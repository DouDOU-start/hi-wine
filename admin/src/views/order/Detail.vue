<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">订单详情</div>
      <el-button @click="goBack">返回列表</el-button>
    </div>
    
    <el-card shadow="never" class="info-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
          <div class="order-status">
            <el-tag :type="getOrderStatusType(orderDetail.status)" size="large">
              {{ getOrderStatusText(orderDetail.status) }}
            </el-tag>
          </div>
        </div>
      </template>
      
      <el-descriptions :column="3" border>
        <el-descriptions-item label="订单号">{{ orderDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ orderDetail.username }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ orderDetail.phone }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">￥{{ orderDetail.totalAmount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ getPayMethodText(orderDetail.payMethod) }}</el-descriptions-item>
        <el-descriptions-item label="商品数量">{{ orderDetail.productCount }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ orderDetail.createTime }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ orderDetail.payTime || '未支付' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ orderDetail.finishTime || '未完成' }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="action-bar" v-if="orderDetail.status !== 2 && orderDetail.status !== 3">
        <el-button 
          v-if="orderDetail.status === 1"
          type="primary" 
          @click="handleCompleteOrder"
        >
          完成订单
        </el-button>
        <el-button 
          v-if="orderDetail.status === 0"
          type="danger" 
          @click="handleCancelOrder"
        >
          取消订单
        </el-button>
      </div>
    </el-card>
    
    <el-card shadow="never" class="info-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>商品信息</span>
        </div>
      </template>
      
      <el-table :data="orderDetail.items" border style="width: 100%">
        <el-table-column prop="productId" label="商品ID" width="100" />
        <el-table-column prop="productImage" label="商品图片" width="100">
          <template #default="scope">
            <el-image 
              v-if="scope.row.productImage" 
              :src="scope.row.productImage" 
              style="width: 60px; height: 60px"
              fit="cover"
              :preview-src-list="[scope.row.productImage]"
            />
            <span v-else>无图片</span>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="商品名称" show-overflow-tooltip />
        <el-table-column prop="price" label="单价" width="120">
          <template #default="scope">
            ￥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="subtotal" label="小计" width="120">
          <template #default="scope">
            ￥{{ (scope.row.price * scope.row.quantity).toFixed(2) }}
          </template>
        </el-table-column>
      </el-table>
      
      <div class="order-total">
        <span>订单总金额：<span class="price">￥{{ orderDetail.totalAmount?.toFixed(2) }}</span></span>
      </div>
    </el-card>
    
    <el-card shadow="never" class="info-card" v-if="orderDetail.remark" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>订单备注</span>
        </div>
      </template>
      
      <div class="remark-content">
        {{ orderDetail.remark }}
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderDetail, updateOrderStatus } from '../../api/order';

const router = useRouter();
const route = useRoute();

// 加载状态
const loading = ref(false);

// 订单详情
const orderDetail = reactive({
  id: '',
  username: '',
  phone: '',
  totalAmount: 0,
  productCount: 0,
  status: 0,
  payMethod: 0,
  createTime: '',
  payTime: '',
  finishTime: '',
  remark: '',
  items: []
});

// 获取订单详情
const fetchOrderDetail = async () => {
  loading.value = true;
  try {
    const response = await getOrderDetail(route.params.id);
    Object.assign(orderDetail, response.data);
  } catch (error) {
    console.error('获取订单详情失败:', error);
    ElMessage.error('获取订单详情失败');
    goBack();
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

// 完成订单
const handleCompleteOrder = () => {
  ElMessageBox.confirm(`确定要将订单 ${orderDetail.id} 标记为已完成吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(orderDetail.id, 2);
      ElMessage.success('订单已完成');
      orderDetail.status = 2;
      orderDetail.finishTime = new Date().toLocaleString();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 取消订单
const handleCancelOrder = () => {
  ElMessageBox.confirm(`确定要取消订单 ${orderDetail.id} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(orderDetail.id, 3);
      ElMessage.success('订单已取消');
      orderDetail.status = 3;
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 返回列表
const goBack = () => {
  router.push('/order/list');
};

// 页面加载时获取数据
onMounted(() => {
  if (route.params.id) {
    fetchOrderDetail();
  } else {
    ElMessage.error('订单ID不能为空');
    goBack();
  }
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
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-status {
  font-size: 16px;
}

.action-bar {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.order-total {
  margin-top: 20px;
  text-align: right;
  font-size: 16px;
  padding-right: 20px;
}

.price {
  color: #f56c6c;
  font-size: 20px;
  font-weight: bold;
}

.remark-content {
  padding: 10px;
  background-color: #f8f8f8;
  border-radius: 4px;
  min-height: 60px;
}
</style> 