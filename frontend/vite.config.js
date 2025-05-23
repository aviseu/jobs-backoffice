import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ command, mode }) => {
  loadEnv(mode, process.cwd());
  return {
    plugins: [react()],
    server: {
      port: 3000,
      host: true,
    },
  };
});
