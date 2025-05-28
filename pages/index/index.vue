<template>
	<view class="ins-container">
		<view class="ins-header">
			<text class="ins-title">HI WINE 酒馆</text>
			<text class="ins-sub">INS风格酒水体验馆</text>
		</view>
		<view class="ins-tabs">
			<view
				v-for="(cat, idx) in categories"
				:key="cat.id || idx"
				:class="['ins-tab', idx === activeTab ? 'active' : '']"
				@tap="changeCategory(idx)"
			>
				{{ idx === 0 ? '全部' : cat.name }}
			</view>
		</view>
		<view class="ins-list">
			<view v-for="item in wines" :key="item.id" class="ins-card">
				<image :src="item.image ? item.image : IMG_BASE_URL + '/wine.png'" class="ins-img" mode="aspectFill" lazy-load="true" :style="{background:'#f3f3f3'}" />
				<view class="ins-info">
					<text class="ins-name">{{ item.name }}</text>
					<text class="ins-price">￥{{ item.price }}</text>
				</view>
				<button class="ins-btn" @tap="addToCart(item)">加入购物车</button>
			</view>
		</view>
		<view v-if="loading" class="loading">加载中...</view>
		<view v-if="!loading && wines.length === 0" class="empty">暂无商品</view>
	</view>
</template>

<script>
import { IMG_BASE_URL } from '@/config.js';
import api from '@/utils/api.js';

function debounce(fn, delay) {
	let timer = null;
	return function(...args) {
		if (timer) clearTimeout(timer);
		timer = setTimeout(() => fn.apply(this, args), delay);
	};
}

export default {
	data() {
		return {
			categories: [{ name: '全部' }],
			activeTab: 0,
			wines: [],
			loading: false,
			page: 1,
			size: 10,
			hasMore: true,
			IMG_BASE_URL
		}
	},
	onLoad() {
		this.loadCategories();
		this.loadProducts();
	},
	methods: {
		// 加载商品分类
		async loadCategories() {
			try {
				const res = await api.getCategoryList();
				if (res && res.list) {
					// 添加"全部"分类
					this.categories = [{ name: '全部' }, ...res.list];
				}
			} catch (err) {
				console.error('加载分类失败', err);
			}
		},
		
		// 切换分类
		changeCategory(index) {
			if (this.activeTab === index) return;
			this.activeTab = index;
			this.wines = [];
			this.page = 1;
			this.hasMore = true;
			this.loadProducts();
		},
		
		// 加载商品列表
		async loadProducts() {
			if (this.loading || !this.hasMore) return;
			
			this.loading = true;
			try {
				const categoryId = this.activeTab === 0 ? 0 : this.categories[this.activeTab].id;
				const res = await api.getProductList(categoryId, '', this.page, this.size);
				
				if (res && res.list) {
					// 如果是第一页，直接替换数据，否则追加数据
					if (this.page === 1) {
						this.wines = res.list;
					} else {
						this.wines = [...this.wines, ...res.list];
					}
					
					// 判断是否还有更多数据
					this.hasMore = res.list.length === this.size;
					this.page++;
				}
			} catch (err) {
				console.error('加载商品失败', err);
				uni.showToast({
					title: '加载商品失败',
					icon: 'none'
				});
			} finally {
				this.loading = false;
			}
		},
		
		// 添加到购物车
		addToCart: debounce(function(item) {
			let cart = uni.getStorageSync('cart') || [];
			const idx = cart.findIndex(i => i.id === item.id);
			let changed = false;
			
			if (idx !== -1) {
				cart[idx].count += 1;
				changed = true;
			} else {
				// 转换API返回的商品格式为购物车格式
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
				uni.setStorageSync('cart', cart);
				uni.showToast({ title: `已加入购物车`, icon: 'success' });
			}
		}, 300),
	},
	// 下拉刷新
	onPullDownRefresh() {
		this.wines = [];
		this.page = 1;
		this.hasMore = true;
		this.loadProducts().then(() => {
			uni.stopPullDownRefresh();
		});
	},
	// 上拉加载更多
	onReachBottom() {
		this.loadProducts();
	}
}
</script>

<style scoped>
	.ins-container {
		min-height: 100vh;
		background: #f8f6f4;
		padding-bottom: 40rpx;
	}
	.ins-header {
		padding: 60rpx 0 10rpx 0;
		text-align: center;
	}
	.ins-title {
		font-size: 52rpx;
		font-weight: bold;
		color: #222;
		letter-spacing: 2rpx;
		font-family: 'PingFang SC', 'Helvetica Neue', Arial, sans-serif;
	}
	.ins-sub {
		font-size: 28rpx;
		color: #b8b8b8;
		margin-top: 10rpx;
		letter-spacing: 1rpx;
	}
	.ins-tabs {
		display: flex;
		justify-content: center;
		margin: 30rpx 0 40rpx 0;
		gap: 30rpx;
		flex-wrap: wrap;
		padding: 0 20rpx;
	}
	.ins-tab {
		font-size: 30rpx;
		color: #b8b8b8;
		padding: 12rpx 38rpx;
		border-radius: 32rpx;
		background: linear-gradient(90deg, #fff 60%, #f8f6f4 100%);
		transition: all 0.2s;
		box-shadow: 0 2rpx 8rpx #f3f3f3;
		font-weight: 500;
		margin-bottom: 10rpx;
	}
	.ins-tab.active {
		color: #fff;
		background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
		font-weight: bold;
		box-shadow: 0 4rpx 16rpx #e0e0e0;
		transform: scale(1.08);
		letter-spacing: 2rpx;
	}
	.ins-list {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 40rpx 32rpx;
		padding: 0 24rpx;
	}
	.ins-card {
		width: 340rpx;
		background: #fff;
		border-radius: 36rpx;
		box-shadow: 0 8rpx 32rpx #eaeaea;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-bottom: 24rpx;
		min-height: 420rpx;
		padding-bottom: 24rpx;
	}
	.ins-img {
		width: 100%;
		height: 220rpx;
		object-fit: cover;
		border-top-left-radius: 36rpx;
		border-top-right-radius: 36rpx;
		background: #f3f3f3;
	}
	.ins-info {
		width: 100%;
		padding: 20rpx 24rpx 0 24rpx;
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: flex-start;
	}
	.ins-name {
		font-size: 34rpx;
		color: #222;
		font-weight: 600;
		margin-bottom: 8rpx;
	}
	.ins-price {
		font-size: 30rpx;
		color: #f7cac9;
		font-weight: bold;
	}
	.ins-btn {
		margin-top: 18rpx;
		width: 80%;
		background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
		color: #fff;
		border: none;
		border-radius: 28rpx;
		font-size: 30rpx;
		font-weight: 500;
		box-shadow: 0 2rpx 8rpx #e0e0e0;
		padding: 16rpx 0;
		letter-spacing: 1rpx;
	}
	.loading, .empty {
		text-align: center;
		color: #b8b8b8;
		margin: 20rpx 0;
		font-size: 28rpx;
	}
</style>
