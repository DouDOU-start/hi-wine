<template>
  <el-card shadow="hover" class="search-card">
    <div class="search-header">
      <div class="search-title">
        <el-icon><Search /></el-icon>
        <span>订单查询</span>
      </div>
      <div class="search-mode-toggle">
        <el-switch
          v-model="isAdvancedMode"
          active-text="高级搜索"
          inactive-text="基础搜索"
        />
      </div>
    </div>

    <el-form :inline="true" :model="searchForm" class="search-form" size="default">
      <el-form-item label="订单号">
        <el-input 
          v-model="searchForm.orderSn" 
          placeholder="请输入订单号" 
          clearable 
          style="width: 200px;"
        />
      </el-form-item>
      
      <el-form-item label="用户名">
        <el-input 
          v-model="searchForm.username" 
          placeholder="请输入用户名" 
          clearable 
          style="width: 200px;"
        />
      </el-form-item>
      
      <el-form-item label="订单状态">
        <el-select 
          v-model="searchForm.orderStatus" 
          placeholder="请选择订单状态" 
          clearable
          style="width: 200px;"
        >
          <el-option label="待支付" value="new" />
          <el-option label="已支付" value="processing" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="下单时间">
        <el-date-picker
          v-model="searchForm.dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 360px;"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>查询
        </el-button>
        <el-button @click="resetSearch">
          <el-icon><Refresh /></el-icon>重置
        </el-button>
        <el-button 
          v-if="showExport" 
          type="success" 
          @click="handleExport"
        >
          <el-icon><Download /></el-icon>导出订单
        </el-button>
        <el-button 
          v-if="isAdvancedMode" 
          type="info" 
          @click="saveSearchCondition"
        >
          <el-icon><Star /></el-icon>保存筛选条件
        </el-button>
      </el-form-item>
    </el-form>

    <!-- 高级搜索选项 -->
    <el-collapse-transition>
      <div v-show="isAdvancedMode" class="advanced-search">
        <el-form :inline="true" :model="searchForm" class="search-form" size="default">
          <el-form-item label="支付方式">
            <el-select 
              v-model="searchForm.paymentMethod" 
              placeholder="请选择支付方式" 
              clearable
              style="width: 200px;"
            >
              <el-option label="微信支付" value="wechat" />
              <el-option label="支付宝" value="alipay" />
              <el-option label="余额支付" value="balance" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="订单金额">
            <el-input-number v-model="searchForm.minAmount" placeholder="最低金额" :min="0" :precision="2" style="width: 120px;" />
            <span class="price-separator">至</span>
            <el-input-number v-model="searchForm.maxAmount" placeholder="最高金额" :min="0" :precision="2" style="width: 120px;" />
          </el-form-item>
          
          <el-form-item label="支付时间">
            <el-date-picker 
              v-model="searchForm.payTimeRange" 
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

    <!-- 保存的筛选条件列表 -->
    <div v-if="isAdvancedMode && savedConditions.length > 0" class="saved-conditions">
      <div class="conditions-title">保存的筛选条件:</div>
      <el-tag
        v-for="(condition, index) in savedConditions"
        :key="index"
        class="condition-tag"
        closable
        @click="applySavedCondition(condition)"
        @close="removeSavedCondition(index)"
      >
        {{ condition.name }}
      </el-tag>
    </div>
    
    <!-- 快速筛选标签 -->
    <div class="quick-filter">
      <span class="quick-filter-label">快速筛选:</span>
      <div class="filter-button-group">
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'all' }]" 
          size="small" 
          @click="applyQuickFilter('all')"
        >
          <el-icon><Tickets /></el-icon>全部订单
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'new' }]" 
          size="small" 
          @click="applyQuickFilter('new')"
          type="warning"
          plain
        >
          <el-icon><Timer /></el-icon>待支付
          <el-badge v-if="pendingCount > 0" :value="pendingCount" class="filter-badge" />
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'processing' }]" 
          size="small" 
          @click="applyQuickFilter('processing')"
          type="success"
          plain
        >
          <el-icon><Check /></el-icon>已支付
          <el-badge v-if="processingCount > 0" :value="processingCount" class="filter-badge" />
        </el-button>
        <el-button 
          :class="['filter-button', { 'active': activeFilter === 'completed' }]" 
          size="small" 
          @click="applyQuickFilter('completed')"
          type="primary"
          plain
        >
          <el-icon><CircleCheck /></el-icon>已完成
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, reactive, watch, defineEmits, defineProps, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  Search, Refresh, Download, User, Star, Timer, Check, CircleCheck, Tickets
} from '@element-plus/icons-vue';

