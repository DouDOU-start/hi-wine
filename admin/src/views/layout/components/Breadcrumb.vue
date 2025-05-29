<template>
  <el-breadcrumb separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="item.path">
        <span v-if="index === breadcrumbs.length - 1 || !item.redirect" class="no-redirect">
          {{ item.meta.title }}
        </span>
        <a v-else @click.prevent="handleLink(item)">{{ item.meta.title }}</a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const breadcrumbs = ref([]);

// 获取面包屑数据
const getBreadcrumbs = () => {
  // 过滤没有标题的路由
  const matched = route.matched.filter(item => item.meta && item.meta.title);
  
  // 如果不是首页，添加首页到面包屑
  if (matched[0].path !== '/dashboard') {
    matched.unshift({
      path: '/dashboard',
      redirect: '/dashboard',
      meta: { title: '首页' }
    });
  }
  
  breadcrumbs.value = matched;
};

// 处理点击面包屑
const handleLink = (item) => {
  const { redirect, path } = item;
  // 如果存在重定向，则优先使用重定向
  if (redirect) {
    router.push(redirect);
    return;
  }
  router.push(path);
};

// 监听路由变化
watch(
  () => route.path,
  () => getBreadcrumbs(),
  { immediate: true }
);
</script>

<style scoped>
.el-breadcrumb {
  font-size: 14px;
  line-height: 1.5;
}

.no-redirect {
  color: #97a8be;
  cursor: text;
}

/* 面包屑过渡动画 */
.breadcrumb-enter-active,
.breadcrumb-leave-active {
  transition: all 0.5s;
}

.breadcrumb-enter-from,
.breadcrumb-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.breadcrumb-leave-active {
  position: absolute;
}
</style> 