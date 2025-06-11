<template>
  <div class="page-header">
    <div class="page-title">{{ title }}</div>
    <div v-if="$slots.actions || $slots.extra" class="page-actions">
      <slot name="actions"></slot>
      <slot name="extra"></slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  }
});
</script>

<style scoped>
.page-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 500;
  color: #303133;
  line-height: 1.5;
}

.page-actions {
  display: flex;
  align-items: center;
}
</style> 