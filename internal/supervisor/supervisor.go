package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/Voltamon/Uca/internal/tidy"
)

func Start() error {
	fmt.Println("Running tidy...")
	err := tidy.Run()
	if err != nil {
		return fmt.Errorf("tidy failed: %w", err)
	}

	fmt.Println("Resolving dependencies...")
	err = runGoModTidy()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	fmt.Println("Building project...")
	err = buildProject()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return runAll()
}

func runGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildProject() error {
	cmd := exec.Command("go", "build", "-o", "server")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runAll() error {
	serverCmd := exec.Command("./server")
	serverCmd.Dir = ".uca"
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	viteCmd := exec.Command("npm", "run", "dev")
	viteCmd.Dir = ".uca"
	viteCmd.Stdout = os.Stdout
	viteCmd.Stderr = os.Stderr

	err := serverCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	err = viteCmd.Start()
	if err != nil {
		serverCmd.Process.Kill()
		return fmt.Errorf("failed to start vite: %w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 2)

	go func() {
		done <- serverCmd.Wait()
	}()

	go func() {
		done <- viteCmd.Wait()
	}()

	select {
	case <-quit:
		fmt.Println("\nShutting down...")
		serverCmd.Process.Kill()
		viteCmd.Process.Kill()
		return nil
	case err := <-done:
		serverCmd.Process.Kill()
		viteCmd.Process.Kill()
		return err
	}
}
