package typegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type TSInterface struct {
	Name   string
	Fields []TSField
}

type TSField struct {
	Name     string
	Type     string
	Optional bool
}

func GenerateTypes(servicesDir string, outputDir string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		srcPath := filepath.Join(servicesDir, entry.Name())
		interfaces, err := extractInterfaces(srcPath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
		}

		if len(interfaces) == 0 {
			continue
		}

		baseName := strings.TrimSuffix(entry.Name(), ".go")
		outPath := filepath.Join(outputDir, baseName+".ts")

		err = writeTypeScript(outPath, interfaces)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", outPath, err)
		}

		fmt.Println("Generated:", outPath)
	}

	return nil
}

func extractInterfaces(filePath string) ([]TSInterface, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var interfaces []TSInterface

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			iface := TSInterface{Name: typeSpec.Name.Name}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				tsField := TSField{
					Name: extractJSONName(field),
					Type: goTypeToTS(field.Type),
				}

				if strings.HasPrefix(fmt.Sprintf("%s", field.Type), "*") {
					tsField.Optional = true
				}

				iface.Fields = append(iface.Fields, tsField)
			}

			interfaces = append(interfaces, iface)
		}

		return true
	})

	return interfaces, nil
}

func extractJSONName(field *ast.Field) string {
	if field.Tag != nil {
		tag := field.Tag.Value
		tag = strings.Trim(tag, "`")
		for _, part := range strings.Split(tag, " ") {
			if strings.HasPrefix(part, "json:") {
				name := strings.TrimPrefix(part, "json:")
				name = strings.Trim(name, "\"")
				name = strings.Split(name, ",")[0]
				if name != "-" {
					return name
				}
			}
		}
	}

	if len(field.Names) > 0 {
		return field.Names[0].Name
	}

	return ""
}

func goTypeToTS(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return identToTS(t.Name)
	case *ast.StarExpr:
		return goTypeToTS(t.X) + " | null"
	case *ast.ArrayType:
		return goTypeToTS(t.Elt) + "[]"
	case *ast.MapType:
		return "Record<" + goTypeToTS(t.Key) + ", " + goTypeToTS(t.Value) + ">"
	case *ast.SelectorExpr:
		return "any"
	default:
		return "any"
	}
}

func identToTS(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "byte":
		return "number"
	default:
		return goType
	}
}

func writeTypeScript(outPath string, interfaces []TSInterface) error {
	var sb strings.Builder

	for i, iface := range interfaces {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("export interface %s {\n", iface.Name))
		for _, field := range iface.Fields {
			if field.Name == "" || field.Name == "-" {
				continue
			}
			optional := ""
			if field.Optional {
				optional = "?"
			}
			sb.WriteString(fmt.Sprintf("    %s%s: %s\n", field.Name, optional, field.Type))
		}
		sb.WriteString("}\n")
	}

	return os.WriteFile(outPath, []byte(sb.String()), 0644)
}
