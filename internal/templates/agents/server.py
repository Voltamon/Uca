from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import urlparse, parse_qs
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

import litellm
import json

def load_env():
    env_path = ".env"
    if os.path.exists(env_path):
        with open(env_path) as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith("#"):
                    k, v = line.split("=", 1)
                    os.environ[k] = v

load_env()

MODEL = "{{MODEL}}"
API_KEY = os.environ.get("GITHUB_PAT_TOKEN", "")
API_BASE = "https://models.inference.ai.azure.com"
TIMEOUT = int(os.environ.get("AI_TIMEOUT", "30"))

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
            message = params.get("message", [""])[0]

            try:
                response = litellm.completion(
                    model=MODEL,
                    messages=[{"role": "user", "content": message}],
                    api_key=API_KEY,
                    api_base=API_BASE,
                    timeout=TIMEOUT
                )

                content = response.choices[0].message.content

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

    def do_POST(self):
        self.send_response(404)
        self.end_headers()

if __name__ == "__main__":
    port = int(os.environ.get("AI_PORT", "8091"))
    server = ThreadingHTTPServer(("127.0.0.1", port), AgentHandler)
    server.serve_forever()
