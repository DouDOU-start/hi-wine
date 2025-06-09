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
        <el-descriptions-item label="订单号">{{ orderDetail.orderSn || orderDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ orderDetail.username }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ orderDetail.phone }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">￥{{ orderDetail.totalAmount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ getPayMethodText(orderDetail.payMethod) }}</el-descriptions-item>
        <el-descriptions-item label="商品数量">{{ orderDetail.productCount }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ orderDetail.createTime }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ orderDetail.payTime || '未支付' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ orderDetail.finishTime || '-' }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="action-bar" v-if="orderDetail.status !== 'completed' && orderDetail.status !== 'cancelled'">
        <el-button 
          v-if="orderDetail.status === 'processing'"
          type="primary" 
          @click="handleCompleteOrder"
        >
          完成订单
        </el-button>
        <el-button 
          v-if="orderDetail.status === 'new'"
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
import { ref, reactive, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderDetail, updateOrderStatus } from '../../api/order';

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
  productCount: 0,
  status: '',
  payMethod: 0,
  createTime: '',
  payTime: '',
  finishTime: '',
  remark: '',
  items: []
});

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
      const orderData = response.data;
      console.log('订单数据:', orderData);
      
      // 映射后端返回的字段到前端展示字段
      orderDetail.id = orderData.id || '';
      orderDetail.orderSn = orderData.order_sn || '';
      orderDetail.username = orderData.user_nickname || orderData.user_name || orderData.username || '';
      orderDetail.phone = orderData.user_phone || orderData.phone || '';
      orderDetail.totalAmount = orderData.total_amount !== undefined 
        ? (typeof orderData.total_amount === 'string' ? parseFloat(orderData.total_amount) : orderData.total_amount) 
        : orderData.totalAmount || 0;
      orderDetail.productCount = orderData.item_count || orderData.productCount || 0;
      orderDetail.status = orderData.order_status || orderData.status || '';
      orderDetail.payMethod = orderData.payment_method || orderData.payMethod || 0;
      orderDetail.createTime = orderData.created_at || orderData.createTime || '';
      orderDetail.payTime = orderData.paid_at || orderData.payTime || '';
      orderDetail.finishTime = orderData.completed_at || orderData.finishTime || '';
      orderDetail.remark = orderData.remark || '';
      
      // 处理订单商品项 - 根据截图调整
      if (orderData.items && Array.isArray(orderData.items)) {
        orderDetail.items = orderData.items.map(item => {
          // 根据截图中的数据结构处理
          return {
            productId: item.product_id || item.id || '',
            productName: item.name || '',
            productImage: item.image_url || item.image || '',
            price: typeof item.item_price === 'string' ? parseFloat(item.item_price) : (item.item_price || 0),
            quantity: typeof item.quantity === 'string' ? parseInt(item.quantity) : (item.quantity || 0),
            subtotal: (item.item_price * item.quantity) || 0,
            isPackageItem: item.is_package_item || false
          };
        });
        
        // 调试输出商品图片信息
        console.log('处理后的商品项:', orderDetail.items);
      } else if (orderData.items) {
        // 如果items不是数组，可能是对象结构
        console.log('商品项不是数组:', orderData.items);
        try {
          // 尝试处理可能的其他结构
          const items = [];
          if (orderData.items[0]) {
            // 如果是对象的数字索引
            for (const key in orderData.items) {
              if (!isNaN(parseInt(key))) {
                const item = orderData.items[key];
                items.push({
                  productId: item.product_id || item.id || '',
                  productName: item.name || '',
                  productImage: item.image_url || item.image || '',
                  price: typeof item.item_price === 'string' ? parseFloat(item.item_price) : (item.item_price || 0),
                  quantity: typeof item.quantity === 'string' ? parseInt(item.quantity) : (item.quantity || 0),
                  subtotal: (item.item_price * item.quantity) || 0,
                  isPackageItem: item.is_package_item || false
                });
              }
            }
            orderDetail.items = items;
            
            // 调试输出商品图片信息
            console.log('处理后的商品项(对象结构):', orderDetail.items);
          }
        } catch (e) {
          console.error('处理商品项出错:', e);
          orderDetail.items = [];
        }
      } else {
        orderDetail.items = [];
      }
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

// 完成订单
const handleCompleteOrder = () => {
  ElMessageBox.confirm(`确定要将订单 ${orderDetail.orderSn || orderDetail.id} 标记为已完成吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 使用数字ID
      await updateOrderStatus(orderDetail.id, 'completed');
      ElMessage.success('订单已完成');
      orderDetail.status = 'completed';
      orderDetail.finishTime = new Date().toLocaleString();
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  }).catch(() => {});
};

// 取消订单
const handleCancelOrder = () => {
  ElMessageBox.confirm(`确定要取消订单 ${orderDetail.orderSn || orderDetail.order_sn} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 使用数字ID
      await updateOrderStatus(orderDetail.id, 'cancelled_cancelled');
      ElMessage.success('订单已取消');
      orderDetail.status = 'cancelled';
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
  { immediate: false } // 设置为false，避免与onMounted重复触发
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

// 格式化价格
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return Number(price).toFixed(2);
};
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