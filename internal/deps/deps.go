package deps

import (
	"encoding/json"
	"os"
)

type Deps struct {
	Pages    map[string]string `json:"pages"`
	Services map[string]string `json:"services"`
	Agents   map[string]string `json:"agents"`
}

func Load() (Deps, error) {
	var d Deps
	data, err := os.ReadFile("deps.json")
	if err != nil {
		if os.IsNotExist(err) {
			return Deps{
				Pages:    make(map[string]string),
				Services: make(map[string]string),
				Agents:   make(map[string]string),
			}, nil
		}
		return d, err
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, err
	}
	if d.Pages == nil { d.Pages = make(map[string]string) }
	if d.Services == nil { d.Services = make(map[string]string) }
	if d.Agents == nil { d.Agents = make(map[string]string) }
	return d, nil
}

func Save(d Deps) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("deps.json", data, 0644)
}
