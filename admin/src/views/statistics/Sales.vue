<template>
  <div class="sales-stats-container">
    <div class="page-header">
      <h2>销售统计</h2>
      <div class="date-picker">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :shortcuts="dateShortcuts"
          @change="handleDateChange"
        />
      </div>
    </div>
    
    <div class="stats-cards">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">总销售额</div>
            <div class="stats-value">¥{{ formatNumber(overviewData.totalSales) }}</div>
            <div class="stats-compare" :class="overviewData.salesGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.salesGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.salesGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">订单数</div>
            <div class="stats-value">{{ formatNumber(overviewData.orderCount) }}</div>
            <div class="stats-compare" :class="overviewData.orderGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.orderGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.orderGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">客单价</div>
            <div class="stats-value">¥{{ formatNumber(overviewData.averageOrderValue) }}</div>
            <div class="stats-compare" :class="overviewData.aovGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.aovGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.aovGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">套餐销售额</div>
            <div class="stats-value">¥{{ formatNumber(overviewData.packageSales) }}</div>
            <div class="stats-compare" :class="overviewData.packageGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.packageGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.packageGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <el-row :gutter="20" class="charts-container">
      <el-col :span="16">
        <el-card class="chart-card">
          <div class="chart-header">
            <div class="chart-title">销售趋势</div>
            <el-radio-group v-model="trendType" size="small" @change="handleTrendTypeChange">
              <el-radio-button label="day">日</el-radio-button>
              <el-radio-button label="week">周</el-radio-button>
              <el-radio-button label="month">月</el-radio-button>
            </el-radio-group>
          </div>
          <div class="chart-content" ref="salesTrendChart"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="chart-card">
          <div class="chart-title">销售分类占比</div>
          <div class="chart-content" ref="salesCategoryChart"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="hot-products-card">
      <div class="card-header">
        <div class="card-title">热销商品</div>
        <el-select v-model="hotProductsType" size="small" @change="getHotProducts">
          <el-option label="按销量" value="quantity" />
          <el-option label="按销售额" value="amount" />
        </el-select>
      </div>
      <el-table
        v-loading="hotProductsLoading"
        :data="hotProducts"
        border
        style="width: 100%"
      >
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="商品" min-width="200">
          <template #default="scope">
            <div class="product-info">
              <el-image
                style="width: 40px; height: 40px"
                :src="scope.row.imageUrl"
                fit="cover"
              />
              <span class="product-name">{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="categoryName" label="分类" width="120" />
        <el-table-column prop="price" label="单价" width="120">
          <template #default="scope">
            ¥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="salesQuantity" label="销量" width="120" />
        <el-table-column prop="salesAmount" label="销售额" width="150">
          <template #default="scope">
            ¥{{ scope.row.salesAmount.toFixed(2) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { getSalesStats, getSalesTrend, getHotProducts as getHotProductsApi } from '../../api/stats';
import * as echarts from 'echarts';

// 日期范围
const dateRange = ref([
  new Date(new Date().getTime() - 30 * 24 * 60 * 60 * 1000), // 30天前
  new Date() // 今天
]);

// 日期快捷选项
const dateShortcuts = [
  {
    text: '最近一周',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 7 * 24 * 60 * 60 * 1000);
      return [start, end];
    }
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 30 * 24 * 60 * 60 * 1000);
      return [start, end];
    }
  },
  {
    text: '最近三个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 90 * 24 * 60 * 60 * 1000);
      return [start, end];
    }
  }
];

// 概览数据
const overviewData = reactive({
  totalSales: 0,
  salesGrowth: 0,
  orderCount: 0,
  orderGrowth: 0,
  averageOrderValue: 0,
  aovGrowth: 0,
  packageSales: 0,
  packageGrowth: 0
});

// 趋势图表
const salesTrendChart = ref(null);
const trendType = ref('day');
let trendChartInstance = null;

// 分类占比图表
const salesCategoryChart = ref(null);
let categoryChartInstance = null;

// 热销商品
const hotProducts = ref([]);
const hotProductsLoading = ref(false);
const hotProductsType = ref('quantity');

// 格式化数字
const formatNumber = (num) => {
  if (num === undefined || num === null) return '0';
  return num.toLocaleString('zh-CN', { maximumFractionDigits: 2 });
};

// 处理日期变化
const handleDateChange = () => {
  getSalesData();
  getSalesTrendData();
  getHotProducts();
};

// 处理趋势类型变化
const handleTrendTypeChange = () => {
  getSalesTrendData();
};

