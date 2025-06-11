<template>
  <div class="loading-state" v-if="loading">
    <div class="loading-content">
      <el-icon class="loading-icon"><Loading /></el-icon>
      <span class="loading-text">{{ text }}</span>
    </div>
  </div>
</template>

<script setup>
import { Loading } from '@element-plus/icons-vue';

defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: '加载中...'
  }
});
</script>

<style scoped>
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  width: 100%;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.loading-icon {
  font-size: 32px;
  color: #409EFF;
  animation: rotating 2s linear infinite;
}

.loading-text {
  margin-top: 10px;
  color: #909399;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 