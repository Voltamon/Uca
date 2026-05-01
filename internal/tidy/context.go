package tidy

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateContextPackage() error {
	err := scaffold.CopyTemplate("uca/context/context.go", ".uca/uca/context/context.go", scaffold.TemplateVars{})
	if err != nil {
		return fmt.Errorf("failed to generate context package: %w", err)
	}

	fmt.Println("Generated: .uca/uca/context/context.go")
	return nil
}
