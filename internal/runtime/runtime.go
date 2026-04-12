package runtime

import (
	"fmt"
	"runtime"
)

const (
	NodeVersion   = "20.11.0"
	PythonVersion = "3.12.3"
)

type Platform struct {
	OS   string
	Arch string
}

func Detect() (Platform, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	default:
		return Platform{}, fmt.Errorf("unsupported architecture: %s", arch)
	}

	switch os {
	case "linux", "darwin", "windows":
	default:
		return Platform{}, fmt.Errorf("unsupported OS: %s", os)
	}

	return Platform{OS: os, Arch: arch}, nil
}

func NodeDir() string {
	return ".uca/runtimes/node"
}

func PythonDir() string {
	return ".uca/runtimes/python"
}

func NodeBin() string {
	return ".uca/runtimes/node/bin/node"
}

func NpmBin() string {
	return ".uca/runtimes/node/bin/npm"
}

func PythonBin() string {
	return ".uca/runtimes/python/bin/python3"
}

func PipBin() string {
	return ".uca/runtimes/python/bin/pip3"
}
