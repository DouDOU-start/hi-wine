<template>
  <div class="package-stats-container">
    <div class="page-header">
      <h2>套餐统计</h2>
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
            <div class="stats-title">套餐销售总额</div>
            <div class="stats-value">¥{{ formatNumber(overviewData.totalSales) }}</div>
            <div class="stats-compare" :class="overviewData.salesGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.salesGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.salesGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">套餐销售数量</div>
            <div class="stats-value">{{ formatNumber(overviewData.totalCount) }}</div>
            <div class="stats-compare" :class="overviewData.countGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.countGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.countGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">活跃套餐数</div>
            <div class="stats-value">{{ formatNumber(overviewData.activeCount) }}</div>
            <div class="stats-compare" :class="overviewData.activeGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.activeGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.activeGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <div class="stats-title">套餐使用次数</div>
            <div class="stats-value">{{ formatNumber(overviewData.usageCount) }}</div>
            <div class="stats-compare" :class="overviewData.usageGrowth >= 0 ? 'positive' : 'negative'">
              <el-icon><component :is="overviewData.usageGrowth >= 0 ? 'ArrowUp' : 'ArrowDown'" /></el-icon>
              {{ Math.abs(overviewData.usageGrowth) }}% 较上期
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    
    <el-row :gutter="20" class="charts-container">
      <el-col :span="12">
        <el-card class="chart-card">
          <div class="chart-title">套餐销售趋势</div>
          <div class="chart-content" ref="packageSalesChart"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <div class="chart-title">套餐使用趋势</div>
          <div class="chart-content" ref="packageUsageChart"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-card class="package-stats-card">
      <div class="card-header">
        <div class="card-title">套餐销售排行</div>
      </div>
      <el-table
        v-loading="packageStatsLoading"
        :data="packageStats"
        border
        style="width: 100%"
      >
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="套餐名称" min-width="150" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="scope">
            ¥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="durationMinutes" label="有效时长" width="120">
          <template #default="scope">
            {{ scope.row.durationMinutes }} 分钟
          </template>
        </el-table-column>
        <el-table-column prop="salesCount" label="销售数量" width="120" />
        <el-table-column prop="salesAmount" label="销售金额" width="150">
          <template #default="scope">
            ¥{{ scope.row.salesAmount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="activeCount" label="活跃数量" width="120" />
        <el-table-column prop="usageCount" label="使用次数" width="120" />
      </el-table>
    </el-card>
    
    <el-card class="usage-stats-card">
      <div class="card-header">
        <div class="card-title">套餐使用商品排行</div>
      </div>
      <el-table
        v-loading="usageStatsLoading"
        :data="usageStats"
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
        <el-table-column prop="usageCount" label="使用次数" width="120" />
        <el-table-column prop="usageAmount" label="使用金额" width="150">
          <template #default="scope">
            ¥{{ scope.row.usageAmount.toFixed(2) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { getPackageSalesStats, getPackageUsageStats } from '../../api/stats';
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
  totalCount: 0,
  countGrowth: 0,
  activeCount: 0,
  activeGrowth: 0,
  usageCount: 0,
  usageGrowth: 0
});

// 套餐销售图表
const packageSalesChart = ref(null);
let salesChartInstance = null;

// 套餐使用图表
const packageUsageChart = ref(null);
let usageChartInstance = null;

// 套餐统计数据
const packageStats = ref([]);
const packageStatsLoading = ref(false);

// 使用统计数据
const usageStats = ref([]);
const usageStatsLoading = ref(false);

// 格式化数字
const formatNumber = (num) => {
  if (num === undefined || num === null) return '0';
  return num.toLocaleString('zh-CN', { maximumFractionDigits: 2 });
};

// 处理日期变化
const handleDateChange = () => {
  getPackageSalesData();
  getPackageUsageData();
};

// 获取套餐销售统计数据
const getPackageSalesData = async () => {
  packageStatsLoading.value = true;
  try {
    const params = {
      startDate: dateRange.value[0].toISOString().split('T')[0],
      endDate: dateRange.value[1].toISOString().split('T')[0]
    };
    
    const res = await getPackageSalesStats(params);
    const data = res.data;
    
    // 更新概览数据
    overviewData.totalSales = data.overview.totalSales || 0;
    overviewData.salesGrowth = data.overview.salesGrowth || 0;
    overviewData.totalCount = data.overview.totalCount || 0;
    overviewData.countGrowth = data.overview.countGrowth || 0;
    overviewData.activeCount = data.overview.activeCount || 0;
    overviewData.activeGrowth = data.overview.activeGrowth || 0;
    
    // 更新套餐统计数据
    packageStats.value = data.packages || [];
    
    // 初始化销售趋势图表
    if (data.trend) {
      initSalesChart(data.trend);
    }
  } catch (error) {
    console.error('获取套餐销售统计数据失败:', error);
  } finally {
    packageStatsLoading.value = false;
  }
};

// 获取套餐使用统计数据
const getPackageUsageData = async () => {
  usageStatsLoading.value = true;
  try {
    const params = {
      startDate: dateRange.value[0].toISOString().split('T')[0],
      endDate: dateRange.value[1].toISOString().split('T')[0]
    };
    
    const res = await getPackageUsageStats(params);
    const data = res.data;
    
    // 更新概览数据
    overviewData.usageCount = data.overview.usageCount || 0;
    overviewData.usageGrowth = data.overview.usageGrowth || 0;
    
    // 更新使用统计数据
    usageStats.value = data.products || [];
    
    // 初始化使用趋势图表
    if (data.trend) {
      initUsageChart(data.trend);
    }
  } catch (error) {
    console.error('获取套餐使用统计数据失败:', error);
  } finally {
    usageStatsLoading.value = false;
  }
};

// 初始化销售趋势图表
const initSalesChart = (data) => {
  if (!packageSalesChart.value) return;
  
  if (!salesChartInstance) {
    salesChartInstance = echarts.init(packageSalesChart.value);
  }
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['销售数量', '销售金额']
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
        name: '销售数量',
        axisLabel: {
          formatter: '{value}'
        }
      },
      {
        type: 'value',
        name: '销售金额',
        axisLabel: {
          formatter: '{value} 元'
        }
      }
    ],
    series: [
      {
        name: '销售数量',
        type: 'bar',
        data: data.counts || []
      },
      {
        name: '销售金额',
        type: 'line',
        yAxisIndex: 1,
        data: data.amounts || []
      }
    ]
  };
  
  salesChartInstance.setOption(option);
};

