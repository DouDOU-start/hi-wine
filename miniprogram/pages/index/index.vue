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
		
		<!-- 调试信息 -->
		<view v-if="debugInfo" class="debug-info">
			<text>分类数量: {{ categories.length }}</text>
			<text>商品数量: {{ wines.length }}</text>
			<text>接口状态: {{ debugInfo }}</text>
		</view>
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
			IMG_BASE_URL,
			debugInfo: ''
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
				console.log('分类响应数据:', JSON.stringify(res));
				if (res && res.data && res.data.list) {
					// 添加"全部"分类
					this.categories = [{ name: '全部' }, ...res.data.list];
					this.debugInfo = '分类加载成功';
				} else {
					this.debugInfo = '分类数据格式不正确: ' + JSON.stringify(res);
				}
			} catch (err) {
				console.error('加载分类失败', err);
				this.debugInfo = '分类加载失败: ' + err.message;
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
				console.log('请求商品列表，分类ID:', categoryId);
				const res = await api.getProductList(categoryId, '', this.page, this.size);
				console.log('商品列表响应:', JSON.stringify(res));
				
				if (res && res.data && res.data.list) {
					// 如果是第一页，直接替换数据，否则追加数据
					if (this.page === 1) {
						this.wines = res.data.list;
					} else {
						this.wines = [...this.wines, ...res.data.list];
					}
					
					// 判断是否还有更多数据
					this.hasMore = res.data.list.length === this.size;
					this.page++;
					this.debugInfo = '商品加载成功，共' + this.wines.length + '条';
				} else {
					this.debugInfo = '商品数据格式不正确: ' + JSON.stringify(res);
				}
			} catch (err) {
				console.error('加载商品失败', err);
				this.debugInfo = '商品加载失败: ' + err.message;
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
		margin-bottom: 10rpx;
		line-height: 1.4;
		max-width: 100%;
		height: 96rpx;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
	}
	.ins-price {
		font-size: 32rpx;
		color: #f7cac9;
		font-weight: bold;
	}
	.ins-btn {
		width: 85%;
		height: 80rpx;
		line-height: 80rpx;
		background: linear-gradient(90deg, #f7cac9 0%, #92a8d1 100%);
		color: #fff;
		border-radius: 40rpx;
		font-size: 30rpx;
		font-weight: bold;
		margin-top: 10rpx;
		border: none;
		position: relative;
		overflow: hidden;
	}
	.ins-btn::after {
		content: '';
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		left: 0;
		background: linear-gradient(90deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.3) 50%, rgba(255,255,255,0.1) 100%);
		opacity: 0;
		transition: opacity 0.3s;
	}
	.ins-btn:active::after {
		opacity: 1;
	}
	.loading {
		text-align: center;
		padding: 30rpx 0;
		color: #b8b8b8;
		font-size: 28rpx;
	}
	.empty {
		text-align: center;
		padding: 60rpx 0;
		color: #999;
		font-size: 30rpx;
	}
	.debug-info {
		padding: 20rpx;
		background-color: rgba(0,0,0,0.05);
		margin: 20rpx;
		border-radius: 10rpx;
		font-size: 24rpx;
		color: #666;
		display: flex;
		flex-direction: column;
	}
</style>
