module.exports = {
  mode: "jit",
  purge: { enabled: false, content: ["./public/**/*.html", "./src/**/*.vue"] },
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      shadow: {
        up: "0 4px 4px 0 rgba(0, 0, 0, 0.25)",
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
