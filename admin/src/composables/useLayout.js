import { ref, provide, inject } from 'vue';

// 定义注入键
const SIDEBAR_COLLAPSED_KEY = Symbol('sidebarCollapsed');
const TOGGLE_SIDEBAR_KEY = Symbol('toggleSidebar');

/**
 * 提供布局状态
 * @returns {Object} 布局状态和方法
 */
export function provideLayout() {
  // 侧边栏折叠状态
  const isCollapsed = ref(false);
  
  // 初始化状态
  const initSidebarStatus = () => {
    const sidebarStatus = localStorage.getItem('sidebarStatus');
    isCollapsed.value = sidebarStatus === '1';
  };
  
  // 切换侧边栏状态
  const toggleSidebar = () => {
    isCollapsed.value = !isCollapsed.value;
    localStorage.setItem('sidebarStatus', isCollapsed.value ? '1' : '0');
  };
  
  // 提供状态和方法
  provide(SIDEBAR_COLLAPSED_KEY, isCollapsed);
  provide(TOGGLE_SIDEBAR_KEY, toggleSidebar);
  
  return {
    isCollapsed,
    toggleSidebar,
    initSidebarStatus
  };
}

/**
 * 使用布局状态
 * @returns {Object} 布局状态和方法
 */
export function useLayout() {
  const isCollapsed = inject(SIDEBAR_COLLAPSED_KEY);
  const toggleSidebar = inject(TOGGLE_SIDEBAR_KEY);
  
  if (isCollapsed === undefined || toggleSidebar === undefined) {
    throw new Error('useLayout must be used within a component that calls provideLayout');
  }
  
  return {
    isCollapsed,
    toggleSidebar
  };
} 