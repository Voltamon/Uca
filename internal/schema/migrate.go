package schema

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ApplyChanges(changes []Change) (bool, error) {
	if len(changes) == 0 {
		return true, nil
	}

	var safeChanges []Change
	var destructiveChanges []Change

	for _, c := range changes {
		if c.Destructive {
			destructiveChanges = append(destructiveChanges, c)
		} else {
			safeChanges = append(safeChanges, c)
		}
	}

	if len(safeChanges) > 0 {
		fmt.Println("Applying schema changes:")
		for _, c := range safeChanges {
			fmt.Println(" +", c.Describe())
		}
	}

	if len(destructiveChanges) > 0 {
		fmt.Println("\nDestructive schema changes detected:")
		for _, c := range destructiveChanges {
			fmt.Println(" !", c.Describe())
		}

		fmt.Print("\nThese changes may cause data loss. Are you sure? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("Schema changes cancelled.")
			return false, nil
		}
	}

	return true, nil
}
