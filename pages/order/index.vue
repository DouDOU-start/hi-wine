<template>
  <view class="order-container">
    <view class="order-header">我的订单</view>
    <view v-if="orders.length === 0" class="order-empty">暂无订单~</view>
    <view v-else class="order-list">
      <view v-for="order in orders" :key="order.id" class="order-item">
        <view class="order-info">
          <text class="order-time">{{ order.time }}</text>
          <text class="order-status">{{ order.status }}</text>
        </view>
        <view class="order-goods">
          <view v-for="item in order.items" :key="item.id" class="order-good">
            <image :src="imgBaseUrl + item.img" class="order-img" mode="aspectFill" />
            <view class="order-good-info">
              <text class="order-good-name">{{ item.name }}</text>
              <text class="order-good-count">x{{ item.count }}</text>
            </view>
          </view>
        </view>
        <view class="order-total">合计：￥{{ order.total }}</view>
      </view>
    </view>
  </view>
</template>

<script>
import ordersData from '../../static/orders.json';
const imgBaseUrl = import.meta.env.VITE_IMG_BASE_URL || '';
export default {
  data() {
    return {
      orders: ordersData,
      imgBaseUrl
    };
  },
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
.order-empty {
  text-align: center;
  color: #b8b8b8;
  margin-top: 120rpx;
  font-size: 32rpx;
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
  color: #f7cac9;
  font-size: 28rpx;
  font-weight: bold;
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
</style> 