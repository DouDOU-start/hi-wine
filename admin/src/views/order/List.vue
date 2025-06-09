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
              <el-input v-model="searchForm.orderId" placeholder="请输入订单号" clearable prefix-icon="Search" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="用户名">
              <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable prefix-icon="User" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="订单状态">
              <el-select v-model="searchForm.status" placeholder="请选择订单状态" clearable style="width: 100%">
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
        <el-table-column prop="order_sn" label="订单号" min-width="120" fixed="left" />
        <el-table-column label="用户名" min-width="150">
          <template #default="scope">
            <div class="user-info">
              <span v-if="scope.row.nickname" class="user-nickname">{{ scope.row.nickname }}</span>
              <span v-if="scope.row.username" class="user-name">{{ scope.row.username }}</span>
              <span v-if="!scope.row.nickname && !scope.row.username" class="user-id">用户ID: {{ scope.row.userId || '-' }}</span>
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
            <el-tag :type="getOrderStatusType(scope.row.status)" effect="light">
              {{ getOrderStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="120">
          <template #default="scope">
            <span class="payment-method">
              {{ getPayMethodText(scope.row.payMethod) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="下单时间" min-width="180">
          <template #default="scope">
            <div class="time-info">
              <el-icon><Calendar /></el-icon>
              {{ formatDateTime(scope.row.createTime) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="支付时间" min-width="180">
          <template #default="scope">
            <div class="time-info">
              <el-icon><Timer /></el-icon>
              {{ formatDateTime(scope.row.payTime) }}
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
                v-if="scope.row.status === 'processing' || scope.row.status === 1"
                type="success" 
                link 
                @click="handleCompleteOrder(scope.row)"
              >
                <el-icon><Check /></el-icon>完成
              </el-button>
              <el-button 
                v-if="scope.row.status === 'new' || scope.row.status === 0"
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
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
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
import { ref, reactive, onMounted, computed, onActivated } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderList, updateOrderStatus, exportOrders } from '../../api/order';
import { 
  Search, Refresh, Calendar, Timer, View, Check, Close, 
  Download, User
} from '@element-plus/icons-vue';

const router = useRouter();

// 防止重复请求的锁
const isRequestLocked = ref(false);
// 记录页面是否已经初始化
const isInitialized = ref(false);

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
  orderId: '',
  username: '',
  status: '',
  dateRange: []
});

// 获取订单列表
const fetchOrderList = async () => {
  if (isRequestLocked.value) return;
  isRequestLocked.value = true;

  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      orderId: searchForm.orderId,
      username: searchForm.username,
      status: searchForm.status,
      startDate: searchForm.dateRange && searchForm.dateRange[0],
      endDate: searchForm.dateRange && searchForm.dateRange[1]
    };
    
    const response = await getOrderList(params);
    console.log('订单列表响应数据:', response);
    
    // 根据后端实际返回的数据结构调整
    if (response && response.data) {
      // 检查response.data.list是否存在，如果不存在，可能数据在response.data中
      if (Array.isArray(response.data.list)) {
        orderList.value = response.data.list;
        total.value = response.data.total || 0;
      } else if (Array.isArray(response.data)) {
        // 如果response.data直接是数组，使用它作为订单列表
        orderList.value = response.data;
        total.value = response.total || response.data.length || 0;
      } else {
        // 尝试从截图中看到的数据结构提取
        if (response.data.list) {
          orderList.value = response.data.list;
          total.value = response.data.total || 0;
        } else {
          console.error('无法识别的订单数据结构:', response);
          orderList.value = [];
          total.value = 0;
        }
      }
    } else {
      orderList.value = [];
      total.value = 0;
    }
    
    // 数据处理：确保每个订单对象的属性都有正确的类型
    orderList.value = orderList.value.map(order => {
      return {
        ...order,
        // 确保totalAmount是数字类型，优先使用total_amount字段
        totalAmount: order.total_amount !== undefined 
          ? (typeof order.total_amount === 'string' ? parseFloat(order.total_amount) : order.total_amount) 
          : (typeof order.totalAmount === 'string' ? parseFloat(order.totalAmount) : (order.totalAmount || 0)),
        // 使用order_status字段作为状态
        status: order.order_status || order.status || 'new',
        // 确保商品数量正确，优先使用item_count字段
        itemCount: order.item_count !== undefined
          ? (typeof order.item_count === 'string' ? parseInt(order.item_count) : order.item_count)
          : (typeof order.productCount === 'string' ? parseInt(order.productCount) : (order.productCount || 0)),
        // 处理创建时间
        createTime: order.created_at || order.createTime || '',
        // 处理支付时间
        payTime: order.paid_at || order.payTime || '',
        // ID
        id: order.id,
        // 订单号
        order_sn: order.order_sn || order.id || '',
        // 用户信息处理
        username: order.user_name || '',
        nickname: order.user_nickname || '',
        userId: order.user_id || '',
        // 支付方式
        payMethod: order.payment_status === 'paid' ? 1 : 0
      };
    });
  } catch (error) {
    console.error('获取订单列表失败:', error);
    ElMessage.error('获取订单列表失败');
  } finally {
    loading.value = false;
    isRequestLocked.value = false;
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  switch (status) {
    case 'new': return '待支付';
    case 'processing': return '已支付';
    case 'completed': return '已完成';
    case 'cancelled': return '已取消';
    // 兼容数字状态码
    case 0: return '待支付';
    case 1: return '已支付';
    case 2: return '已完成';
    case 3: return '已取消';
    default: return '未知状态';
  }
};

// 获取订单状态类型
const getOrderStatusType = (status) => {
  switch (status) {
    case 'new': return 'warning';
    case 'processing': return 'success';
    case 'completed': return 'primary';
    case 'cancelled': return 'info';
    // 兼容数字状态码
    case 0: return 'warning';
    case 1: return 'success';
    case 2: return 'primary';
    case 3: return 'info';
    default: return 'info';
  }
};

// 获取支付方式文本
const getPayMethodText = (method) => {
  switch (method) {
    case 1: return '微信支付';
    case 2: return '支付宝';
    case 3: return '余额支付';
    default: return '-';
  }
};

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchOrderList();
};

// 重置搜索
const resetSearch = () => {
  searchForm.orderId = '';
  searchForm.username = '';
  searchForm.status = '';
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

  // 从截图可以看出，后端期望接收的是数字ID（如红框中的18）
  // 从订单号或ID中提取数字ID
  let numericId;
  
  if (row.id) {
    // 如果已经是数字，直接使用
    if (typeof row.id === 'number') {
      numericId = row.id;
    } else {
      // 尝试从id字段中提取数字
      const match = String(row.id).match(/\d+/);
      numericId = match ? parseInt(match[0]) : null;
    }
  }
  
  // 如果没有找到数字ID，则尝试从order_sn中提取
  if (!numericId && row.order_sn) {
    const match = String(row.order_sn).match(/\d+/);
    numericId = match ? parseInt(match[0]) : null;
  }
  
  if (!numericId) {
    ElMessage.warning('无法获取有效的订单ID');
    return;
  }
  
  console.log('查看订单详情，数字ID:', numericId);
  router.push(`/order/detail/${numericId}`);
};

// 完成订单
const handleCompleteOrder = (row) => {
  ElMessageBox.confirm(`确定要将订单 ${row.id} 标记为已完成吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 根据后端API调整状态值
      await updateOrderStatus(row.id, 'completed');
      ElMessage.success('订单已完成');
      fetchOrderList();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 取消订单
const handleCancelOrder = (row) => {
  ElMessageBox.confirm(`确定要取消订单 ${row.order_sn} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 根据后端API调整状态值
      await updateOrderStatus(row.id, 'cancelled_cancelled');
      ElMessage.success('订单已取消');
      fetchOrderList();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 导出订单数据
const exportOrderData = () => {
  ElMessageBox.confirm('确定要导出订单数据吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const params = {
        orderId: searchForm.orderId,
        username: searchForm.username,
        status: searchForm.status,
        startDate: searchForm.dateRange && searchForm.dateRange[0],
        endDate: searchForm.dateRange && searchForm.dateRange[1]
      };
      
      await exportOrders(params);
      ElMessage.success('导出成功，文件已开始下载');
    } catch (error) {
      console.error('导出失败:', error);
      ElMessage.error('导出失败');
    }
  }).catch(() => {});
};

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '-';
  
  // 尝试解析不同格式的日期时间
  try {
    // 如果是ISO格式或标准日期格式，直接使用Date对象格式化
    const date = new Date(dateTimeStr);
    if (!isNaN(date.getTime())) {
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      });
    }
    
    // 如果是类似"2023-06-06 11:32:24"的格式
    if (typeof dateTimeStr === 'string' && dateTimeStr.includes('-')) {
      return dateTimeStr;
    }
  } catch (e) {
    console.error('日期格式化错误:', e);
  }
  
  return dateTimeStr || '-';
};

