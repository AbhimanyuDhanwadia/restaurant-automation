import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath, URL } from "node:url";

const isGitHubActions = (globalThis as { process?: { env?: Record<string, string | undefined> } }).process?.env?.GITHUB_ACTIONS === "true";

export default defineConfig({
  plugins: [
    tailwindcss(),
    react(),
  ],
  base: isGitHubActions ? "/restaurant-automation/" : "/",
  resolve: {
    alias: {
      // Allows imports like: import { X } from "@/components/..."
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
