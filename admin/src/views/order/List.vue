<template>
  <div class="page-container">
    <div class="page-title">订单管理</div>
    
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.orderId" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchForm.status" placeholder="请选择订单状态" clearable>
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="已完成" :value="2" />
            <el-option label="已取消" :value="3" />
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
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <div class="left">
          <el-button type="primary" @click="exportOrderData">导出订单</el-button>
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="orderList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="订单号" width="120" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="totalAmount" label="订单金额" width="120">
          <template #default="scope">
            ￥{{ scope.row.totalAmount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="productCount" label="商品数量" width="100" />
        <el-table-column prop="status" label="订单状态" width="100">
          <template #default="scope">
            <el-tag :type="getOrderStatusType(scope.row.status)">
              {{ getOrderStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="payMethod" label="支付方式" width="100">
          <template #default="scope">
            {{ getPayMethodText(scope.row.payMethod) }}
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="下单时间" width="180" />
        <el-table-column prop="payTime" label="支付时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              @click="handleViewDetail(scope.row)"
            >
              查看详情
            </el-button>
            <el-button 
              v-if="scope.row.status === 1"
              type="success" 
              link 
              @click="handleCompleteOrder(scope.row)"
            >
              完成订单
            </el-button>
            <el-button 
              v-if="scope.row.status === 0"
              type="danger" 
              link 
              @click="handleCancelOrder(scope.row)"
            >
              取消订单
            </el-button>
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
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderList, updateOrderStatus, exportOrders } from '../../api/order';

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
  orderId: '',
  username: '',
  status: '',
  dateRange: []
});

// 获取订单列表
const fetchOrderList = async () => {
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
    orderList.value = response.data.list || [];
    total.value = response.data.total || 0;
  } catch (error) {
    console.error('获取订单列表失败:', error);
    ElMessage.error('获取订单列表失败');
  } finally {
    loading.value = false;
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  switch (status) {
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
    default: return '未支付';
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
  router.push(`/order/detail/${row.id}`);
};

// 完成订单
const handleCompleteOrder = (row) => {
  ElMessageBox.confirm(`确定要将订单 ${row.id} 标记为已完成吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(row.id, 2);
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
  ElMessageBox.confirm(`确定要取消订单 ${row.id} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(row.id, 3);
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

// 页面加载时获取数据
onMounted(() => {
  fetchOrderList();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.page-title {
  margin-bottom: 20px;
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 