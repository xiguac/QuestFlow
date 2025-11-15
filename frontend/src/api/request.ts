import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

// 创建 Axios 实例
const service = axios.create({
  baseURL: '/api/v1',
  timeout: 10000 // 导出可能耗时较长，增加超时时间
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    console.log(error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    // 【核心修复】
    // 检查响应数据是否是 Blob 类型（文件下载）
    // 如果是，直接构造并返回一个包含 blob 和 fileName 的对象
    if (response.data instanceof Blob) {
      const contentDisposition = response.headers['content-disposition'];
      let fileName = 'download.xlsx'; // 默认文件名
      if (contentDisposition) {
        // 优先匹配 RFC 5987 格式的编码文件名
        const fileNameMatch = contentDisposition.match(/filename\*=UTF-8''(.+)/);
        if (fileNameMatch && fileNameMatch.length > 1) {
          fileName = decodeURIComponent(fileNameMatch[1]);
        } else {
          // 回退到匹配普通 filename="xxx" 格式
          const fileNameMatchFallback = contentDisposition.match(/filename="(.+)"/);
          if (fileNameMatchFallback && fileNameMatchFallback.length > 1) {
            fileName = decodeURIComponent(fileNameMatchFallback[1]);
          }
        }
      }
      return {
        blob: response.data,
        fileName: fileName,
      };
    }

    // 如果不是文件，则按原有的 JSON API 格式处理
    const res = response.data;
    if (res.code !== 0) {
      ElMessage({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      });
      return Promise.reject(new Error(res.message || 'Error'));
    } else {
      return res.data;
    }
  },
  (error) => {
    console.log('err:', error);

    // 401 登录过期处理
    if (error.response && error.response.status === 401) {
      const userStore = useUserStore();
      if (router.currentRoute.value.name !== 'login') {
        ElMessage.error('登录状态已过期，请重新登录');
        userStore.logout();
        router.push(`/login?redirect=${router.currentRoute.value.fullPath}`);
      }
    }
    // 【新增】处理后端返回的业务错误（例如 404 Not Found），此时响应体可能是 JSON 格式的 Blob
    else if (error.response && error.response.data instanceof Blob && error.response.data.type.toLowerCase().includes('application/json')) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
          const errorJson = JSON.parse(reader.result as string);
          ElMessage.error(errorJson.message || '操作失败');
          reject(new Error(errorJson.message || '操作失败'));
        };
        reader.onerror = () => {
          ElMessage.error('无法解析错误信息');
          reject(new Error('无法解析错误信息'));
        };
        reader.readAsText(error.response.data);
      });
    }
    // 其他网络错误等
    else {
      ElMessage({
        message: error.message || '网络错误',
        type: 'error',
        duration: 5 * 1000
      });
    }

    return Promise.reject(error);
  }
);

export default service;
