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
                accent: {
                    blue: "#1D4ED8",
                    fuchsia: {
                        100: "#E879F9",
                        200: "#A21CAF",
                    },
                },
            },
        },
    },
    plugins: [require("daisyui")],
    daisyui: {
        themes: ["dark"]
        // extends: {
        //     themes: [
        //         {
        //             envhub: {
        //                 primary: "#0d9488",
        //             },
        //         },
        //     ],
        // }
    },
};
