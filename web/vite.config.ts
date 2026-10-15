import { fileURLToPath } from 'node:url'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

const path = (relative: string) => fileURLToPath(new URL(relative, import.meta.url))

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    // Mirrors the path aliases of tsconfig.json (ADR web core 0007).
    alias: [
      { find: /^@bc\/([^/]+)$/, replacement: path('./src/bounded-contexts/$1/index.ts') },
      { find: '@app', replacement: path('./src/app') },
      { find: '@core', replacement: path('./src/core') },
      { find: '@shared', replacement: path('./src/shared') },
      { find: '@api-client', replacement: path('./var/api-client') },
    ],
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test-setup.ts'],
    css: false,
    env: { VITE_API_URL: 'http://api.test' },
  },
})
