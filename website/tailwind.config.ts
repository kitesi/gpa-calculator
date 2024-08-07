import type { Config } from "tailwindcss";

const config: Config = {
    content: [
        "./pages/**/*.{js,ts,jsx,tsx,mdx}",
        "./components/**/*.{js,ts,jsx,tsx,mdx}",
        "./app/**/*.{js,ts,jsx,tsx,mdx}",
    ],
    theme: {
        extend: {
            backgroundImage: {
                "gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
                "gradient-conic":
                    "conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
            },
            colors: {
                midnight: {
                    700: "#30363d",
                    800: "#0d1117",
                    900: "#010409",
                },
                "my-green": "#357541",
                "my-red": "#8e3b3b",
                "my-neutral": "#30363d",
                "my-blue": "#295359",
            },
        },
    },
    plugins: [],
};
export default config;
