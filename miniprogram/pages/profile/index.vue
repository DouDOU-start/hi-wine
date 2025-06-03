<template>
  <view class="profile-container">
    <view class="profile-header">
      <button 
        v-if="showLoginBtn" 
        open-type="chooseAvatar" 
        @chooseavatar="onChooseAvatar"
        class="avatar-wrapper"
      >
        <image class="profile-avatar" :src="userInfo.avatar || 'https://randomuser.me/api/portraits/men/32.jpg'" />
      </button>
      <image v-else class="profile-avatar" :src="userInfo.avatar || 'https://randomuser.me/api/portraits/men/32.jpg'" />
      
      <view v-if="showLoginBtn" class="nickname-input">
        <input 
          type="nickname" 
          placeholder="请输入昵称" 
          @blur="onNicknameInput" 
          class="nickname-field" 
        />
      </view>
      <text v-else class="profile-nick">Hi，{{ userInfo.nickname || '酒友' }}</text>
      
      <button
        v-if="showLoginBtn && tempUserInfo.avatar && tempUserInfo.nickname"
        @tap="doWechatLogin"
        type="primary"
        class="main-login-btn ins-btn"
      >点击登录</button>
      <view v-if="showLoginBtn && (!tempUserInfo.avatar || !tempUserInfo.nickname)" class="login-tip">
        请先选择头像并输入昵称
      </view>
      <button
        v-if="!showLoginBtn"
        @tap="logout"
        type="default"
        class="logout-btn"
      >退出登录</button>
    </view>
    <view class="profile-list">
      <view class="profile-item" @tap="goOrder">
        <text class="profile-icon">📦</text>
        <text class="profile-label">我的订单</text>
      </view>
    </view>
    
    <!-- 授权失败提示 -->
    <button v-if="showSettingBtn" @tap="openSetting" type="warn" class="ins-btn ins-warn-btn">去设置重新授权</button>
    <view v-if="showSettingBtn" class="ins-warn-text">请在设置中授权后，再点击"微信一键登录"</view>
  </view>
</template>

<script>
import api from '@/utils/api.js';

export default {
  data() {
    return {
      userInfo: uni.getStorageSync('userInfo') || {},
      showSettingBtn: false,
      showLoginBtn: !uni.getStorageSync('token'),
      tempUserInfo: {
        avatar: '',
        nickname: ''
      }
    };
  },
  onLoad() {
    this.initUserStatus();
  },
  onShow() {
    // 优先用本地userInfo渲染
    const userInfo = uni.getStorageSync('userInfo') || {};
    this.userInfo = userInfo;
    this.showLoginBtn = !uni.getStorageSync('token');
    
    // 如果未登录，显示提示信息
    if (this.showLoginBtn) {
      uni.showToast({
        title: '请选择头像和昵称后登录',
        icon: 'none',
        duration: 2000
      });
    }
  },
  methods: {
    initUserStatus() {
      const token = uni.getStorageSync('token');
      if (!token) {
        this.showLoginBtn = true;
        this.showSettingBtn = false;
      } else {
        this.showLoginBtn = false;
        this.showSettingBtn = false;
      }
    },
    goOrder() {
      uni.switchTab({ url: '/pages/order/index' });
    },
    
    // 获取头像
    onChooseAvatar(e) {
      const { avatarUrl } = e.detail;
      this.tempUserInfo.avatar = avatarUrl;
      
      // 将临时头像显示到界面
      this.userInfo.avatar = avatarUrl;
      
      // 提示用户
      if (this.tempUserInfo.nickname) {
        uni.showToast({
          title: '头像已选择，请点击登录',
          icon: 'none',
          duration: 1500
        });
      } else {
        uni.showToast({
          title: '请输入昵称',
          icon: 'none',
          duration: 1500
        });
      }
    },
    
    // 获取昵称
    onNicknameInput(e) {
      const nickname = e.detail.value;
      this.tempUserInfo.nickname = nickname;
      
      // 将临时昵称显示到界面
      this.userInfo.nickname = nickname;
      
      // 提示用户
      if (this.tempUserInfo.avatar) {
        uni.showToast({
          title: '昵称已输入，请点击登录',
          icon: 'none',
          duration: 1500
        });
      } else {
        uni.showToast({
          title: '请选择头像',
          icon: 'none',
          duration: 1500
        });
      }
    },
    
    // 确认登录
    doWechatLogin() {
      // 检查是否已选择头像和填写昵称
      if (!this.tempUserInfo.avatar) {
        uni.showToast({
          title: '请选择头像',
          icon: 'none'
        });
        return;
      }
      
      if (!this.tempUserInfo.nickname) {
        uni.showToast({
          title: '请输入昵称',
          icon: 'none'
        });
        return;
      }
      
      // 显示加载中
      uni.showLoading({
        title: '登录中...',
        mask: true
      });
      
      // 1. 调用 wx.login() 获取临时登录凭证 code
      uni.login({
        success: (res) => {
          if (res.code) {
            console.log('获取到临时登录凭证code:', res.code);
            
            // 2. 将code发送到开发者服务器
            api.login(res.code)
              .then(resp => {
                // 隐藏加载
                uni.hideLoading();
                console.log('登录成功，响应:', JSON.stringify(resp));
                
                if (!resp.data.token) {
                  throw new Error('服务器未返回有效的登录凭证');
                }
                
                // 3. 保存自定义登录态
                uni.setStorageSync('token', resp.data.token);
                
                // 4. 保存用户信息
                const userInfo = {
                  ...resp.data.userInfo,
                  nickname: this.tempUserInfo.nickname,
                  avatar: this.tempUserInfo.avatar
                };
                uni.setStorageSync('userInfo', userInfo);
                this.userInfo = userInfo;
                this.showLoginBtn = false;
                
                // 5. 更新用户信息到服务器
                this.updateUserInfo(resp.data.token, this.tempUserInfo.nickname, this.tempUserInfo.avatar);
                
                uni.showToast({
                  title: '登录成功',
                  icon: 'success'
                });
              })
              .catch(err => {
                // 隐藏加载
                uni.hideLoading();
                
                console.error('登录失败', err);
                // 显示更详细的错误信息
                uni.showModal({
                  title: '登录失败',
                  content: `服务器返回错误: ${err.message || '未知错误'}`,
                  showCancel: false
                });
              });
          } else {
            // 隐藏加载
            uni.hideLoading();
            
            console.error('获取code失败', res);
            uni.showToast({ 
              title: '获取临时登录凭证失败', 
              icon: 'none' 
            });
          }
        },
        fail: (err) => {
          // 隐藏加载
          uni.hideLoading();
          
          console.error('wx.login 失败', err);
          uni.showToast({ 
            title: '微信登录失败', 
            icon: 'none' 
          });
        }
      });
    },
    
    // 更新用户信息到服务器
    updateUserInfo(token, nickname, avatar) {
      api.updateUserInfo(nickname, avatar)
        .then(resp => {
          console.log('用户信息更新成功', resp);
        })
        .catch(err => {
          console.error('用户信息更新失败', err);
          // 不影响用户体验，静默失败
        });
    },
    
    // 上传头像到服务器
    uploadAvatar(tempFilePath, callback) {
      // 微信临时文件路径可能在某些情况下无法直接使用
      // 检查是否是临时文件路径，如果是则需要处理
      if (tempFilePath.startsWith('wxfile://') || tempFilePath.startsWith('http://tmp/')) {
        // 将临时文件上传到服务器或转为base64
        uni.getFileSystemManager().readFile({
          filePath: tempFilePath,
          encoding: 'base64',
          success: (res) => {
            // 将base64编码的图片传递给回调
            const base64Img = 'data:image/jpeg;base64,' + res.data;
            callback(base64Img);
          },
          fail: (err) => {
            console.error('读取头像文件失败', err);
            // 读取失败时，尝试直接使用临时路径
            callback(tempFilePath);
          }
        });
      } else {
        // 如果不是临时文件路径，直接使用
        callback(tempFilePath);
      }
    },
    
    openSetting() {
      uni.openSetting({
        success: (res) => {
          this.showSettingBtn = false;
          uni.showToast({ title: '请重新登录', icon: 'none' });
        }
      });
    },
    
    // 退出登录
    logout() {
      uni.showModal({
        title: '提示',
        content: '确定要退出登录吗？',
        success: (res) => {
          if (res.confirm) {
            // 清除token和用户信息
            uni.removeStorageSync('token');
            uni.removeStorageSync('userInfo');
            
            // 更新状态
            this.userInfo = {};
            this.showLoginBtn = true;
            this.tempUserInfo = {
              avatar: '',
              nickname: ''
            };
            
            uni.showToast({
              title: '已退出登录',
              icon: 'success'
            });
          }
        }
      });
    }
  }
};
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: #f8f6f4;
  padding-bottom: 40rpx;
}
.profile-header {
  padding: 60rpx 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}
