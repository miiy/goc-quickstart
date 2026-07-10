import { fileURLToPath, URL } from "node:url";

import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

const entries = {
  app: "frontend/src/app.ts",
  "post/create": "frontend/src/features/post/entries/create.tsx",
  "post/detail": "frontend/src/features/post/entries/detail.tsx",
  "post/edit": "frontend/src/features/post/entries/edit.tsx",
  "user/profile": "frontend/src/features/profile/entries/profile.tsx",
  "user/show": "frontend/src/features/user/entries/show.tsx"
};

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [react(), tailwindcss()],
    publicDir: "resources/static",
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./frontend/src", import.meta.url))
      }
    },
    server: {
      host: env.VITE_HOST || '127.0.0.1',
      port: parseInt(env.VITE_PORT, 10),
      strictPort: true,
      cors: true
    },
    preview: {
      host: env.VITE_HOST || '127.0.0.1',
      port: parseInt(env.VITE_PORT, 10) || 4173,
      strictPort: true
    },
    build: {
      outDir: "dist",
      emptyOutDir: true,
      manifest: true,
      rollupOptions: {
        input: entries
      }
    }
  };
});
