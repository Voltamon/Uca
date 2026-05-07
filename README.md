# 📐 Uca

**The Hermetic, Polyglot Microframework for AI-Native Applications.**

Uca (Unified Component Architecture) is a specialized framework designed to eliminate the friction of building full-stack applications that require high-performance backends (**Go**), reactive frontends (**TypeScript/Preact**), and sophisticated AI logic (**Python**).

---

## 🌟 Why Uca?

Traditional full-stack development requires managing multiple toolchains, environment variables, and type definitions across languages. Uca automates this by providing a **hermetic build sandbox** that manages its own runtimes and synchronizes your architecture via a single manifest.

- **Zero-Config Runtimes:** Uca downloads and manages Node.js and Python internally. No more `nvm` or `pyenv` conflicts.
- **Manifest-as-Code:** Your `uca.yaml` defines your database, API, and AI tools in one place.
- **Type-Safe Polyglotism:** Types defined in Go are automatically available in TypeScript. Python agents automatically understand your Go services.
- **AI-First:** Agents aren't an afterthought; they are first-class citizens with automatic tool-spec generation.

---

## 🚀 Quick Start (60 Seconds)

### 1. Install & Initialize
```bash
# Assuming you have the 'uca' binary in your path
uca init my-project
cd my-project
```

### 2. Define your Schema
Edit `uca.yaml` to define a service:
```yaml
services:
  - name: Task
    methods: [GET, POST]
    schema:
      title: string | required
      done: bool
```

### 3. Tidy & Run
```bash
uca tidy  # Generates types, DB migrations, and tests
uca dev   # Starts Go, Python, and Vite servers
```
Your app is now live at `http://localhost:5173`.

---

## 🏗 Framework Anatomy

Uca splits your project into a **User Space** (which you own) and an **Engine Room** (managed by Uca).

```text
.
├── uca.yaml           # The Source of Truth
├── pages/             # Frontend (TypeScript + Preact)
├── services/          # Backend (Go + PocketBase)
├── agents/            # AI (Python + LiteLLM)
├── assets/            # Global Static Assets
└── .uca/              # [MANAGED] Build sandbox, Runtimes, and DB
```

### The "Tidy" Lifecycle
Whenever you run `uca tidy`, the framework performs a full reconciliation:
1.  **Sync Manifest:** Updates the internal registry of services and agents.
2.  **Schema Migration:** Automatically updates the SQLite/PocketBase schema to match your `uca.yaml`.
3.  **Type Projection:** Exports Go structs as TypeScript interfaces to `.uca/types/`.
4.  **Auto-Testing:** Generates "Smoke Tests" for every page, service, and agent in `.uca/tests/autogen/`.

---

## 🧪 Testing Philosophy

Uca believes testing should be **Invisible by Default, Extensible by Choice.**

### 1. Auto-Generated Baseline
Every time you run `uca test`, the framework generates and executes:
- **Mount Tests:** Ensures every Preact component in `pages/` renders without crashing.
- **Wiring Tests:** Verifies all API endpoints are correctly registered in the Go backend.
- **Import Tests:** Confirms Python agents load their models and tools correctly.

### 2. Custom Extensions
Simply add `*_test.go`, `*.test.tsx`, or `*_test.py` in your user-space directories. Uca's test runner will merge them into the execution flow.

---

## 🔗 Polyglot Communication

Uca provides "Magic Imports" that bridge the language gap:

### Frontend → Backend (TS to Go)
```tsx
import { Task } from "uca/srv" // Auto-generated based on uca.yaml

// Reactive signal: automatically updates the UI on fetch
useEffect(() => { Task.GET.fetch() }, [])
```

### Agent → Backend (Python to Go)
```python
from uca.srv import Task # Agents can call Go logic directly

def check_tasks():
    pending = Task.GET()
    return f"You have {len(pending)} tasks left."
```

---

## ⌨️ CLI Command Reference

| Command | Action |
| :--- | :--- |
| `init <name>` | Create a new project from a template. |
| `tidy` | Reconcile manifest, generate types, and update DB. |
| `dev` | Start the unified development server with hot-reload. |
| `test` | Run the integrated auto-generated and custom test suite. |
| `env` | Interactively manage environment variables and API keys. |
| `export` | Package the application for production deployment. |

---

## 🛠 System Requirements

- **Go 1.22+** (for building the core engine)
- **Git** (for version control and dependency management)
- *Note: Node.js and Python are managed automatically by Uca.*