.avatar-wrapper {
  padding: 0;
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  background-color: transparent;
}
.avatar-wrapper::after {
  border: none;
}
.profile-avatar {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  border: 6rpx solid #fff;
  box-shadow: 0 8rpx 32rpx rgba(0,0,0,0.1);
}
.nickname-input {
  margin-top: 20rpx;
  width: 80%;
  text-align: center;
}
.nickname-field {
  border-bottom: 1px solid #f7cac9;
  padding: 10rpx 0;
  font-size: 36rpx;
  text-align: center;
}
.profile-nick {
  margin-top: 20rpx;
  font-size: 40rpx;
  font-weight: bold;
  color: #222;
}
.login-tip {
  margin-top: 20rpx;
  font-size: 28rpx;
  color: #f7cac9;
  text-align: center;
}
.main-login-btn {
  margin-top: 30rpx;
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  border: none;
  color: #fff;
  font-size: 32rpx;
  padding: 16rpx 60rpx;
  border-radius: 40rpx;
  box-shadow: 0 8rpx 16rpx rgba(0,0,0,0.1);
}
.logout-btn {
  margin-top: 20rpx;
  font-size: 28rpx;
  color: #999;
  background: #f5f5f5;
  border: none;
  border-radius: 30rpx;
  padding: 10rpx 40rpx;
}
.profile-list {
  margin-top: 40rpx;
  padding: 0 30rpx;
}
.profile-item {
  background: #fff;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.05);
}
.profile-icon {
  font-size: 48rpx;
  margin-right: 20rpx;
}
.profile-label {
  font-size: 32rpx;
  color: #333;
}
.ins-btn {
  background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
  border: none;
  color: #fff;
  font-size: 32rpx;
  padding: 16rpx 60rpx;
  border-radius: 40rpx;
  box-shadow: 0 8rpx 16rpx rgba(0,0,0,0.1);
}
.ins-warn-btn {
  margin: 30rpx auto;
  display: block;
  background: #ff6b6b;
}
.ins-warn-text {
  text-align: center;
  color: #ff6b6b;
  font-size: 24rpx;
}
</style> 