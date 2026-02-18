import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    server: {
        proxy: {
            '/v1/monitors': { // Ловим все запросы на /monitor
                target: 'http://localhost:8080',
                changeOrigin: true,
                secure: false,

            },
        },
    },
    plugins: [react()],
})
