import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'https://fast-order-xvkq.onrender.com',
        changeOrigin: true,
        secure: true,
        rewrite: (path) => path, // Keep the /api prefix
        configure: (proxy, options) => {
          proxy.on('error', (err, req, res) => {
            console.log('proxy error', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            console.log('proxying:', req.method, req.url, '->', options.target + req.url);
          });
        },
      },
    },
  },
})
