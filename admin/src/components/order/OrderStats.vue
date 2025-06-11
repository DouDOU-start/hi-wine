<template>
  <div class="order-stats-container">
    <el-card shadow="hover" class="stats-card">
      <template #header>
        <div class="card-header">
          <span>订单统计</span>
          <el-button v-if="refreshable" type="primary" link @click="$emit('refresh')">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </template>
      
      <div class="stats-overview">
        <div class="stat-item total">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">总订单</div>
          <el-icon class="stat-icon"><Document /></el-icon>
        </div>
        <div class="stat-item completed">
          <div class="stat-value">{{ stats.completed }}</div>
          <div class="stat-label">已完成</div>
          <el-icon class="stat-icon"><CircleCheck /></el-icon>
        </div>
        <div class="stat-item pending">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">待支付</div>
          <el-icon class="stat-icon"><Clock /></el-icon>
        </div>
        <div class="stat-item cancelled">
          <div class="stat-value">{{ stats.cancelled }}</div>
          <div class="stat-label">已取消</div>
          <el-icon class="stat-icon"><CircleClose /></el-icon>
        </div>
      </div>
      
      <div class="chart-container" v-if="showChart">
        <div ref="pieChartRef" class="pie-chart"></div>
      </div>
      
      <div class="stats-details" v-if="showDetails">
        <el-divider content-position="center">详细数据</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="detail-item">
              <span class="detail-label">今日订单:</span>
              <span class="detail-value">{{ stats.todayOrders || 0 }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="detail-item">
              <span class="detail-label">今日销售额:</span>
              <span class="detail-value price">¥{{ formatPrice(stats.todaySales) }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="detail-item">
              <span class="detail-label">本月订单:</span>
              <span class="detail-value">{{ stats.monthOrders || 0 }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="detail-item">
              <span class="detail-label">本月销售额:</span>
              <span class="detail-value price">¥{{ formatPrice(stats.monthSales) }}</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, onMounted, watch } from 'vue';
import { Document, CircleCheck, Clock, CircleClose, Refresh } from '@element-plus/icons-vue';
import * as echarts from 'echarts/core';
import { PieChart } from 'echarts/charts';
import { 
  TitleComponent, 
  TooltipComponent, 
  LegendComponent 
} from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';

// 注册必要的组件
echarts.use([
  TitleComponent, 
  TooltipComponent, 
  LegendComponent, 
  PieChart, 
  CanvasRenderer
]);

const props = defineProps({
  stats: {
    type: Object,
    required: true,
    default: () => ({
      total: 0,
      completed: 0,
      pending: 0,
      cancelled: 0,
      todayOrders: 0,
      todaySales: 0,
      monthOrders: 0,
      monthSales: 0
    })
  },
  showChart: {
    type: Boolean,
    default: true
  },
  showDetails: {
    type: Boolean,
    default: true
  },
  refreshable: {
    type: Boolean,
    default: true
  }
});

defineEmits(['refresh']);

const pieChartRef = ref(null);
let pieChart = null;

// 格式化价格
const formatPrice = (price) => {
  if (price === undefined || price === null) return '0.00';
  return parseFloat(price).toFixed(2);
};

// 初始化饼图
const initPieChart = () => {
  if (!pieChartRef.value || !props.showChart) return;
  
  // 创建图表实例
  pieChart = echarts.init(pieChartRef.value);
  
  // 更新图表
  updatePieChart();
  
  // 监听窗口大小变化，调整图表大小
  window.addEventListener('resize', () => {
    pieChart && pieChart.resize();
  });
};

// 更新饼图数据
const updatePieChart = () => {
  if (!pieChart) return;
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      data: ['已完成', '待支付', '已取消']
    },
    series: [
      {
        name: '订单状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
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
        data: [
          { 
            value: props.stats.completed, 
            name: '已完成',
            itemStyle: { color: '#409EFF' }
          },
          { 
            value: props.stats.pending, 
            name: '待支付',
            itemStyle: { color: '#E6A23C' }
          },
          { 
            value: props.stats.cancelled, 
            name: '已取消',
            itemStyle: { color: '#909399' }
          }
        ]
      }
    ]
  };
  
  pieChart.setOption(option);
};

// 监听数据变化，更新图表
watch(() => props.stats, () => {
  if (pieChart) {
    updatePieChart();
  }
}, { deep: true });

// 组件挂载后初始化图表
onMounted(() => {
  initPieChart();
});
</script>

<style scoped>
.order-stats-container {
  margin-bottom: 20px;
}

.stats-card {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-overview {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 10px;
}

.stat-item {
  flex: 1;
  min-width: 120px;
  padding: 15px;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 5px;
  position: relative;
  z-index: 1;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  position: relative;
  z-index: 1;
}

.stat-icon {
  position: absolute;
  right: 10px;
  bottom: 10px;
  font-size: 40px;
  opacity: 0.2;
}

.stat-item.total {
  background-color: #f0f5ff;
  color: #409EFF;
}

.stat-item.completed {
  background-color: #f0f9eb;
  color: #67C23A;
}

.stat-item.pending {
  background-color: #fdf6ec;
  color: #E6A23C;
}

.stat-item.cancelled {
  background-color: #f4f4f5;
  color: #909399;
}

.chart-container {
  margin-top: 20px;
}

.pie-chart {
  height: 300px;
  width: 100%;
}

.stats-details {
  margin-top: 20px;
}

.detail-item {
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
}

.detail-label {
  color: #606266;
}

.detail-value {
  font-weight: 500;
}

.detail-value.price {
  color: #f56c6c;
}

@media (max-width: 768px) {
  .stats-overview {
    flex-direction: column;
  }
  
  .stat-item {
    margin-bottom: 10px;
  }
}
</style> 
