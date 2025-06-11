<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <page-header title="商品管理">
      <template #extra>
        <div class="data-summary">
          <el-tag type="info">总商品: {{ total }}</el-tag>
          <el-tag v-if="lowStockCount > 0" type="warning">库存不足: {{ lowStockCount }}</el-tag>
          <el-tag v-if="outOfStockCount > 0" type="danger">缺货: {{ outOfStockCount }}</el-tag>
        </div>
      </template>
    </page-header>
    
    <!-- 搜索表单 -->
    <product-search-form
      :category-options="categoryOptions"
      :low-stock-count="lowStockCount"
      :out-of-stock-count="outOfStockCount"
      @search="handleSearch"
      @reset="handleReset"
      @filter="handleFilter"
    />
    
    <!-- 商品表格 -->
    <product-table
      :data="productList"
      :loading="loading"
      :total="total"
      :current-page="pagination.page"
      :page-size="pagination.limit"
      @refresh="fetchProductList"
      @add="handleAdd"
      @edit="handleEdit"
      @delete="handleDelete"
      @toggle-status="handleToggleStatus"
      @export="handleExport"
      @batch-action="handleBatchAction"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
      @sort-change="handleSortChange"
      @update:current-page="pagination.page = $event"
      @update:page-size="pagination.limit = $event"
    />
    
    <!-- 加载状态 -->
    <loading-state :loading="loading" />
    
    <!-- 空数据状态 -->
    <empty-state v-if="!loading && productList.length === 0" description="暂无商品数据">
      <el-button type="primary" @click="handleAdd">添加商品</el-button>
    </empty-state>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, provide } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getProductList, deleteProduct, updateProductStatus, batchUpdateProductStatus } from '@/api/product';
import { getCategoryList } from '@/api/category';
import ProductSearchForm from './components/ProductSearchForm.vue';
import ProductTable from './components/ProductTable.vue';

const router = useRouter();

// 数据加载状态
const loading = ref(false);

// 商品列表数据
const productList = ref([]);

// 分类选项
const categoryOptions = ref([]);

// 提供分类选项给子组件
provide('categoryOptions', categoryOptions);

// 总数量
const total = ref(0);

// 库存不足数量
const lowStockCount = ref(0);

// 缺货数量
const outOfStockCount = ref(0);

// 分页参数
const pagination = ref({
  page: 1,
  limit: 10,
  sort: 'createTime',
  order: 'desc'
});

// 搜索参数
const searchParams = ref({});

// 获取商品列表
const fetchProductList = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.value.page,
      limit: pagination.value.limit,
      sort_field: pagination.value.sort,
      sort_order: pagination.value.order,
      ...searchParams.value
    };
    
    console.log('请求参数:', params);
    const response = await getProductList(params);
    console.log('API响应:', response);
    
    // 确保响应数据格式正确
    if (response && response.data) {
      // 检查数据格式
      if (Array.isArray(response.data)) {
        // 直接是数组格式
        productList.value = response.data;
        total.value = response.data.length; // 或者使用响应头中的total
      } else if (response.data.list && Array.isArray(response.data.list)) {
        // 包含list字段的对象格式
        productList.value = response.data.list;
        total.value = response.data.total || response.data.list.length;
      } else if (typeof response.data === 'object') {
        // 其他对象格式，尝试提取数据
        const dataArray = Object.values(response.data).filter(item => typeof item === 'object');
        if (dataArray.length > 0) {
          productList.value = dataArray;
          total.value = dataArray.length;
        } else {
          productList.value = [];
          total.value = 0;
        }
      } else {
        productList.value = [];
        total.value = 0;
      }
    } else {
      productList.value = [];
      total.value = 0;
    }
    
    console.log('处理后的商品列表:', productList.value);
    
    // 计算库存相关数量
    lowStockCount.value = productList.value.filter(item => item.stock > 0 && item.stock <= 10).length;
    outOfStockCount.value = productList.value.filter(item => item.stock <= 0).length;
  } catch (error) {
    console.error('获取商品列表失败:', error);
    ElMessage.error('获取商品列表失败');
    productList.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
};

