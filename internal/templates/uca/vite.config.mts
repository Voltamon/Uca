import { defineConfig } from "vite"
import preact from "@preact/preset-vite"
import { resolve } from "path"

export default defineConfig({
  plugins: [preact()],
  root: ".",
  resolve: {
    alias: {
      "uca/types": resolve(__dirname, "types/index.ts")
    }
  },
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
