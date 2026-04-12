package runtime

import "fmt"

func EnsureAll() error {
	fmt.Println("Checking runtimes...")

	err := EnsureNode()
	if err != nil {
		return fmt.Errorf("node runtime error: %w", err)
	}

	err = EnsurePython()
	if err != nil {
		return fmt.Errorf("python runtime error: %w", err)
	}

	return nil
}
