<template>
  <div class="dashboard-container">
    <el-row :gutter="20" class="data-overview">
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-header">
            <div class="card-title">总销售额</div>
            <el-icon :size="24" color="#409EFF"><Money /></el-icon>
          </div>
          <div class="card-data">￥{{ stats.totalSales.toLocaleString() }}</div>
          <div class="card-footer">
            <span>今日 ￥{{ stats.todaySales.toLocaleString() }}</span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-header">
            <div class="card-title">订单数量</div>
            <el-icon :size="24" color="#67C23A"><Tickets /></el-icon>
          </div>
          <div class="card-data">{{ stats.totalOrders }}</div>
          <div class="card-footer">
            <span>今日 {{ stats.todayOrders }}</span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-header">
            <div class="card-title">商品数量</div>
            <el-icon :size="24" color="#E6A23C"><Goods /></el-icon>
          </div>
          <div class="card-data">{{ stats.totalProducts }}</div>
          <div class="card-footer">
            <span>上架中 {{ stats.activeProducts }}</span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6" :lg="6" :xl="6">
        <el-card class="stats-card" shadow="hover">
          <div class="card-header">
            <div class="card-title">用户数量</div>
            <el-icon :size="24" color="#F56C6C"><User /></el-icon>
          </div>
          <div class="card-data">{{ stats.totalUsers }}</div>
          <div class="card-footer">
            <span>新增 {{ stats.newUsers }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span>商品销量排行</span>
            </div>
          </template>
          <div class="ranking-list">
            <div v-for="(item, index) in productRanking" :key="item.id" class="ranking-item">
              <div class="ranking-index" :class="{ 'top-three': index < 3 }">{{ index + 1 }}</div>
              <div class="ranking-info">
                <div class="ranking-name">{{ item.name }}</div>
                <div class="ranking-data">{{ item.sales }} 件</div>
              </div>
            </div>
            <div v-if="productRanking.length === 0" class="empty-data">
              暂无数据
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span>最近订单</span>
              <router-link to="/order/list" class="view-more">查看更多</router-link>
            </div>
          </template>
          <el-table :data="recentOrders" stripe style="width: 100%">
            <el-table-column prop="id" label="订单号" width="100" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="totalAmount" label="金额">
              <template #default="scope">
                ￥{{ scope.row.totalAmount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getOrderStatusType(scope.row.status)">
                  {{ getOrderStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" />
          </el-table>
          <div v-if="recentOrders.length === 0" class="empty-data">
            暂无数据
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getDashboardStats, getProductRanking } from '../../api/stats';
import { getOrderList } from '../../api/order';

// 统计数据
const stats = reactive({
  totalSales: 0,
  todaySales: 0,
  totalOrders: 0,
  todayOrders: 0,
  totalProducts: 0,
  activeProducts: 0,
  totalUsers: 0,
  newUsers: 0
});

// 商品销量排行
const productRanking = ref([]);

// 最近订单
const recentOrders = ref([]);

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    const response = await getDashboardStats();
    if (response.data) {
      Object.assign(stats, response.data);
    }
  } catch (error) {
    console.error('获取仪表盘数据失败:', error);
  }
};

// 获取商品销量排行
const fetchProductRanking = async () => {
  try {
    const response = await getProductRanking({ limit: 10 });
    if (response.data && response.data.list) {
      productRanking.value = response.data.list;
    }
  } catch (error) {
    console.error('获取商品销量排行失败:', error);
  }
};

// 获取最近订单
const fetchRecentOrders = async () => {
  try {
    const response = await getOrderList({ 
      page: 1, 
      limit: 5,
      sort: 'create_time',
      order: 'desc'
    });
    if (response.data && response.data.list) {
      recentOrders.value = response.data.list;
    }
  } catch (error) {
    console.error('获取最近订单失败:', error);
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  const statusMap = {
    'pending': '待支付',
    'paid': '已支付',
    'completed': '已完成',
    'cancelled': '已取消'
  };
  return statusMap[status] || '未知状态';
};

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const typeMap = {
    'pending': 'warning',
    'paid': 'success',
    'completed': 'success',
    'cancelled': 'danger'
  };
  return typeMap[status] || 'info';
};

// 初始化
onMounted(() => {
  fetchDashboardData();
  fetchProductRanking();
  fetchRecentOrders();
});
</script>

<style scoped>
.dashboard-container {
  padding: 20px 0;
}

.data-overview {
  margin-bottom: 20px;
}

.stats-card {
  height: 100%;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.card-title {
  font-size: 16px;
  color: #606266;
}

.card-data {
  font-size: 28px;
  font-weight: bold;
  margin: 10px 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 13px;
}

.chart-card {
  margin-bottom: 20px;
  height: 100%;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-more {
  font-size: 12px;
  color: #409EFF;
  text-decoration: none;
}

.chart {
  height: 300px;
}

.ranking-list {
  padding: 0 10px;
}

.ranking-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #EBEEF5;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-index {
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  border-radius: 50%;
  background-color: #F2F6FC;
  margin-right: 12px;
  font-size: 12px;
  font-weight: bold;
  color: #909399;
}

.ranking-index.top-three {
  background-color: #409EFF;
  color: white;
}

.ranking-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.ranking-name {
  font-size: 14px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 70%;
}

.ranking-data {
  font-size: 14px;
  color: #606266;
  font-weight: bold;
}

.empty-data {
  padding: 30px 0;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style> 