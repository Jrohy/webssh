const path = require('path')
const webpack = require('webpack')

const proxyTargetMap = {
    prod: 'https://xxx.xxx.com/',
    dev: 'http://127.0.0.1:5032'
}
const proxyTarget = proxyTargetMap[process.env.API_TYPE] || proxyTargetMap.prod
const wsTarget = proxyTarget.replace('http', 'ws')
const publicPath = process.env.NODE_ENV === 'production' ? '/' : '/'

module.exports = {
    publicPath: publicPath,
    outputDir: 'dist',

    // 放置静态资源的地方 (js/css/img/font/...)
    assetsDir: 'static',

    lintOnSave: true,

    transpileDependencies: [
        /* string or regex */
    ],

    // 是否为生产环境构建生成 source map？
    productionSourceMap: false,

    // 调整内部的 webpack 配置。
    // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli/webpack.md
    chainWebpack: config => {
        // 移除 prefetch 插件,解决组件懒加载失效的问题
        config.plugins.delete('prefetch')
    },

    // 在生产环境下为 Babel 和 TypeScript 使用 `thread-loader`
    // 在多核机器下会默认开启。
    parallel: require('os').cpus().length > 1,

    // PWA 插件的选项。
    // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli-plugin-pwa/README.md
    pwa: {},

    // 配置 webpack-dev-server 行为。
    devServer: {
        disableHostCheck: false,
        open: process.platform === 'darwin',
        host: '0.0.0.0',
        port: 8257,
        https: false,
        hotOnly: false,
        // eslint-disable-next-line no-dupe-keys
        open: true,
        // 查阅 https://github.com/vuejs/vue-docs-zh-cn/blob/master/vue-cli/cli-service.md#配置代理
        proxy: {
            '/api': {
                target: proxyTarget,
                changeOrigin: true,
                ws: true,
                pathRewrite: {
                    '^/api': ''
                }
            },
            '/ws': {
                target: wsTarget,
                changeOrigin: true,
                ws: true,
                pathRewrite: {
                    '^/ws': ''
                }
            }
        },
        before: app => {}
    },

    // 第三方插件的选项
    pluginOptions: {}
}
