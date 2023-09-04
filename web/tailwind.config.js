/** @type {import('tailwindcss').Config} */
export default {
	content: ["./src/**/*.{html,ts,js,svelte}"],
	theme: {
		extend: {
			fontFamily: {
				accent: ["Alegreya SC", "serifs"],
				regular: ["Inter", "sans-serifs"],
			},
			colors: {
				main: "#0d9488",
				dark: {
					100: "#343135",
					200: "#262626",
				},
				light: "#d4d4d4",
                accent: "#F250A3"
			},
		},
	},
	plugins: [require("daisyui")],
	daisyui: {
		themes: ["night"],
	},
};
