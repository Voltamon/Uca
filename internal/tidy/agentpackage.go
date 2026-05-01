package tidy

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateAgentPythonPackage() error {
	err := scaffold.CopyTemplate("uca/python/__init__.py", ".uca/uca/__init__.py", scaffold.TemplateVars{})
	if err != nil {
		return fmt.Errorf("failed to generate agent python package: %w", err)
	}

	fmt.Println("Generated: .uca/uca/__init__.py")
	return nil
}
