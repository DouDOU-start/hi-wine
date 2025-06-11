import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const useThemeStore = defineStore('theme', () => {
  // 主题模式：'light' 或 'dark'
  const theme = ref(localStorage.getItem('theme') || 'light');
  
  // 监听主题变化
  watch(theme, (newTheme) => {
    // 保存到本地存储
    localStorage.setItem('theme', newTheme);
    
    // 应用主题
    applyTheme(newTheme);
  });
  
  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light';
  };
  
  // 设置主题
  const setTheme = (newTheme) => {
    if (newTheme === 'light' || newTheme === 'dark') {
      theme.value = newTheme;
    }
  };
  
  // 应用主题
  const applyTheme = (themeName) => {
    // 移除旧主题类
    document.documentElement.classList.remove('light-theme', 'dark-theme');
    
    // 添加新主题类
    document.documentElement.classList.add(`${themeName}-theme`);
    
    // 设置Element Plus主题
    if (themeName === 'dark') {
      document.documentElement.setAttribute('data-theme', 'dark');
    } else {
      document.documentElement.removeAttribute('data-theme');
    }
  };
  
  // 初始化主题
  const initTheme = () => {
    applyTheme(theme.value);
  };
  
  return {
    theme,
    toggleTheme,
    setTheme,
    initTheme
  };
}); 