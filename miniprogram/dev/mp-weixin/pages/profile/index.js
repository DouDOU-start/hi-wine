"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_api = require("../../utils/api.js");
const _sfc_main = {
  data() {
    return {
      userInfo: common_vendor.index.getStorageSync("userInfo") || {},
      showSettingBtn: false,
      showLoginBtn: !common_vendor.index.getStorageSync("token"),
      tempUserInfo: {
        avatar: "",
        nickname: ""
      }
    };
  },
  onLoad() {
    this.initUserStatus();
  },
  onShow() {
    const userInfo = common_vendor.index.getStorageSync("userInfo") || {};
    this.userInfo = userInfo;
    this.showLoginBtn = !common_vendor.index.getStorageSync("token");
    if (this.showLoginBtn) {
      common_vendor.index.showToast({
        title: "请选择头像和昵称后登录",
        icon: "none",
        duration: 2e3
      });
    }
  },
  methods: {
    initUserStatus() {
      const token = common_vendor.index.getStorageSync("token");
      if (!token) {
        this.showLoginBtn = true;
        this.showSettingBtn = false;
      } else {
        this.showLoginBtn = false;
        this.showSettingBtn = false;
      }
    },
    goOrder() {
      common_vendor.index.switchTab({ url: "/pages/order/index" });
    },
    // 获取头像
    onChooseAvatar(e) {
      const { avatarUrl } = e.detail;
      this.tempUserInfo.avatar = avatarUrl;
      this.userInfo.avatar = avatarUrl;
      if (this.tempUserInfo.nickname) {
        common_vendor.index.showToast({
          title: "头像已选择，请点击登录",
          icon: "none",
          duration: 1500
        });
      } else {
        common_vendor.index.showToast({
          title: "请输入昵称",
          icon: "none",
          duration: 1500
        });
      }
    },
    // 获取昵称
    onNicknameInput(e) {
      const nickname = e.detail.value;
      this.tempUserInfo.nickname = nickname;
      this.userInfo.nickname = nickname;
      if (this.tempUserInfo.avatar) {
        common_vendor.index.showToast({
          title: "昵称已输入，请点击登录",
          icon: "none",
          duration: 1500
        });
      } else {
        common_vendor.index.showToast({
          title: "请选择头像",
          icon: "none",
          duration: 1500
        });
      }
    },
    // 确认登录
    doWechatLogin() {
      if (!this.tempUserInfo.avatar) {
        common_vendor.index.showToast({
          title: "请选择头像",
          icon: "none"
        });
        return;
      }
      if (!this.tempUserInfo.nickname) {
        common_vendor.index.showToast({
          title: "请输入昵称",
          icon: "none"
        });
        return;
      }
      common_vendor.index.showLoading({
        title: "登录中...",
        mask: true
      });
      common_vendor.index.login({
        success: (res) => {
          if (res.code) {
            common_vendor.index.__f__("log", "at pages/profile/index.vue:178", "获取到临时登录凭证code:", res.code);
            utils_api.api.login(res.code).then((resp) => {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("log", "at pages/profile/index.vue:185", "登录成功，响应:", JSON.stringify(resp));
              if (!resp.data.token) {
                throw new Error("服务器未返回有效的登录凭证");
              }
              common_vendor.index.setStorageSync("token", resp.data.token);
              const userInfo = {
                ...resp.data.userInfo,
                nickname: this.tempUserInfo.nickname,
                avatar: this.tempUserInfo.avatar
              };
              common_vendor.index.setStorageSync("userInfo", userInfo);
              this.userInfo = userInfo;
              this.showLoginBtn = false;
              this.updateUserInfo(resp.data.token, this.tempUserInfo.nickname, this.tempUserInfo.avatar);
              common_vendor.index.showToast({
                title: "登录成功",
                icon: "success"
              });
            }).catch((err) => {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("error", "at pages/profile/index.vue:216", "登录失败", err);
              common_vendor.index.showModal({
                title: "登录失败",
                content: `服务器返回错误: ${err.message || "未知错误"}`,
                showCancel: false
              });
            });
          } else {
            common_vendor.index.hideLoading();
            common_vendor.index.__f__("error", "at pages/profile/index.vue:228", "获取code失败", res);
            common_vendor.index.showToast({
              title: "获取临时登录凭证失败",
              icon: "none"
            });
          }
        },
        fail: (err) => {
          common_vendor.index.hideLoading();
          common_vendor.index.__f__("error", "at pages/profile/index.vue:239", "wx.login 失败", err);
          common_vendor.index.showToast({
            title: "微信登录失败",
            icon: "none"
          });
        }
      });
    },
    // 更新用户信息到服务器
    updateUserInfo(token, nickname, avatar) {
      utils_api.api.updateUserInfo(nickname, avatar).then((resp) => {
        common_vendor.index.__f__("log", "at pages/profile/index.vue:252", "用户信息更新成功", resp);
      }).catch((err) => {
        common_vendor.index.__f__("error", "at pages/profile/index.vue:255", "用户信息更新失败", err);
      });
    },
    // 上传头像到服务器
    uploadAvatar(tempFilePath, callback) {
      if (tempFilePath.startsWith("wxfile://") || tempFilePath.startsWith("http://tmp/")) {
        common_vendor.index.getFileSystemManager().readFile({
          filePath: tempFilePath,
          encoding: "base64",
          success: (res) => {
            const base64Img = "data:image/jpeg;base64," + res.data;
            callback(base64Img);
          },
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/profile/index.vue:275", "读取头像文件失败", err);
            callback(tempFilePath);
          }
        });
      } else {
        callback(tempFilePath);
      }
    },
    openSetting() {
      common_vendor.index.openSetting({
        success: (res) => {
          this.showSettingBtn = false;
          common_vendor.index.showToast({ title: "请重新登录", icon: "none" });
        }
      });
    },
    // 退出登录
    logout() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要退出登录吗？",
        success: (res) => {
          if (res.confirm) {
            common_vendor.index.removeStorageSync("token");
            common_vendor.index.removeStorageSync("userInfo");
            this.userInfo = {};
            this.showLoginBtn = true;
            this.tempUserInfo = {
              avatar: "",
              nickname: ""
            };
            common_vendor.index.showToast({
              title: "已退出登录",
              icon: "success"
            });
          }
        }
      });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.showLoginBtn
  }, $data.showLoginBtn ? {
    b: $data.userInfo.avatar || "https://randomuser.me/api/portraits/men/32.jpg",
    c: common_vendor.o((...args) => $options.onChooseAvatar && $options.onChooseAvatar(...args))
  } : {
    d: $data.userInfo.avatar || "https://randomuser.me/api/portraits/men/32.jpg"
  }, {
    e: $data.showLoginBtn
  }, $data.showLoginBtn ? {
    f: common_vendor.o((...args) => $options.onNicknameInput && $options.onNicknameInput(...args))
  } : {
    g: common_vendor.t($data.userInfo.nickname || "酒友")
  }, {
    h: $data.showLoginBtn && $data.tempUserInfo.avatar && $data.tempUserInfo.nickname
  }, $data.showLoginBtn && $data.tempUserInfo.avatar && $data.tempUserInfo.nickname ? {
    i: common_vendor.o((...args) => $options.doWechatLogin && $options.doWechatLogin(...args))
  } : {}, {
    j: $data.showLoginBtn && (!$data.tempUserInfo.avatar || !$data.tempUserInfo.nickname)
  }, $data.showLoginBtn && (!$data.tempUserInfo.avatar || !$data.tempUserInfo.nickname) ? {} : {}, {
    k: !$data.showLoginBtn
  }, !$data.showLoginBtn ? {
    l: common_vendor.o((...args) => $options.logout && $options.logout(...args))
  } : {}, {
    m: common_vendor.o((...args) => $options.goOrder && $options.goOrder(...args)),
    n: $data.showSettingBtn
  }, $data.showSettingBtn ? {
    o: common_vendor.o((...args) => $options.openSetting && $options.openSetting(...args))
  } : {}, {
    p: $data.showSettingBtn
  }, $data.showSettingBtn ? {} : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-201c0da5"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/profile/index.js.map
