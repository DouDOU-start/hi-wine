<template>
  <div class="table-list-container">
    <div class="page-header">
      <h2>桌号管理</h2>
      <el-button type="primary" size="large" icon="Plus" @click="handleAddTable">添加桌号</el-button>
    </div>
    
    <el-card class="filter-container" shadow="hover">
      <el-form :inline="true" :model="queryParams" class="filter-form">
        <el-form-item label="桌号">
          <el-input v-model="queryParams.tableNumber" placeholder="请输入桌号" clearable @keyup.enter="handleQuery">
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
            <el-option label="空闲" value="idle" />
            <el-option label="占用" value="occupied" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery" icon="Search">查询</el-button>
          <el-button @click="resetQuery" icon="Refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card class="list-container" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">桌号列表</span>
          <div class="header-operations">
            <el-button type="primary" icon="Refresh" circle @click="getList" />
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="tableList"
        border
        stripe
        highlight-current-row
        style="width: 100%"
        table-layout="auto"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="tableNumber" label="桌号" min-width="120" align="center" />
        <el-table-column label="二维码" width="120" align="center">
          <template #default="scope">
            <div class="qrcode-container">
              <el-image
                class="qrcode-image"
                :src="scope.row.qrcodeUrl"
                :preview-src-list="[scope.row.qrcodeUrl]"
                fit="cover"
                preview-teleported
              >
                <template #error>
                  <div class="qrcode-error">
                    <el-icon><Picture /></el-icon>
                    <span>加载失败</span>
                  </div>
                </template>
              </el-image>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'idle' ? 'success' : 'danger'" effect="dark" round>
              {{ scope.row.status === 'idle' ? '空闲' : '占用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isActive" label="是否激活" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.isActive ? 'success' : 'info'" effect="plain">
              {{ scope.row.isActive ? '已激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180" align="center">
          <template #default="scope">
            <div class="time-info">
              <el-icon><Calendar /></el-icon>
              <span>{{ formatDate(scope.row.createdAt) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right" align="center">
          <template #default="scope">
            <div class="table-actions">
              <el-button
                type="primary"
                size="small"
                @click="handleEdit(scope.row)"
              >
                <el-icon><Edit /></el-icon>
                <span class="button-text">编辑</span>
              </el-button>
              
              <el-button
                type="success"
                size="small"
                @click="handleDownloadQrcode(scope.row)"
              >
                <el-icon><Download /></el-icon>
                <span class="button-text">下载</span>
              </el-button>
              
              <el-dropdown trigger="hover" @command="(command) => handleCommand(command, scope.row)">
                <el-button type="info" size="small">
                  <el-icon><MoreFilled /></el-icon>
                  <span class="button-text">更多</span>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="regenerate">
                      <el-icon><RefreshRight /></el-icon>
                      <span>重新生成二维码</span>
                    </el-dropdown-item>
                    <el-dropdown-item divided command="delete">
                      <el-icon><Delete /></el-icon>
                      <span style="color: #f56c6c;">删除桌号</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getTableList, deleteTable, regenerateQrcode, downloadQrcode } from '../../api/table';
import { Search, Refresh, Calendar, Edit, RefreshRight, Download, Delete, Picture, Plus, MoreFilled } from '@element-plus/icons-vue';
import { formatDate as formatDateTime } from '../../utils/format';

const router = useRouter();

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  tableNumber: '',
  status: ''
});

// 桌号列表数据
const tableList = ref([]);
const total = ref(0);
const loading = ref(false);

// 格式化日期
const formatDate = (timestamp) => {
  if (!timestamp) return '-';
  return formatDateTime(timestamp);
};

// 获取桌号列表
const getList = async () => {
  loading.value = true;
  try {
    const response = await getTableList(queryParams);
    tableList.value = response.data.list || [];
    total.value = response.data.total || 0;
  } catch (error) {
    ElMessage.error('获取桌号列表失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 查询操作
const handleQuery = () => {
  queryParams.pageNum = 1;
  getList();
};

// 重置查询
const resetQuery = () => {
  queryParams.tableNumber = '';
  queryParams.status = '';
  handleQuery();
};

// 分页大小变化
const handleSizeChange = (size) => {
  queryParams.pageSize = size;
  getList();
};

// 页码变化
const handleCurrentChange = (page) => {
  queryParams.pageNum = page;
  getList();
};

// 添加桌号
const handleAddTable = () => {
  router.push('/table/add');
};

// 编辑桌号
const handleEdit = (row) => {
  router.push(`/table/edit/${row.id}`);
};

// 下载二维码
const handleDownloadQrcode = async (row) => {
  try {
    await downloadQrcode(row.id);
    ElMessage.success('二维码下载成功');
  } catch (error) {
    ElMessage.error('下载二维码失败');
    console.error(error);
  }
};

// 处理更多操作
const handleCommand = (command, row) => {
  switch (command) {
    case 'regenerate':
      handleRegenerateQrcode(row);
      break;
    case 'delete':
      handleDelete(row);
      break;
  }
};

// 重新生成二维码
const handleRegenerateQrcode = (row) => {
  ElMessageBox.confirm('确认重新生成二维码吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await regenerateQrcode(row.id);
      ElMessage.success('二维码重新生成成功');
      getList();
    } catch (error) {
      ElMessage.error('重新生成二维码失败');
      console.error(error);
    }
  }).catch(() => {});
};

// 删除桌号
const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除桌号 ${row.tableNumber} 吗?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteTable(row.id);
      ElMessage.success('删除桌号成功');
      getList();
    } catch (error) {
      ElMessage.error('删除桌号失败');
      console.error(error);
    }
  }).catch(() => {});
};

// 初始化
onMounted(() => {
  getList();
});
</script>

<style scoped>
.table-list-container {
  padding: 0 20px 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  font-weight: 500;
  color: #303133;
  margin: 0;
}

.filter-container {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.list-container {
  margin-bottom: 20px;
}

.qrcode-container {
  display: flex;
  justify-content: center;
  align-items: center;
}

.qrcode-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
}

.qrcode-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 12px;
}

.time-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}

.table-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.empty-data {
  padding: 40px 0;
}
</style> 