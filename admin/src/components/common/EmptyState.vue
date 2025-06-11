<template>
  <div class="empty-state">
    <el-empty 
      :image="image" 
      :image-size="imageSize" 
      :description="description"
    >
      <slot></slot>
    </el-empty>
  </div>
</template>

<script setup>
defineProps({
  image: {
    type: String,
    default: undefined
  },
  imageSize: {
    type: Number,
    default: 100
  },
  description: {
    type: String,
    default: '暂无数据'
  }
});
</script>

<style scoped>
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  width: 100%;
}
</style> 