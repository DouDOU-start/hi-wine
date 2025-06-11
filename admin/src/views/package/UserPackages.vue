<template>
  <div class="user-package-container">
    <div class="page-header">
      <h2>用户套餐管理</h2>
    </div>
    
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="filter-form">
        <el-form-item label="用户昵称">
          <el-input v-model="queryParams.nickname" placeholder="请输入用户昵称" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="套餐名称">
          <el-input v-model="queryParams.packageName" placeholder="请输入套餐名称" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="queryParams.status" 
            placeholder="全部状态" 
            clearable
          >
            <el-option 
              v-for="(value, key) in STATUS_MAP" 
              :key="key" 
              :label="value.text" 
              :value="key" 
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card class="list-container">
      <el-table
        v-loading="loading"
        :data="userPackageList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="180">
          <template #default="scope">
            <div class="user-info">
              <el-avatar :size="40" :src="scope.row.avatarUrl || ''" />
              <div class="user-detail">
                <div>{{ scope.row.userName || '未知用户' }}</div>
                <div class="user-phone">{{ scope.row.userPhone || '未绑定手机' }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="packageName" label="套餐名称" min-width="150" />
        <el-table-column label="有效期" width="280">
          <template #default="scope">
            <div>开始: {{ formatDate(scope.row.startTime) }}</div>
            <div>结束: {{ formatDate(scope.row.endTime) }}</div>
            <div v-if="scope.row.status === 'active'" class="remaining-time">
              <el-tag size="small" effect="plain" type="success">
                剩余: {{ scope.row.remainingDays || 0 }} 天
              </el-tag>
            </div>
            <div v-else class="status-desc">
              <el-tag size="small" effect="plain" :type="getStatusType(scope.row.status)">
                {{ scope.row.statusDesc || getStatusText(scope.row.status) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="订单信息" width="160">
          <template #default="scope">
            <div v-if="scope.row.orderSn">
              <div class="order-info">订单号: {{ scope.row.orderSn }}</div>
              <div class="order-status" v-if="scope.row.orderStatus">
                <el-tag size="small" :type="scope.row.orderStatus === 'paid' ? 'success' : 'warning'">
                  {{ scope.row.orderStatus === 'paid' ? '已支付' : scope.row.orderStatus === 'pending' ? '未支付' : scope.row.orderStatus || '未知状态' }}
                </el-tag>
              </div>
            </div>
            <div v-else>-</div>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 'active'"
              type="danger"
              size="small"
              @click="handleStatusChange(scope.row, 'expired')"
            >设为过期</el-button>
            <el-button
              v-if="scope.row.status === 'pending'"
              type="success"
              size="small"
              @click="handleStatusChange(scope.row, 'active')"
            >设为生效</el-button>
            <el-button
              type="primary"
              size="small"
              @click="handleViewDetail(scope.row)"
            >详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-size="queryParams.pageSize"
          :current-page="queryParams.pageNum"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 用户套餐详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="用户套餐详情"
      width="700px"
      destroy-on-close
    >
      <div v-if="!detailLoading && packageDetail.id" v-loading="detailLoading">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h3 class="section-title">基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="套餐ID">{{ packageDetail.id }}</el-descriptions-item>
            <el-descriptions-item label="套餐名称">{{ packageDetail.packageName }}</el-descriptions-item>
            <el-descriptions-item label="用户昵称">{{ packageDetail.user?.nickname || packageDetail.userName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="用户手机">{{ packageDetail.user?.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(packageDetail.status)">
                {{ getStatusText(packageDetail.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="有效期">
              <el-tag size="small" effect="plain" :type="getStatusType(packageDetail.status)">
                {{ packageDetail.validPeriod || '-' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ formatDate(packageDetail.startTime) }}</el-descriptions-item>
            <el-descriptions-item label="结束时间">{{ formatDate(packageDetail.endTime) }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(packageDetail.createdAt) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 订单信息 -->
        <div v-if="packageDetail.order" class="detail-section">
          <h3 class="section-title">订单信息</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="订单号">{{ packageDetail.order.orderSn || '-' }}</el-descriptions-item>
            <el-descriptions-item label="支付状态">
              <el-tag v-if="packageDetail.order.payStatus" size="small" :type="packageDetail.order.payStatus === 'paid' ? 'success' : 'warning'">
                {{ packageDetail.order.payStatus === 'paid' ? '已支付' : '未支付' }}
              </el-tag>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item label="支付时间">{{ formatDate(packageDetail.order.payTime) }}</el-descriptions-item>
            <el-descriptions-item label="订单金额">
              <span class="price-value">¥{{ packageDetail.order.totalFee || 0 }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
        
        <!-- 套餐详情 -->
        <div v-if="packageDetail.package" class="detail-section">
          <h3 class="section-title">套餐详情</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="套餐名称">{{ packageDetail.package.name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="套餐价格">
              <span class="price-value">¥{{ packageDetail.package.price || 0 }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="时长(分钟)">{{ packageDetail.package.durationMinutes || 0 }}</el-descriptions-item>
            <el-descriptions-item label="套餐描述" :span="2">{{ packageDetail.package.description || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>
        
        <!-- 使用记录 -->
        <div v-if="hasUsageRecords" class="detail-section">
          <h3 class="section-title">使用记录</h3>
          <div class="usage-stats">
            <el-tag type="info">总使用次数: {{ packageDetail.usage?.totalUsedTimes || 0 }}</el-tag>
            <el-tag type="success" v-if="packageDetail.usage?.lastUsedTime">最后使用: {{ formatDate(packageDetail.usage.lastUsedTime) }}</el-tag>
          </div>
          <el-table :data="packageDetail.usage?.records || []" border style="width: 100%; margin-top: 10px;">
            <el-table-column prop="use_time" label="使用时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.use_time) }}</template>
            </el-table-column>
            <el-table-column prop="use_duration" label="使用时长" width="120">
              <template #default="scope">{{ scope.row.use_duration || 0 }} 分钟</template>
            </el-table-column>
            <el-table-column prop="table_name" label="桌号" width="120" />
            <el-table-column prop="note" label="备注" />
          </el-table>
        </div>
        
        <!-- 操作按钮 -->
        <div class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button v-if="packageDetail.status === 'active'" type="danger" @click="handleStatusChangeInDialog('expired')">设为过期</el-button>
          <el-button v-if="packageDetail.status === 'pending'" type="success" @click="handleStatusChangeInDialog('active')">设为生效</el-button>
        </div>
      </div>
      <div v-else v-loading="detailLoading" style="min-height: 200px;"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getUserPackages, getUserPackageDetail, updateUserPackageStatus } from '../../api/package';
import { formatDate as formatDateTime } from '../../utils/format';

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  nickname: '',
  packageName: '',
  status: ''
});

// 列表数据
const userPackageList = ref([]);
const total = ref(0);
const loading = ref(false);

// 详情相关
const detailDialogVisible = ref(false);
const detailLoading = ref(false);
const packageDetail = reactive({});

// 状态映射表
const STATUS_MAP = {
  'active': { text: '生效中', type: 'success' },
  'expired': { text: '已过期', type: 'info' },
  'pending': { text: '待支付', type: 'warning' },
  'refunded': { text: '已退款', type: 'danger' }
};

// 日期格式化
const formatDate = (timestamp) => {
  if (!timestamp) return '-';
  return formatDateTime(timestamp);
};

// 获取状态文本
const getStatusText = (status) => {
  return STATUS_MAP[status]?.text || '未知状态';
};

// 获取状态类型
const getStatusType = (status) => {
  return STATUS_MAP[status]?.type || 'info';
};

// 获取用户套餐列表
const getList = async () => {
  loading.value = true;
  try {
    const response = await getUserPackages(queryParams);
    console.log('API原始响应:', response);
    
    // 确保返回的数据是数组，并对每个项目进行处理
    let list = [];
    
    // 检查响应结构，兼容不同的数据格式
    if (response.data?.list) {
      // 标准格式：{data: {list: [], total: 0}}
      list = response.data.list;
      total.value = response.data.total || 0;
    } else if (Array.isArray(response.data)) {
      // 数组格式：{data: []}
      list = response.data;
      total.value = list.length;
    } else if (response.list) {
      // 直接包含list的格式：{list: [], total: 0}
      list = response.list;
      total.value = response.total || 0;
    } else {
      console.error('未知的数据格式:', response);
      list = [];
      total.value = 0;
    }
    
    console.log('识别的列表数据:', list);
    
    // 确保list是数组
    if (!Array.isArray(list)) {
      console.error('列表数据不是数组:', list);
      list = [];
    }
    
    // 处理每一项
    userPackageList.value = list.map(item => {
      // 直接返回原始数据，假设API已经处理好
      return item;
    });
    
  } catch (error) {
    ElMessage.error('获取用户套餐列表失败');
    console.error('获取用户套餐列表错误详情:', error);
    userPackageList.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
};

// 查询
const handleQuery = () => {
  queryParams.pageNum = 1;
  getList();
};

// 重置
const resetQuery = () => {
  Object.assign(queryParams, {
    nickname: '',
    packageName: '',
    status: '',
    pageNum: 1
  });
  getList();
};

// 分页改变
const handleSizeChange = (size) => {
  queryParams.pageSize = size;
  getList();
};

const handleCurrentChange = (page) => {
  queryParams.pageNum = page;
  getList();
};

// 状态更改通用处理函数
const changeStatus = async (id, status) => {
  try {
    await updateUserPackageStatus(id, { status });
    ElMessage.success(`状态已更新为${getStatusText(status)}`);
    getList();
    return true;
  } catch (error) {
    ElMessage.error('更新状态失败');
    console.error('更新状态错误详情:', error);
    return false;
  }
};

// 修改状态
const handleStatusChange = (row, status) => {
  const statusText = getStatusText(status);
  ElMessageBox.confirm(`确认要将该用户套餐状态修改为"${statusText}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => changeStatus(row.id, status))
    .catch(() => {});
};

// 详情对话框中的状态更改
const handleStatusChangeInDialog = (status) => {
  const statusText = getStatusText(status);
  ElMessageBox.confirm(`确认要将该用户套餐状态修改为"${statusText}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    if (await changeStatus(packageDetail.id, status)) {
      // 刷新详情
      handleViewDetail({ id: packageDetail.id });
    }
  }).catch(() => {});
};

// 查看详情
const handleViewDetail = async (row) => {
  detailDialogVisible.value = true;
  detailLoading.value = true;
  // 清空旧数据以确保响应性
  Object.keys(packageDetail).forEach(key => delete packageDetail[key]);
  
  try {
    const response = await getUserPackageDetail(row.id);
    console.log('详情API响应:', response);
    if (response.data) {
      // 将获取的数据赋值给响应式对象
      Object.assign(packageDetail, response.data);
    }
  } catch (error) {
    ElMessage.error('获取套餐详情失败');
    console.error('获取套餐详情错误:', error);
  } finally {
    detailLoading.value = false;
  }
};

// 判断是否有使用记录
const hasUsageRecords = computed(() => {
  const usage = packageDetail?.usage;
  if (!usage) return false;
  
  return usage.totalUsedTimes > 0 || 
         (usage.records && usage.records.length > 0) || 
         usage.lastUsedTime;
});

// 初始化
onMounted(() => {
  getList();
});
</script>

<style scoped>
.user-package-container {
  padding: 0 20px 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  font-weight: 500;
  color: #303133;
}

.filter-container {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.list-container {
  margin-bottom: 20px;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-detail {
  margin-left: 10px;
}

.user-phone {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.remaining-time, .status-desc {
  margin-top: 5px;
}

.order-info {
  font-size: 13px;
  color: #606266;
}

.order-status {
  margin-top: 5px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.detail-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  color: #303133;
  margin: 15px 0;
  padding-left: 10px;
  border-left: 3px solid #409EFF;
  font-weight: 500;
}

.usage-stats {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.price-value {
  color: #f56c6c;
  font-weight: bold;
}

.dialog-footer {
  margin-top: 20px;
  text-align: right;
}
</style> 