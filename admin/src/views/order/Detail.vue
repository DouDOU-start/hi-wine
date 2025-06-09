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
            <el-tag :type="getOrderStatusType(orderDetail.orderStatus)" size="large">
              {{ getOrderStatusText(orderDetail.orderStatus) }}
            </el-tag>
          </div>
        </div>
      </template>
      
      <el-descriptions :column="3" border>
        <el-descriptions-item label="订单号">{{ orderDetail.orderSn }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ orderDetail.username }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ orderDetail.phone }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">￥{{ formatPrice(orderDetail.totalAmount) }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ getPayMethodText(orderDetail.paymentMethod) }}</el-descriptions-item>
        <el-descriptions-item label="商品数量">{{ orderDetail.itemCount }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ orderDetail.createdAt ? formatDate(orderDetail.createdAt) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ orderDetail.paidAt ? formatDate(orderDetail.paidAt) : '未支付' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ orderDetail.completedAt ? formatDate(orderDetail.completedAt) : '-' }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="action-bar" v-if="orderDetail.orderStatus !== 'completed' && orderDetail.orderStatus !== 'cancelled'">
        <el-button 
          v-if="orderDetail.orderStatus === 'processing'"
          type="primary" 
          @click="handleCompleteOrder"
        >
          完成订单
        </el-button>
        <el-button 
          v-if="orderDetail.orderStatus === 'new'"
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
              style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;"
              fit="cover"
              :preview-src-list="[scope.row.productImage]"
              preview-teleported
            />
            <div v-else class="no-image">无图片</div>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="商品名称" show-overflow-tooltip>
          <template #default="scope">
            <div class="product-name">
              {{ scope.row.productName }}
              <el-tag v-if="scope.row.isPackageItem" size="small" type="success">套餐</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="单价" width="120">
          <template #default="scope">
            <span class="price-value">￥{{ formatPrice(scope.row.price) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" align="center" />
        <el-table-column label="小计" width="120">
          <template #default="scope">
            <span class="price-value">￥{{ formatPrice(scope.row.price * scope.row.quantity) }}</span>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="order-total">
        <span>订单总金额：<span class="price">￥{{ formatPrice(orderDetail.totalAmount) }}</span></span>
      </div>
    </el-card>
    
    <el-card shadow="never" class="info-card" v-if="orderDetail.message" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>订单备注</span>
        </div>
      </template>
      
      <div class="remark-content">
        {{ orderDetail.message }}
      </div>
    </el-card>

    <pre v-if="false">{{ orderDetail }}</pre>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderDetail, updateOrderStatus } from '../../api/order';
import { formatDate } from '../../utils/format';

const router = useRouter();
const route = useRoute();

// 加载状态
const loading = ref(false);

// 数据加载标志，防止重复请求
const isLoading = ref(false);

// 订单详情
const orderDetail = reactive({
  id: '',
  orderSn: '',
  username: '',
  phone: '',
  totalAmount: 0,
  itemCount: 0,
  orderStatus: '',
  paymentMethod: '',
  paymentStatus: '',
  createdAt: '',
  paidAt: '',
  completedAt: '',
  message: '',
  items: []
});

// 格式化价格
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return Number(price).toFixed(2);
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  switch (status) {
    case 'new': return '待支付';
    case 'processing': return '已支付';
    case 'completed': return '已完成';
    case 'cancelled': return '已取消';
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
    default: return 'info';
  }
};

// 获取支付方式文本
const getPayMethodText = (method) => {
  switch (method) {
    case 'wechat': return '微信支付';
    case 'alipay': return '支付宝';
    case 'balance': return '余额支付';
    case 'paid': return '已支付'; // 兼容payment_status作为支付方式
    default: return '-';
  }
};

