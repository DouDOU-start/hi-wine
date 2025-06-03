"use strict";
const common_vendor = require("../common/vendor.js");
const config = require("../config.js");
const getToken = () => {
  return common_vendor.index.getStorageSync("token") || "";
};
const request = (url, method, data, needAuth = true) => {
  return new Promise((resolve, reject) => {
    common_vendor.index.request({
      url: config.BASE_URL + url,
      method,
      data,
      header: {
        "Content-Type": "application/json",
        "Authorization": needAuth ? `Bearer ${getToken()}` : ""
      },
      success: (res) => {
        if (res.statusCode === 200) {
          common_vendor.index.__f__("log", "at utils/api.js:23", `API响应(${url}):`, JSON.stringify(res.data));
          if (res.data.hasOwnProperty("code")) {
            if (res.data.code === 0) {
              resolve(res.data);
            } else {
              common_vendor.index.showToast({
                title: res.data.message || "请求失败",
                icon: "none"
              });
              reject(new Error(res.data.message || "请求失败"));
            }
          } else {
            resolve({ code: 0, message: "操作成功", data: res.data });
          }
        } else if (res.statusCode === 401) {
          common_vendor.index.showToast({
            title: "请先登录",
            icon: "none"
          });
          common_vendor.index.removeStorageSync("token");
          setTimeout(() => {
            common_vendor.index.switchTab({
              url: "/pages/profile/index"
            });
          }, 1500);
          reject(new Error("未登录或登录已过期"));
        } else {
          common_vendor.index.showToast({
            title: res.data.message || "请求失败",
            icon: "none"
          });
          reject(new Error(res.data.message || "请求失败"));
        }
      },
      fail: (err) => {
        common_vendor.index.showToast({
          title: "网络错误",
          icon: "none"
        });
        reject(err);
      }
    });
  });
};
const api = {
  // 用户登录
  login(code) {
    return request("/wechat/login", "POST", {
      code
      // 临时登录凭证
    }, false);
  },
  // 更新用户信息
  updateUserInfo(nickname, avatar) {
    return request("/api/user/update", "POST", {
      nickname,
      avatar
    }, true);
  },
  // 获取商品分类列表
  getCategoryList() {
    return request("/api/category/list", "GET");
  },
  // 获取商品列表
  getProductList(categoryId = 0, keyword = "", page = 1, size = 10) {
    return request("/api/product/list", "GET", { categoryId, keyword, page, size });
  },
  // 获取商品详情
  getProductDetail(id) {
    return request("/api/product/detail", "GET", { id });
  },
  // 创建订单
  createOrder(tableId, items) {
    return request("/api/order/create", "POST", { tableId, items });
  },
  // 获取订单列表
  getOrderList(status = -1, page = 1, size = 10) {
    return request("/api/order/list", "GET", { status, page, size });
  },
  // 获取订单详情
  getOrderDetail(id) {
    return request("/api/order/detail", "GET", { id });
  },
  // 更新订单状态
  updateOrderStatus(id, status) {
    return request("/api/order/updateStatus", "POST", { id, status });
  }
};
exports.api = api;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/api.js.map
