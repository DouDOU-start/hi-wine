// Token相关操作
const TokenKey = 'Admin-Token';

// 获取Token
export function getToken() {
  return localStorage.getItem(TokenKey);
}

// 设置Token
export function setToken(token) {
  return localStorage.setItem(TokenKey, token);
}

// 移除Token
export function removeToken() {
  return localStorage.removeItem(TokenKey);
}

// 管理员信息相关操作
const AdminInfoKey = 'Admin-Info';

// 获取管理员信息
export function getAdminInfo() {
  const adminInfo = localStorage.getItem(AdminInfoKey);
  return adminInfo ? JSON.parse(adminInfo) : null;
}

// 设置管理员信息
export function setAdminInfo(adminInfo) {
  return localStorage.setItem(AdminInfoKey, JSON.stringify(adminInfo));
}

// 移除管理员信息
export function removeAdminInfo() {
  return localStorage.removeItem(AdminInfoKey);
}

// 清除所有认证信息
export function clearAuth() {
  removeToken();
  removeAdminInfo();
}

// 判断是否已登录
export function isLoggedIn() {
  return !!getToken();
} 