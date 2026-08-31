import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import wails from "@wailsio/runtime/plugins/vite";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, "./src/lib"),
      $shared: path.resolve(__dirname, "./src/lib/shared"),
      $navigation: path.resolve(__dirname, "./src/lib/shared/navigation"),
      $session: path.resolve(__dirname, "./src/lib/shared/session"),
      $bindings: path.resolve(__dirname, "./bindings/go-walis/internal"),
    },
  },
  plugins: [tailwindcss(), svelte(), wails("./bindings")],
});
