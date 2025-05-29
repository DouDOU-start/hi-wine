import request from '../utils/request';

// 用户登录
export function login(data) {
  return request({
    url: '/api/admin/login',
    method: 'post',
    data
  });
}

// 获取当前用户信息
export function getUserInfo() {
  return request({
    url: '/api/admin/info',
    method: 'get'
  });
}

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/api/user/list',
    method: 'get',
    params
  });
}

// 更新用户状态
export function updateUserStatus(id, status) {
  return request({
    url: '/api/user/status',
    method: 'post',
    data: { id, status }
  });
}

// 获取用户详情
export function getUserDetail(id) {
  return request({
    url: '/api/user/detail',
    method: 'get',
    params: { id }
  });
} 