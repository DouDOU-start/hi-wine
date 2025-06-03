"use strict";
const common_vendor = require("../../common/vendor.js");
const config = require("../../config.js");
const utils_api = require("../../utils/api.js");
const _sfc_main = {
  data() {
    return {
      cart: [],
      tableList: [],
      selectedTableId: null,
      IMG_BASE_URL: config.IMG_BASE_URL,
      submitting: false
    };
  },
  computed: {
    totalPrice() {
      return this.cart.reduce((sum, item) => sum + item.price * item.count, 0);
    }
  },
  onShow() {
    this.loadCartData();
  },
  onLoad() {
    this.fetchTableList();
  },
  methods: {
    async fetchTableList() {
      const res = await common_vendor.index.request({ url: "http://你的后端地址/api/table/list" });
      if (res.data.code === 0) {
        this.tableList = res.data.data || [];
      }
    },
    onTableChange(e) {
      this.selectedTableId = this.tableList[e.detail.value].id;
    },
    // 加载购物车数据
    loadCartData() {
      const cartData = common_vendor.index.getStorageSync("cart") || [];
      this.cart = cartData;
    },
    // 修改商品数量
    changeCount(item, delta) {
      const idx = this.cart.findIndex((i) => i.id === item.id);
      if (idx !== -1) {
        this.cart[idx].count += delta;
        if (this.cart[idx].count < 1)
          this.cart[idx].count = 1;
        this.updateCartStorage();
      }
    },
    // 移除商品
    removeItem(item) {
      this.cart = this.cart.filter((i) => i.id !== item.id);
      this.updateCartStorage();
    },
    // 更新本地存储
    updateCartStorage() {
      common_vendor.index.setStorageSync("cart", this.cart);
    },
    // 显示下单确认弹窗
    showOrderConfirm() {
      if (this.cart.length === 0) {
        common_vendor.index.showToast({
          title: "购物车为空",
          icon: "none"
        });
        return;
      }
      const token = common_vendor.index.getStorageSync("token");
      if (!token) {
        common_vendor.index.showToast({
          title: "请先登录",
          icon: "none"
        });
        setTimeout(() => {
          common_vendor.index.switchTab({
            url: "/pages/profile/index"
          });
        }, 1500);
        return;
      }
      common_vendor.index.showModal({
        title: "确认下单",
        content: "确定要提交订单吗？",
        success: (res) => {
          if (res.confirm) {
            this.submitOrder();
          }
        }
      });
    },
    // 提交订单
    async submitOrder() {
      if (this.cart.length === 0) {
        common_vendor.index.showToast({ title: "购物车为空", icon: "none" });
        return;
      }
      this.submitting = true;
      common_vendor.index.showLoading({ title: "订单提交中...", mask: true });
      try {
        const items = this.cart.map((item) => ({
          productId: item.id,
          quantity: item.count
        }));
        const res = await utils_api.api.createOrder(this.selectedTableId || 0, items);
        if (res && res.code === 0 && res.data && res.data.id) {
          this.cart = [];
          this.updateCartStorage();
          common_vendor.index.hideLoading();
          common_vendor.index.showToast({ title: "下单成功", icon: "success" });
          setTimeout(() => {
            common_vendor.index.switchTab({ url: "/pages/order/index" });
          }, 1500);
        } else {
          throw new Error("创建订单失败");
        }
      } catch (err) {
        common_vendor.index.showToast({ title: "下单失败，请重试", icon: "none" });
      } finally {
        this.submitting = false;
        common_vendor.index.hideLoading();
      }
    }
  }
};
if (!Array) {
  const _component_uni_popup_dialog = common_vendor.resolveComponent("uni-popup-dialog");
  const _component_uni_popup = common_vendor.resolveComponent("uni-popup");
  (_component_uni_popup_dialog + _component_uni_popup)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  var _a;
  return common_vendor.e({
    a: $data.cart.length === 0
  }, $data.cart.length === 0 ? {} : {
    b: common_vendor.f($data.cart, (item, k0, i0) => {
      return {
        a: item.image ? item.image : $data.IMG_BASE_URL + "/wine.png",
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.price),
        d: common_vendor.o(($event) => $options.changeCount(item, -1), item.id),
        e: item.count <= 1,
        f: common_vendor.t(item.count),
        g: common_vendor.o(($event) => $options.changeCount(item, 1), item.id),
        h: common_vendor.o(($event) => $options.removeItem(item), item.id),
        i: item.id
      };
    }),
    c: common_vendor.t($options.totalPrice),
    d: common_vendor.o((...args) => $options.showOrderConfirm && $options.showOrderConfirm(...args))
  }, {
    e: common_vendor.p({
      title: "处理中",
      content: "订单提交中，请稍候...",
      ["show-cancel"]: false,
      ["before-close"]: true
    }),
    f: common_vendor.sr("loading", "8039fbf1-0"),
    g: common_vendor.p({
      type: "dialog",
      ["mask-click"]: false
    }),
    h: common_vendor.t($data.selectedTableId ? ((_a = $data.tableList.find((t) => t.id === $data.selectedTableId)) == null ? void 0 : _a.name) || "" : "未选择"),
    i: $data.tableList,
    j: common_vendor.o((...args) => $options.onTableChange && $options.onTableChange(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-8039fbf1"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/cart/index.js.map
