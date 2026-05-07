# Uca - A Hermetic, Polyglot Microframework

Uca is a modern web framework designed for building full-stack applications with AI agents. It abstracts away the complexity of managing multiple languages and environments, allowing you to focus on behavior and business logic while the framework handles the infrastructure "magic."

---

## 🚀 Core Philosophy

*   **Manifest-Driven:** Your `uca.yaml` is the single source of truth for your entire application architecture.
*   **Hermetic Runtimes:** Uca manages its own Node.js and Python runtimes within the project directory, ensuring zero-dependency, reproducible builds.
*   **Polyglot by Design:** Write your UI in **TypeScript (Preact)**, your business logic in **Go**, and your AI agents in **Python**.
*   **Zero Boilerplate:** Hidden CRUD layers and reactive Signals mean you only write the code that is unique to your application.

---

## 📁 Project Structure

```text
my-app/
├── uca.yaml          # The application manifest (Source of Truth)
├── .env              # Environment secrets (managed interactively)
├── agents/           # Python AI Agent definitions
├── pages/            # Preact frontend components
├── services/         # Go backend business logic
├── assets/           # Static files (images, audio, etc.) shared across all layers
└── .uca/             # The "Magic" glue layer (Build sandbox, runtimes, DB)
```

---

## 🛠️ The Manifest (`uca.yaml`)

The manifest defines your application's DNA.

```yaml
app:
  name: my-app
  version: "1.0"
  keys: [GITHUB_PAT_TOKEN] # Required secrets
  port:
    frontend: 5173
    backend: 8090
    ai: 8091

services:
  - name: Todo
    methods: [GET, POST, PUT, DELETE]
    schema:
      title: string | required
      completed: bool

agents:
  - name: Assistant
    model: github/gpt-4o
    tools: [Todo.All] # Automatically registers all CRUD tools
```

---

## 💻 Backend Logic (Go)

Uca uses **PocketBase** as its core engine.

### Hidden CRUD
If you define a service in `uca.yaml` but don't write any Go code, Uca automatically provides a high-quality API for that service.

### The Decorator Pattern
To add custom logic (e.g., validation, notifications), simply implement the function in `services/*.go` and use the provided helpers.

```go
func UserPOST(e *context.RequestEvent) error {
    // 1. Custom logic before
    log.Println("New user signing up!")

    // 2. Call the magic helper
    err := uca.DefaultPOST(e, "User")

    // 3. Custom logic after
    return err
}
```

---

## 🤖 AI Agents (Python)

Uca provides a simplified, object-oriented API for Python agents using **LiteLLM**.

### The Agent Class
*   **Auto-Sync:** Changing `agent.model` in Python automatically updates your `uca.yaml`.
*   **Introspective Tools:** Pass any Python function to `agent.tools`, and Uca automatically generates the JSON tool-spec for the AI.

```python
from uca.ai import Agent, Message
from uca.srv import Todo

def get_weather(city: str):
    """Returns the weather."""
    return "Sunny"

agent = Agent(model="gpt-4o", tools=[Todo.All, get_weather])
agent.prompt = f"System: Be helpful.\nUser: {Message}"
```

---

## ⚛️ Frontend (Preact + Signals)

Uca implements a "Zero-Boilerplate" frontend using **Preact Signals**.

### Modular Imports
Import your services and agents directly by name—no more manual `fetch` calls.

```tsx
import { Todo } from "uca/srv"
import { TaskBot } from "uca/ai"

export default function Home() {
    // 1. Reactive Data: Automatically fetches and updates UI
    useEffect(() => { Todo.GET.fetch() }, [])

    // 2. No State Boilerplate: Just render the signal values
    return (
        <div>
            {Todo.GET.data.value.map(t => <p>{t.title}</p>)}
            <button onClick={() => TaskBot.chat("Hello!", console.log)}>
                Chat
            </button>
        </div>
    )
}
```

---

## ⌨️ CLI Commands

*   **`uca init <name>`**: Scaffold a new project.
*   **`uca tidy`**: Synchronize the manifest with the code (generates boilerplate, types, and DB schema).
*   **`uca dev`**: Start the development environment with hot-reloading for Go, Python, and TypeScript.
*   **`uca agent add <name> <model>`**: Interactively add a new agent to your manifest.

---

## 🔒 Secret Management

Uca manages your `.env` file interactively. If a key defined in `uca.yaml` is missing during `uca dev`, the framework will pause and prompt you to enter the value in the terminal, saving it securely for future runs.

---

## 📦 Assets

Any file in the `assets/` directory is served globally under the `/assets/` URL prefix. It is accessible to your frontend components (`<img src="/assets/logo.png" />`), your Go services, and your Python agents.
