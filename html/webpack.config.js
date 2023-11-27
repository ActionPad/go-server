const HtmlWebpackPlugin = require('html-webpack-plugin');
const TerserPlugin = require("terser-webpack-plugin");

module.exports = {
  mode: 'development',
  entry: './src/index.js',
  output: {
	filename: 'bundle.min.js',
  },
  optimization: {
  	minimize: true,
  	minimizer: [new TerserPlugin()],
  },
  module: {
	rules: [
	  {
		test: /\.css$/,
		use: ['style-loader', 'css-loader'],
	  },
	],
  },
  plugins: [
	new HtmlWebpackPlugin({
	  template: './src/index.html',
	  inlineSource: '.(js|css)$',
	  minify: true,
	}),
  ],
};