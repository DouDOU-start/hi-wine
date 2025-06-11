<template>
  <search-form 
    :model="searchForm" 
    :initial-model="initialSearchForm"
    @search="handleSearch"
    @reset="handleReset"
  >
    <!-- 基础搜索项 -->
    <el-row :gutter="20">
      <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
        <form-item
          type="input"
          label="商品名称"
          v-model="searchForm.name"
          placeholder="请输入商品名称"
          clearable
        />
      </el-col>
      <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
        <form-item
          type="select"
          label="分类"
          v-model="searchForm.categoryId"
          placeholder="请选择分类"
          clearable
          filterable
          :options="categoryOptions"
        />
      </el-col>
      <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
        <form-item
          type="select"
          label="状态"
          v-model="searchForm.status"
          placeholder="请选择状态"
          clearable
          :options="statusOptions"
        />
      </el-col>
    </el-row>
    
    <!-- 高级搜索项 -->
    <template #advanced>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
          <el-form-item label="价格区间">
            <el-input-number v-model="searchForm.minPrice" placeholder="最低价" :min="0" :precision="2" style="width: 120px;" />
            <span class="price-separator">至</span>
            <el-input-number v-model="searchForm.maxPrice" placeholder="最高价" :min="0" :precision="2" style="width: 120px;" />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
          <form-item
            type="select"
            label="库存"
            v-model="searchForm.stockStatus"
            placeholder="库存状态"
            clearable
            :options="stockOptions"
          />
        </el-col>
        <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
          <form-item
            type="date"
            label="创建时间"
            v-model="searchForm.dateRange"
            date-type="daterange"
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
        </el-col>
      </el-row>
    </template>
    
    <!-- 快速筛选标签 -->
    <template #filters>
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
    </template>
  </search-form>
</template>

<script setup>
import { ref, defineEmits, defineProps } from 'vue';
import { Tickets, WarningFilled, RemoveFilled, Sell } from '@element-plus/icons-vue';

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
  status: '',
  minPrice: '',
  maxPrice: '',
  stockStatus: '',
  dateRange: []
};

// 搜索表单
const searchForm = ref({ ...initialSearchForm });

// 当前激活的快速筛选
const activeFilter = ref('all');

// 处理搜索
const handleSearch = () => {
  emit('search', searchForm.value);
};

// 处理重置
const handleReset = () => {
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
      searchForm.value.status = 1;
      break;
  }
  
  emit('filter', filter, searchForm.value);
  handleSearch();
};
</script>

<style scoped>
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
</style> 