module.exports = {
    outputDir: 'dist',
    transpileDependencies: [
        'vuetify'
    ],
    devServer: {
        hot: true,
        compress: true,
        disableHostCheck: true, //webpack4.0 开启热更新
    },
}
