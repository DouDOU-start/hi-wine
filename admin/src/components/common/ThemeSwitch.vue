<template>
  <div class="theme-switch">
    <el-tooltip :content="theme === 'light' ? '切换到暗黑模式' : '切换到明亮模式'" placement="bottom">
      <el-button circle @click="toggleTheme">
        <el-icon>
          <component :is="theme === 'light' ? 'Moon' : 'Sunny'" />
        </el-icon>
      </el-button>
    </el-tooltip>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useThemeStore } from '@/stores/theme';

// 获取主题状态
const themeStore = useThemeStore();
const { theme, toggleTheme, initTheme } = themeStore;

// 初始化主题
onMounted(() => {
  initTheme();
});
</script>

<style scoped>
.theme-switch {
  display: inline-flex;
  align-items: center;
}
</style> 