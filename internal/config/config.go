package config

type Config struct {
    App      AppConfig       `yaml:"app"`
    Pages    []PageConfig    `yaml:"pages,omitempty"`
    Services []ServiceConfig `yaml:"services,omitempty"`
    Agents   []AgentConfig   `yaml:"agents,omitempty"`
}

type AppConfig struct {
    Name    string   `yaml:"name"`
    Version string   `yaml:"version"`
    Icon    string   `yaml:"icon,omitempty"`
    Keys    []string `yaml:"keys,omitempty"`
    Port    PortConfig `yaml:"port,omitempty"`
}

type PortConfig struct {
    Frontend int `yaml:"frontend,omitempty"`
    Backend  int `yaml:"backend,omitempty"`
    AI       int `yaml:"ai,omitempty"`
}

type PageConfig struct {
    Name   string `yaml:"name"`
    Route  string `yaml:"route"`
    Role   string `yaml:"role,omitempty"`
    Render string `yaml:"render,omitempty"`
    Loader string `yaml:"loader,omitempty"`
}

type ServiceConfig struct {
    Name    string            `yaml:"name"`
    Desc    string            `yaml:"desc,omitempty"`
    Methods []string          `yaml:"methods"`
    Schema  map[string]string `yaml:"schema,omitempty"`
    Actions []string          `yaml:"actions,omitempty"`
}

type AgentConfig struct {
    Name    string   `yaml:"name"`
    Model   string   `yaml:"model"`
    Tools   []string `yaml:"tools,omitempty"`
    Context string   `yaml:"context,omitempty"`
    Timeout int      `yaml:"timeout,omitempty"`
}
