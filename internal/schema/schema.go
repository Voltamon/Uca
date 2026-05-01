package schema

type Schema struct {
	Collections []Collection `json:"collections"`
	Roles       []string     `json:"roles,omitempty"`
	DefaultRole string       `json:"default_role,omitempty"`
}

type Collection struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
}
