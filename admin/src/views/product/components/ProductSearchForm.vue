<template>
  <el-card shadow="hover" class="search-card">
    <el-form :inline="true" :model="searchForm" class="search-form" size="default">
      <!-- 基础搜索项 -->
      <el-form-item label="商品名称">
        <el-input 
          v-model="searchForm.name" 
          placeholder="请输入商品名称" 
          clearable
          style="width: 200px;"
        />
      </el-form-item>
      
      <el-form-item label="分类">
        <el-select 
          v-model="searchForm.categoryId" 
          placeholder="请选择分类" 
          clearable 
          filterable
          style="width: 200px;"
        >
          <el-option 
            v-for="item in categoryOptions" 
            :key="item.value" 
            :label="item.label" 
            :value="item.value" 
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="状态">
        <el-select 
          v-model="searchForm.isActive" 
          placeholder="请选择状态" 
          clearable
          style="width: 200px;"
        >
          <el-option 
            v-for="item in statusOptions" 
            :key="item.value" 
            :label="item.label" 
            :value="item.value" 
          />
        </el-select>
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>查询
        </el-button>
        <el-button @click="handleReset">
          <el-icon><Refresh /></el-icon>重置
        </el-button>
        <el-button type="text" @click="showAdvanced = !showAdvanced">
          {{ showAdvanced ? '收起' : '高级筛选' }}
          <el-icon>
            <component :is="showAdvanced ? 'ArrowUp' : 'ArrowDown'" />
          </el-icon>
        </el-button>
      </el-form-item>
    </el-form>
    
    <!-- 高级搜索区域 -->
    <el-collapse-transition>
      <div v-show="showAdvanced" class="advanced-search">
        <el-form :inline="true" :model="searchForm" class="search-form" size="default">
          <el-form-item label="价格区间">
            <el-input-number v-model="searchForm.minPrice" placeholder="最低价" :min="0" :precision="2" style="width: 120px;" />
            <span class="price-separator">至</span>
            <el-input-number v-model="searchForm.maxPrice" placeholder="最高价" :min="0" :precision="2" style="width: 120px;" />
          </el-form-item>
          
          <el-form-item label="库存">
            <el-select 
              v-model="searchForm.stockStatus" 
              placeholder="库存状态" 
              clearable
              style="width: 200px;"
            >
              <el-option 
                v-for="item in stockOptions" 
                :key="item.value" 
                :label="item.label" 
                :value="item.value" 
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="创建时间">
            <el-date-picker 
              v-model="searchForm.dateRange" 
              type="daterange" 
              value-format="YYYY-MM-DD" 
              range-separator="至" 
              start-placeholder="开始日期" 
              end-placeholder="结束日期"
              style="width: 360px;"
            />
          </el-form-item>
        </el-form>
      </div>
    </el-collapse-transition>
    
    <!-- 快速筛选标签 -->
    <div class="quick-filter">
      <span class="quick-filter-label">快速筛选:</span>
      <div class="filter-button-group">
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'all' }]" 
          size="small" 
          @click="applyQuickFilter('all')"
        >
          <el-icon><Tickets /></el-icon>全部商品
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'lowStock' }]" 
          size="small" 
          @click="applyQuickFilter('lowStock')"
          type="warning"
          plain
        >
          <el-icon><WarningFilled /></el-icon>库存不足
          <el-badge v-if="lowStockCount > 0" :value="lowStockCount" class="filter-badge" />
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'outOfStock' }]" 
          size="small" 
          @click="applyQuickFilter('outOfStock')"
          type="danger"
          plain
        >
          <el-icon><RemoveFilled /></el-icon>缺货商品
          <el-badge v-if="outOfStockCount > 0" :value="outOfStockCount" class="filter-badge" />
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'onSale' }]" 
          size="small" 
          @click="applyQuickFilter('onSale')"
          type="success"
          plain
        >
          <el-icon><Sell /></el-icon>在售商品
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, defineEmits, defineProps } from 'vue';
import { Tickets, WarningFilled, RemoveFilled, Sell, Search, Refresh, ArrowUp, ArrowDown } from '@element-plus/icons-vue';

const props = defineProps({
  categoryOptions: {
    type: Array,
    default: () => []
  },
  lowStockCount: {
    type: Number,
    default: 0
  },
  outOfStockCount: {
    type: Number,
    default: 0
  }
});

const emit = defineEmits(['search', 'reset', 'filter']);

// 状态选项
const statusOptions = [
  { label: '上架', value: 1 },
  { label: '下架', value: 0 }
];

// 库存选项
const stockOptions = [
  { label: '充足', value: 'normal' },
  { label: '不足', value: 'low' },
  { label: '缺货', value: 'out' }
];

// 初始搜索表单
const initialSearchForm = {
  name: '',
  categoryId: '',
  isActive: '',
  minPrice: '',
  maxPrice: '',
  stockStatus: '',
  dateRange: []
};

// 搜索表单
const searchForm = ref({ ...initialSearchForm });

// 当前激活的快速筛选
const activeFilter = ref('all');

// 高级搜索显示状态
const showAdvanced = ref(false);

// 处理搜索
const handleSearch = () => {
  emit('search', searchForm.value);
};

// 处理重置
const handleReset = () => {
  // 重置表单
  Object.keys(searchForm.value).forEach(key => {
    searchForm.value[key] = initialSearchForm[key];
  });
  
  activeFilter.value = 'all';
  emit('reset');
};

// 应用快速筛选
const applyQuickFilter = (filter) => {
  activeFilter.value = filter;
  
  // 重置搜索表单
  Object.keys(searchForm.value).forEach(key => {
    searchForm.value[key] = initialSearchForm[key];
  });
  
  // 根据筛选类型设置搜索条件
  switch (filter) {
    case 'lowStock':
      searchForm.value.stockStatus = 'low';
      break;
    case 'outOfStock':
      searchForm.value.stockStatus = 'out';
      break;
    case 'onSale':
      searchForm.value.isActive = 1;
      break;
  }
  
  emit('filter', filter, searchForm.value);
  handleSearch();
};
</script>

<style scoped>
.search-card {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
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

.price-separator {
  margin: 0 5px;
}

.filter-button {
  display: flex;
  align-items: center;
  gap: 5px;
}

.filter-button.active {
  font-weight: bold;
}

.filter-badge {
  margin-left: 5px;
}

@media (max-width: 768px) {
  .search-form {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .el-form-item {
    margin-right: 0;
    margin-bottom: 10px;
    width: 100%;
  }
}
</style> 