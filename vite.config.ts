import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/

const server = 'http://localhost:8080'
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/s': server,
      '^/raw/.*': server,
      '^/info/.*': server,
      '/config': server,
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
