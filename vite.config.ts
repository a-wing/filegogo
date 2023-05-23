import { defineConfig } from 'vite'
import UnoCSS from 'unocss/vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  base: "./",
  server: {
    port: 3000,
    proxy: {
      '^.*/api/.*': 'http://localhost:8080',
      '^.*/api/signal/.*': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
  publicDir: "webapp/public",
  build: {
    outDir: "server/build"
  },
  plugins: [react(), UnoCSS()]
})
