import request from '../utils/request';

// 管理员登录
export function login(data) {
  return request({
    url: '/admin/auth/login',
    method: 'post',
    data
  });
}

// 获取当前管理员信息
export function getAdminInfo() {
  return request({
    url: '/admin/profile',
    method: 'get'
  });
}

// 获取管理员列表
export function getAdminList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  });
}

// 创建管理员
export function createAdmin(data) {
  return request({
    url: '/admin/users',
    method: 'post',
    data
  });
}

// 更新管理员信息
export function updateAdmin(id, data) {
  return request({
    url: `/admin/users/${id}`,
    method: 'put',
    data
  });
}

// 重置管理员密码
export function resetAdminPassword(id, newPassword) {
  return request({
    url: `/admin/users/${id}/reset-password`,
    method: 'post',
    data: { newPassword }
  });
}

// 删除管理员
export function deleteAdmin(id) {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete'
  });
}

// 获取普通用户列表
export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  });
}

// 获取用户详情
export function getUserDetail(id) {
  return request({
    url: `/admin/users/${id}`,
    method: 'get'
  });
}

// 更新用户状态
export function updateUserStatus(id, status) {
  return request({
    url: `/api/v1/admin/users/${id}/status`,
    method: 'put',
    data: { status }
  });
} 