const props = defineProps({
  showExport: {
    type: Boolean,
    default: true
  },
  initialConditions: {
    type: Object,
    default: () => ({})
  },
  pendingCount: {
    type: Number,
    default: 0
  },
  processingCount: {
    type: Number,
    default: 0
  },
  completedCount: {
    type: Number,
    default: 0
  }
});

const emit = defineEmits(['search', 'reset', 'export', 'filter']);

// 搜索模式
const isAdvancedMode = ref(false);

// 当前激活的快速筛选
const activeFilter = ref('all');

// 订单状态选项
const orderStatusOptions = {
  'new': '待支付',
  'processing': '已支付',
  'completed': '已完成',
  'cancelled': '已取消'
};

// 搜索表单
const searchForm = reactive({
  orderSn: '',
  username: '',
  orderStatus: '',
  dateRange: [],
  paymentMethod: '',
  minAmount: '',
  maxAmount: '',
  payTimeRange: []
});

// 保存的筛选条件
const savedConditions = ref([]);

// 初始化表单数据
const initFormData = () => {
  if (props.initialConditions) {
    Object.keys(props.initialConditions).forEach(key => {
      if (searchForm.hasOwnProperty(key)) {
        searchForm[key] = props.initialConditions[key];
      }
    });
  }
};

// 处理搜索
const handleSearch = () => {
  // 构建搜索参数
  const params = {
    orderSn: searchForm.orderSn,
    username: searchForm.username,
    orderStatus: searchForm.orderStatus,
    startDate: searchForm.dateRange && searchForm.dateRange[0],
    endDate: searchForm.dateRange && searchForm.dateRange[1]
  };

  // 如果是高级搜索模式，添加额外参数
  if (isAdvancedMode.value) {
    params.paymentMethod = searchForm.paymentMethod;
    params.minAmount = searchForm.minAmount;
    params.maxAmount = searchForm.maxAmount;
    params.payStartDate = searchForm.payTimeRange && searchForm.payTimeRange[0];
    params.payEndDate = searchForm.payTimeRange && searchForm.payTimeRange[1];
  }

  emit('search', params);
};

// 重置搜索
const resetSearch = () => {
  searchForm.orderSn = '';
  searchForm.username = '';
  searchForm.orderStatus = '';
  searchForm.dateRange = [];
  
  // 如果是高级搜索模式，重置额外字段
  if (isAdvancedMode.value) {
    searchForm.paymentMethod = '';
    searchForm.minAmount = '';
    searchForm.maxAmount = '';
    searchForm.payTimeRange = [];
  }
  
  activeFilter.value = 'all';
  emit('reset');
};

// 应用快速筛选
const applyQuickFilter = (filter) => {
  activeFilter.value = filter;
  
  // 重置当前表单
  resetSearch();
  
  // 根据筛选类型设置相应的搜索条件
  if (filter !== 'all') {
    searchForm.orderStatus = filter;
  }
  
  // 发送筛选事件
  emit('filter', { filter, orderStatus: filter });
  
  // 执行搜索
  handleSearch();
};

