<template>
  <div class="table-toolbar">
    <div class="left-actions">
      <slot name="left"></slot>
    </div>
    <div class="right-actions">
      <slot name="right"></slot>
      <el-tooltip content="刷新" placement="top">
        <el-button 
          circle 
          plain 
          @click="$emit('refresh')" 
          :loading="loading"
          class="refresh-btn"
        >
          <el-icon><Refresh /></el-icon>
        </el-button>
      </el-tooltip>
      <el-tooltip content="密度" placement="top">
        <el-dropdown trigger="click" @command="handleSizeChange">
          <el-button circle plain class="size-btn">
            <el-icon><Operation /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="'default'" :class="{ 'active': size === 'default' }">
                <el-icon v-if="size === 'default'"><Check /></el-icon>
                默认
              </el-dropdown-item>
              <el-dropdown-item :command="'medium'" :class="{ 'active': size === 'medium' }">
                <el-icon v-if="size === 'medium'"><Check /></el-icon>
                中等
              </el-dropdown-item>
              <el-dropdown-item :command="'small'" :class="{ 'active': size === 'small' }">
                <el-icon v-if="size === 'small'"><Check /></el-icon>
                紧凑
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-tooltip>
      <el-tooltip content="列设置" placement="top">
        <el-popover 
          placement="bottom" 
          trigger="click" 
          :width="200" 
          popper-class="column-popover"
        >
          <template #reference>
            <el-button circle plain class="column-btn">
              <el-icon><SetUp /></el-icon>
            </el-button>
          </template>
          <div class="column-list">
            <div class="column-list-header">
              <span>列展示</span>
              <el-button type="text" @click="resetColumns">重置</el-button>
            </div>
            <el-divider style="margin: 8px 0" />
            <el-checkbox
              v-model="allColumnsSelected"
              :indeterminate="isIndeterminate"
              @change="handleCheckAllChange"
            >
              全选
            </el-checkbox>
            <el-divider style="margin: 8px 0" />
            <el-checkbox-group v-model="selectedColumns" @change="handleColumnChange">
              <div v-for="col in columns" :key="col.prop" class="column-item">
                <el-checkbox :label="col.prop">{{ col.label }}</el-checkbox>
              </div>
            </el-checkbox-group>
          </div>
        </el-popover>
      </el-tooltip>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { Refresh, Operation, SetUp, Check } from '@element-plus/icons-vue';

const props = defineProps({
  columns: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  defaultSize: {
    type: String,
    default: 'default'
  }
});

const emit = defineEmits(['refresh', 'size-change', 'column-change']);

// 表格尺寸
const size = ref(props.defaultSize);

// 选中的列
const selectedColumns = ref(props.columns.map(col => col.prop));

// 是否全选
const allColumnsSelected = computed(() => {
  return selectedColumns.value.length === props.columns.length;
});

// 是否半选
const isIndeterminate = computed(() => {
  return selectedColumns.value.length > 0 && selectedColumns.value.length < props.columns.length;
});

// 处理全选
const handleCheckAllChange = (val) => {
  selectedColumns.value = val ? props.columns.map(col => col.prop) : [];
  emitColumnChange();
};

// 处理列选择变化
const handleColumnChange = () => {
  emitColumnChange();
};

// 发出列变化事件
const emitColumnChange = () => {
  const visibleColumns = props.columns.filter(col => 
    selectedColumns.value.includes(col.prop)
  );
  emit('column-change', visibleColumns);
};

// 重置列设置
const resetColumns = () => {
  selectedColumns.value = props.columns.map(col => col.prop);
  emitColumnChange();
};

// 处理表格尺寸变化
const handleSizeChange = (command) => {
  size.value = command;
  emit('size-change', command);
};

// 监听列变化
watch(() => props.columns, (newColumns) => {
  // 保持当前选中的列，并添加新列
  const currentSelected = [...selectedColumns.value];
  const newColumnProps = newColumns.map(col => col.prop);
  
  // 添加新列
  newColumnProps.forEach(prop => {
    if (!currentSelected.includes(prop)) {
      currentSelected.push(prop);
    }
  });
  
  // 过滤掉不存在的列
  selectedColumns.value = currentSelected.filter(prop => 
    newColumnProps.includes(prop)
  );
}, { deep: true });
</script>

<style scoped>
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.left-actions, .right-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.column-list {
  padding: 8px;
}

.column-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.column-item {
  margin: 8px 0;
}

:deep(.el-dropdown-menu__item.active) {
  color: #409EFF;
  font-weight: bold;
}
</style> 