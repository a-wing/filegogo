import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

import { createHtmlPlugin } from 'vite-plugin-html'

// https://vitejs.dev/config/
export default defineConfig(({command, mode}) => {
  const env = loadEnv(mode, process.cwd());
  return {
  base: "./",
  server: {
    port: 3000,
    proxy: {
      '/s': 'http://localhost:8080',
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
  plugins: [react(),
    createHtmlPlugin({
      inject: {
        data: {
          // ...env,
          BASE_URL: env.VITE_BASE_URL,
        },
      },
    }),
  ]
  };
})
