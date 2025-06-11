<template>
  <el-card shadow="hover" class="search-card">
    <el-form :inline="true" :model="formModel" class="search-form" size="default">
      <slot></slot>
      
      <div class="search-buttons">
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
          <el-button v-if="$slots.advanced" type="text" @click="toggleAdvanced">
            {{ showAdvanced ? '收起' : '高级筛选' }}
            <el-icon>
              <component :is="showAdvanced ? 'ArrowUp' : 'ArrowDown'" />
            </el-icon>
          </el-button>
        </el-form-item>
      </div>
      
      <!-- 高级搜索区域 -->
      <el-collapse-transition>
        <div v-show="showAdvanced" class="advanced-search">
          <slot name="advanced"></slot>
        </div>
      </el-collapse-transition>
    </el-form>
    
    <!-- 快速筛选标签 -->
    <div v-if="$slots.filters" class="quick-filter">
      <span class="quick-filter-label">快速筛选:</span>
      <div class="filter-button-group">
        <slot name="filters"></slot>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref } from 'vue';
import { Search, Refresh, ArrowUp, ArrowDown } from '@element-plus/icons-vue';

const props = defineProps({
  model: {
    type: Object,
    required: true
  },
  initialModel: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits(['search', 'reset']);

// 表单数据
const formModel = ref(props.model);

// 高级搜索显示状态
const showAdvanced = ref(false);

// 切换高级搜索
const toggleAdvanced = () => {
  showAdvanced.value = !showAdvanced.value;
};

// 处理搜索
const handleSearch = () => {
  emit('search', formModel.value);
};

// 处理重置
const handleReset = () => {
  // 重置为初始值
  Object.keys(formModel.value).forEach(key => {
    formModel.value[key] = props.initialModel[key] !== undefined 
      ? props.initialModel[key] 
      : undefined;
  });
  
  emit('reset', formModel.value);
};
</script>

<style scoped>
.search-card {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.search-buttons {
  margin-left: auto;
  display: flex;
  align-items: center;
}

.advanced-search {
  width: 100%;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #dcdfe6;
}

.quick-filter {
  display: flex;
  align-items: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.quick-filter-label {
  margin-right: 10px;
  color: #606266;
  font-size: 14px;
}

.filter-button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

@media (max-width: 768px) {
  .search-form {
    flex-direction: column;
  }
  
  .search-buttons {
    margin-left: 0;
    margin-top: 10px;
  }
}
</style> 