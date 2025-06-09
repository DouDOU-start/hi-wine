import { get, post, put, del } from '../utils/request';

// 管理员登录
export function login(data) {
  return post('/admin/auth/login', data);
}

// 获取当前管理员信息
export function getAdminInfo() {
  return get('/admin/profile');
}

// 获取管理员列表
export function getAdminList(params) {
  return get('/admin/users', params);
}

// 创建管理员
export function createAdmin(data) {
  return post('/admin/users', data);
}

// 更新管理员信息
export function updateAdmin(id, data) {
  return put(`/admin/users/${id}`, data);
}

// 重置管理员密码
export function resetAdminPassword(id, newPassword) {
  return post(`/admin/users/${id}/reset-password`, { newPassword });
}

// 删除管理员
export function deleteAdmin(id) {
  return del(`/admin/users/${id}`);
}

// 获取普通用户列表
export function getUserList(params) {
  return get('/admin/users', params);
}

// 获取用户详情
export function getUserDetail(id) {
  return get(`/admin/users/${id}`);
}

// 更新用户状态
export function updateUserStatus(id, status) {
  return put(`/admin/users/${id}/status`, { status });
} 