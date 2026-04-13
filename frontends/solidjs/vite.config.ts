import { defineConfig } from 'vite';
import solid from 'vite-plugin-solid';
import { fileURLToPath, URL } from 'node:url';

export default defineConfig({
  plugins: [solid()],
  server: {
    port: 3000,
    host: '127.0.0.1',
    fs: {
      allow: [fileURLToPath(new URL('../../..', import.meta.url))],
    },
  },
  test: {
    environment: 'jsdom',
    include: ['tests/**/*.test.{ts,tsx}'],
    exclude: ['tests/e2e/**', 'node_modules/**'],
  },
});
