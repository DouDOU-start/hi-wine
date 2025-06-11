<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">订单管理</div>
      <div class="order-stats">
        <el-tag type="info">总订单: {{ total }}</el-tag>
        <el-tag type="success">已完成: {{ completedCount }}</el-tag>
        <el-tag type="warning">待支付: {{ pendingCount }}</el-tag>
        <el-tag type="danger">已取消: {{ cancelledCount }}</el-tag>
      </div>
    </div>
    
    <el-card shadow="hover" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="订单号">
              <el-input v-model="searchForm.orderSn" placeholder="请输入订单号" clearable prefix-icon="Search" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="用户名">
              <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable prefix-icon="User" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="订单状态">
              <el-select v-model="searchForm.orderStatus" placeholder="请选择订单状态" clearable style="width: 100%">
                <el-option label="待支付" value="new" />
                <el-option label="已支付" value="processing" />
                <el-option label="已完成" value="completed" />
                <el-option label="已取消" value="cancelled" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="下单时间">
              <el-date-picker
                v-model="searchForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <div class="search-buttons">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>查询
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
          <el-button type="success" @click="exportOrderData">
            <el-icon><Download /></el-icon>导出订单
          </el-button>
        </div>
      </el-form>
    </el-card>
    
    <el-card shadow="hover" class="table-card">
      <el-table
        v-loading="loading"
        :data="orderList"
        border
        stripe
        style="width: 100%"
        :header-cell-style="{ background: '#f5f7fa' }"
      >
        <el-table-column prop="orderSn" label="订单号" min-width="120" fixed="left" />
        <el-table-column label="用户名" min-width="150">
          <template #default="scope">
            <div class="user-info">
              <span v-if="scope.row.userNickname" class="user-nickname">{{ scope.row.userNickname }}</span>
              <span v-else-if="scope.row.userName" class="user-name">{{ scope.row.userName }}</span>
              <span v-else class="user-id">用户ID: {{ scope.row.userId || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="订单金额" width="120">
          <template #default="scope">
            <span class="price-value" :class="{'price-zero': scope.row.totalAmount <= 0}">
              ￥{{ formatPrice(scope.row.totalAmount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="商品数量" width="100" align="center">
          <template #default="scope">
            <el-tag size="small" effect="plain" type="info">
              {{ scope.row.itemCount || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="订单状态" width="120">
          <template #default="scope">
            <el-tag :type="getOrderStatusType(scope.row.orderStatus)" effect="light">
              {{ getOrderStatusText(scope.row.orderStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="120">
          <template #default="scope">
            <span class="payment-method">
              {{ getPayMethodText(scope.row.paymentMethod) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="下单时间" min-width="180">
          <template #default="scope">
            <div class="time-info">
              <el-icon><Calendar /></el-icon>
              {{ formatDate(scope.row.createdAt) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="支付时间" min-width="180">
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
      
      <div class="pagination-container">
        <el-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderList, updateOrderStatus, exportOrders } from '../../api/order';
import { formatDate } from '../../utils/format';
import { 
  Search, Refresh, Calendar, Timer, View, Check, Close, 
  Download, User
} from '@element-plus/icons-vue';

const router = useRouter();

// 加载状态
const loading = ref(false);

// 订单列表数据
const orderList = ref([]);

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 搜索表单
const searchForm = reactive({
  orderSn: '',
  username: '',
  orderStatus: '',
  dateRange: []
});

// 格式化金额
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return parseFloat(price).toFixed(2);
};

// 统计数据
const completedCount = computed(() => {
  return orderList.value.filter(order => order.orderStatus === 'completed').length;
});

const pendingCount = computed(() => {
  return orderList.value.filter(order => order.orderStatus === 'new').length;
});

const cancelledCount = computed(() => {
  return orderList.value.filter(order => order.orderStatus === 'cancelled').length;
});

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      orderSn: searchForm.orderSn,
      username: searchForm.username,
      orderStatus: searchForm.orderStatus,
      startDate: searchForm.dateRange && searchForm.dateRange[0],
      endDate: searchForm.dateRange && searchForm.dateRange[1]
    };
    
    const response = await getOrderList(params);
    
    // 确保响应中包含列表数据和总数
    orderList.value = response.data.list || [];
    total.value = response.data.total || 0;
    
    // 对订单数据进行额外处理，确保字段一致性
    orderList.value = orderList.value.map(order => {
      // 统一订单状态字段名称
      if (order.status !== undefined && order.orderStatus === undefined) {
        order.orderStatus = order.status;
      }
      
      // 统一创建时间字段
      if (order.createTime !== undefined && order.createdAt === undefined) {
        order.createdAt = order.createTime;
      }
      
      // 统一支付时间字段
      if (order.payTime !== undefined && order.paidAt === undefined) {
        order.paidAt = order.payTime;
      }
      
      // 统一完成时间字段
      if (order.finishTime !== undefined && order.completedAt === undefined) {
        order.completedAt = order.finishTime;
      }
      
      return order;
    });
  } catch (error) {
    console.error('获取订单列表失败:', error);
    ElMessage.error('获取订单列表失败');
  } finally {
    loading.value = false;
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  const statusMap = {
    'new': '待支付',
    'processing': '已支付',
    'completed': '已完成',
    'cancelled': '已取消'
  };
  return statusMap[status] || '未知状态';
};

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const typeMap = {
    'new': 'warning',
    'processing': 'success',
    'completed': 'primary',
    'cancelled': 'info'
  };
  return typeMap[status] || 'info';
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

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchOrderList();
};

// 重置搜索
const resetSearch = () => {
  searchForm.orderSn = '';
  searchForm.username = '';
  searchForm.orderStatus = '';
  searchForm.dateRange = [];
  currentPage.value = 1;
  fetchOrderList();
};

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchOrderList();
};

// 页码变化
const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchOrderList();
};

// 查看订单详情
const handleViewDetail = (row) => {
  router.push(`/order/detail/${row.id}`);
};

// 完成订单
const handleCompleteOrder = (row) => {
  ElMessageBox.confirm('确认将此订单标记为已完成?', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(row.id, 'completed');
      ElMessage.success('订单已完成');
      fetchOrderList();
    } catch (error) {
      console.error('更新订单状态失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 取消订单
const handleCancelOrder = (row) => {
  ElMessageBox.confirm('确认取消此订单?', '警告', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(row.id, 'cancelled');
      ElMessage.success('订单已取消');
      fetchOrderList();
    } catch (error) {
      console.error('取消订单失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 导出订单数据
const exportOrderData = () => {
  ElMessageBox.confirm('确定要导出订单数据吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      const params = {
        orderSn: searchForm.orderSn,
        username: searchForm.username,
        orderStatus: searchForm.orderStatus,
        startDate: searchForm.dateRange && searchForm.dateRange[0],
        endDate: searchForm.dateRange && searchForm.dateRange[1]
      };
      
      await exportOrders(params);
      ElMessage.success('导出成功');
    } catch (error) {
      console.error('导出订单数据失败:', error);
      ElMessage.error('导出失败');
    }
  }).catch(() => {});
};

// 初始化
onMounted(() => {
  fetchOrderList();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 500;
  color: #303133;
}

.order-stats {
  display: flex;
  gap: 15px;
}

.search-card {
  margin-bottom: 20px;
}

.search-buttons {
  display: flex;
  justify-content: center;
  margin-top: 10px;
  gap: 10px;
}

.table-card {
  margin-bottom: 20px;
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
</style> 