<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useThemeStore } from './stores/theme';

// 初始化主题
onMounted(() => {
  const themeStore = useThemeStore();
  themeStore.initTheme();
});
</script>

<style>
/* 全局样式已移至 assets/css/main.css */
</style> 