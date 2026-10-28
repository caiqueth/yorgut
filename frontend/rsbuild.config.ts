import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";

// Without it the login button can't work; each developer uses their own
// OAuth client (README.md, "Local development").
if (!process.env.GOOGLE_CLIENT_ID) {
  throw new Error("GOOGLE_CLIENT_ID is not set: put it in frontend/.env.local, or pass it to the Docker build.");
}

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
  source: {
    define: {
      "__GOOGLE_CLIENT_ID__": JSON.stringify(process.env.GOOGLE_CLIENT_ID),
      "__API_BASE_URL__": JSON.stringify(process.env.API_BASE_URL),
    },
  },
});
