package tidy

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateUIPackage() error {
	err := scaffold.CopyTemplate("uca/ui/index.ts", ".uca/ui/index.ts", scaffold.TemplateVars{})
	if err != nil {
		return fmt.Errorf("failed to generate ui package: %w", err)
	}

	fmt.Println("Generated: .uca/ui/index.ts")
	return nil
}
