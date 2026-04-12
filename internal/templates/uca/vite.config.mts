import { defineConfig } from "vite"
import preact from "@preact/preset-vite"

export default defineConfig({
  plugins: [preact()],
  root: ".",
  server: {
    port: {{FRONTEND_PORT}},
    proxy: {
      "/api": "http://localhost:{{BACKEND_PORT}}"
    }
  },
  build: {
    outDir: ".vite"
  }
})
