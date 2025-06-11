<template>
  <el-tag 
    :type="getOrderStatusType" 
    :effect="effect" 
    :size="size"
    :class="{ 'status-tag': true, 'with-icon': showIcon }"
  >
    <el-icon v-if="showIcon" class="status-icon">
      <component :is="getOrderStatusIcon"></component>
    </el-icon>
    {{ getOrderStatusText }}
  </el-tag>
</template>

<script setup>
import { computed, defineProps } from 'vue';
import { 
  Clock, Check, CircleCheck, CircleClose 
} from '@element-plus/icons-vue';

const props = defineProps({
  status: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: 'default'
  },
  effect: {
    type: String,
    default: 'light'
  },
  showIcon: {
    type: Boolean,
    default: false
  }
});

// 获取订单状态文本
const getOrderStatusText = computed(() => {
  const statusMap = {
    'new': '待支付',
    'processing': '已支付',
    'completed': '已完成',
    'cancelled': '已取消'
  };
  return statusMap[props.status] || '未知状态';
});

// 获取订单状态类型
const getOrderStatusType = computed(() => {
  const typeMap = {
    'new': 'warning',
    'processing': 'success',
    'completed': 'primary',
    'cancelled': 'info'
  };
  return typeMap[props.status] || 'info';
});

// 获取订单状态图标
const getOrderStatusIcon = computed(() => {
  const iconMap = {
    'new': Clock,
    'processing': Check,
    'completed': CircleCheck,
    'cancelled': CircleClose
  };
  return iconMap[props.status] || Clock;
});
</script>

<style scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
}

.status-tag.with-icon {
  padding-left: 8px;
}

.status-icon {
  margin-right: 4px;
}
</style> 
