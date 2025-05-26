<template>
  <view class="profile-container">
    <view class="profile-header">
      <image class="profile-avatar" src="https://randomuser.me/api/portraits/men/32.jpg" />
      <text class="profile-nick">Hi，酒友</text>
    </view>
    <view class="profile-list">
      <view class="profile-item" @tap="goOrder">
        <text class="profile-icon">📦</text>
        <text class="profile-label">我的订单</text>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  methods: {
    goOrder() {
      uni.navigateTo({ url: '/pages/order/index' });
    },
  },
};
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: #f8f6f4;
  padding-bottom: 40rpx;
}
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0 40rpx 0;
}
.profile-avatar {
  width: 140rpx;
  height: 140rpx;
  border-radius: 50%;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 24rpx #eaeaea;
}
.profile-nick {
  font-size: 36rpx;
  color: #222;
  font-weight: bold;
}
.profile-list {
  margin: 0 40rpx;
}
.profile-item {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 32rpx;
  box-shadow: 0 4rpx 24rpx #eaeaea;
  padding: 32rpx 24rpx;
  margin-bottom: 32rpx;
  font-size: 32rpx;
  color: #222;
  font-weight: 500;
  transition: box-shadow 0.2s;
}
.profile-icon {
  font-size: 40rpx;
  margin-right: 24rpx;
}
</style> 