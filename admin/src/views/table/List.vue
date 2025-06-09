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
            <el-tooltip content="刷新数据" placement="top">
              <el-button type="primary" icon="Refresh" circle @click="getList" />
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <div v-if="loading" class="table-loading">
        <el-skeleton :rows="5" animated />
      </div>
      
      <div v-else-if="tableList.length === 0" class="empty-data">
        <el-empty description="暂无桌号数据" />
      </div>
      
      <el-table
        v-else
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
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="scope">
            <div class="table-actions">
              <el-tooltip content="编辑" placement="top" :hide-after="1000">
                <el-button
                  type="primary"
                  size="small"
                  @click="handleEdit(scope.row)"
                >
                  <el-icon><Edit /></el-icon>
                  <span class="button-text">编辑</span>
                </el-button>
              </el-tooltip>
              
              <el-tooltip content="下载二维码" placement="top" :hide-after="1000">
                <el-button
                  type="success"
                  size="small"
                  @click="handleDownloadQrcode(scope.row)"
                >
                  <el-icon><Download /></el-icon>
                  <span class="button-text">下载</span>
                </el-button>
              </el-tooltip>
              
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
import { ref, reactive, onMounted, onActivated, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getTableList, deleteTable, regenerateQrcode, downloadQrcode } from '../../api/table';
import { Search, Refresh, Calendar, Edit, RefreshRight, Download, Delete, Picture, Plus, MoreFilled } from '@element-plus/icons-vue';

const router = useRouter();
const route = useRoute();

// 防止重复请求的锁
const isRequestLocked = ref(false);

// 记录页面是否已经初始化
const isInitialized = ref(false);

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

// 处理单个桌号项的数据
const processTableItem = (item) => {
  // 从截图中看到字段名是下划线形式
  return {
    id: item.id,
    tableNumber: item.table_number || '',
    qrcodeUrl: item.qrcode_url || '',
    status: item.status || 'idle',
    isActive: item.is_active !== undefined ? item.is_active : false,
    createdAt: item.created_at || ''
  };
};

// 获取桌号列表
const getList = async (fromRoute = false) => {
  // 如果已经在加载中或请求被锁定且不是来自路由变化，则跳过
  if (loading.value || (isRequestLocked.value && !fromRoute)) {
    console.log('请求被锁定或正在加载中，跳过此次请求');
    return;
  }
  
  // 锁定请求，防止短时间内重复调用
  isRequestLocked.value = true;
  loading.value = true;
  
  try {
    const res = await getTableList(queryParams);
    console.log('桌号列表响应数据:', res);
    
    // 根据截图中的数据结构进行处理
    if (res.data && res.data.list && Array.isArray(res.data.list)) {
      // 处理list数组形式的数据
      tableList.value = res.data.list.map(item => {
        return processTableItem(item);
      });
      total.value = res.data.total || res.data.list.length;
    } else if (res.data && Array.isArray(res.data)) {
      // 处理直接返回数组的情况
      tableList.value = res.data.map(item => {
        return processTableItem(item);
      });
      total.value = res.data.length;
    } else if (res.data) {
      // 根据截图，后端返回的是对象结构，需要特殊处理
      console.log('非标准数据结构:', res.data);
      
      // 从截图看，数据在res.data.list中，且list是一个对象，不是数组
      const listData = res.data.list || {};
      const extractedList = [];
      
      // 遍历对象中的每个键，如果是数字索引，则添加到列表中
      for (const key in listData) {
        if (!isNaN(parseInt(key))) {
          extractedList.push(listData[key]);
        }
      }
      
      console.log('提取的列表数据:', extractedList);
      
      // 处理每个桌号项
      tableList.value = extractedList.map(item => {
        return processTableItem(item);
      });
      
      // 设置总数
      total.value = res.data.total || extractedList.length;
    } else {
      tableList.value = [];
      total.value = 0;
      console.error('获取桌号列表数据格式异常', res);
    }
    
    // 调试输出处理后的数据
    console.log('处理后的桌号列表:', tableList.value);
  } catch (error) {
    console.error('获取桌号列表失败:', error);
    ElMessage.error('获取桌号列表失败');
    tableList.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
    // 延迟解锁，防止短时间内重复请求
    setTimeout(() => {
      isRequestLocked.value = false;
    }, 300);
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

// 添加处理下拉菜单命令的方法
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

// 初始化
onMounted(() => {
  console.log('桌号列表页面挂载');
  isInitialized.value = true;
  getList(true);
});

// 监听路由变化，从添加或编辑页面返回时刷新数据
watch(
  () => route.fullPath,
  (newPath, oldPath) => {
    console.log('路由变化:', oldPath, '->', newPath);
    if (newPath.includes('/table/list') && (oldPath.includes('/table/add') || oldPath.includes('/table/edit'))) {
      console.log('从添加/编辑页面返回，刷新数据');
      // 设置请求标志，避免onActivated重复请求
      isRequestLocked.value = true;
      getList(true);
    }
  }
);

// 添加页面激活时刷新数据（用于keep-alive场景）
onActivated(() => {
  console.log('桌号列表页面激活');
  if (isInitialized.value && !isRequestLocked.value) {
    getList(false); // 不强制刷新，因为可能已经在路由监听中触发了刷新
  }
});
</script>

<style scoped>
.table-list-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.filter-container {
  margin-bottom: 20px;
  border-radius: 8px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.list-container {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.header-operations {
  display: flex;
  gap: 10px;
}

.table-loading {
  padding: 20px;
}

.empty-data {
  padding: 40px 0;
}

.qrcode-container {
  display: flex;
  justify-content: center;
}

.qrcode-image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  transition: all 0.3s;
}

.qrcode-image:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.qrcode-error {
  width: 80px;
  height: 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  color: #909399;
  font-size: 12px;
  border-radius: 4px;
}

.qrcode-error .el-icon {
  font-size: 24px;
  margin-bottom: 5px;
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
  flex-wrap: wrap;
  gap: 8px;
}

.button-text {
  margin-left: 4px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .filter-form {
    flex-direction: column;
  }
  
  .table-actions {
    flex-direction: column;
  }
  
  .pagination-container {
    justify-content: center;
  }
}
</style> 