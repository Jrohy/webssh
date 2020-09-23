import axios from 'axios'
import { Message } from 'element-ui'

function validateStatus(status) {
    switch (status) {
    case 400:
        Message.error('请求出错')
        break
    case 403:
        Message.warning({
            message: '拒绝访问'
        })
        break
    case 404:
        Message.warning({
            message: '请求错误,未找到该资源'
        })
        break
    case 500:
        Message.warning({
            message: '服务端错误'
        })
        break
    }
    return status >= 200 && status < 300
}

var instance = axios.create({
    timeout: 8000,
    baseURL: process.env.NODE_ENV === 'production' ? '/' : '/api',
    validateStatus
})

// 响应拦截器即异常处理
instance.interceptors.response.use(
    response => {
        return response.data
    },
    err => {
        if (err.response === undefined) {
            Message.error('连接服务器失败')
        }
        return Promise.reject(err.response)
    }
)

export default instance