// 获取销售统计数据
const getSalesData = async () => {
  try {
    const params = {
      startDate: dateRange.value[0].toISOString().split('T')[0],
      endDate: dateRange.value[1].toISOString().split('T')[0]
    };
    
    const res = await getSalesStats(params);
    const data = res.data;
    
    overviewData.totalSales = data.totalSales || 0;
    overviewData.salesGrowth = data.salesGrowth || 0;
    overviewData.orderCount = data.orderCount || 0;
    overviewData.orderGrowth = data.orderGrowth || 0;
    overviewData.averageOrderValue = data.averageOrderValue || 0;
    overviewData.aovGrowth = data.aovGrowth || 0;
    overviewData.packageSales = data.packageSales || 0;
    overviewData.packageGrowth = data.packageGrowth || 0;
    
    // 初始化分类占比图表
    if (data.categorySales) {
      initCategoryChart(data.categorySales);
    }
  } catch (error) {
    console.error('获取销售统计数据失败:', error);
  }
};

// 获取销售趋势数据
const getSalesTrendData = async () => {
  try {
    const params = {
      startDate: dateRange.value[0].toISOString().split('T')[0],
      endDate: dateRange.value[1].toISOString().split('T')[0],
      type: trendType.value
    };
    
    const res = await getSalesTrend(params);
    initTrendChart(res.data);
  } catch (error) {
    console.error('获取销售趋势数据失败:', error);
  }
};

// 获取热销商品
const getHotProducts = async () => {
  hotProductsLoading.value = true;
  try {
    const params = {
      startDate: dateRange.value[0].toISOString().split('T')[0],
      endDate: dateRange.value[1].toISOString().split('T')[0],
      type: hotProductsType.value,
      limit: 10
    };
    
    const res = await getHotProductsApi(params);
    hotProducts.value = res.data.list || [];
  } catch (error) {
    console.error('获取热销商品失败:', error);
  } finally {
    hotProductsLoading.value = false;
  }
};

// 初始化趋势图表
const initTrendChart = (data) => {
  if (!salesTrendChart.value) return;
  
  if (!trendChartInstance) {
    trendChartInstance = echarts.init(salesTrendChart.value);
  }
  
  const option = {
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
      data: data.dates || []
    },
    yAxis: [
      {
        type: 'value',
        name: '销售额',
        axisLabel: {
          formatter: '{value} 元'
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
        type: 'bar',
        data: data.sales || []
      },
      {
        name: '订单数',
        type: 'line',
        yAxisIndex: 1,
        data: data.orders || []
      }
    ]
  };
  
  trendChartInstance.setOption(option);
};

// 初始化分类占比图表
const initCategoryChart = (data) => {
  if (!salesCategoryChart.value) return;
  
  if (!categoryChartInstance) {
    categoryChartInstance = echarts.init(salesCategoryChart.value);
  }
  
  const categories = data.map(item => item.name);
  const values = data.map(item => item.value);
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10,
      data: categories
    },
    series: [
      {
        name: '销售额',
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
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data
      }
    ]
  };
  
  categoryChartInstance.setOption(option);
};

// 窗口大小变化时重新调整图表大小
const handleResize = () => {
  if (trendChartInstance) {
    trendChartInstance.resize();
  }
  if (categoryChartInstance) {
    categoryChartInstance.resize();
  }
};

// 初始化
onMounted(() => {
  getSalesData();
  getSalesTrendData();
  getHotProducts();
  
  window.addEventListener('resize', handleResize);
});

// 清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  
  if (trendChartInstance) {
    trendChartInstance.dispose();
    trendChartInstance = null;
  }
  
  if (categoryChartInstance) {
    categoryChartInstance.dispose();
    categoryChartInstance = null;
  }
});
</script>

<style scoped>
.sales-stats-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.stats-cards {
  margin-bottom: 20px;
}

.stats-card {
  height: 120px;
}

.stats-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 10px;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 10px;
}

.stats-compare {
  font-size: 12px;
  display: flex;
  align-items: center;
}

.stats-compare.positive {
  color: #67C23A;
}

.stats-compare.negative {
  color: #F56C6C;
}

.charts-container {
  margin-bottom: 20px;
}

.chart-card {
  margin-bottom: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.chart-title {
  font-size: 16px;
  font-weight: bold;
}

.chart-content {
  height: 300px;
}

.hot-products-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: bold;
}

.product-info {
  display: flex;
  align-items: center;
}

.product-name {
  margin-left: 10px;
}
</style> 