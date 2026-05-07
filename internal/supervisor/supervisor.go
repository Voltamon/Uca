package supervisor

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/tidy"
	"github.com/Voltamon/Uca/internal/env"
	"github.com/Voltamon/Uca/internal/runtime"
	"github.com/Voltamon/Uca/internal/auth"
)

var defaultRole string
var logStore *LogStore
var serverLogger *TaggedLogger
var viteLogger *TaggedLogger
var agentLogger *TaggedLogger

var serverProcess *exec.Cmd
var agentProcess *exec.Cmd
var currentAiPort string

func Start() error {
	defaultRole = auth.DefaultRole

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

	// Check for missing keys
	err = checkMissingKeys(cfg)
	if err != nil {
		return err
	}

	aiPort := fmt.Sprintf("%d", cfg.App.Port.AI)
	currentAiPort = aiPort
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

	return runAll(cfg, aiPort)
}

func checkMissingKeys(cfg *config.Config) error {
	missing := []string{}
	for _, key := range cfg.App.Keys {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	fmt.Printf("\n[UCA] Some required keys are missing from your .env file.\n")
	
	for _, key := range missing {
		fmt.Printf("Please enter value for %s: ", key)
		var val string
		fmt.Scanln(&val)
		if val != "" {
			os.Setenv(key, val)
			// Append to .env
			f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("%s=%s\n", key, val))
				f.Close()
			}
		}
	}

	return nil
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

func runAll(cfg *config.Config, aiPort string) error {
	logStore = NewLogStore(1000)

	logFile, err := ensureLogDir()
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	serverLogger = NewTaggedLogger("server", color.New(color.FgCyan), logStore, logFile)
	viteLogger = NewTaggedLogger("vite", color.New(color.FgYellow), logStore, logFile)
	agentLogger = NewTaggedLogger("agent", color.New(color.FgGreen), logStore, logFile)

	done := make(chan error, 3)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 1. Start Backend Server (Always needed for API or Static Files)
	serverProcess = exec.Command(".uca/server")
	serverProcess.Env = append(os.Environ(), "UCA_DEFAULT_ROLE="+defaultRole)
	serverProcess.Stdout = serverLogger
	serverProcess.Stderr = serverLogger
	err = serverProcess.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	go func() {
		err := serverProcess.Wait()
		if err != nil && err.Error() == "signal: killed" {
			return
		}
		done <- err
	}()
	go watchServices(done)

	// 2. Start Vite (Only if pages exist)
	var viteProcess *exec.Cmd
	if len(cfg.Pages) > 0 {
		absNodeBin, err := filepath.Abs(runtime.NodeBin())
		if err != nil {
			return fmt.Errorf("failed to resolve node path: %w", err)
		}

		viteProcess = exec.Command(absNodeBin, "node_modules/.bin/vite")
		viteProcess.Dir = ".uca"
		viteProcess.Stdout = viteLogger
		viteProcess.Stderr = viteLogger
		err = viteProcess.Start()
		if err != nil {
			fmt.Printf("[UCA] Warning: failed to start vite: %v\n", err)
		} else {
			go func() { done <- viteProcess.Wait() }()
		}
	}

	// 3. Start Agent Server (Only if agents exist)
	if len(cfg.Agents) > 0 {
		agentProcess, err = startAgent(aiPort)
		if err != nil {
			fmt.Printf("[UCA] Warning: failed to start agent: %v\n", err)
		} else {
			go watchAgent(agentProcess, aiPort, done)
		}
	}

	select {
	case <-quit:
		fmt.Println("\nShutting down...")
		if serverProcess != nil { serverProcess.Process.Kill() }
		if viteProcess != nil { viteProcess.Process.Kill() }
		if agentProcess != nil { agentProcess.Process.Kill() }
		return nil
	case err := <-done:
		if serverProcess != nil { serverProcess.Process.Kill() }
		if viteProcess != nil { viteProcess.Process.Kill() }
		if agentProcess != nil { agentProcess.Process.Kill() }
		return err
	}
}

func restartAgent() {
	if agentProcess != nil && agentProcess.Process != nil {
		agentProcess.Process.Kill()
	}
	// The watchAgent loop will handle the restart
}

func startAgent(aiPort string) (*exec.Cmd, error) {
	absUca, _ := filepath.Abs(".uca")
	pycache := filepath.Join(absUca, "pycache")
	os.MkdirAll(pycache, 0755)

	cmd := exec.Command("venv/bin/python3", "venv/server.py")
	cmd.Dir = ".uca"
	cmd.Env = append(os.Environ(), 
		"AI_PORT="+aiPort,
		"PYTHONPYCACHEPREFIX="+pycache,
		"PYTHONDONTWRITEBYTECODE=1",
	)
	cmd.Stdout = agentLogger
	cmd.Stderr = agentLogger
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
		time.Sleep(time.Second)
	}

	cmd.Process.Kill()
	return nil, fmt.Errorf("agent failed to start within 5 seconds")
}

func watchAgent(cmd *exec.Cmd, aiPort string, done chan error) {
	for {
		err := cmd.Wait()
		if err != nil && err.Error() != "signal: killed" {
			fmt.Println("Agent crashed, restarting...")
		}

		time.Sleep(time.Second)

		newCmd, err := startAgent(aiPort)
		if err != nil {
			done <- fmt.Errorf("failed to restart agent: %w", err)
			return
		}
		agentProcess = newCmd
		cmd = newCmd
		fmt.Println("Agent restarted")
	}
}
