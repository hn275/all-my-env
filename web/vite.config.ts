import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import path from "path";
import dns from "dns";

dns.setDefaultResultOrder("verbatim");
export default defineConfig({
	plugins: [sveltekit()],
	server: {
		host: "localhost",
		port: 3000,
    proxy: {
      "/app": {
        target: "http://localhost:8080",
        changeOrigin: false,
      }
    }
	},
	resolve: {
		alias: {
			"@lib": path.resolve("./src/lib"),
			"@assets": path.resolve("./src/assets"),
		},
	},
});
