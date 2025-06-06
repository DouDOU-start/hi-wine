<template>
  <div class="package-list-container">
    <div class="page-header">
      <h2>畅饮套餐列表</h2>
      <el-button type="primary" @click="handleAddPackage">添加套餐</el-button>
    </div>
    
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="filter-form">
        <el-form-item label="套餐名称">
          <el-input v-model="queryParams.name" placeholder="请输入套餐名称" clearable @keyup.enter="handleQuery" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.isActive" placeholder="全部状态" clearable>
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
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
        :data="packageList"
        border
        style="width: 100%"
      >
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
        <el-table-column prop="isActive" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isActive ? 'success' : 'danger'">
              {{ scope.row.isActive ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="success" size="small" @click="handleViewProducts(scope.row)">查看商品</el-button>
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
    
    <!-- 套餐商品对话框 -->
    <el-dialog
      v-model="productDialogVisible"
      title="套餐包含商品"
      width="650px"
    >
      <el-table
        v-loading="productLoading"
        :data="packageProducts"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="商品名称" min-width="150" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="scope">
            ¥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="图片" width="100">
          <template #default="scope">
            <el-image
              style="width: 50px; height: 50px"
              :src="scope.row.imageUrl"
              :preview-src-list="[scope.row.imageUrl]"
              fit="cover"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getPackageList, deletePackage, getPackageProducts } from '../../api/package';

const router = useRouter();

// 查询参数
const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  isActive: null
});

// 套餐列表数据
const packageList = ref([]);
const total = ref(0);
const loading = ref(false);

// 套餐商品相关
const productDialogVisible = ref(false);
const packageProducts = ref([]);
const productLoading = ref(false);
const currentPackageId = ref(null);

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
};

// 获取套餐列表
const getList = async () => {
  loading.value = true;
  try {
    const res = await getPackageList(queryParams);
    packageList.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取套餐列表失败:', error);
    ElMessage.error('获取套餐列表失败');
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
  queryParams.name = '';
  queryParams.isActive = null;
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

// 添加套餐
const handleAddPackage = () => {
  router.push('/package/add');
};

// 编辑套餐
const handleEdit = (row) => {
  router.push(`/package/edit/${row.id}`);
};

// 查看套餐商品
const handleViewProducts = async (row) => {
  productDialogVisible.value = true;
  currentPackageId.value = row.id;
  productLoading.value = true;
  
  try {
    const res = await getPackageProducts(row.id);
    packageProducts.value = res.data.products || [];
  } catch (error) {
    console.error('获取套餐商品失败:', error);
    ElMessage.error('获取套餐商品失败');
  } finally {
    productLoading.value = false;
  }
};

// 删除套餐
const handleDelete = (row) => {
  ElMessageBox.confirm(`确认删除套餐 "${row.name}" 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePackage(row.id);
      ElMessage.success('删除成功');
      getList();
    } catch (error) {
      console.error('删除套餐失败:', error);
      ElMessage.error('删除套餐失败');
    }
  }).catch(() => {});
};

// 初始化
onMounted(() => {
  getList();
});
</script>

<style scoped>
.package-list-container {
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