<template>
  <div class="error-page">
    <div class="error-code">404</div>
    <div class="error-desc">抱歉，您访问的页面不存在</div>
    <div class="error-actions">
      <el-button type="primary" @click="goHome">返回首页</el-button>
      <el-button @click="goBack">返回上一页</el-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

// 返回首页
const goHome = () => {
  router.push('/');
};

// 返回上一页
const goBack = () => {
  router.go(-1);
};
</script>

<style scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #f5f7fa;
}

.error-code {
  font-size: 160px;
  font-weight: bold;
  color: #409eff;
  text-shadow: 0 4px 8px rgba(64, 158, 255, 0.3);
  margin-bottom: 20px;
}

.error-desc {
  font-size: 24px;
  color: #606266;
  margin-bottom: 40px;
}

.error-actions {
  display: flex;
  gap: 20px;
}
</style> 