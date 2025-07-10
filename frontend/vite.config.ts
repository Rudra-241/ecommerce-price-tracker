import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'node:path'

// The Go server (Gin) serves the built SPA from `src/public` and exposes the
// JSON API under `/api`. In dev we run Vite on :5173 and proxy `/api` to the
// Go server on :3000 so the browser treats everything as same-origin — that
// keeps the httpOnly auth cookies working without any CORS setup.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: path.resolve(__dirname, '../src/public'),
    emptyOutDir: true,
  },
})
