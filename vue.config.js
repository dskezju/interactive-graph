const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: "/interactive-graph/",
  chainWebpack: (config) => {
    // webgl Loader
    config.module
      .rule("glsl")
      .test(/\.glsl$/)
      .use("raw-loader")
      .loader("raw-loader")
      .end();
  },
  configureWebpack: {
    devServer: {
      allowedHosts: "all",
    },
  },
});
