import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import WindiCSS from 'vite-plugin-windicss'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    watch: {
      usePolling: true
    }
  },
  plugins: [
      vue(),
      WindiCSS(),
  ],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
})
