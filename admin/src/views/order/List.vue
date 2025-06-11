<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">订单管理</div>
    </div>
    
    <!-- 订单统计组件 -->
    <order-stats 
      :stats="statsData" 
      @refresh="fetchOrderStats"
    />
    
    <!-- 订单搜索组件 -->
    <order-search-form 
      :initial-conditions="searchForm"
      :pending-count="statsData.pending"
      :processing-count="statsData.processing"
      :completed-count="statsData.completed"
      @search="handleSearch"
      @reset="resetSearch"
      @export="exportOrderData"
      @filter="handleFilter"
    />
    
    <!-- 订单表格组件 -->
    <el-card shadow="hover" class="table-card">
      <order-table
        :orders="orderList"
        :loading="loading"
        :total="total"
        :current-page="currentPage"
        :page-size="pageSize"
        @view-detail="handleViewDetail"
        @complete-order="handleCompleteOrder"
        @cancel-order="handleCancelOrder"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        @refresh="fetchOrderList"
        @batch-complete="handleBatchComplete"
        @batch-cancel="handleBatchCancel"
        @batch-export="handleBatchExport"
        @update:current-page="currentPage = $event"
        @update:page-size="pageSize = $event"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderList, updateOrderStatus, exportOrders, getOrderStats } from '../../api/order';
import OrderSearchForm from '../../components/order/OrderSearchForm.vue';
import OrderTable from '../../components/order/OrderTable.vue';
import OrderStats from '../../components/order/OrderStats.vue';

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

// 统计数据
const statsData = reactive({
  total: 0,
  completed: 0,
  pending: 0,
  processing: 0,
  cancelled: 0,
  todayOrders: 0,
  todaySales: 0,
  monthOrders: 0,
  monthSales: 0
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
      
      // 确保订单状态值与组件中使用的一致
      if (order.orderStatus) {
        // 将可能的其他状态值映射到我们使用的状态值
        const statusMap = {
          '待支付': 'new',
          '已支付': 'processing',
          '已完成': 'completed',
          '已取消': 'cancelled',
          'pending': 'new',
          'paid': 'processing',
          'complete': 'completed',
          'cancel': 'cancelled'
        };
        
        if (statusMap[order.orderStatus]) {
          order.orderStatus = statusMap[order.orderStatus];
        }
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

// 获取订单统计数据
const fetchOrderStats = async () => {
  try {
    const response = await getOrderStats();
    if (response && response.data) {
      const data = response.data;
      
      // 更新统计数据
      statsData.total = data.total || 0;
      statsData.completed = data.completed || 0;
      statsData.pending = data.new || 0; // 后端可能使用new表示待支付
      statsData.processing = data.processing || 0;
      statsData.cancelled = data.cancelled || 0;
      statsData.todayOrders = data.todayOrders || 0;
      statsData.todaySales = data.todaySales || 0;
      statsData.monthOrders = data.monthOrders || 0;
      statsData.monthSales = data.monthSales || 0;
    }
  } catch (error) {
    console.error('获取订单统计数据失败:', error);
  }
};

// 搜索
const handleSearch = (params) => {
  // 更新搜索表单
  Object.keys(params).forEach(key => {
    if (searchForm.hasOwnProperty(key)) {
      searchForm[key] = params[key];
    }
  });
  
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

// 处理快速筛选
const handleFilter = (filterData) => {
  searchForm.orderStatus = filterData.orderStatus;
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
      fetchOrderStats(); // 更新统计数据
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
      fetchOrderStats(); // 更新统计数据
    } catch (error) {
      console.error('取消订单失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 批量完成订单
const handleBatchComplete = (rows) => {
  if (rows.length === 0) {
    ElMessage.warning('请选择要操作的订单');
    return;
  }
  
  ElMessageBox.confirm(`确认将选中的 ${rows.length} 个订单标记为已完成?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const promises = rows.map(row => updateOrderStatus(row.id, 'completed'));
      await Promise.all(promises);
      ElMessage.success(`已成功完成 ${rows.length} 个订单`);
      fetchOrderList();
      fetchOrderStats(); // 更新统计数据
    } catch (error) {
      console.error('批量更新订单状态失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 批量取消订单
const handleBatchCancel = (rows) => {
  if (rows.length === 0) {
    ElMessage.warning('请选择要操作的订单');
    return;
  }
  
  ElMessageBox.confirm(`确认取消选中的 ${rows.length} 个订单?`, '警告', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const promises = rows.map(row => updateOrderStatus(row.id, 'cancelled'));
      await Promise.all(promises);
      ElMessage.success(`已成功取消 ${rows.length} 个订单`);
      fetchOrderList();
      fetchOrderStats(); // 更新统计数据
    } catch (error) {
      console.error('批量取消订单失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 批量导出订单
const handleBatchExport = (rows) => {
  if (rows.length === 0) {
    ElMessage.warning('请选择要导出的订单');
    return;
  }
  
  ElMessageBox.confirm(`确定要导出选中的 ${rows.length} 个订单数据吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      const orderIds = rows.map(row => row.id);
      await exportOrders({ ids: orderIds.join(',') });
      ElMessage.success('导出成功');
    } catch (error) {
      console.error('导出订单数据失败:', error);
      ElMessage.error('导出失败');
    }
  }).catch(() => {});
};

// 导出订单数据
const exportOrderData = (params) => {
  ElMessageBox.confirm('确定要导出订单数据吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
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
  fetchOrderStats();
});
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
</style>

<style scoped>
.page-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 500;
  color: #303133;
}

.table-card {
  margin-bottom: 20px;
}
</style> 