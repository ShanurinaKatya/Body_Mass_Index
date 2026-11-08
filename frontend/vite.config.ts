import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    strictPort: true,
    proxy: {
      // Обратный прокси к Go-сервису: снимает ограничение CORS.
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // Второй маршрут — к объектному хранилищу MinIO (S3).
      '/s3': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/s3/, '/body-mass-index'),
      },
    },
  },
  plugins: [react()],
})