// 获取订单详情
const fetchOrderDetail = async () => {
  // 如果已经在加载中，则不重复请求
  if (isLoading.value) {
    console.log('数据正在加载中，跳过重复请求');
    return;
  }
  
  isLoading.value = true;
  loading.value = true;
  
  try {
    const orderId = route.params.id; // 从路由参数获取订单ID
    console.log('获取订单详情，订单ID:', orderId);
    
    // 确保ID是数字
    const numericId = parseInt(orderId);
    if (isNaN(numericId)) {
      throw new Error(`无效的订单ID: ${orderId}`);
    }
    
    const response = await getOrderDetail(numericId);
    console.log('订单详情响应:', response);
    
    if (response && response.data) {
      // 处理订单基本信息
      // 由于请求拦截器会将下划线字段转为驼峰，我们需要处理可能的两种情况
      const orderData = response.data;
      console.log('订单数据:', orderData);
      
      // 映射后端返回的字段到前端展示字段
      orderDetail.id = orderData.id || '';
      orderDetail.orderSn = orderData.orderSn || orderData.order_sn || '';
      orderDetail.username = orderData.userNickname || orderData.userName || orderData.user_nickname || orderData.user_name || '';
      orderDetail.phone = orderData.userPhone || orderData.user_phone || '';
      
      // 处理金额字段，确保是数字类型
      const amount = orderData.totalAmount || orderData.total_amount || 0;
      orderDetail.totalAmount = parseFloat(amount);
      
      // 处理商品数量
      orderDetail.itemCount = parseInt(orderData.itemCount || orderData.item_count || 0);
      
      // 处理订单状态
      orderDetail.orderStatus = orderData.orderStatus || orderData.order_status || '';
      
      // 处理支付方式
      orderDetail.paymentMethod = orderData.paymentMethod || orderData.payment_method || '';
      
      // 处理支付状态
      orderDetail.paymentStatus = orderData.paymentStatus || orderData.payment_status || '';
      
      // 处理时间字段
      orderDetail.createdAt = orderData.createdAt || orderData.created_at || '';
      orderDetail.paidAt = orderData.paidAt || orderData.paid_at || '';
      orderDetail.completedAt = orderData.completedAt || orderData.completed_at || '';
      
      // 处理备注信息
      orderDetail.message = orderData.message || '';
      
      // 处理订单商品项
      if (orderData.items && Array.isArray(orderData.items)) {
        orderDetail.items = orderData.items.map(item => {
          const productId = item.productId || item.product_id || '';
          const productName = item.name || '';
          const productImage = item.imageUrl || item.image_url || '';
          const price = parseFloat(item.itemPrice || item.item_price || 0);
          const quantity = parseInt(item.quantity || 0);
          const isPackageItem = item.isPackageItem || item.is_package_item || false;
          
          return {
            productId,
            productName,
            productImage,
            price,
            quantity,
            subtotal: price * quantity,
            isPackageItem
          };
        });
      } else {
        orderDetail.items = [];
      }

      // 打印处理后的数据，方便调试
      console.log('处理后的订单数据:', orderDetail);
    } else {
      throw new Error('订单数据格式异常');
    }
  } catch (error) {
    console.error('获取订单详情失败:', error);
    ElMessage.error(`获取订单详情失败: ${error.message}`);
    goBack();
  } finally {
    loading.value = false;
    isLoading.value = false;
  }
};

// 完成订单
const handleCompleteOrder = () => {
  ElMessageBox.confirm(`确定要将订单 ${orderDetail.orderSn} 标记为已完成吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(orderDetail.id, 'completed');
      ElMessage.success('订单已完成');
      orderDetail.orderStatus = 'completed';
      orderDetail.completedAt = new Date().toISOString();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 取消订单
const handleCancelOrder = () => {
  ElMessageBox.confirm(`确定要取消订单 ${orderDetail.orderSn} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await updateOrderStatus(orderDetail.id, 'cancelled');
      ElMessage.success('订单已取消');
      orderDetail.orderStatus = 'cancelled';
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

// 监听路由参数变化，重新获取数据
watch(
  () => route.params.id,
  (newId, oldId) => {
    if (newId && newId !== oldId) {
      console.log('订单ID变化，重新获取数据:', newId);
      fetchOrderDetail();
    }
  },
  { immediate: false }
);

// 页面加载时获取数据
onMounted(() => {
  console.log('页面挂载，获取数据');
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

.no-image {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  color: #909399;
  border-radius: 4px;
}

.product-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.price-value {
  font-weight: bold;
  color: #f56c6c;
}
</style> 