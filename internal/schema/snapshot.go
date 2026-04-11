package schema

import (
	"encoding/json"
	"os"
)

const snapshotPath = ".uca/schema.json"

func LoadSnapshot() (Schema, error) {
	var s Schema

	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Schema{}, nil
		}
		return s, err
	}

	err = json.Unmarshal(data, &s)
	return s, err
}

func SaveSnapshot(s Schema) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(snapshotPath, data, 0644)
}
