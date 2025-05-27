<template>
  <view class="profile-container">
    <view class="profile-header">
      <image class="profile-avatar" :src="userInfo.avatar || 'https://randomuser.me/api/portraits/men/32.jpg'" />
      <text class="profile-nick">Hi，{{ userInfo.nickname || '酒友' }}</text>
      <button
        v-if="showLoginBtn && !showAuthModal"
        @tap="showAuthModal = true"
        type="primary"
        class="main-login-btn ins-btn"
      >微信一键登录</button>
    </view>
    <view class="profile-list">
      <view class="profile-item" @tap="goOrder">
        <text class="profile-icon">📦</text>
        <text class="profile-label">我的订单</text>
      </view>
    </view>
    <!-- ins风格授权引导弹窗 -->
    <view v-if="showAuthModal" class="auth-modal ins-modal">
      <view class="auth-modal-content ins-modal-content">
        <text class="auth-modal-title ins-modal-title">🌸 微信授权登录</text>
        <text class="auth-modal-desc ins-modal-desc">
          为了更好地为你服务，我们需要获取你的微信头像和昵称。
          <br />
          <text style="color:#f7cac9;font-weight:bold;">请放心，信息仅用于完善你的个人资料。</text>
        </text>
        <button type="primary" class="main-login-btn ins-btn" @tap="doWechatLogin">一键授权，开启美好体验</button>
        <button class="ins-cancel-btn" @tap="closeAuthModal">暂不授权</button>
      </view>
    </view>
    <!-- ins风格去设置重新授权 -->
    <button v-if="showSettingBtn" @tap="openSetting" type="warn" class="ins-btn ins-warn-btn">去设置重新授权</button>
    <view v-if="showSettingBtn" class="ins-warn-text">请在设置中授权后，再点击"微信一键登录"</view>
  </view>
</template>

