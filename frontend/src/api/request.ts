import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建 Axios 实例
const service = axios.create({
  // 注意：因为我们在 vite.config.ts 中配置了 proxy，所以这里的 baseURL 是 /api
  // 请求会通过代理转发到 http://localhost:8080/api
  baseURL: '/api/v1',
  timeout: 5000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么
    // 例如，从 Pinia store 中获取 token 并添加到请求头
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    // 对请求错误做些什么
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    // 对响应数据做点什么
    const res = response.data

    // 如果业务 code 不是 0，则判断为错误。
    if (res.code !== 0) {
      ElMessage({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      })
      // 可以根据不同的 code 做不同的处理
      // 例如：code 4001 代表 token 问题，可以触发登出操作
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      // 业务 code 为 0，直接返回 data 部分
      return res.data
    }
  },
  (error) => {
    // 对响应错误做点什么
    console.log('err' + error) // for debug
    ElMessage({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