// 导出订单
const handleExport = () => {
  emit('export', {
    orderSn: searchForm.orderSn,
    username: searchForm.username,
    orderStatus: searchForm.orderStatus,
    startDate: searchForm.dateRange && searchForm.dateRange[0],
    endDate: searchForm.dateRange && searchForm.dateRange[1],
    paymentMethod: searchForm.paymentMethod,
    minAmount: searchForm.minAmount,
    maxAmount: searchForm.maxAmount,
    payStartDate: searchForm.payTimeRange && searchForm.payTimeRange[0],
    payEndDate: searchForm.payTimeRange && searchForm.payTimeRange[1]
  });
};

// 保存筛选条件
const saveSearchCondition = () => {
  ElMessageBox.prompt('请输入筛选条件名称', '保存筛选条件', {
    confirmButtonText: '保存',
    cancelButtonText: '取消',
    inputValidator: (value) => {
      if (!value) {
        return '名称不能为空';
      }
      return true;
    }
  }).then(({ value }) => {
    const condition = {
      name: value,
      data: JSON.parse(JSON.stringify(searchForm))
    };
    savedConditions.value.push(condition);
    
    // 保存到本地存储
    try {
      const existingConditions = JSON.parse(localStorage.getItem('orderSearchConditions') || '[]');
      existingConditions.push(condition);
      localStorage.setItem('orderSearchConditions', JSON.stringify(existingConditions));
      ElMessage.success('筛选条件保存成功');
    } catch (error) {
      console.error('保存筛选条件失败:', error);
    }
  }).catch(() => {});
};

// 应用保存的筛选条件
const applySavedCondition = (condition) => {
  Object.keys(condition.data).forEach(key => {
    searchForm[key] = condition.data[key];
  });
  handleSearch();
};

// 删除保存的筛选条件
const removeSavedCondition = (index) => {
  savedConditions.value.splice(index, 1);
  
  // 更新本地存储
  try {
    localStorage.setItem('orderSearchConditions', JSON.stringify(savedConditions.value));
  } catch (error) {
    console.error('更新筛选条件失败:', error);
  }
};

// 从本地存储加载保存的筛选条件
const loadSavedConditions = () => {
  try {
    const conditions = JSON.parse(localStorage.getItem('orderSearchConditions') || '[]');
    savedConditions.value = conditions;
  } catch (error) {
    console.error('加载筛选条件失败:', error);
  }
};

// 初始化
onMounted(() => {
  initFormData();
  loadSavedConditions();
});

// 监听初始条件变化
watch(() => props.initialConditions, (newVal) => {
  if (newVal) {
    initFormData();
  }
}, { deep: true });
</script>

<style>
/* 全局样式，确保下拉框正确显示 */
.el-select {
  width: 100%;
}

.el-select .el-input {
  width: 100%;
}

.el-select-dropdown {
  min-width: 150px !important;
}

.select-wrapper {
  width: 100%;
}

.select-wrapper .el-select {
  width: 100%;
  display: block;
}

.select-wrapper .el-input {
  width: 100%;
}
</style>

<style scoped>
.search-card {
  margin-bottom: 20px;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.search-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
}

.search-form {
  margin-bottom: 10px;
}

.search-buttons {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  gap: 10px;
  flex-wrap: wrap;
}

.saved-conditions {
  margin-top: 15px;
  border-top: 1px dashed #dcdfe6;
  padding-top: 10px;
  margin-bottom: 15px;
}

.conditions-title {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.condition-tag {
  margin-right: 8px;
  margin-bottom: 8px;
  cursor: pointer;
}

.advanced-search {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #ebeef5;
}

.price-separator {
  margin: 0 5px;
}

/* 快速筛选样式 */
.quick-filter {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px dashed #ebeef5;
}

.quick-filter-label {
  font-size: 14px;
  color: #606266;
  margin-right: 10px;
}

.filter-button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}

.filter-button {
  display: flex;
  align-items: center;
  gap: 5px;
}

.filter-button.active {
  font-weight: bold;
  border-color: currentColor;
}

.filter-badge {
  margin-left: 5px;
}

@media (max-width: 768px) {
  .search-buttons {
    justify-content: flex-start;
  }
  
  .filter-button-group {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style> 
