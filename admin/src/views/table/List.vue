<template>
  <div class="table-list-container">
    <div class="page-header">
      <h2>桌号管理</h2>
      <el-button type="primary" @click="handleAddTable">添加桌号</el-button>
    </div>
    
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="filter-form">
        <el-form-item label="桌号">
          <el-input v-model="queryParams.tableNumber" placeholder="请输入桌号" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部状态" clearable>
            <el-option label="空闲" value="idle" />
            <el-option label="占用" value="occupied" />
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
        :data="tableList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="tableNumber" label="桌号" min-width="120" />
        <el-table-column label="二维码" width="120">
          <template #default="scope">
            <el-image
              style="width: 80px; height: 80px"
              :src="scope.row.qrcodeUrl"
              :preview-src-list="[scope.row.qrcodeUrl]"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'idle' ? 'success' : 'danger'">
              {{ scope.row.status === 'idle' ? '空闲' : '占用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="isActive" label="是否激活" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isActive ? 'success' : 'info'">
              {{ scope.row.isActive ? '已激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEdit(scope.row)"
            >编辑</el-button>
            <el-button
              type="success"
              size="small"
              @click="handleRegenerateQrcode(scope.row)"
            >重新生成二维码</el-button>
            <el-button
              type="info"
              size="small"
              @click="handleDownloadQrcode(scope.row)"
            >下载二维码</el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row)"
            >删除</el-button>
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
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
};

// 获取桌号列表
const getList = async () => {
  loading.value = true;
  try {
    const res = await getTableList(queryParams);
    tableList.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取桌号列表失败:', error);
    ElMessage.error('获取桌号列表失败');
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
  queryParams.tableNumber = '';
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

// 添加桌号
const handleAddTable = () => {
  router.push('/table/add');
};

// 编辑桌号
const handleEdit = (row) => {
  router.push(`/table/edit/${row.id}`);
};

// 重新生成二维码
const handleRegenerateQrcode = async (row) => {
  ElMessageBox.confirm(`确认重新生成桌号 "${row.tableNumber}" 的二维码吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await regenerateQrcode(row.id);
      ElMessage.success('二维码重新生成成功');
      getList(); // 刷新列表
    } catch (error) {
      console.error('重新生成二维码失败:', error);
      ElMessage.error('重新生成二维码失败');
    }
  }).catch(() => {});
};

// 下载二维码
const handleDownloadQrcode = async (row) => {
  try {
    const response = await downloadQrcode(row.id);
    
    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `qrcode-${row.tableNumber}.png`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    ElMessage.success('二维码下载成功');
  } catch (error) {
    console.error('下载二维码失败:', error);
    ElMessage.error('下载二维码失败');
  }
};

// 删除桌号
const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除桌号 "${row.tableNumber}" 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteTable(row.id);
      ElMessage.success('删除成功');
      getList();
    } catch (error) {
      console.error('删除桌号失败:', error);
      ElMessage.error('删除桌号失败');
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
</style> 