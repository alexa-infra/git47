const postcssImport = require('postcss-import');
const tailwind = require('tailwindcss');
const autoprefixer = require('autoprefixer');
const cssnano = require('cssnano');

const plugins = [postcssImport, tailwind, autoprefixer];

if (process.env.NODE_ENV === 'production') {
  plugins.push(cssnano);
}

module.exports = { plugins };
