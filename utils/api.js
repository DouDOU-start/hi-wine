// api.js - API请求工具类
import { BASE_URL, APP_ID } from '@/config.js';

// 获取存储的token
const getToken = () => {
  return uni.getStorageSync('token') || '';
};

// 请求封装
const request = (url, method, data, needAuth = true) => {
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + url,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': needAuth ? `Bearer ${getToken()}` : ''
      },
      success: (res) => {
        if (res.statusCode === 200) {
          resolve(res.data);
        } else if (res.statusCode === 401) {
          // token失效，跳转到登录页
          uni.showToast({
            title: '请先登录',
            icon: 'none'
          });
          // 清除token
          uni.removeStorageSync('token');
          // 跳转到个人中心页面
          setTimeout(() => {
            uni.switchTab({
              url: '/pages/profile/index'
            });
          }, 1500);
          reject(new Error('未登录或登录已过期'));
        } else {
          uni.showToast({
            title: res.data.message || '请求失败',
            icon: 'none'
          });
          reject(new Error(res.data.message || '请求失败'));
        }
      },
      fail: (err) => {
        uni.showToast({
          title: '网络错误',
          icon: 'none'
        });
        reject(err);
      }
    });
  });
};

// API方法
export default {
  // 用户登录
  login(code) {
    // 按照微信官方API参数格式发送请求
    return request('/wechat/login', 'POST', { 
      code: code  // 临时登录凭证
    }, false);
  },
  
  // 更新用户信息
  updateUserInfo(nickname, avatar) {
    return request('/api/user/update', 'POST', {
      nickname: nickname,
      avatar: avatar
    }, true);
  },
  
  // 获取商品分类列表
  getCategoryList() {
    return request('/api/category/list', 'GET');
  },
  
  // 获取商品列表
  getProductList(categoryId = 0, keyword = '', page = 1, size = 10) {
    return request('/api/product/list', 'GET', { categoryId, keyword, page, size });
  },
  
  // 获取商品详情
  getProductDetail(id) {
    return request('/api/product/detail', 'GET', { id });
  },
  
  // 创建订单
  createOrder(tableId, items) {
    return request('/api/order/create', 'POST', { tableId, items });
  },
  
  // 获取订单列表
  getOrderList(status = -1, page = 1, size = 10) {
    return request('/api/order/list', 'GET', { status, page, size });
  },
  
  // 获取订单详情
  getOrderDetail(id) {
    return request('/api/order/detail', 'GET', { id });
  },
  
  // 更新订单状态
  updateOrderStatus(id, status) {
    return request('/api/order/updateStatus', 'POST', { id, status });
  }
}; 