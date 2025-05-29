<template>
  <view class="order-container">
    <view class="order-header">我的订单</view>
    
    <!-- 订单状态选项卡 -->
    <view class="order-tabs">
      <view 
        v-for="(tab, index) in statusTabs" 
        :key="index" 
        :class="['order-tab', activeTab === index ? 'active' : '']"
        @tap="changeTab(index)"
      >
        {{ tab.name }}
      </view>
    </view>
    
    <view v-if="loading" class="order-loading">加载中...</view>
    <view v-else-if="orders.length === 0" class="order-empty">暂无订单~</view>
    <view v-else class="order-list">
      <view v-for="order in orders" :key="order.id" class="order-item">
        <view class="order-info">
          <text class="order-time">{{ formatTime(order.createTime) }}</text>
          <text :class="['order-status', 'status-' + order.status]">{{ getStatusText(order.status) }}</text>
        </view>
        <view class="order-goods">
          <view v-for="item in orderItems[order.id]" :key="item.id" class="order-good">
            <image :src="item.productImage || IMG_BASE_URL + '/wine.png'" class="order-img" mode="aspectFill" />
            <view class="order-good-info">
              <text class="order-good-name">{{ item.productName }}</text>
              <text class="order-good-count">x{{ item.quantity }}</text>
            </view>
          </view>
        </view>
        <view class="order-total">合计：￥{{ order.totalAmount }}</view>
        
        <!-- 订单操作按钮 -->
        <view class="order-actions" v-if="order.status === 0">
          <button class="order-btn pay-btn" @tap="payOrder(order)">立即支付</button>
          <button class="order-btn cancel-btn" @tap="cancelOrder(order)">取消订单</button>
        </view>
      </view>
    </view>
    
    <!-- 上拉加载更多提示 -->
    <view v-if="orders.length > 0 && !hasMore" class="order-no-more">没有更多订单了</view>
    
    <!-- 确认弹窗 -->
    <uni-popup ref="confirmPopup" type="dialog">
      <uni-popup-dialog
        :title="popupConfig.title"
        :content="popupConfig.content"
        :before-close="true"
        @confirm="popupConfig.onConfirm"
        @close="closePopup"
      ></uni-popup-dialog>
    </uni-popup>
  </view>
</template>

<script>
import { IMG_BASE_URL } from '@/config.js';
import api from '@/utils/api.js';

export default {
  data() {
    return {
      orders: [],
      orderItems: {}, // 订单商品映射表，key为订单ID
      loading: false,
      page: 1,
      size: 10,
      hasMore: true,
      IMG_BASE_URL,
      activeTab: 0,
      statusTabs: [
        { name: '全部', value: -1 },
        { name: '待支付', value: 0 },
        { name: '已支付', value: 1 },
        { name: '已完成', value: 2 },
        { name: '已取消', value: 3 }
      ],
      popupConfig: {
        title: '',
        content: '',
        onConfirm: null
      }
    };
  },
  onLoad() {
    this.loadOrders();
  },
  onShow() {
    // 每次显示页面时刷新数据
    this.refreshOrders();
  },
  methods: {
    // 切换状态选项卡
    changeTab(index) {
      if (this.activeTab === index) return;
      this.activeTab = index;
      this.refreshOrders();
    },
    
    // 刷新订单列表
    refreshOrders() {
      this.orders = [];
      this.orderItems = {};
      this.page = 1;
      this.hasMore = true;
      this.loadOrders();
    },
    
    // 加载订单列表
    async loadOrders() {
      if (this.loading || !this.hasMore) return;
      
      this.loading = true;
      try {
        // 获取当前选中的状态值
        const status = this.statusTabs[this.activeTab].value;
        
        const res = await api.getOrderList(status, this.page, this.size);
        console.log('订单列表响应:', JSON.stringify(res));
        
        if (res && res.code === 0 && res.data && res.data.list) {
          if (this.page === 1) {
            this.orders = res.data.list;
          } else {
            this.orders = [...this.orders, ...res.data.list];
          }
          
          // 判断是否还有更多数据
          this.hasMore = res.data.list.length === this.size;
          this.page++;
          
          // 加载每个订单的商品详情
          for (const order of res.data.list) {
            this.loadOrderItems(order.id);
          }
        }
      } catch (err) {
        console.error('加载订单失败', err);
        uni.showToast({
          title: '加载订单失败',
          icon: 'none'
        });
      } finally {
        this.loading = false;
      }
    },
    
    // 加载订单商品
    async loadOrderItems(orderId) {
      try {
        const res = await api.getOrderDetail(orderId);
        console.log('订单详情响应:', JSON.stringify(res));
        if (res && res.code === 0 && res.data && res.data.orderItems) {
          // 使用Vue.set确保响应式更新
          this.$set(this.orderItems, orderId, res.data.orderItems);
        }
      } catch (err) {
        console.error(`加载订单${orderId}商品失败`, err);
      }
    },
    
    // 支付订单
    payOrder(order) {
      this.popupConfig = {
        title: '支付订单',
        content: `确认支付订单金额 ￥${order.totalAmount}？`,
        onConfirm: async () => {
          try {
            // 调用更新订单状态API，将状态改为已支付(1)
            const res = await api.updateOrderStatus(order.id, 1);
            if (res && res.code === 0) {
              uni.showToast({
                title: '支付成功',
                icon: 'success'
              });
              // 刷新订单列表
              this.refreshOrders();
            } else {
              throw new Error('支付失败');
            }
          } catch (err) {
            console.error('支付失败', err);
            uni.showToast({
              title: '支付失败',
              icon: 'none'
            });
          }
        }
      };
      this.$refs.confirmPopup.open();
    },
    
    // 取消订单
    cancelOrder(order) {
      this.popupConfig = {
        title: '取消订单',
        content: '确认取消该订单？',
        onConfirm: async () => {
          try {
            // 调用更新订单状态API，将状态改为已取消(3)
            const res = await api.updateOrderStatus(order.id, 3);
            if (res && res.code === 0) {
              uni.showToast({
                title: '订单已取消',
                icon: 'success'
              });
              // 刷新订单列表
              this.refreshOrders();
            } else {
              throw new Error('取消订单失败');
            }
          } catch (err) {
            console.error('取消订单失败', err);
            uni.showToast({
              title: '取消订单失败',
              icon: 'none'
            });
          }
        }
      };
      this.$refs.confirmPopup.open();
    },
    
    // 关闭弹窗
    closePopup() {
      this.$refs.confirmPopup.close();
    },
    
    // 格式化时间
    formatTime(time) {
      if (!time) return '';
      // 简单处理，实际项目中可以使用日期格式化库
      return time.replace('T', ' ').split('.')[0];
    },
    
    // 获取状态文本
    getStatusText(status) {
      switch (status) {
        case 0: return '待支付';
        case 1: return '已支付';
        case 2: return '已完成';
        case 3: return '已取消';
        default: return '未知状态';
      }
    }
  },
  // 下拉刷新
  onPullDownRefresh() {
    this.refreshOrders();
    uni.stopPullDownRefresh();
  },
  // 上拉加载更多
  onReachBottom() {
    this.loadOrders();
  }
};
</script>

