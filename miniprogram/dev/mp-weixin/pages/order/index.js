"use strict";
const common_vendor = require("../../common/vendor.js");
const config = require("../../config.js");
const utils_api = require("../../utils/api.js");
const _sfc_main = {
  data() {
    return {
      orders: [],
      orderItems: {},
      // 订单商品映射表，key为订单ID
      loading: false,
      page: 1,
      size: 10,
      hasMore: true,
      IMG_BASE_URL: config.IMG_BASE_URL,
      activeTab: 0,
      statusTabs: [
        { name: "全部", value: -1 },
        { name: "待支付", value: 0 },
        { name: "已支付", value: 1 },
        { name: "已完成", value: 2 },
        { name: "已取消", value: 3 }
      ],
      popupConfig: {
        title: "",
        content: "",
        onConfirm: null
      }
    };
  },
  onLoad() {
    this.loadOrders();
  },
  onShow() {
    this.refreshOrders();
  },
  methods: {
    // 切换状态选项卡
    changeTab(index) {
      if (this.activeTab === index)
        return;
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
      if (this.loading || !this.hasMore)
        return;
      this.loading = true;
      try {
        const status = this.statusTabs[this.activeTab].value;
        const res = await utils_api.api.getOrderList(status, this.page, this.size);
        common_vendor.index.__f__("log", "at pages/order/index.vue:123", "订单列表响应:", JSON.stringify(res));
        if (res && res.code === 0 && res.data && res.data.list) {
          if (this.page === 1) {
            this.orders = res.data.list;
          } else {
            this.orders = [...this.orders, ...res.data.list];
          }
          this.hasMore = res.data.list.length === this.size;
          this.page++;
          for (const order of res.data.list) {
            this.loadOrderItems(order.id);
          }
        }
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/order/index.vue:142", "加载订单失败", err);
        common_vendor.index.showToast({
          title: "加载订单失败",
          icon: "none"
        });
      } finally {
        this.loading = false;
      }
    },
    // 加载订单商品
    async loadOrderItems(orderId) {
      try {
        const res = await utils_api.api.getOrderDetail(orderId);
        common_vendor.index.__f__("log", "at pages/order/index.vue:156", "订单详情响应:", JSON.stringify(res));
        if (res && res.code === 0 && res.data && res.data.orderItems) {
          this.$set(this.orderItems, orderId, res.data.orderItems);
        }
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/order/index.vue:162", `加载订单${orderId}商品失败`, err);
      }
    },
    // 支付订单
    payOrder(order) {
      this.popupConfig = {
        title: "支付订单",
        content: `确认支付订单金额 ￥${order.totalAmount}？`,
        onConfirm: async () => {
          try {
            const res = await utils_api.api.updateOrderStatus(order.id, 1);
            if (res && res.code === 0) {
              common_vendor.index.showToast({
                title: "支付成功",
                icon: "success"
              });
              this.refreshOrders();
            } else {
              throw new Error("支付失败");
            }
          } catch (err) {
            common_vendor.index.__f__("error", "at pages/order/index.vue:186", "支付失败", err);
            common_vendor.index.showToast({
              title: "支付失败",
              icon: "none"
            });
          }
        }
      };
      this.$refs.confirmPopup.open();
    },
    // 取消订单
    cancelOrder(order) {
      this.popupConfig = {
        title: "取消订单",
        content: "确认取消该订单？",
        onConfirm: async () => {
          try {
            const res = await utils_api.api.updateOrderStatus(order.id, 3);
            if (res && res.code === 0) {
              common_vendor.index.showToast({
                title: "订单已取消",
                icon: "success"
              });
              this.refreshOrders();
            } else {
              throw new Error("取消订单失败");
            }
          } catch (err) {
            common_vendor.index.__f__("error", "at pages/order/index.vue:217", "取消订单失败", err);
            common_vendor.index.showToast({
              title: "取消订单失败",
              icon: "none"
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
      if (!time)
        return "";
      return time.replace("T", " ").split(".")[0];
    },
    // 获取状态文本
    getStatusText(status) {
      switch (status) {
        case 0:
          return "待支付";
        case 1:
          return "已支付";
        case 2:
          return "已完成";
        case 3:
          return "已取消";
        default:
          return "未知状态";
      }
    }
  },
  // 下拉刷新
  onPullDownRefresh() {
    this.refreshOrders();
    common_vendor.index.stopPullDownRefresh();
  },
  // 上拉加载更多
  onReachBottom() {
    this.loadOrders();
  }
};
if (!Array) {
  const _component_uni_popup_dialog = common_vendor.resolveComponent("uni-popup-dialog");
  const _component_uni_popup = common_vendor.resolveComponent("uni-popup");
  (_component_uni_popup_dialog + _component_uni_popup)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.f($data.statusTabs, (tab, index, i0) => {
      return {
        a: common_vendor.t(tab.name),
        b: index,
        c: common_vendor.n($data.activeTab === index ? "active" : ""),
        d: common_vendor.o(($event) => $options.changeTab(index), index)
      };
    }),
    b: $data.loading
  }, $data.loading ? {} : $data.orders.length === 0 ? {} : {
    d: common_vendor.f($data.orders, (order, k0, i0) => {
      return common_vendor.e({
        a: common_vendor.t($options.formatTime(order.createTime)),
        b: common_vendor.t($options.getStatusText(order.status)),
        c: common_vendor.n("status-" + order.status),
        d: common_vendor.f($data.orderItems[order.id], (item, k1, i1) => {
          return {
            a: item.productImage || $data.IMG_BASE_URL + "/wine.png",
            b: common_vendor.t(item.productName),
            c: common_vendor.t(item.quantity),
            d: item.id
          };
        }),
        e: common_vendor.t(order.totalAmount),
        f: order.status === 0
      }, order.status === 0 ? {
        g: common_vendor.o(($event) => $options.payOrder(order), order.id),
        h: common_vendor.o(($event) => $options.cancelOrder(order), order.id)
      } : {}, {
        i: order.id
      });
    })
  }, {
    c: $data.orders.length === 0,
    e: $data.orders.length > 0 && !$data.hasMore
  }, $data.orders.length > 0 && !$data.hasMore ? {} : {}, {
    f: common_vendor.o($data.popupConfig.onConfirm),
    g: common_vendor.o($options.closePopup),
    h: common_vendor.p({
      title: $data.popupConfig.title,
      content: $data.popupConfig.content,
      ["before-close"]: true
    }),
    i: common_vendor.sr("confirmPopup", "17a44f9d-0"),
    j: common_vendor.p({
      type: "dialog"
    })
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-17a44f9d"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/order/index.js.map
