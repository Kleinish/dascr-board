const plugin = require('tailwindcss/plugin');

module.exports = {
  content: [
    './src/**/*.svelte', // Include all Svelte files in src
    './public/*.html', // Include any HTML files in public
  ],
  theme: {
    extend: {
      gradientColorStops: {
        primary: '#1d3c4c',
        secondary: '#3c5b6c',
      },
    },
  },
  darkMode: 'media', // or 'class'
  plugins: [],
};
