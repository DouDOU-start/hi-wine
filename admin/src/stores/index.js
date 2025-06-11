import { createPinia } from 'pinia';

// 导出Pinia实例
export const pinia = createPinia();

// 导出所有store
export * from './user';
export * from './theme'; 