// 获取分类列表
const fetchCategoryList = async () => {
  try {
    const response = await getCategoryList();
    console.log('分类数据:', response);
    
    if (response && response.data) {
      let categories = [];
      
      // 处理不同格式的分类数据
      if (Array.isArray(response.data)) {
        categories = response.data;
      } else if (response.data.list && Array.isArray(response.data.list)) {
        categories = response.data.list;
      }
      
      categoryOptions.value = categories.map(item => ({
        label: item.name,
        value: item.id
      }));
      
      console.log('处理后的分类选项:', categoryOptions.value);
    }
  } catch (error) {
    console.error('获取分类列表失败:', error);
  }
};

// 处理搜索
const handleSearch = (formData) => {
  // 处理日期范围
  const params = { ...formData };
  if (params.dateRange && params.dateRange.length === 2) {
    params.start_date = params.dateRange[0];
    params.end_date = params.dateRange[1];
    delete params.dateRange;
  }
  
  searchParams.value = params;
  pagination.value.page = 1;
  fetchProductList();
};

// 处理重置
const handleReset = () => {
  searchParams.value = {};
  pagination.value.page = 1;
  fetchProductList();
};

// 处理筛选
const handleFilter = (filter, formData) => {
  // 已在 ProductSearchForm 组件中处理
};

// 处理添加
const handleAdd = () => {
  router.push('/product/add');
};

// 处理编辑
const handleEdit = (row) => {
  router.push(`/product/edit/${row.id}`);
};

// 处理删除
const handleDelete = async (row) => {
  try {
    await deleteProduct(row.id);
    ElMessage.success('删除成功');
    fetchProductList();
  } catch (error) {
    console.error('删除商品失败:', error);
    ElMessage.error('删除商品失败');
  }
};

// 处理状态切换
const handleToggleStatus = async (row) => {
  const newStatus = row.status === 1 ? 0 : 1;
  try {
    await updateProductStatus(row.id, newStatus);
    ElMessage.success(`${newStatus === 1 ? '上架' : '下架'}成功`);
    fetchProductList();
  } catch (error) {
    console.error('更新商品状态失败:', error);
    ElMessage.error('更新商品状态失败');
  }
};

// 处理导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...');
};

// 处理批量操作
const handleBatchAction = async (actionType, ids) => {
  if (!ids || ids.length === 0) return;
  
  try {
    switch (actionType) {
      case 'online':
        await batchUpdateProductStatus(ids, 1);
        ElMessage.success('批量上架成功');
        break;
      case 'offline':
        await batchUpdateProductStatus(ids, 0);
        ElMessage.success('批量下架成功');
        break;
      case 'delete':
        // 这里需要实现批量删除API
        ElMessage.info('批量删除功能开发中...');
        break;
    }
    fetchProductList();
  } catch (error) {
    console.error('批量操作失败:', error);
    ElMessage.error('批量操作失败');
  }
};

// 处理页码变化
const handlePageChange = (page) => {
  pagination.value.page = page;
  fetchProductList();
};

// 处理每页条数变化
const handleSizeChange = (page, size) => {
  pagination.value.page = page;
  pagination.value.limit = size;
  fetchProductList();
};

// 处理排序变化
const handleSortChange = (sort) => {
  if (sort.prop) {
    pagination.value.sort = sort.prop;
    pagination.value.order = sort.order === 'ascending' ? 'asc' : 'desc';
  } else {
    pagination.value.sort = 'createTime';
    pagination.value.order = 'desc';
  }
  fetchProductList();
};

// 页面初始化
onMounted(() => {
  fetchProductList();
  fetchCategoryList();
});
</script>

<style scoped>
.data-summary {
  display: flex;
  gap: 10px;
}
</style> 