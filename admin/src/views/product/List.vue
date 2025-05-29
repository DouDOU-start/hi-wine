<template>
  <div class="page-container">
    <div class="page-title">商品管理</div>
    
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="商品名称">
          <el-input v-model="searchForm.name" placeholder="请输入商品名称" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.categoryId" placeholder="请选择分类" clearable>
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="上架" :value="1" />
            <el-option label="下架" :value="0" />
          </el-select>
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
          <el-button type="primary" @click="handleAdd">添加商品</el-button>
          <el-button @click="handleBatchDelete" :disabled="selectedIds.length === 0">批量删除</el-button>
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="productList"
        border
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="image" label="图片" width="100">
          <template #default="scope">
            <el-image 
              v-if="scope.row.image" 
              :src="scope.row.image" 
              style="width: 60px; height: 60px"
              fit="cover"
              :preview-src-list="[scope.row.image]"
              @error="handleImageError(scope.row)"
            />
            <span v-else>无图片</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" show-overflow-tooltip />
        <el-table-column prop="categoryName" label="分类" width="120" />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="scope">
            ￥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="sales" label="销量" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              @click="handleEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(scope.row)"
            >
              删除
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getProductList, deleteProduct, updateProductStatus } from '../../api/product';
import { getCategoryList } from '../../api/category';

const router = useRouter();

// 加载状态
const loading = ref(false);

// 商品列表数据
const productList = ref([]);

// 分类选项
const categoryOptions = ref([]);

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 搜索表单
const searchForm = reactive({
  name: '',
  categoryId: '',
  status: ''
});

// 选中的商品ID
const selectedIds = ref([]);

// 获取商品列表
const fetchProductList = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      ...searchForm
    };
    
    const response = await getProductList(params);
    
    // 处理返回的数据，确保每个商品有分类名称
    if (response.data && response.data.list) {
      productList.value = response.data.list.map(item => {
        // 查找对应的分类名称
        const category = categoryOptions.value.find(cat => cat.id === item.category_id);
        return {
          ...item,
          categoryName: category ? category.name : '未分类',
          // 确保状态是数字类型
          status: typeof item.status === 'string' ? parseInt(item.status) : item.status
        };
      });
      total.value = response.data.total || 0;
    } else {
      productList.value = [];
      total.value = 0;
    }
  } catch (error) {
    console.error('获取商品列表失败:', error);
    ElMessage.error('获取商品列表失败');
  } finally {
    loading.value = false;
  }
};

// 获取分类列表
const fetchCategoryList = async () => {
  try {
    const response = await getCategoryList({ size: 100 });
    categoryOptions.value = response.data.list || [];
  } catch (error) {
    console.error('获取分类列表失败:', error);
  }
};

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchProductList();
};

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = '';
  });
  currentPage.value = 1;
  fetchProductList();
};

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchProductList();
};

// 页码变化
const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchProductList();
};

// 选择变化
const handleSelectionChange = (selection) => {
  selectedIds.value = selection.map(item => item.id);
};

// 添加商品
const handleAdd = () => {
  router.push('/product/add');
};

// 编辑商品
const handleEdit = (row) => {
  router.push(`/product/edit/${row.id}`);
};

// 删除商品
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除商品"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteProduct(row.id);
      ElMessage.success('删除成功');
      fetchProductList();
    } catch (error) {
      console.error('删除商品失败:', error);
      ElMessage.error('删除商品失败');
    }
  }).catch(() => {});
};

// 批量删除商品
const handleBatchDelete = () => {
  if (selectedIds.value.length === 0) {
    return;
  }
  
  ElMessageBox.confirm(`确定要删除选中的${selectedIds.value.length}个商品吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 这里需要后端提供批量删除接口
      for (const id of selectedIds.value) {
        await deleteProduct(id);
      }
      ElMessage.success('批量删除成功');
      fetchProductList();
      selectedIds.value = [];
    } catch (error) {
      console.error('批量删除商品失败:', error);
      ElMessage.error('批量删除商品失败');
    }
  }).catch(() => {});
};

// 修改商品状态
const handleStatusChange = async (row) => {
  try {
    await updateProductStatus(row.id, row.status);
    ElMessage.success(`商品已${row.status === 1 ? '上架' : '下架'}`);
  } catch (error) {
    console.error('更新商品状态失败:', error);
    ElMessage.error('更新商品状态失败');
    // 恢复原状态
    row.status = row.status === 1 ? 0 : 1;
  }
};

// 图片错误处理
const handleImageError = (row) => {
  console.error(`图片加载失败: ${row.name}`, row.image);
  // 避免多次提示，不显示消息提示
};

// 页面加载时获取数据
onMounted(() => {
  fetchCategoryList();
  fetchProductList();
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