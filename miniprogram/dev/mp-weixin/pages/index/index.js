"use strict";
const common_vendor = require("../../common/vendor.js");
const config = require("../../config.js");
const utils_api = require("../../utils/api.js");
function debounce(fn, delay) {
  let timer = null;
  return function(...args) {
    if (timer)
      clearTimeout(timer);
    timer = setTimeout(() => fn.apply(this, args), delay);
  };
}
const _sfc_main = {
  data() {
    return {
      categories: [{ name: "全部" }],
      activeTab: 0,
      wines: [],
      loading: false,
      page: 1,
      size: 10,
      hasMore: true,
      IMG_BASE_URL: config.IMG_BASE_URL,
      debugInfo: ""
    };
  },
  onLoad() {
    this.loadCategories();
    this.loadProducts();
  },
  methods: {
    // 加载商品分类
    async loadCategories() {
      try {
        const res = await utils_api.api.getCategoryList();
        common_vendor.index.__f__("log", "at pages/index/index.vue:74", "分类响应数据:", JSON.stringify(res));
        if (res && res.data && res.data.list) {
          this.categories = [{ name: "全部" }, ...res.data.list];
          this.debugInfo = "分类加载成功";
        } else {
          this.debugInfo = "分类数据格式不正确: " + JSON.stringify(res);
        }
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/index/index.vue:83", "加载分类失败", err);
        this.debugInfo = "分类加载失败: " + err.message;
      }
    },
    // 切换分类
    changeCategory(index) {
      if (this.activeTab === index)
        return;
      this.activeTab = index;
      this.wines = [];
      this.page = 1;
      this.hasMore = true;
      this.loadProducts();
    },
    // 加载商品列表
    async loadProducts() {
      if (this.loading || !this.hasMore)
        return;
      this.loading = true;
      try {
        const categoryId = this.activeTab === 0 ? 0 : this.categories[this.activeTab].id;
        common_vendor.index.__f__("log", "at pages/index/index.vue:105", "请求商品列表，分类ID:", categoryId);
        const res = await utils_api.api.getProductList(categoryId, "", this.page, this.size);
        common_vendor.index.__f__("log", "at pages/index/index.vue:107", "商品列表响应:", JSON.stringify(res));
        if (res && res.data && res.data.list) {
          if (this.page === 1) {
            this.wines = res.data.list;
          } else {
            this.wines = [...this.wines, ...res.data.list];
          }
          this.hasMore = res.data.list.length === this.size;
          this.page++;
          this.debugInfo = "商品加载成功，共" + this.wines.length + "条";
        } else {
          this.debugInfo = "商品数据格式不正确: " + JSON.stringify(res);
        }
      } catch (err) {
        common_vendor.index.__f__("error", "at pages/index/index.vue:125", "加载商品失败", err);
        this.debugInfo = "商品加载失败: " + err.message;
        common_vendor.index.showToast({
          title: "加载商品失败",
          icon: "none"
        });
      } finally {
        this.loading = false;
      }
    },
    // 添加到购物车
    addToCart: debounce(function(item) {
      let cart = common_vendor.index.getStorageSync("cart") || [];
      const idx = cart.findIndex((i) => i.id === item.id);
      let changed = false;
      if (idx !== -1) {
        cart[idx].count += 1;
        changed = true;
      } else {
        cart.push({
          id: item.id,
          name: item.name,
          price: item.price,
          image: item.image,
          count: 1
        });
        changed = true;
      }
      if (changed) {
        common_vendor.index.setStorageSync("cart", cart);
        common_vendor.index.showToast({ title: `已加入购物车`, icon: "success" });
      }
    }, 300)
  },
  // 下拉刷新
  onPullDownRefresh() {
    this.wines = [];
    this.page = 1;
    this.hasMore = true;
    this.loadProducts().then(() => {
      common_vendor.index.stopPullDownRefresh();
    });
  },
  // 上拉加载更多
  onReachBottom() {
    this.loadProducts();
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.f($data.categories, (cat, idx, i0) => {
      return {
        a: common_vendor.t(idx === 0 ? "全部" : cat.name),
        b: cat.id || idx,
        c: common_vendor.n(idx === $data.activeTab ? "active" : ""),
        d: common_vendor.o(($event) => $options.changeCategory(idx), cat.id || idx)
      };
    }),
    b: common_vendor.f($data.wines, (item, k0, i0) => {
      return {
        a: item.image ? item.image : $data.IMG_BASE_URL + "/wine.png",
        b: common_vendor.t(item.name),
        c: common_vendor.t(item.price),
        d: common_vendor.o(($event) => $options.addToCart(item), item.id),
        e: item.id
      };
    }),
    c: $data.loading
  }, $data.loading ? {} : {}, {
    d: !$data.loading && $data.wines.length === 0
  }, !$data.loading && $data.wines.length === 0 ? {} : {}, {
    e: $data.debugInfo
  }, $data.debugInfo ? {
    f: common_vendor.t($data.categories.length),
    g: common_vendor.t($data.wines.length),
    h: common_vendor.t($data.debugInfo)
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-1cf27b2a"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/index/index.js.map
