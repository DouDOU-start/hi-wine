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
            <span>日同比 {{ stats.daySalesGrowth > 0 ? '+' : '' }}{{ stats.daySalesGrowth }}%</span>
            <span>周同比 {{ stats.weekSalesGrowth > 0 ? '+' : '' }}{{ stats.weekSalesGrowth }}%</span>
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
            <span>昨日 {{ stats.yesterdayOrders }}</span>
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
            <span>已下架 {{ stats.inactiveProducts }}</span>
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
            <span>活跃 {{ stats.activeUsers }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span>销售趋势</span>
              <el-radio-group v-model="timeRange" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">本年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="salesChartRef" class="chart"></div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
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
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="chart-header">
              <span>分类销售占比</span>
            </div>
          </template>
          <div ref="categoryChartRef" class="chart"></div>
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
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, onUnmounted } from 'vue';
import { getDashboardStats } from '../../api/stats';
import { getProductRanking } from '../../api/stats';
import { getCategorySales } from '../../api/stats';
import { getSalesStats } from '../../api/stats';
import { getOrderList } from '../../api/order';
import * as echarts from 'echarts/core';
import { BarChart, LineChart, PieChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, TitleComponent, LegendComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';

// 注册 ECharts 组件
echarts.use([
  BarChart,
  LineChart,
  PieChart,
  GridComponent,
  TooltipComponent,
  TitleComponent,
  LegendComponent,
  CanvasRenderer
]);

// 统计数据
const stats = reactive({
  totalSales: 0,
  daySalesGrowth: 0,
  weekSalesGrowth: 0,
  totalOrders: 0,
  yesterdayOrders: 0,
  todayOrders: 0,
  totalProducts: 0,
  activeProducts: 0,
  inactiveProducts: 0,
  totalUsers: 0,
  newUsers: 0,
  activeUsers: 0
});

// 销售趋势图表
const salesChartRef = ref(null);
let salesChart = null;

// 分类销售占比图表
const categoryChartRef = ref(null);
let categoryChart = null;

// 时间范围
const timeRange = ref('week');

// 商品销量排行
const productRanking = ref([]);

// 最近订单
const recentOrders = ref([]);

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    const response = await getDashboardStats();
    Object.assign(stats, response.data);
  } catch (error) {
    console.error('获取仪表盘数据失败:', error);
  }
};

// 获取商品销量排行
const fetchProductRanking = async () => {
  try {
    const response = await getProductRanking();
    productRanking.value = response.data.list || [];
  } catch (error) {
    console.error('获取商品销量排行失败:', error);
  }
};

// 获取分类销售占比
const fetchCategorySales = async () => {
  try {
    const response = await getCategorySales();
    const data = response.data;
    
    // 初始化分类销售占比图表
    if (categoryChart) {
      categoryChart.setOption({
        series: [{
          data: data.map(item => ({
            name: item.name,
            value: item.sales
          }))
        }]
      });
    }
  } catch (error) {
    console.error('获取分类销售占比失败:', error);
  }
};

// 获取销售趋势
const fetchSalesTrend = async () => {
  try {
    const response = await getSalesStats({ timeRange: timeRange.value });
    const data = response.data;
    
    // 初始化销售趋势图表
    if (salesChart) {
      salesChart.setOption({
        xAxis: {
          data: data.dates
        },
        series: [
          {
            name: '销售额',
            data: data.sales
          },
          {
            name: '订单数',
            data: data.orders
          }
        ]
      });
    }
  } catch (error) {
    console.error('获取销售趋势失败:', error);
  }
};

// 获取最近订单
const fetchRecentOrders = async () => {
  try {
    const response = await getOrderList({ page: 1, size: 5 });
    recentOrders.value = response.data.list || [];
  } catch (error) {
    console.error('获取最近订单失败:', error);
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  switch (status) {
    case 0: return '待支付';
    case 1: return '已支付';
    case 2: return '已完成';
    case 3: return '已取消';
    default: return '未知状态';
  }
};

// 获取订单状态类型
const getOrderStatusType = (status) => {
  switch (status) {
    case 0: return 'warning';
    case 1: return 'success';
    case 2: return 'primary';
    case 3: return 'info';
    default: return 'info';
  }
};

// 初始化销售趋势图表
const initSalesChart = () => {
  if (salesChartRef.value) {
    salesChart = echarts.init(salesChartRef.value);
    
    salesChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        data: ['销售额', '订单数']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: []
      },
      yAxis: [
        {
          type: 'value',
          name: '销售额',
          axisLabel: {
            formatter: '￥{value}'
          }
        },
        {
          type: 'value',
          name: '订单数',
          axisLabel: {
            formatter: '{value}'
          }
        }
      ],
      series: [
        {
          name: '销售额',
          type: 'line',
          data: []
        },
        {
          name: '订单数',
          type: 'bar',
          yAxisIndex: 1,
          data: []
        }
      ]
    });
  }
};

// 初始化分类销售占比图表
const initCategoryChart = () => {
  if (categoryChartRef.value) {
    categoryChart = echarts.init(categoryChartRef.value);
    
    categoryChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 10,
        data: []
      },
      series: [
        {
          name: '销售占比',
          type: 'pie',
          radius: ['50%', '70%'],
          avoidLabelOverlap: false,
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '14',
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: []
        }
      ]
    });
  }
};

// 监听窗口大小变化，重绘图表
const handleResize = () => {
  if (salesChart) {
    salesChart.resize();
  }
  if (categoryChart) {
    categoryChart.resize();
  }
};

// 监听时间范围变化
watch(timeRange, () => {
  fetchSalesTrend();
});

// 页面挂载时执行
onMounted(async () => {
  // 获取数据
  await fetchDashboardData();
  await fetchProductRanking();
  await fetchRecentOrders();
  
  // 初始化图表
  initSalesChart();
  initCategoryChart();
  
  // 获取图表数据
  fetchSalesTrend();
  fetchCategorySales();
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize);
});

// 页面卸载时执行
onUnmounted(() => {
  // 移除事件监听
  window.removeEventListener('resize', handleResize);
  
  // 销毁图表实例
  if (salesChart) {
    salesChart.dispose();
    salesChart = null;
  }
  if (categoryChart) {
    categoryChart.dispose();
    categoryChart = null;
  }
});
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.data-overview {
  margin-bottom: 20px;
}

.stats-card {
  height: 100%;
  cursor: pointer;
  transition: all 0.3s;
}

.stats-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
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
  color: #303133;
  margin-bottom: 16px;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  color: #909399;
  font-size: 13px;
}

.chart-card {
  margin-bottom: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-more {
  font-size: 14px;
  color: #409EFF;
  text-decoration: none;
}

.chart {
  height: 300px;
}

.ranking-list {
  height: 300px;
  overflow-y: auto;
}

.ranking-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.ranking-index {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #f5f7fa;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-size: 14px;
  font-weight: bold;
}

.ranking-index.top-three {
  background-color: #409EFF;
  color: #fff;
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
}

.ranking-data {
  font-size: 14px;
  color: #606266;
  font-weight: bold;
}
</style> 