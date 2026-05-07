package deps

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/Voltamon/Uca/internal/runtime"
)

func SyncFromFile(path string) error {
	fmt.Printf("Syncing dependencies from %s...\n", path)
	d, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load deps file: %w", err)
	}

	for pkg := range d.Pages {
		if err := AddPagesDep(pkg); err != nil {
			return err
		}
	}
	for pkg := range d.Services {
		if err := AddServicesDep(pkg); err != nil {
			return err
		}
	}
	for pkg := range d.Agents {
		if err := AddAgentsDep(pkg); err != nil {
			return err
		}
	}
	return nil
}

func AddPagesDep(pkg string) error {
	fmt.Printf("Installing %s for pages...\n", pkg)
	cmd := exec.Command(runtime.NodeBin(), runtime.NpmBin(), "install", "--prefix", ".uca", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AddServicesDep(pkg string) error {
	fmt.Printf("Installing %s for services...\n", pkg)
	cmd := exec.Command("go", "get", pkg)
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AddAgentsDep(pkg string) error {
	fmt.Printf("Installing %s for agents...\n", pkg)
	cmd := exec.Command(".uca/venv/bin/pip", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
