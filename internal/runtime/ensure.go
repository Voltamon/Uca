package runtime

import (
	"fmt"
	"github.com/Voltamon/Uca/internal/config"
)

func EnsureAll(cfg *config.Config) error {
	fmt.Println("Checking runtimes...")

	if len(cfg.Pages) > 0 {
		err := EnsureNode()
		if err != nil {
			return fmt.Errorf("node runtime error: %w", err)
		}
	}

	if len(cfg.Agents) > 0 {
		err := EnsurePython()
		if err != nil {
			return fmt.Errorf("python runtime error: %w", err)
		}
	}

	return nil
}
