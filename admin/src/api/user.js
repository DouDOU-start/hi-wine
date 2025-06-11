import { createApi } from '../utils/apiFactory';
import { get, post, put } from '../utils/request';

// 创建管理员API
const adminApi = createApi('/admin/users');

// 导出管理员API方法
export const {
  getList: getAdminList,
  create: createAdmin,
  update: updateAdmin,
  delete: deleteAdmin
} = adminApi;

// 管理员登录
export function login(data) {
  return post('/admin/auth/login', data);
}

// 获取当前管理员信息
export function getAdminInfo() {
  return get('/admin/profile');
}

// 重置管理员密码
export function resetAdminPassword(id, newPassword) {
  return post(`/admin/users/${id}/reset-password`, { newPassword });
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