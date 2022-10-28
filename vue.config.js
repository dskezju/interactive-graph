const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  chainWebpack: (config) => {
    // webgl Loader
    config.module
      .rule("glsl")
      .test(/\.glsl$/)
      .use("raw-loader")
      .loader("raw-loader")
      .end();
  },
});
