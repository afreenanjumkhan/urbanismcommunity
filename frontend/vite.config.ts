import { defineConfig } from "vite";

// The UI is the prototype bundle served verbatim from public/legacy.
// Vite only compiles src/ (the sync adapter) and serves/builds the shell.
export default defineConfig({
  server: { port: 5173 },
  build: { outDir: "dist", sourcemap: false },
});