<style scoped>
.order-container {
  min-height: 100vh;
  background: #f8f6f4;
  padding-bottom: 40rpx;
}
.order-header {
  padding: 60rpx 0 20rpx 0;
  text-align: center;
  font-size: 48rpx;
  font-weight: bold;
  color: #222;
  letter-spacing: 2rpx;
}
.order-tabs {
  display: flex;
  justify-content: space-around;
  padding: 0 20rpx;
  margin-bottom: 30rpx;
}
.order-tab {
  font-size: 28rpx;
  color: #666;
  padding: 10rpx 20rpx;
  position: relative;
}
.order-tab.active {
  color: #f7cac9;
  font-weight: bold;
}
.order-tab.active::after {
  content: '';
  position: absolute;
  bottom: -6rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 40rpx;
  height: 6rpx;
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  border-radius: 3rpx;
}
.order-loading, .order-empty, .order-no-more {
  text-align: center;
  color: #b8b8b8;
  margin: 40rpx 0;
  font-size: 28rpx;
}
.order-list {
  padding: 0 24rpx;
}
.order-item {
  background: #fff;
  border-radius: 32rpx;
  box-shadow: 0 4rpx 24rpx #eaeaea;
  margin-bottom: 32rpx;
  padding: 24rpx;
}
.order-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16rpx;
}
.order-time {
  color: #b8b8b8;
  font-size: 28rpx;
}
.order-status {
  font-size: 28rpx;
  font-weight: bold;
}
.status-0 {
  color: #f7cac9;
}
.status-1 {
  color: #92a8d1;
}
.status-2 {
  color: #67c23a;
}
.status-3 {
  color: #909399;
}
.order-goods {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 12rpx;
}
.order-good {
  display: flex;
  align-items: center;
  background: #f8f6f4;
  border-radius: 16rpx;
  padding: 8rpx 16rpx;
  margin-right: 12rpx;
}
.order-img {
  width: 60rpx;
  height: 60rpx;
  border-radius: 12rpx;
  object-fit: cover;
  margin-right: 12rpx;
}
.order-good-info {
  display: flex;
  flex-direction: column;
}
.order-good-name {
  font-size: 28rpx;
  color: #222;
}
.order-good-count {
  font-size: 24rpx;
  color: #b8b8b8;
}
.order-total {
  text-align: right;
  color: #f7cac9;
  font-size: 32rpx;
  font-weight: bold;
  margin-top: 8rpx;
}
.order-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 20rpx;
  gap: 20rpx;
}
.order-btn {
  font-size: 24rpx;
  border-radius: 24rpx;
  padding: 6rpx 24rpx;
}
.pay-btn {
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  color: #fff;
  border: none;
}
.cancel-btn {
  background: #f5f5f5;
  color: #909399;
  border: none;
}
</style> 