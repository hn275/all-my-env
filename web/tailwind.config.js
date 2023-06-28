/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{ts,tsx,js,jsx,html}"],
  theme: {
    extend: {
      fontFamily: {},
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