// 初始化使用趋势图表
const initUsageChart = (data) => {
  if (!packageUsageChart.value) return;
  
  if (!usageChartInstance) {
    usageChartInstance = echarts.init(packageUsageChart.value);
  }
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: data.dates || []
    },
    yAxis: {
      type: 'value',
      name: '使用次数'
    },
    series: [
      {
        name: '使用次数',
        type: 'bar',
        data: data.counts || [],
        itemStyle: {
          color: '#91cc75'
        }
      }
    ]
  };
  
  usageChartInstance.setOption(option);
};

// 窗口大小变化时重新调整图表大小
const handleResize = () => {
  if (salesChartInstance) {
    salesChartInstance.resize();
  }
  if (usageChartInstance) {
    usageChartInstance.resize();
  }
};

// 初始化
onMounted(() => {
  getPackageSalesData();
  getPackageUsageData();
  
  window.addEventListener('resize', handleResize);
});

// 清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  
  if (salesChartInstance) {
    salesChartInstance.dispose();
    salesChartInstance = null;
  }
  
  if (usageChartInstance) {
    usageChartInstance.dispose();
    usageChartInstance = null;
  }
});
</script>

<style scoped>
.package-stats-container {
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

.chart-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 20px;
}

.chart-content {
  height: 300px;
}

.package-stats-card,
.usage-stats-card {
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