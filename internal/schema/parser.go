package schema

import (
	"strings"

	"github.com/Voltamon/Uca/internal/config"
)

func ParseFromConfig(cfg *config.Config) Schema {
	var s Schema

	for _, service := range cfg.Services {
		if len(service.Schema) == 0 {
			continue
		}

		collection := Collection{
			Name: service.Name,
		}

		for fieldName, fieldDef := range service.Schema {
			field := parseField(fieldName, fieldDef)
			collection.Fields = append(collection.Fields, field)
		}

		s.Collections = append(s.Collections, collection)
	}

	return s
}

func parseField(name string, def string) Field {
	field := Field{Name: name}

	parts := strings.SplitN(def, "|", 2)
	typePart := strings.TrimSpace(parts[0])

	if len(parts) == 2 {
		constraint := strings.TrimSpace(parts[1])
		if constraint == "required" {
			field.Required = true
		}
	}

	if strings.HasPrefix(typePart, "select:") {
		field.Type = "select"
		optionStr := strings.TrimPrefix(typePart, "select:")
		for _, opt := range strings.Split(optionStr, ",") {
			field.Options = append(field.Options, strings.TrimSpace(opt))
		}
	} else {
		field.Type = typePart
	}

	return field
}
