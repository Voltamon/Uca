package schema

import "fmt"

type ChangeType string

const (
	ChangeAdd        ChangeType = "ADD"
	ChangeRemove     ChangeType = "REMOVE"
	ChangeTypeChange ChangeType = "TYPE_CHANGE"
	ChangeNewCollection ChangeType = "NEW_COLLECTION"
)

type Change struct {
	Collection string
	Field      string
	ChangeType ChangeType
	From       string
	To         string
	Destructive bool
}

func (c Change) Describe() string {
	switch c.ChangeType {
	case ChangeNewCollection:
		return fmt.Sprintf("CREATE collection %q", c.Collection)
	case ChangeAdd:
		return fmt.Sprintf("ADD field %q (%s) to %q", c.Field, c.To, c.Collection)
	case ChangeRemove:
		return fmt.Sprintf("REMOVE field %q from %q", c.Field, c.Collection)
	case ChangeTypeChange:
		return fmt.Sprintf("CHANGE field %q in %q from %s to %s", c.Field, c.Collection, c.From, c.To)
	}
	return ""
}

func Diff(current Schema, desired Schema) []Change {
	var changes []Change

	currentMap := make(map[string]Collection)
	for _, c := range current.Collections {
		currentMap[c.Name] = c
	}

	for _, desired := range desired.Collections {
		current, exists := currentMap[desired.Name]
		if !exists {
			changes = append(changes, Change{
				Collection:  desired.Name,
				ChangeType:  ChangeNewCollection,
				Destructive: false,
			})
			continue
		}

		currentFields := make(map[string]Field)
		for _, f := range current.Fields {
			currentFields[f.Name] = f
		}

		desiredFields := make(map[string]Field)
		for _, f := range desired.Fields {
			desiredFields[f.Name] = f
		}

		for _, df := range desired.Fields {
			cf, exists := currentFields[df.Name]
			if !exists {
				changes = append(changes, Change{
					Collection:  desired.Name,
					Field:       df.Name,
					ChangeType:  ChangeAdd,
					To:          df.Type,
					Destructive: false,
				})
				continue
			}

			if cf.Type != df.Type {
				changes = append(changes, Change{
					Collection:  desired.Name,
					Field:       df.Name,
					ChangeType:  ChangeTypeChange,
					From:        cf.Type,
					To:          df.Type,
					Destructive: true,
				})
			}
		}

		for _, cf := range current.Fields {
			if _, exists := desiredFields[cf.Name]; !exists {
				changes = append(changes, Change{
					Collection:  desired.Name,
					Field:       cf.Name,
					ChangeType:  ChangeRemove,
					From:        cf.Type,
					Destructive: true,
				})
			}
		}
	}

	return changes
}
