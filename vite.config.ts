import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/s': 'http://localhost:8080',
      '/config': 'http://localhost:8080',
      '^/s/.*': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
  build: {
    outDir: "server/build"
  },
  plugins: [react()]
})
