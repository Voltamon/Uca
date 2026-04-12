package runtime

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func EnsureNode() error {
	if _, err := os.Stat(NodeBin()); err == nil {
		return nil
	}

	platform, err := Detect()
	if err != nil {
		return err
	}

	url := nodeURL(platform)
	fmt.Printf("Downloading Node.js %s...\n", NodeVersion)

	archive, err := download(url)
	if err != nil {
		return fmt.Errorf("failed to download Node.js: %w", err)
	}
	defer os.Remove(archive)

	fmt.Println("Extracting Node.js...")
	err = extractNodeTarGz(archive, NodeDir())
	if err != nil {
		return fmt.Errorf("failed to extract Node.js: %w", err)
	}

	fmt.Printf("Node.js %s ready\n", NodeVersion)
	return nil
}

func nodeURL(p Platform) string {
	os := p.OS
	if os == "windows" {
		os = "win"
	}
	return fmt.Sprintf(
		"https://nodejs.org/dist/v%s/node-v%s-%s-%s.tar.gz",
		NodeVersion, NodeVersion, os, p.Arch,
	)
}

func download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "uca-download-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return "", err
	}

	return tmp.Name(), nil
}

func extractNodeTarGz(archive string, destDir string) error {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parts := strings.SplitN(header.Name, "/", 2)
		if len(parts) < 2 {
			continue
		}
		relPath := parts[1]
		if relPath == "" {
			continue
		}

		target := filepath.Join(destDir, relPath)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			io.Copy(out, tr)
			out.Close()
			os.Chmod(target, os.FileMode(header.Mode))
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, target)
		}
	}

	return nil
}
