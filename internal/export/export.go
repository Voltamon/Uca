package export

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/tidy"
	"github.com/klauspost/compress/zstd"
)

func Run(includeEnv bool) error {
	fmt.Println("Running tidy...")
	cfg, err := tidy.Run()
	if err != nil {
		return fmt.Errorf("tidy failed: %w", err)
	}

	fmt.Println("Building frontend...")
	err = buildFrontend()
	if err != nil {
		return fmt.Errorf("frontend build failed: %w", err)
	}

	fmt.Println("Building server...")
	err = buildServer()
	if err != nil {
		return fmt.Errorf("server build failed: %w", err)
	}

	fmt.Println("Generating .env.example...")
	err = generateEnvExample(cfg)
	if err != nil {
		return fmt.Errorf("failed to generate .env.example: %w", err)
	}

	fmt.Println("Creating export artifact...")
	err = createArtifact(includeEnv)
	if err != nil {
		return fmt.Errorf("failed to create artifact: %w", err)
	}

	fmt.Println("Export complete: app.tar.zst")
	return nil
}

func generateEnvExample(cfg *config.Config) error {
	var sb strings.Builder
	for _, key := range cfg.App.Keys {
		sb.WriteString(key + "=\n")
	}
	return os.WriteFile(".uca/.env.example", []byte(sb.String()), 0644)
}

func buildFrontend() error {
	absNodeBin, err := filepath.Abs(".uca/runtimes/node/bin/node")
	if err != nil {
		return err
	}

	cmd := exec.Command(absNodeBin, "node_modules/.bin/vite", "build")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return stripCrossOrigin(".uca/.vite/index.html")
}

func stripCrossOrigin(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(data), ` crossorigin`, "")
	return os.WriteFile(path, []byte(content), 0644)
}

func buildServer() error {
	fmt.Println("Resolving dependencies...")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = ".uca"
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	cmd := exec.Command("go", "build", "-o", "server")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createArtifact(includeEnv bool) error {
	outFile, err := os.Create("app.tar.zst")
	if err != nil {
		return err
	}
	defer outFile.Close()

	zw, err := zstd.NewWriter(outFile)
	if err != nil {
		return err
	}
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	items := []struct {
		src  string
		dest string
	}{
		{".uca/server", "app/server"},
		{".uca/venv", "app/venv"},
		{".uca/.vite", "app/dist"},
		{".uca/schema.json", "app/schema.json"},
		{".uca/.env.example", "app/.env.example"},
	}

	if includeEnv {
		if _, err := os.Stat(".uca/.env"); err == nil {
			items = append(items, struct{ src, dest string }{".uca/.env", "app/.env"})
		}
	}

	for _, item := range items {
		info, err := os.Stat(item.src)
		if err != nil {
			return fmt.Errorf("missing %s: %w", item.src, err)
		}

		if info.IsDir() {
			err = addDir(tw, item.src, item.dest)
		} else {
			err = addFile(tw, item.src, item.dest)
		}

		if err != nil {
			return fmt.Errorf("failed to add %s: %w", item.src, err)
		}

		fmt.Println("Packed:", item.dest)
	}

	return nil
}

func addFile(tw *tar.Writer, src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: dest,
		Mode: int64(info.Mode()),
		Size: info.Size(),
	}

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, f)
	return err
}

func addDir(tw *tar.Writer, src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeDir,
				Name:     targetPath + "/",
				Mode:     int64(info.Mode()),
			})
		}

		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeSymlink,
				Name:     targetPath,
				Linkname: linkTarget,
			})
		}

		if strings.Contains(path, "__pycache__") {
			return nil
		}

		return addFile(tw, path, targetPath)
	})
}
