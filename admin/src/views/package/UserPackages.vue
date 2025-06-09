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
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
            <el-option label="生效中" value="active" />
            <el-option label="已过期" value="expired" />
            <el-option label="待支付" value="pending" />
            <el-option label="已退款" value="refunded" />
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
              <el-avatar :size="40" :src="scope.row.user.avatarUrl" />
              <div class="user-detail">
                <div>{{ scope.row.user.nickname }}</div>
                <div class="user-phone">{{ scope.row.user.phone || '未绑定手机' }}</div>
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
                剩余: {{ scope.row.remainingMinutes }} 分钟
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
        <el-table-column label="订单信息" width="150">
          <template #default="scope">
            <div v-if="scope.row.order && scope.row.order.sn">
              <div class="order-info">订单号: {{ scope.row.order.sn }}</div>
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
      <el-descriptions v-loading="detailLoading" :column="2" border>
        <el-descriptions-item label="套餐ID">{{ packageDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="套餐名称">{{ packageDetail.packageName }}</el-descriptions-item>
        <el-descriptions-item label="用户昵称">{{ packageDetail.user?.nickname }}</el-descriptions-item>
        <el-descriptions-item label="用户手机">{{ packageDetail.user?.phone || '未绑定手机' }}</el-descriptions-item>
        <el-descriptions-item label="套餐价格">
          <span class="price-value">¥{{ packageDetail.package?.price || 0 }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="有效期">
          <el-tag size="small" effect="plain" :type="getStatusType(packageDetail.status)">
            {{ packageDetail.validPeriod || '未知' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ formatDate(packageDetail.startTime) }}</el-descriptions-item>
        <el-descriptions-item label="结束时间">{{ formatDate(packageDetail.endTime) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(packageDetail.status)">
            {{ getStatusText(packageDetail.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态描述">{{ packageDetail.statusDesc || '-' }}</el-descriptions-item>
        <el-descriptions-item label="订单号">
          <el-link v-if="packageDetail.order?.sn" type="primary" :underline="false">
            {{ packageDetail.order?.sn || '-' }}
          </el-link>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(packageDetail.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatDate(packageDetail.updatedAt) }}</el-descriptions-item>
      </el-descriptions>

      <!-- 订单信息 -->
      <div v-if="packageDetail.order" class="detail-section">
        <h3 class="section-title">订单信息</h3>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="订单ID">{{ packageDetail.order.id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单号">{{ packageDetail.order.sn || '-' }}</el-descriptions-item>
          <el-descriptions-item label="支付状态">
            <el-tag size="small" :type="packageDetail.order.status === 'paid' ? 'success' : 'warning'">
              {{ packageDetail.order.status === 'paid' ? '已支付' : '未支付' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ formatDate(packageDetail.order.payTime) || '-' }}</el-descriptions-item>
          <el-descriptions-item label="订单金额" :span="2">
            <span class="price-value">¥{{ packageDetail.order.totalFee || 0 }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      
      <!-- 套餐信息 -->
      <div v-if="packageDetail.package" class="detail-section">
        <h3 class="section-title">套餐详情</h3>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="套餐ID">{{ packageDetail.package.id || '-' }}</el-descriptions-item>
          <el-descriptions-item label="套餐名称">{{ packageDetail.package.name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="套餐价格">
            <span class="price-value">¥{{ packageDetail.package.price || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="时长(分钟)">{{ packageDetail.package.durationMinutes || 0 }}</el-descriptions-item>
          <el-descriptions-item label="套餐描述" :span="2">{{ packageDetail.package.description || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <!-- 使用记录 -->
      <div v-if="packageDetail.usage && packageDetail.usage.records && packageDetail.usage.records.length > 0" class="detail-section">
        <h3 class="section-title">使用记录</h3>
        <div class="usage-stats">
          <el-tag type="info">总使用次数: {{ packageDetail.usage.totalUsedTimes || 0 }}</el-tag>
          <el-tag type="success" v-if="packageDetail.usage.lastUsedTime">最后使用: {{ formatDate(packageDetail.usage.lastUsedTime) }}</el-tag>
        </div>
        <el-table :data="packageDetail.usage.records" border style="width: 100%">
          <el-table-column prop="orderId" label="订单ID" width="100" />
          <el-table-column prop="productName" label="商品名称" min-width="150" />
          <el-table-column prop="quantity" label="数量" width="80" align="center" />
          <el-table-column label="使用时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else class="no-records">暂无使用记录</div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button 
            v-if="packageDetail.status === 'active'"
            type="danger"
            @click="handleStatusChange(packageDetail, 'expired')"
          >设为过期</el-button>
          <el-button 
            v-if="packageDetail.status === 'pending'"
            type="success"
            @click="handleStatusChange(packageDetail, 'active')"
          >设为生效</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getUserPackageList, getUserPackageDetail, updateUserPackageStatus } from '../../api/package';

// 防止重复请求的锁
const isRequestLocked = ref(false);
// 记录页面是否已经初始化
const isInitialized = ref(false);
// 详情请求锁
const isDetailLocked = ref(false);

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  nickname: '',
  packageName: '',
  status: ''
});

// 用户套餐列表数据
const userPackageList = ref([]);
const total = ref(0);
const loading = ref(false);

// 用户套餐详情
const detailDialogVisible = ref(false);
const detailLoading = ref(false);
const packageDetail = ref({});

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
};

// 获取状态类型
const getStatusType = (status) => {
  switch (status) {
    case 'active':
      return 'success';
    case 'expired':
      return 'info';
    case 'pending':
      return 'warning';
    case 'refunded':
      return 'danger';
    default:
      return 'info'; // 默认为info类型
  }
};

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'active':
      return '生效中';
    case 'expired':
      return '已过期';
    case 'pending':
      return '待支付';
    case 'refunded':
      return '已退款';
    default:
      return status || '未知'; // 直接显示状态值或未知
  }
};

// 获取用户套餐列表
const getList = async () => {
  if (isRequestLocked.value) return;
  isRequestLocked.value = true;
  loading.value = true;
  try {
    // 准备API请求参数
    const apiParams = {
      page: queryParams.pageNum,
      limit: queryParams.pageSize
    };
    
    // 添加筛选条件
    if (queryParams.nickname) {
      apiParams.user_id = queryParams.nickname; // 后端可能支持通过昵称或ID查询
    }
    if (queryParams.packageName) {
      apiParams.package_id = queryParams.packageName; // 后端可能支持通过套餐名称或ID查询
    }
    if (queryParams.status) {
      apiParams.status = queryParams.status;
    }
    
    console.log('请求参数:', apiParams);
    const res = await getUserPackageList(apiParams);
    console.log('用户套餐列表响应:', res);
    
    if (res.code === 200 && res.data) {
      // 处理后端返回的数据
      if (Array.isArray(res.data.list)) {
        userPackageList.value = res.data.list.map(item => processUserPackage(item));
      } else if (res.data.list && typeof res.data.list === 'object') {
        // 如果是对象形式的列表，转换为数组
        const listArray = [];
        for (const key in res.data.list) {
          if (Object.prototype.hasOwnProperty.call(res.data.list, key)) {
            listArray.push(processUserPackage(res.data.list[key]));
          }
        }
        userPackageList.value = listArray;
      } else if (Array.isArray(res.data)) {
        userPackageList.value = res.data.map(item => processUserPackage(item));
      } else {
        userPackageList.value = [];
      }
      
      total.value = res.data.total || userPackageList.value.length;
    } else {
      userPackageList.value = [];
      total.value = 0;
      console.error('获取用户套餐列表响应格式异常:', res);
    }
  } catch (error) {
    console.error('获取用户套餐列表失败:', error);
    ElMessage.error('获取用户套餐列表失败');
    userPackageList.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
    isRequestLocked.value = false;
  }
};

// 处理用户套餐数据
const processUserPackage = (item) => {
  return {
    id: item.id,
    packageName: item.package_name || '',
    startTime: item.start_time || '',
    endTime: item.end_time || '',
    status: item.status || 'unknown',
    remainingMinutes: item.remaining_minutes || 0,
    createdAt: item.created_at || '',
    // 用户信息
    user: {
      nickname: item.user_name || '',
      phone: item.user_phone || '',
      avatarUrl: item.user_avatar || ''
    },
    // 套餐信息
    package: {
      id: item.package_id || 0,
      name: item.package_name || '',
      price: item.package_price || 0
    },
    // 订单信息
    order: {
      id: item.order_id || '',
      sn: item.order_sn || ''
    }
  };
};

// 查询
const handleQuery = () => {
  queryParams.pageNum = 1;
  getList();
};

// 重置查询
const resetQuery = () => {
  queryParams.nickname = '';
  queryParams.packageName = '';
  queryParams.status = '';
  handleQuery();
};

// 处理分页大小变化
const handleSizeChange = (size) => {
  queryParams.pageSize = size;
  getList();
};

// 处理页码变化
const handleCurrentChange = (page) => {
  queryParams.pageNum = page;
  getList();
};

// 查看详情
const handleViewDetail = async (row) => {
  if (isDetailLocked.value) return;
  isDetailLocked.value = true;
  detailDialogVisible.value = true;
  detailLoading.value = true;
  
  try {
    console.log('获取用户套餐详情，ID:', row.id);
    const res = await getUserPackageDetail(row.id);
    console.log('用户套餐详情原始数据:', res.data);
    
    if (res.code !== 200 || !res.data) {
      throw new Error('获取详情失败或数据为空');
    }
    
    // 处理返回的数据
    const detailData = res.data;
    
    // 构建套餐详情对象
    packageDetail.value = {
      id: detailData.id,
      packageName: detailData.package_name || '',
      user: {
        nickname: detailData.user_name || '',
        phone: detailData.user_phone || '',
        avatarUrl: detailData.user_avatar || detailData.avatar_url || ''
      },
      startTime: detailData.start_time || '',
      endTime: detailData.end_time || '',
      status: detailData.status || 'unknown',
      remainingMinutes: detailData.remaining_minutes || 0,
      createdAt: detailData.created_at || '',
      updatedAt: detailData.updated_at || '',
      validPeriod: detailData.valid_period || '',
      statusDesc: detailData.status_desc || '',
      
      // 套餐信息
      package: {
        id: detailData.package_id || 0,
        name: detailData.package_name || '',
        price: detailData.package_price || 0,
        durationMinutes: detailData.duration_minutes || 60,
        description: detailData.package?.description || detailData.description || ''
      },
      
      // 订单信息
      order: detailData.order ? {
        id: detailData.order.id || detailData.order_id || '',
        sn: detailData.order.order_sn || detailData.order_sn || '',
        status: detailData.order.pay_status || detailData.pay_status || '',
        payTime: detailData.order.pay_time || detailData.pay_time || '',
        totalFee: detailData.order.total_fee || detailData.total_fee || 0
      } : null,
      
      // 使用记录
      usage: {
        totalUsedTimes: detailData.total_used_times || detailData.usage?.total_used_times || 0,
        lastUsedTime: detailData.last_used_time || detailData.usage?.last_used_time || '',
        records: Array.isArray(detailData.usage_records) 
          ? detailData.usage_records.map(record => ({
              orderId: record.order_id || '',
              productName: record.product_name || '',
              quantity: record.quantity || 1,
              createdAt: record.created_at || ''
            }))
          : []
      }
    };
    
    console.log('处理后的套餐详情:', packageDetail.value);
  } catch (error) {
    console.error('获取用户套餐详情失败:', error);
    ElMessage.error('获取用户套餐详情失败');
    packageDetail.value = {}; // 重置详情对象
  } finally {
    detailLoading.value = false;
    isDetailLocked.value = false;
  }
};

// 修改套餐状态
const handleStatusChange = (row, status) => {
  const statusText = getStatusText(status);
  
  ElMessageBox.confirm(`确认将该套餐状态修改为"${statusText}"吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 发送状态更新请求
      await updateUserPackageStatus(row.id, status);
      ElMessage.success('状态修改成功');
      getList(); // 刷新列表
    } catch (error) {
      console.error('修改套餐状态失败:', error);
      ElMessage.error('修改套餐状态失败');
    }
  }).catch(() => {});
};

// 初始化
onMounted(() => {
  if (!isInitialized.value) {
    getList();
    isInitialized.value = true;
  }
});

// 当页面被激活时（从缓存中恢复）重新加载数据
onActivated(() => {
  // 避免重复请求数据
  if (!isRequestLocked.value) {
    getList();
  }
});
</script>

<style scoped>
.user-package-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
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
  margin-top: 3px;
}

.usage-title {
  margin-top: 20px;
  margin-bottom: 10px;
}

.no-records {
  text-align: center;
  color: #909399;
  padding: 20px 0;
}

.remaining-time {
  margin-top: 5px;
}

.status-desc {
  margin-top: 5px;
}

.order-info {
  font-size: 12px;
  color: #606266;
  margin-top: 3px;
}

.detail-section {
  margin-top: 20px;
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 15px;
  padding-left: 10px;
  border-left: 3px solid #409EFF;
}

.usage-stats {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.price-value {
  color: #f56c6c;
  font-weight: 600;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 