<script>
import { BASE_URL } from '@/config.js'
export default {
  data() {
    return {
      userInfo: uni.getStorageSync('userInfo') || {},
      showSettingBtn: false,
      showLoginBtn: !uni.getStorageSync('token'),
      showAuthModal: false,
    };
  },
  onLoad() {
    this.initUserStatus();
  },
  onShow() {
    // 优先用本地userInfo渲染
    const userInfo = uni.getStorageSync('userInfo') || {};
    this.userInfo = userInfo;
  },
  methods: {
    initUserStatus() {
      const token = uni.getStorageSync('token');
      if (!token) {
        this.showLoginBtn = true;
        this.showAuthModal = false;
        this.showSettingBtn = false;
      } else {
        this.showLoginBtn = false;
        this.showAuthModal = false;
        this.showSettingBtn = false;
      }
    },
    goOrder() {
      uni.navigateTo({ url: '/pages/order/index' });
    },
    // 检查用户信息授权状态（仅首次登录或token失效时调用）
    checkUserInfoAuth() {
      const token = uni.getStorageSync('token');
      if (token) {
        // 已有token，无需重复授权
        this.showLoginBtn = false;
        this.showAuthModal = false;
        this.showSettingBtn = false;
        return;
      }
      wx.getSetting({
        success: (res) => {
          if (res.authSetting && res.authSetting['scope.userInfo']) {
            // 已授权，自动获取微信用户信息
            wx.getUserInfo({
              success: (userRes) => {
                const userInfo = {
                  nickname: userRes.userInfo.nickName,
                  avatar: userRes.userInfo.avatarUrl
                };
                uni.setStorageSync('userInfo', userInfo);
                this.userInfo = userInfo;
                this.showLoginBtn = false;
                this.showAuthModal = false;
                this.showSettingBtn = false;
              },
              fail: () => {
                this.showLoginBtn = true;
                this.showAuthModal = false;
                this.showSettingBtn = false;
              }
            })
          } else {
            // 未授权，需弹窗引导
            this.showLoginBtn = true;
            this.showAuthModal = false;
            this.showSettingBtn = false;
          }
        },
        fail: () => {
          this.showLoginBtn = true;
          this.showAuthModal = false;
          this.showSettingBtn = false;
        }
      })
    },
    // 登录弹窗按钮事件（仅无token或token失效时才调用）
    doWechatLogin() {
      wx.login({
        success: (res) => {
          if (res.code) {
            wx.getUserProfile({
              desc: '用于完善会员资料',
              success: (userRes) => {
                this.showSettingBtn = false;
                this.showAuthModal = false;
                const wxUserInfo = {
                  nickname: userRes.userInfo.nickName,
                  avatar: userRes.userInfo.avatarUrl
                };
                console.log('准备请求后端', BASE_URL + '/api/wechat/login', wxUserInfo, res.code);
                if (!BASE_URL) {
                  uni.showToast({ title: '后端地址未配置', icon: 'none' });
                  return;
                }
                uni.request({
                  url: BASE_URL + '/api/wechat/login',
                  method: 'POST',
                  data: {
                    code: res.code,
                    nickname: wxUserInfo.nickname,
                    avatar: wxUserInfo.avatar
                  },
                  success: (resp) => {
                    console.log('后端响应', resp);
                    if (resp.data.code === 0) {
                      uni.setStorageSync('token', resp.data.data.token)
                      const userInfo = {
                        ...resp.data.data.userInfo,
                        nickname: wxUserInfo.nickname,
                        avatar: wxUserInfo.avatar
                      }
                      uni.setStorageSync('userInfo', userInfo)
                      this.userInfo = userInfo
                      this.showLoginBtn = false;
                      this.$nextTick(() => {
                        this.onShow && this.onShow();
                      });
                    } else {
                      uni.showToast({ title: resp.data.message || '登录失败', icon: 'none' })
                    }
                  },
                  fail: (err) => {
                    console.log('请求后端失败', err);
                    uni.showToast({ title: '无法连接后端服务', icon: 'none' })
                  }
                })
              },
              fail: () => {
                this.showSettingBtn = true;
                this.showAuthModal = false;
                uni.showModal({
                  title: '授权失败',
                  content: '如需正常使用，请在小程序设置中授权"用户信息"',
                  confirmText: '去设置',
                  success: (res) => {
                    if (res.confirm) {
                      this.openSetting();
                    }
                  }
                });
              }
            })
          } else {
            uni.showToast({ title: '获取code失败', icon: 'none' })
          }
        },
        fail: (err) => {
          console.log('wx.login 失败', err);
          uni.showToast({ title: '微信登录失败', icon: 'none' })
        }
      })
    },
    openSetting() {
      uni.openSetting({
        success: (res) => {
          this.showSettingBtn = false;
          uni.showToast({ title: '请授权后再次点击微信一键登录', icon: 'none' });
          this.checkUserInfoAuth();
        }
      })
    },
    closeAuthModal() {
      this.showAuthModal = false;
    },
  },
  // 全局拦截401，自动清除token并跳转登录
  onLoad() {
    this.initUserStatus();
    uni.addInterceptor && uni.addInterceptor('request', {
      complete: (res) => {
        if (res.statusCode === 401) {
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          this.showLoginBtn = true;
          this.showAuthModal = false;
          this.showSettingBtn = false;
          uni.showToast({ title: '登录已失效，请重新登录', icon: 'none' });
        }
      }
    })
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
.main-login-btn {
  width: 80vw;
  font-size: 32rpx;
  margin: 0 auto;
  margin-top: 32rpx;
}
.ins-btn {
  width: 80vw;
  font-size: 32rpx;
  margin: 0 auto;
  margin-top: 32rpx;
  border-radius: 32rpx;
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  color: #fff;
  font-weight: 500;
  box-shadow: 0 4rpx 24rpx #eaeaea;
  border: none;
  letter-spacing: 1rpx;
  transition: background 0.3s;
}
.ins-btn:active {
  background: linear-gradient(90deg, #92a8d1 0%, #f7cac9 100%);
}
.ins-modal {
  position: fixed;
  left: 0; top: 0; right: 0; bottom: 0;
  background: rgba(247,202,201,0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}
.ins-modal-content {
  background: #fff;
  border-radius: 32rpx;
  padding: 56rpx 36rpx 40rpx 36rpx;
  width: 82vw;
  text-align: center;
  box-shadow: 0 8rpx 32rpx #f7cac9;
  border: 2rpx solid #f7cac9;
}
.ins-modal-title {
  font-size: 40rpx;
  font-weight: bold;
  margin-bottom: 28rpx;
  color: #92a8d1;
  display: block;
  letter-spacing: 2rpx;
}
.ins-modal-desc {
  font-size: 30rpx;
  color: #666;
  margin-bottom: 36rpx;
  display: block;
  line-height: 1.7;
}
.ins-cancel-btn {
  width: 80vw;
  font-size: 30rpx;
  margin: 0 auto;
  margin-top: 24rpx;
  border-radius: 32rpx;
  background: #fff;
  color: #92a8d1;
  border: 2rpx solid #92a8d1;
  font-weight: 500;
  box-shadow: 0 2rpx 8rpx #eaeaea;
  letter-spacing: 1rpx;
}
.ins-warn-btn {
  background: linear-gradient(90deg, #f7cac9 0%, #f67280 100%);
  color: #fff;
  border: none;
  margin-top: 40rpx;
}
.ins-warn-text {
  color: #f67280;
  margin-top: 16rpx;
  font-size: 28rpx;
  text-align: center;
  letter-spacing: 1rpx;
}
</style> 