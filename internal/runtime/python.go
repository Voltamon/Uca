package runtime

import (
	"fmt"
	"os"
)

func EnsurePython() error {
	if _, err := os.Stat(PythonBin()); err == nil {
		return nil
	}

	platform, err := Detect()
	if err != nil {
		return err
	}

	url := pythonURL(platform)
	fmt.Printf("Downloading Python %s...\n", PythonVersion)

	archive, err := download(url)
	if err != nil {
		return fmt.Errorf("failed to download Python: %w", err)
	}
	defer os.Remove(archive)

	fmt.Println("Extracting Python...")
	err = extractNodeTarGz(archive, PythonDir())
	if err != nil {
		return fmt.Errorf("failed to extract Python: %w", err)
	}

	fmt.Printf("Python %s ready\n", PythonVersion)
	return nil
}

func pythonURL(p Platform) string {
	arch := p.Arch
	if arch == "x64" {
		arch = "x86_64"
	}

	os := p.OS
	if os == "darwin" {
		os = "apple-darwin"
	} else if os == "linux" {
		os = "unknown-linux-gnu"
	} else if os == "windows" {
		os = "pc-windows-msvc"
	}

	return fmt.Sprintf(
		"https://github.com/indygreg/python-build-standalone/releases/download/20240415/cpython-%s+20240415-%s-%s-install_only.tar.gz",
		PythonVersion, arch, os,
	)
}
