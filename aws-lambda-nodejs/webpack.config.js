const webpack = require('webpack');

module.exports = {
  mode: "development",
  resolve: {
    alias: {
      'node-fetch$': "node-fetch/lib/index.js"
    }
  },
  plugins: [
    new webpack.IgnorePlugin(/^pg-native$/),
    new webpack.IgnorePlugin(/^tedious$/),
  ],
}
