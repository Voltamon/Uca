package supervisor

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/Voltamon/Uca/internal/tidy"
	"github.com/Voltamon/Uca/internal/env"
	"github.com/Voltamon/Uca/internal/runtime"
)

func Start() error {
	envVars, err := env.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	fmt.Println("Running tidy...")
	cfg, err := tidy.Run()
	if err != nil {
		return fmt.Errorf("tidy failed: %w", err)
	}

	aiPort := fmt.Sprintf("%d", cfg.App.Port.AI)
	backendPort := fmt.Sprintf("%d", cfg.App.Port.Backend)

	fmt.Println("Resolving dependencies...")
	err = runGoModTidy()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	fmt.Println("Building project...")
	err = buildProject(backendPort)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return runAll(aiPort)
}

func runGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildProject(backendPort string) error {
	cmd := exec.Command("go", "build", "-o", "server")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runAll(aiPort string) error {
	serverCmd := exec.Command("./server")
	serverCmd.Dir = ".uca"
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	viteCmd := exec.Command("../"+runtime.NodeBin(), "node_modules/.bin/vite")
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

	agentCmd, err := startAgent(aiPort)
	if err != nil {
		serverCmd.Process.Kill()
		viteCmd.Process.Kill()
		return fmt.Errorf("failed to start agent: %w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 3)

	go func() { done <- serverCmd.Wait() }()
	go func() { done <- viteCmd.Wait() }()
	go watchAgent(agentCmd, aiPort, done)

	select {
	case <-quit:
		fmt.Println("\nShutting down...")
		serverCmd.Process.Kill()
		viteCmd.Process.Kill()
		agentCmd.Process.Kill()
		return nil
	case err := <-done:
		serverCmd.Process.Kill()
		viteCmd.Process.Kill()
		agentCmd.Process.Kill()
		return err
	}
}

func startAgent(aiPort string) (*exec.Cmd, error) {
	cmd := exec.Command("venv/bin/python3", "server.py")
	cmd.Dir = ".uca"
	cmd.Env = append(os.Environ(), "AI_PORT="+aiPort)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	fmt.Println("Waiting for agent to be ready...")
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://127.0.0.1:" + aiPort + "/health")
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("Agent ready on http://127.0.0.1:" + aiPort)
			return cmd, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	cmd.Process.Kill()
	return nil, fmt.Errorf("agent failed to start within 5 seconds")
}

func watchAgent(cmd *exec.Cmd, aiPort string, done chan error) {
	for {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("Agent crashed, restarting...")
		}

		time.Sleep(time.Second)

		cmd := exec.Command("venv/bin/python3", "server.py")
		cmd.Dir = ".uca"
		cmd.Env = append(os.Environ(), "AI_PORT="+aiPort)
		err = cmd.Start()
		if err != nil {
			done <- fmt.Errorf("failed to restart agent: %w", err)
			return
		}

		fmt.Println("Agent restarted")
	}
}
