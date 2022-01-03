export default {
    /**
     * @description 配置显示在浏览器标签的title
     */
    title: process.env.VUE_APP_SYSTEM_TITLE,
    /**
     * @description token在Cookie中存储的天数，默认1天
     */
    cookieExpires: process.env.VUE_APP_COOKIE_EXPIRES,

    baseUrl: process.env.VUE_APP_BACKEND_URL,
}
