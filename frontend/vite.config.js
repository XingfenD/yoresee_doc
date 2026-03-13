import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  optimizeDeps: {
    include: ['codemirror', 'easymde']
  },
  server: {
    host: true,
    port: 80,
    strictPort: true,
    hmr: {
      protocol: 'ws',
      host: 'localhost',
      port: 80,
      clientPort: 8080
    }
  }
})
