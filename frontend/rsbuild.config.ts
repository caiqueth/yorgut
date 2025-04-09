import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";

export default defineConfig({
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        pathRewrite: { "^/api": "/api/v1" }
      },
    }
  },
  html: {
    template: "./public/index.html",
  },
  plugins: [pluginReact()],
});
