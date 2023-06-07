/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.{svelte,ts,js,html,astro}"],
	theme: {
		extend: {
			fontFamily: {
				main: ["Monomaniac One", "sans-serif"],
				accent: ["Exo 2", "sans-serif"],
				regular: ["Cabin", "sans-serif"],
				mono: ["JetBrains Mono", "monospace"],
			},
			colors: {
				main: "#0d9488",
				dark: "#262626",
				light: "#d4d4d4",
				accent: {
					blue: "#1D4ED8",
					fuchsia: "#A21CAF",
				},
			},
		},
	},
	plugins: [],
};
