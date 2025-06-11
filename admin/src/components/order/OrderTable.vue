<template>
  <div class="order-table-container">
    <div class="table-header">
      <div class="view-switch">
        <el-radio-group v-model="viewMode" size="small">
          <el-radio-button label="table">
            <el-icon><Grid /></el-icon>表格视图
          </el-radio-button>
          <el-radio-button label="card">
            <el-icon><Menu /></el-icon>卡片视图
          </el-radio-button>
        </el-radio-group>
      </div>
      
      <div class="table-actions">
        <el-tooltip content="列设置" placement="top">
          <el-button @click="showColumnSelector = true" circle>
            <el-icon><Setting /></el-icon>
          </el-button>
        </el-tooltip>
        
        <el-tooltip content="刷新" placement="top">
          <el-button @click="$emit('refresh')" circle>
            <el-icon><Refresh /></el-icon>
          </el-button>
        </el-tooltip>
        
        <el-tooltip v-if="selectedRows.length > 0" content="批量操作" placement="top">
          <el-dropdown @command="handleBatchCommand" trigger="click">
            <el-button type="primary">
              批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="complete" v-if="hasProcessingOrders">标记为已完成</el-dropdown-item>
                <el-dropdown-item command="cancel" v-if="hasNewOrders">标记为已取消</el-dropdown-item>
                <el-dropdown-item command="export">导出所选订单</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-tooltip>
      </div>
    </div>
    
    <!-- 表格视图 -->
    <div v-show="viewMode === 'table'">
      <el-table
        v-loading="loading"
        :data="orders"
        border
        stripe
        style="width: 100%"
        :header-cell-style="{ background: '#f5f7fa' }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" fixed="left" />
        
        <el-table-column v-if="columns.orderSn.visible" prop="orderSn" label="订单号" min-width="120" fixed="left" sortable />
        
        <el-table-column v-if="columns.username.visible" label="用户名" min-width="150">
          <template #default="scope">
            <div class="user-info">
              <span v-if="scope.row.userNickname" class="user-nickname">{{ scope.row.userNickname }}</span>
              <span v-else-if="scope.row.userName" class="user-name">{{ scope.row.userName }}</span>
              <span v-else class="user-id">用户ID: {{ scope.row.userId || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.totalAmount.visible" label="订单金额" width="120" sortable>
          <template #default="scope">
            <span class="price-value" :class="{'price-zero': scope.row.totalAmount <= 0}">
              ￥{{ formatPrice(scope.row.totalAmount) }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.itemCount.visible" label="商品数量" width="100" align="center" sortable>
          <template #default="scope">
            <el-tag size="small" effect="plain" type="info">
              {{ scope.row.itemCount || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.orderStatus.visible" label="订单状态" width="120" sortable>
          <template #default="scope">
            <order-status-tag :status="scope.row.orderStatus" show-icon />
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.paymentMethod.visible" label="支付方式" width="120">
          <template #default="scope">
            <span class="payment-method">
              {{ getPayMethodText(scope.row.paymentMethod) }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.createdAt.visible" label="下单时间" min-width="180" sortable>
          <template #default="scope">
            <div class="time-info">
              <el-icon><Calendar /></el-icon>
              {{ formatDate(scope.row.createdAt) }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column v-if="columns.paidAt.visible" label="支付时间" min-width="180" sortable>
          <template #default="scope">
            <div class="time-info">
              <el-icon><Timer /></el-icon>
              {{ formatDate(scope.row.paidAt) }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                type="primary" 
                link 
                @click="handleViewDetail(scope.row)"
              >
                <el-icon><View /></el-icon>详情
              </el-button>
              <el-button 
                v-if="scope.row.orderStatus === 'processing'"
                type="success" 
                link 
                @click="handleCompleteOrder(scope.row)"
              >
                <el-icon><Check /></el-icon>完成
              </el-button>
              <el-button 
                v-if="scope.row.orderStatus === 'new'"
                type="danger" 
                link 
                @click="handleCancelOrder(scope.row)"
              >
                <el-icon><Close /></el-icon>取消
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 卡片视图 -->
    <div v-show="viewMode === 'card'" class="card-view">
      <el-checkbox-group v-model="selectedCardIds" class="card-checkbox-group">
        <el-card 
          v-for="order in orders" 
          :key="order.id" 
          shadow="hover" 
          class="order-card"
          :class="{ 'selected': selectedCardIds.includes(order.id) }"
        >
          <template #header>
            <div class="card-header">
              <el-checkbox 
                :value="order.id" 
                :label="order.id"
                @change="updateCardSelection"
              />
              <div class="order-sn">订单号: {{ order.orderSn }}</div>
              <order-status-tag :status="order.orderStatus" show-icon />
            </div>
          </template>
          
          <div class="card-body">
            <div class="card-row">
              <span class="card-label">用户:</span>
              <span class="card-value">{{ order.userNickname || order.userName || `用户ID: ${order.userId || '-'}` }}</span>
            </div>
            
            <div class="card-row">
              <span class="card-label">金额:</span>
              <span class="card-value price">￥{{ formatPrice(order.totalAmount) }}</span>
            </div>
            
            <div class="card-row">
              <span class="card-label">商品数量:</span>
              <span class="card-value">{{ order.itemCount || 0 }}</span>
            </div>
            
            <div class="card-row">
              <span class="card-label">下单时间:</span>
              <span class="card-value">{{ formatDate(order.createdAt) }}</span>
            </div>
            
            <div class="card-row">
              <span class="card-label">支付方式:</span>
              <span class="card-value">{{ getPayMethodText(order.paymentMethod) }}</span>
            </div>
          </div>
          
          <div class="card-footer">
            <el-button 
              type="primary" 
              size="small" 
              @click="handleViewDetail(order)"
            >
              <el-icon><View /></el-icon>详情
            </el-button>
            <el-button 
              v-if="order.orderStatus === 'processing'"
              type="success" 
              size="small" 
              @click="handleCompleteOrder(order)"
            >
              <el-icon><Check /></el-icon>完成
            </el-button>
            <el-button 
              v-if="order.orderStatus === 'new'"
              type="danger" 
              size="small" 
              @click="handleCancelOrder(order)"
            >
              <el-icon><Close /></el-icon>取消
            </el-button>
          </div>
        </el-card>
      </el-checkbox-group>
    </div>
    
    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        @update:current-page="emit('update:currentPage', $event)"
        @update:page-size="emit('update:pageSize', $event)"
        background
      />
    </div>
    
    <!-- 列选择器对话框 -->
    <el-dialog
      v-model="showColumnSelector"
      title="列设置"
      width="300px"
    >
      <el-checkbox-group v-model="visibleColumns">
        <div class="column-item" v-for="(col, key) in columns" :key="key">
          <el-checkbox :label="key">{{ col.label }}</el-checkbox>
        </div>
      </el-checkbox-group>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="resetColumns">重置</el-button>
          <el-button type="primary" @click="showColumnSelector = false">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, defineProps, defineEmits } from 'vue';
import { formatDate } from '../../utils/format';
import OrderStatusTag from './OrderStatusTag.vue';
import { 
  Grid, Menu, Setting, Refresh, ArrowDown, Calendar, Timer, 
  View, Check, Close 
} from '@element-plus/icons-vue';

const props = defineProps({
  orders: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  }
});

const emit = defineEmits([
  'view-detail', 
  'complete-order', 
  'cancel-order', 
  'size-change', 
  'current-change', 
  'refresh',
  'batch-complete',
  'batch-cancel',
  'batch-export',
  'update:currentPage',
  'update:pageSize'
]);

// 视图模式
const viewMode = ref('table');

// 列设置
const columns = reactive({
  orderSn: { visible: true, label: '订单号' },
  username: { visible: true, label: '用户名' },
  totalAmount: { visible: true, label: '订单金额' },
  itemCount: { visible: true, label: '商品数量' },
  orderStatus: { visible: true, label: '订单状态' },
  paymentMethod: { visible: true, label: '支付方式' },
  createdAt: { visible: true, label: '下单时间' },
  paidAt: { visible: true, label: '支付时间' }
});

// 可见列
const visibleColumns = ref(Object.keys(columns).filter(key => columns[key].visible));

// 显示列选择器
const showColumnSelector = ref(false);

// 选中的行
const selectedRows = ref([]);

// 卡片视图选中的ID
const selectedCardIds = ref([]);

// 格式化价格
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return parseFloat(price).toFixed(2);
};

// 获取支付方式文本
const getPayMethodText = (method) => {
  const methodMap = {
    'wechat': '微信支付',
    'alipay': '支付宝',
    'balance': '余额支付'
  };
  return methodMap[method] || '-';
};

// 是否有待支付订单
const hasNewOrders = computed(() => {
  return selectedRows.value.some(row => row.orderStatus === 'new');
});

// 是否有已支付订单
const hasProcessingOrders = computed(() => {
  return selectedRows.value.some(row => row.orderStatus === 'processing');
});

// 查看订单详情
const handleViewDetail = (row) => {
  emit('view-detail', row);
};

// 完成订单
const handleCompleteOrder = (row) => {
  emit('complete-order', row);
};

// 取消订单
const handleCancelOrder = (row) => {
  emit('cancel-order', row);
};

// 分页大小变化
const handleSizeChange = (size) => {
  emit('update:pageSize', size);
  emit('size-change', size);
};

// 页码变化
const handleCurrentChange = (page) => {
  emit('update:currentPage', page);
  emit('current-change', page);
};

// 表格选择变化
const handleSelectionChange = (rows) => {
  selectedRows.value = rows;
  // 同步卡片视图选择
  selectedCardIds.value = rows.map(row => row.id);
};

// 更新卡片选择
const updateCardSelection = () => {
  // 根据selectedCardIds更新selectedRows
  selectedRows.value = props.orders.filter(order => 
    selectedCardIds.value.includes(order.id)
  );
};

// 批量操作
const handleBatchCommand = (command) => {
  switch (command) {
    case 'complete':
      emit('batch-complete', selectedRows.value);
      break;
    case 'cancel':
      emit('batch-cancel', selectedRows.value);
      break;
    case 'export':
      emit('batch-export', selectedRows.value);
      break;
  }
};

// 重置列设置
const resetColumns = () => {
  Object.keys(columns).forEach(key => {
    columns[key].visible = true;
  });
  visibleColumns.value = Object.keys(columns);
};

// 监听可见列变化
watch(visibleColumns, (newVal) => {
  // 更新列可见性
  Object.keys(columns).forEach(key => {
    columns[key].visible = newVal.includes(key);
  });
}, { deep: true });
</script>

<style scoped>
.order-table-container {
  width: 100%;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.table-actions {
  display: flex;
  gap: 10px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-nickname {
  font-weight: 500;
}

.user-name, .user-id {
  font-size: 12px;
  color: #909399;
}

.price-value {
  color: #f56c6c;
  font-weight: bold;
}

.price-zero {
  color: #909399;
}

.payment-method {
  color: #606266;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #606266;
}

.action-buttons {
  display: flex;
  justify-content: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.column-item {
  margin: 10px 0;
}

/* 卡片视图样式 */
.card-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
}

.order-card {
  transition: all 0.3s;
}

.order-card.selected {
  border: 1px solid #409EFF;
  box-shadow: 0 0 10px rgba(64, 158, 255, 0.3);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.order-sn {
  flex: 1;
  margin: 0 10px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-body {
  margin-bottom: 15px;
}

.card-row {
  display: flex;
  margin-bottom: 8px;
  font-size: 14px;
}

.card-label {
  color: #606266;
  width: 80px;
  flex-shrink: 0;
}

.card-value {
  flex: 1;
}

.card-value.price {
  color: #f56c6c;
  font-weight: bold;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  border-top: 1px solid #ebeef5;
  padding-top: 10px;
}

@media (max-width: 768px) {
  .table-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .table-actions {
    align-self: flex-end;
  }
  
  .card-view {
    grid-template-columns: 1fr;
  }
}
</style> 
