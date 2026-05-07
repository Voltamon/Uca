package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/Voltamon/Uca/internal/tidy"
)

func watchServices(done chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		done <- fmt.Errorf("failed to create watcher: %w", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add("services")
	if err != nil {
		done <- fmt.Errorf("failed to watch services: %w", err)
		return
	}

	err = watcher.Add("agents")
	if err != nil {
		done <- fmt.Errorf("failed to watch agents: %w", err)
		return
	}

	err = watcher.Add("uca.yaml")
	if err != nil {
		fmt.Printf("Warning: failed to watch uca.yaml: %v\n", err)
	}

	fmt.Println("Watching services/, agents/ and uca.yaml for changes...")

	var debounce *time.Timer
	var lastFile string

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !strings.HasSuffix(event.Name, ".go") && !strings.HasSuffix(event.Name, ".py") && !strings.HasSuffix(event.Name, "uca.yaml") {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			if event.Name == lastFile {
				continue
			}

			if debounce != nil {
				debounce.Stop()
			}

			lastFile = event.Name
			fileName := event.Name
			debounce = time.AfterFunc(time.Second, func() {
				lastFile = ""
				fmt.Printf("\nChange detected in %s — refreshing...\n", filepath.Base(fileName))
				
				if strings.HasSuffix(fileName, "uca.yaml") {
					_, err := tidy.Run()
					if err != nil {
						fmt.Printf("[UCA] Tidy failed: %v\n", err)
						return
					}
					rebuildAndRestart()
					restartAgent()
				} else if strings.HasSuffix(fileName, ".go") {
					rebuildAndRestart()
				} else {
					restartAgent()
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watcher error:", err)
		}
	}
}

func rebuildAndRestart() {
	if serverProcess != nil && serverProcess.Process != nil {
		serverProcess.Process.Kill()
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("Building...")
	build := exec.Command("go", "build", "-o", "server")
	build.Dir = ".uca"
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		fmt.Println("Build failed:", err)
		return
	}

	fmt.Println("Restarting server...")
	serverProcess = exec.Command("./server")
	serverProcess.Dir = ".uca"
	serverProcess.Env = append(os.Environ(), "UCA_DEFAULT_ROLE="+defaultRole)
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr

	if err := serverProcess.Start(); err != nil {
		fmt.Println("Failed to restart server:", err)
		return
	}

	go func() {
		err := serverProcess.Wait()
		if err != nil && err.Error() == "signal: killed" {
			return
		}
	}()

	fmt.Println("Server restarted")
}
