import { defineConfig } from "vite"
import preact from "@preact/preset-vite"
import { resolve } from "path"

export default defineConfig({
  plugins: [preact()],
  root: ".",
  resolve: {
    alias: {
      "uca/types": resolve(__dirname, "types/index.ts"),
      "uca/roles": resolve(__dirname, "roles/index.ts"),
      "preact/jsx-runtime": resolve(__dirname, "node_modules/preact/jsx-runtime"),
      "preact/hooks": resolve(__dirname, "node_modules/preact/hooks"),
      "preact": resolve(__dirname, "node_modules/preact")
    }
  },
  server: {
    port: {{FRONTEND_PORT}},
    proxy: {
      "/api": "http://localhost:{{BACKEND_PORT}}"
    }
  },
  build: {
    outDir: ".vite",
    modulePreload: false,
    rollupOptions: {
      input: resolve(__dirname, "index.html"),
      output: {
        crossOriginLoading: false
      }
    }
  },
  optimizeDeps: {
    include: ["preact", "preact/hooks", "preact/jsx-runtime"]
  }
})
