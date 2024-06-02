/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.{go,js,templ,html}",
    "./components/**/*.{go,js,templ,html}",
  ],
  theme: {
    extend: {
      colors: {
        platform: {
          red: "#E11D48",
          green: "#26c03d",
          border: "#305ba0",
          accent: "#0682D4",
          background: "#0E123E",
        },
      },
      fontSize: { "2xs": "0.7rem", "3xs": "0.65rem" },
      backgroundImage: { gradient: "var(--main-bg)" },
    },
  },
  plugins: [],
};
