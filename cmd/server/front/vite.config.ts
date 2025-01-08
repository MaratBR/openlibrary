import path from 'node:path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

const SRC_DIR = path.resolve(__dirname, './src')
const COMPONENTS_DIR = path.resolve(SRC_DIR, './components')
const NODE_MODULES = path.resolve(__dirname, './node_modules')

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': SRC_DIR,
    },
  },

  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.startsWith(COMPONENTS_DIR)) {
            return 'components'
          }

          if (id.startsWith(NODE_MODULES)) {
            return 'vendor'
          }

          return undefined
        },
      },
    },
  },
})
