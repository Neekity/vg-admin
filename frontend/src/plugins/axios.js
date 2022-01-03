"use strict";

import axios from "axios";
import Vue from 'vue';
import config from '../config';
import {getToken} from "../store/module/user";

const _axios = axios.create({
    baseUrl: config.baseUrl,
    timeout: 30000,
});
const errorResponseHandle = (response) => {
    if (response === undefined) {
        alert('服务器出错！')
    } else {
        if (response.status === 401) {
            window.location.href = '/login';
        }
        if (response.status === 403) {
            alert('请求方式错误！')
        }
        if (response.status === 500 || response.status > 504) {
            alert('服务器发生错误')
        }

        if (response.status === 502) {
            alert('服务器网关错误');
        }
        if (response.status === 503) {
            alert('服务器不可用');
        }
        if (response.status === 504) {
            alert('网关超时');
        }
    }
}
export const apiUrlDomain = config.baseUrl;

function parseParams(data) {
    try {
        let paramsArray = [];
        for (let i in data) {
            let key = encodeURIComponent(i);
            let value = encodeURIComponent(data[i]);
            paramsArray.push(key + '=' + value);
        }
        return paramsArray.join('&');
    } catch (e) {
        return '';
    }
}

_axios.interceptors.request.use(function (config) {
    if (getToken()) {
        config.headers.Authorization = 'Bearer ' + getToken();
    } else {
        window.location.href = '/login';
    }
    return config;
}, function (error) {
    return Promise.reject(error);
})

// Add a response interceptor
_axios.interceptors.response.use(
    response => {
        errorResponseHandle(response)
        return response;
    },
    error => {
        const response = error.response;
        errorResponseHandle(response);

        return Promise.reject(error);
    }
);

export const http = (method, uri, data = {}) => {
    let url = apiUrlDomain + (uri.indexOf('/') === 0 ? uri : '/' + uri);
    if (method === 'get') {
        url = url + '?' + parseParams(data);
    }
    return _axios[method](url, data);
}

Plugin.install = function (Vue) {
    Vue.http = http;
    window.http = http;
    Object.defineProperties(Vue.prototype, {
        $http: {
            get() {
                return http;
            }
        },
        http: {
            get() {
                return http;
            }
        },
    });
};

Vue.use(Plugin)

export default Plugin;
