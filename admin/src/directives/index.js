import permission from './permission';

// 注册所有指令
export function registerDirectives(app) {
  // 权限指令
  app.directive('permission', permission);
} 