<template>
  <div class="page-container">
    <div class="page-title">用户管理</div>
    
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/手机号" clearable />
        </el-form-item>
        <el-form-item label="注册时间">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <div class="left">
          <el-button type="primary" @click="exportUserData">导出数据</el-button>
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="userList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="avatar_url" label="头像" width="80">
          <template #default="scope">
            <el-avatar :src="scope.row.avatar_url || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
          </template>
        </el-table-column>
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="openid" label="OpenID" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column prop="updated_at" label="最后更新" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              @click="handleViewDetail(scope.row)"
            >
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="用户详情"
      width="600px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item label="OpenID">{{ currentUser.openid }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ currentUser.nickname }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentUser.phone }}</el-descriptions-item>
        <el-descriptions-item label="注册时间">{{ currentUser.created_at }}</el-descriptions-item>
        <el-descriptions-item label="最后更新">{{ currentUser.updated_at }}</el-descriptions-item>
        <el-descriptions-item label="头像" :span="2">
          <el-avatar 
            :size="100" 
            :src="currentUser.avatar_url || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" 
          />
        </el-descriptions-item>
      </el-descriptions>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getUserList } from '../../api/user';

// 防止重复请求的锁
const isRequestLocked = ref(false);

// 记录页面是否已经初始化
const isInitialized = ref(false);

// 加载状态
const loading = ref(false);

// 用户列表数据
const userList = ref([]);

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 搜索表单
const searchForm = reactive({
  keyword: '',
  dateRange: [],
  startDate: '',
  endDate: ''
});

// 详情对话框
const detailDialogVisible = ref(false);
const currentUser = reactive({
  id: '',
  username: '',
  nickname: '',
  phone: '',
  gender: '',
  status: 1,
  created_at: '',
  updated_at: '',
  avatar_url: ''
});

// 获取用户列表
const fetchUserList = async () => {
  if (isRequestLocked.value) return;
  isRequestLocked.value = true;
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      keyword: searchForm.keyword || '',
      start_date: searchForm.startDate || '',
      end_date: searchForm.endDate || ''
    };
    
    const response = await getUserList(params);
    if (response && response.code === 200) {
      userList.value = response.data.list || [];
      total.value = response.data.total || 0;
    } else {
      ElMessage.error(response?.message || '获取用户列表失败');
    }
  } catch (error) {
    console.error('获取用户列表失败:', error);
    ElMessage.error('获取用户列表失败');
  } finally {
    loading.value = false;
    isRequestLocked.value = false;
  }
};

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchUserList();
};

// 处理日期范围变化
const handleDateRangeChange = (val) => {
  if (val && val.length === 2) {
    searchForm.startDate = val[0];
    searchForm.endDate = val[1];
  } else {
    searchForm.startDate = '';
    searchForm.endDate = '';
  }
};

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = '';
  searchForm.dateRange = [];
  searchForm.startDate = '';
  searchForm.endDate = '';
  currentPage.value = 1;
  fetchUserList();
};

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchUserList();
};

// 页码变化
const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchUserList();
};

// 查看用户详情
const handleViewDetail = async (row) => {
  try {
    // 直接使用列表中的数据，不再额外请求
    Object.assign(currentUser, row);
    detailDialogVisible.value = true;
  } catch (error) {
    console.error('获取用户详情失败:', error);
    ElMessage.error('获取用户详情失败');
  }
};

// 导出用户数据
const exportUserData = () => {
  ElMessageBox.confirm('确定要导出用户数据吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage({
      type: 'success',
      message: '导出成功，文件已开始下载'
    });
    // 这里实现导出逻辑
  }).catch(() => {});
};

// 页面加载时获取数据
onMounted(() => {
  console.log('用户列表页面已挂载');
  if (!isInitialized.value) {
    fetchUserList();
    isInitialized.value = true;
  }
});

// 当页面从缓存中激活时触发（切换tab时）
onActivated(() => {
  console.log('用户列表页面已激活');
  // 避免重复请求数据
  if (!isRequestLocked.value) {
    fetchUserList();
  }
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.page-title {
  margin-bottom: 20px;
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 