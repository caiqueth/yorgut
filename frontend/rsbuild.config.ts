import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";

export default defineConfig({
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        pathRewrite: { "^/api": "/api/v1" }
      },
      "/storage/uploads": {
        target: "http://localhost:8080",
      },
    }
  },
  html: {
    template: "./public/index.html",
  },
  plugins: [pluginReact()],
  source: {
    define: {
      "__GOOGLE_CLIENT_ID__": JSON.stringify(process.env.GOOGLE_CLIENT_ID),
      "__API_BASE_URL__": JSON.stringify(process.env.API_BASE_URL),
    },
  },
});
