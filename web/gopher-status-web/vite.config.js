import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    server: {
        proxy: {
            '/v1/monitors': {
                target: 'http://127.0.0.1:8080', // ИСПРАВЛЕНО ЗДЕСЬ
                changeOrigin: true,
                secure: false,
            },
            '/v1/auth': {
                target: 'http://127.0.0.1:8080', // И ИСПРАВЛЕНО ЗДЕСЬ
                changeOrigin: true,
                secure: false,
            },
        },
    },
    plugins: [react()],
})
