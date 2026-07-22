import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const isGitHubActions = (globalThis as { process?: { env?: Record<string, string | undefined> } }).process?.env?.GITHUB_ACTIONS === "true";

export default defineConfig({
  plugins: [react()],
  base: isGitHubActions ? "/restaurant-automation/" : "/",
});