// 格式化价格
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return Number(price).toFixed(2);
};

// 计算不同状态的订单数量
const completedCount = computed(() => {
  return orderList.value.filter(order => 
    order.status === 'completed' || order.status === 2
  ).length;
});

const pendingCount = computed(() => {
  return orderList.value.filter(order => 
    order.status === 'new' || order.status === 0
  ).length;
});

const cancelledCount = computed(() => {
  return orderList.value.filter(order => 
    order.status === 'cancelled' || order.status === 3
  ).length;
});

// 页面加载时获取数据
onMounted(() => {
  console.log('订单列表页面已挂载');
  if (!isInitialized.value) {
    fetchOrderList();
    isInitialized.value = true;
  }
});

// 当页面从缓存中激活时触发（切换tab时）
onActivated(() => {
  console.log('订单列表页面已激活');
  // 避免重复请求数据
  if (!isRequestLocked.value) {
    fetchOrderList();
  }
});
</script>

<style scoped>
.page-container {
  padding: 20px;
  background-color: #f5f7fa;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.order-stats {
  display: flex;
  gap: 10px;
}

.search-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.search-form {
  display: flex;
  flex-direction: column;
}

.search-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.table-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-nickname {
  font-weight: bold;
  color: #303133;
}

.user-name {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.user-id {
  color: #909399;
  font-size: 12px;
}

.price-value {
  font-weight: bold;
  color: #f56c6c;
}

.price-zero {
  color: #909399;
  font-weight: normal;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #606266;
}

.payment-method {
  color: #606266;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .order-stats {
    flex-wrap: wrap;
  }
  
  .search-buttons {
    justify-content: flex-start;
  }
}
</style> 