package scaffold

import (
    "os"

    "gopkg.in/yaml.v3"
    "github.com/Voltamon/Uca/internal/config"
)

func generateManifest(appName string) error {
    cfg := config.Config{
        App: config.AppConfig{
            Name:    appName,
            Version: "1.0",
            Keys:    []string{"GITHUB_PAT_TOKEN"},
        },
        Pages: []config.PageConfig{
            {
                Name:   "Welcome",
                Route:  "/",
                Loader: "User.GET",
            },
            {
                Name:   "Chat",
                Route:  "/chat",
                Render: "ssr",
                Loader: "History.GetChatHistory",
            },
        },
        Services: []config.ServiceConfig{
            {
                Name:    "User",
                Desc:    "Manages user identity and onboarding",
                Methods: []string{"GET", "POST"},
                Schema: map[string]string{
                    "username": "string | required",
                },
            },
            {
                Name:    "History",
                Desc:    "Persistence layer for agent interactions",
                Methods: []string{"GET", "POST"},
                Schema: map[string]string{
                    "sender":  "select:user,agent | required",
                    "content": "string | required",
                },
                Actions: []string{"GetChatHistory"},
            },
        },
        Agents: []config.AgentConfig{
            {
                Name:    "Assistant",
                Model:   "openai/gpt-4o",
                Timeout: 30,
                Context: "History.GetChatHistory",
                Tools:   []string{"History.GetChatHistory", "User.GET", "DuckDuckGo"},
            },
        },
    }

    data, err := yaml.Marshal(&cfg)
    if err != nil {
        return err
    }

    return os.WriteFile(appName+"/uca.yaml", data, 0644)
}
