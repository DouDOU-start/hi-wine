import axios from 'axios';
import { ElMessage } from 'element-plus';
import { getToken, clearAuth } from './auth';
import router from '../router';
import { toCamelCase, toSnakeCase } from './format';

// 创建axios实例
const service = axios.create({
  baseURL: '/api', // 使用Vite配置的代理路径
  timeout: 30000, // 请求超时时间
});

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 设置token
    const token = getToken();
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    
    // 将请求数据从驼峰格式转换为下划线格式
    if (config.method === 'post' || config.method === 'put' || config.method === 'patch') {
      if (config.data && typeof config.data === 'object' && !(config.data instanceof FormData)) {
        config.data = toSnakeCase(config.data);
      }
    }
    
    // 将GET请求参数从驼峰格式转换为下划线格式
    if (config.params && typeof config.params === 'object') {
      config.params = toSnakeCase(config.params);
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
    // 检查是否需要显示全局成功消息
    if (response.config.showSuccessMessage) {
      ElMessage.success(response.data.message || '操作成功');
    }
    
    // 如果是下载文件请求，直接返回响应
    if (response.config.responseType === 'blob') {
      return response;
    }
    
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
      // 将整个响应数据从下划线格式转换为驼峰格式
      if (res && typeof res === 'object') {
        // 先转换res.data
        if (res.data && typeof res.data === 'object') {
          res.data = toCamelCase(res.data);
        }
      }
      return res;
    }
  },
  error => {
    // 处理网络错误
    if (!error.response) {
      ElMessage.error('网络错误，请检查您的网络连接');
      return Promise.reject(error);
    }
    
    // 处理HTTP错误
    const { status, data } = error.response;
    
    // 处理401未授权错误
    if (status === 401) {
      ElMessage.error('登录已过期，请重新登录');
      clearAuth();
      router.push('/login');
      return Promise.reject(error);
    }
    
    // 处理其他错误
    const errorMsg = (data && data.message) || '请求失败';
    ElMessage.error(errorMsg);
    
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

/**
 * 通用GET请求方法
 * @param {string} url - 请求URL
 * @param {Object} params - 请求参数
 * @returns {Promise} 请求Promise
 */
export function get(url, params = {}) {
  return service({
    url,
    method: 'get',
    params
  });
}

/**
 * 通用POST请求方法
 * @param {string} url - 请求URL
 * @param {Object} data - 请求数据
 * @returns {Promise} 请求Promise
 */
export function post(url, data = {}) {
  return service({
    url,
    method: 'post',
    data
  });
}

/**
 * 通用PUT请求方法
 * @param {string} url - 请求URL
 * @param {Object} data - 请求数据
 * @returns {Promise} 请求Promise
 */
export function put(url, data = {}) {
  return service({
    url,
    method: 'put',
    data
  });
}

/**
 * 通用DELETE请求方法
 * @param {string} url - 请求URL
 * @param {Object} params - 请求参数
 * @returns {Promise} 请求Promise
 */
export function del(url, params = {}) {
  return service({
    url,
    method: 'delete',
    params
  });
}

/**
 * 下载文件请求方法
 * @param {string} url - 请求URL
 * @param {Object} params - 请求参数
 * @param {string} filename - 文件名
 * @returns {Promise} 请求Promise
 */
export function download(url, params = {}, filename = '') {
  return service({
    url,
    method: 'get',
    params,
    responseType: 'blob'
  }).then(response => {
    const blob = new Blob([response.data]);
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = filename || getFilenameFromResponse(response);
    document.body.appendChild(link);
    link.click();
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
    return response;
  });
}

/**
 * 从响应头中获取文件名
 * @param {Object} response - 响应对象
 * @returns {string} 文件名
 */
function getFilenameFromResponse(response) {
  const contentDisposition = response.headers['content-disposition'];
  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename="?(.+)"?/);
    if (filenameMatch && filenameMatch.length > 1) {
      return filenameMatch[1];
    }
  }
  return 'downloaded-file';
}

export default service; 