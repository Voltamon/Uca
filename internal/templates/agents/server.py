import importlib.util
import json
import os
import sys
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import urlparse, parse_qs

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

def load_env():
    env_path = ".env"
    if os.path.exists(env_path):
        with open(env_path) as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith("#"):
                    if "=" in line:
                        k, v = line.split("=", 1)
                        os.environ[k] = v

load_env()

API_KEY = os.environ.get("GITHUB_PAT_TOKEN", "")
API_BASE = "https://models.inference.ai.azure.com"

# Cache for loaded agents
agents_cache = {}

def get_agent(name):
    # For dev, we might want to reload it if file changed, but for now cache is fine
    # Actually, the supervisor will restart the server if files change.
    if name in agents_cache:
        return agents_cache[name]
    
    # Try relative to root
    path = f"agents/{name}.py"
    if not os.path.exists(path):
        # Try relative to .uca/venv/
        path = os.path.join(os.path.dirname(__file__), "../../agents", f"{name}.py")
        if not os.path.exists(path):
            return None
    
    try:
        spec = importlib.util.spec_from_file_location(name, path)
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        agent = getattr(module, "agent", None)
        if agent:
            agents_cache[name] = agent
        return agent
    except Exception as e:
        print(f"Error loading agent {name}: {e}")
        return None

class AgentHandler(BaseHTTPRequestHandler):

    def log_message(self, format, *args):
        pass

    def do_GET(self):
        parsed = urlparse(self.path)

        if parsed.path == "/health":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({"status": "ok"}).encode())
            return

        if parsed.path == "/chat":
            params = parse_qs(parsed.query)
            agent_name = params.get("agent", [""])[0]
            message = params.get("message", [""])[0]

            if not agent_name:
                self.send_response(400)
                self.end_headers()
                self.wfile.write(b"Missing agent parameter")
                return

            agent = get_agent(agent_name)
            if not agent:
                self.send_response(404)
                self.end_headers()
                self.wfile.write(f"Agent {agent_name} not found".encode())
                return

            try:
                content = agent.run(message)

                self.send_response(200)
                self.send_header("Content-Type", "text/event-stream")
                self.send_header("Cache-Control", "no-cache")
                self.send_header("Access-Control-Allow-Origin", "*")
                self.end_headers()

                self.wfile.write(f"data: {content}\n\n".encode())
                self.wfile.write(b"data: [DONE]\n\n")
                self.wfile.flush()

            except Exception as e:
                self.send_response(200)
                self.send_header("Content-Type", "text/event-stream")
                self.end_headers()
                self.wfile.write(f"data: Error: {str(e)}\n\n".encode())
                self.wfile.write(b"data: [DONE]\n\n")
                self.wfile.flush()

            return

        self.send_response(404)
        self.end_headers()

if __name__ == "__main__":
    port = int(os.environ.get("AI_PORT", "8091"))
    server = ThreadingHTTPServer(("127.0.0.1", port), AgentHandler)
    server.serve_forever()
