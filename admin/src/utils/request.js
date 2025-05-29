import axios from 'axios';
import { ElMessage } from 'element-plus';

// 创建axios实例
const service = axios.create({
  baseURL: '', // 不设置baseURL，使用相对路径
  timeout: 30000, // 请求超时时间增加到30秒
});

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 从本地存储获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data;
    
    // 如果响应码不为0，则表示有错误
    if (res.code !== 0) {
      ElMessage({
        message: res.message || '服务器错误',
        type: 'error',
        duration: 5 * 1000
      });
      
      // 如果是401，则清除token并重定向到登录页
      if (res.code === 401) {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
      
      return Promise.reject(new Error(res.message || '服务器错误'));
    } else {
      return res;
    }
  },
  error => {
    console.error('响应错误:', error);
    let message = '连接服务器失败';
    
    if (error.response) {
      switch (error.response.status) {
        case 401:
          message = '未授权，请重新登录';
          localStorage.removeItem('token');
          window.location.href = '/login';
          break;
        case 403:
          message = '拒绝访问';
          break;
        case 404:
          message = '请求的资源不存在';
          break;
        case 500:
          message = '服务器内部错误';
          break;
        default:
          message = `请求错误 (${error.response.status})`;
      }
      
      // 尝试从响应中获取更详细的错误信息
      if (error.response.data && error.response.data.message) {
        message = error.response.data.message;
      }
    } else if (error.request) {
      message = '服务器无响应';
    }
    
    ElMessage({
      message: message,
      type: 'error',
      duration: 5 * 1000
    });
    
    return Promise.reject(error);
  }
);

export default service; 