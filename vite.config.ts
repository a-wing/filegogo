import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/s': 'http://localhost:8080',
      '^/raw/.*': 'http://localhost:8080',
      '/config': 'http://localhost:8080',
      '^/s/.*': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
  publicDir: "webapp/public",
  build: {
    outDir: "server/build"
  },
  plugins: [react()]
})
