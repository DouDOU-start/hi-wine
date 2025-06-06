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
            <div v-if="scope.row.status === 'active'">
              剩余: {{ scope.row.remainingMinutes }} 分钟
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
    >
      <el-descriptions v-loading="detailLoading" :column="2" border>
        <el-descriptions-item label="套餐ID">{{ packageDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="套餐名称">{{ packageDetail.packageName }}</el-descriptions-item>
        <el-descriptions-item label="用户昵称">{{ packageDetail.user?.nickname }}</el-descriptions-item>
        <el-descriptions-item label="用户手机">{{ packageDetail.user?.phone || '未绑定手机' }}</el-descriptions-item>
        <el-descriptions-item label="开始时间">{{ formatDate(packageDetail.startTime) }}</el-descriptions-item>
        <el-descriptions-item label="结束时间">{{ formatDate(packageDetail.endTime) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(packageDetail.status)">
            {{ getStatusText(packageDetail.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="剩余时间" v-if="packageDetail.status === 'active'">
          {{ packageDetail.remainingMinutes }} 分钟
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatDate(packageDetail.createdAt) }}</el-descriptions-item>
      </el-descriptions>
      
      <template v-if="packageDetail.usageRecords && packageDetail.usageRecords.length > 0">
        <h3 class="usage-title">使用记录</h3>
        <el-table :data="packageDetail.usageRecords" border style="width: 100%">
          <el-table-column prop="orderId" label="订单ID" width="100" />
          <el-table-column prop="productName" label="商品名称" min-width="150" />
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column label="使用时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>
      </template>
      <template v-else>
        <div class="no-records">暂无使用记录</div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getUserPackageList, getUserPackageDetail, updateUserPackageStatus } from '../../api/package';

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
      return '';
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
      return '未知';
  }
};

// 获取用户套餐列表
const getList = async () => {
  loading.value = true;
  try {
    const res = await getUserPackageList(queryParams);
    userPackageList.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取用户套餐列表失败:', error);
    ElMessage.error('获取用户套餐列表失败');
  } finally {
    loading.value = false;
  }
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
  detailDialogVisible.value = true;
  detailLoading.value = true;
  
  try {
    const res = await getUserPackageDetail(row.id);
    packageDetail.value = res.data;
  } catch (error) {
    console.error('获取用户套餐详情失败:', error);
    ElMessage.error('获取用户套餐详情失败');
  } finally {
    detailLoading.value = false;
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
      await updateUserPackageStatus(row.id, status);
      ElMessage.success('状态修改成功');
      getList();
    } catch (error) {
      console.error('修改套餐状态失败:', error);
      ElMessage.error('修改套餐状态失败');
    }
  }).catch(() => {});
};

// 初始化
onMounted(() => {
  getList();
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
</style> 