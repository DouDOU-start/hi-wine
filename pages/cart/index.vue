<template>
  <view class="cart-container">
    <view class="cart-header">购物车</view>
    <view v-if="cart.length === 0" class="cart-empty">购物车空空如也~</view>
    <view v-else class="cart-list">
      <view v-for="item in cart" :key="item.id" class="cart-item">
        <image :src="item.image ? item.image : IMG_BASE_URL + '/wine.png'" class="cart-img" mode="aspectFill" />
        <view class="cart-info">
          <text class="cart-name">{{ item.name }}</text>
          <text class="cart-price">￥{{ item.price }}</text>
          <view class="cart-ctrl">
            <button class="cart-btn" @tap="changeCount(item, -1)" :disabled="item.count<=1">-</button>
            <text class="cart-count">{{ item.count }}</text>
            <button class="cart-btn" @tap="changeCount(item, 1)">+</button>
            <button class="cart-del" @tap="removeItem(item)">删除</button>
          </view>
        </view>
      </view>
      <view class="cart-footer">
        <text>合计：</text>
        <text class="cart-total">￥{{ totalPrice }}</text>
        <button class="cart-order" @tap="showOrderConfirm">下单</button>
      </view>
    </view>
    
    <!-- 下单确认弹窗 -->
    <uni-popup ref="popup" type="dialog">
      <uni-popup-dialog
        title="确认下单"
        content="确定要提交订单吗？"
        :before-close="true"
        @confirm="submitOrder"
        @close="closePopup"
      ></uni-popup-dialog>
    </uni-popup>
    
    <!-- 加载中 -->
    <uni-popup ref="loading" type="dialog" :mask-click="false">
      <uni-popup-dialog
        title="处理中"
        content="订单提交中，请稍候..."
        :show-cancel="false"
        :before-close="true"
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
      cart: [],
      IMG_BASE_URL,
      submitting: false
    };
  },
  computed: {
    totalPrice() {
      return this.cart.reduce((sum, item) => sum + item.price * item.count, 0);
    },
  },
  onShow() {
    // 每次显示页面时，从本地存储读取购物车数据
    this.loadCartData();
  },
  methods: {
    // 加载购物车数据
    loadCartData() {
      const cartData = uni.getStorageSync('cart') || [];
      this.cart = cartData;
    },
    
    // 修改商品数量
    changeCount(item, delta) {
      const idx = this.cart.findIndex(i => i.id === item.id);
      if (idx !== -1) {
        this.cart[idx].count += delta;
        if (this.cart[idx].count < 1) this.cart[idx].count = 1;
        this.updateCartStorage();
      }
    },
    
    // 移除商品
    removeItem(item) {
      this.cart = this.cart.filter(i => i.id !== item.id);
      this.updateCartStorage();
    },
    
    // 更新本地存储
    updateCartStorage() {
      uni.setStorageSync('cart', this.cart);
    },
    
    // 显示下单确认弹窗
    showOrderConfirm() {
      if (this.cart.length === 0) {
        uni.showToast({
          title: '购物车为空',
          icon: 'none'
        });
        return;
      }
      
      // 检查是否登录
      const token = uni.getStorageSync('token');
      if (!token) {
        uni.showToast({
          title: '请先登录',
          icon: 'none'
        });
        setTimeout(() => {
          uni.switchTab({
            url: '/pages/profile/index'
          });
        }, 1500);
        return;
      }
      
      this.$refs.popup.open();
    },
    
    // 关闭弹窗
    closePopup() {
      this.$refs.popup.close();
    },
    
    // 提交订单
    async submitOrder() {
      if (this.submitting) return;
      
      this.submitting = true;
      this.$refs.loading.open();
      
      try {
        // 将购物车商品转换为订单项
        const items = this.cart.map(item => ({
          productId: item.id,
          quantity: item.count
        }));
        
        // 调用创建订单API
        const res = await api.createOrder(0, items);
        
        if (res && res.id) {
          // 清空购物车
          this.cart = [];
          this.updateCartStorage();
          
          // 关闭加载弹窗
          this.$refs.loading.close();
          
          // 显示成功提示
          uni.showToast({
            title: '下单成功',
            icon: 'success'
          });
          
          // 跳转到订单页面
          setTimeout(() => {
            uni.switchTab({
              url: '/pages/order/index'
            });
          }, 1500);
        }
      } catch (err) {
        console.error('下单失败', err);
        uni.showToast({
          title: '下单失败，请重试',
          icon: 'none'
        });
      } finally {
        this.submitting = false;
        this.$refs.loading.close();
      }
    }
  }
};
</script>

<style scoped>
.cart-container {
  min-height: 100vh;
  background: #f8f6f4;
  padding-bottom: 40rpx;
}
.cart-header {
  padding: 60rpx 0 20rpx 0;
  text-align: center;
  font-size: 48rpx;
  font-weight: bold;
  color: #222;
  letter-spacing: 2rpx;
}
.cart-empty {
  text-align: center;
  color: #b8b8b8;
  margin-top: 120rpx;
  font-size: 32rpx;
}
.cart-list {
  padding: 0 24rpx;
}
.cart-item {
  display: flex;
  background: #fff;
  border-radius: 32rpx;
  box-shadow: 0 4rpx 24rpx #eaeaea;
  margin-bottom: 24rpx;
  overflow: hidden;
}
.cart-img {
  width: 180rpx;
  height: 140rpx;
  object-fit: cover;
  border-top-left-radius: 32rpx;
  border-bottom-left-radius: 32rpx;
}
.cart-info {
  flex: 1;
  padding: 20rpx 24rpx 0 24rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.cart-name {
  font-size: 32rpx;
  color: #222;
  font-weight: 500;
}
.cart-price {
  font-size: 28rpx;
  color: #f7cac9;
  font-weight: bold;
  margin-bottom: 8rpx;
}
.cart-ctrl {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
}
.cart-btn {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #f7cac9;
  color: #fff;
  border: none;
  font-size: 32rpx;
  font-weight: bold;
  line-height: 48rpx;
  text-align: center;
}
.cart-del {
  background: #eee;
  color: #b8b8b8;
  border: none;
  border-radius: 24rpx;
  font-size: 24rpx;
  padding: 0 16rpx;
  margin-left: 10rpx;
}
.cart-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 24rpx;
  margin-top: 30rpx;
}
.cart-total {
  font-size: 36rpx;
  color: #f7cac9;
  font-weight: bold;
}
.cart-order {
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  color: #fff;
  border: none;
  border-radius: 24rpx;
  font-size: 28rpx;
  font-weight: 500;
  padding: 8rpx 36rpx;
  box-shadow: 0 2rpx 8rpx #e0e0e0;
}
</style> 