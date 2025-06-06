import axios from 'axios';
import { ElMessage } from 'element-plus';
import { getToken, clearAuth } from './auth';
import router from '../router';

// 创建axios实例
const service = axios.create({
  baseURL: '/api', // 使用Vite配置的代理路径
  timeout: 30000, // 请求超时时间增加到30秒
});

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 从认证工具获取token
    const token = getToken();
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
    
    // 如果响应码不为200，则表示有错误
    if (res.code !== 200) {
      ElMessage({
        message: res.message || '服务器错误',
        type: 'error',
        duration: 5 * 1000
      });
      
      // 如果是401，则清除认证信息并重定向到登录页
      if (res.code === 401) {
        handleUnauthorized();
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
          handleUnauthorized();
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

// 处理未授权情况
function handleUnauthorized() {
  // 清除所有认证信息
  clearAuth();
  
  // 获取当前路径，用于登录后跳转回来
  const currentPath = router.currentRoute.value.fullPath;
  
  // 跳转到登录页面
  router.push(`/login?redirect=${encodeURIComponent(currentPath)}`);
}

export default service; 