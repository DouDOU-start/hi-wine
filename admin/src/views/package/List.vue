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
            {{ formatDuration(scope.row.durationMinutes || scope.row.duration_minutes || 0) }}
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
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="图片" width="100">
          <template #default="scope">
            <el-image
              v-if="scope.row.imageUrl"
              style="width: 50px; height: 50px"
              :src="scope.row.imageUrl"
              :preview-src-list="[scope.row.imageUrl]"
              fit="cover"
              @error="() => { scope.row.imageLoadError = true }"
            >
              <template #error>
                <div style="width: 50px; height: 50px; display: flex; align-items: center; justify-content: center; background: #f5f7fa; color: #909399; border-radius: 4px;">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div 
              v-else 
              style="width: 50px; height: 50px; display: flex; align-items: center; justify-content: center; background: #f5f7fa; color: #909399; border-radius: 4px;"
            >
              <el-icon><Picture /></el-icon>
            </div>
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
              v-if="scope.row.imageUrl"
              style="width: 50px; height: 50px"
              :src="scope.row.imageUrl"
              :preview-src-list="[scope.row.imageUrl]"
              fit="cover"
              @error="() => { scope.row.imageLoadError = true }"
            >
              <template #error>
                <div style="width: 50px; height: 50px; display: flex; align-items: center; justify-content: center; background: #f5f7fa; color: #909399; border-radius: 4px;">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div 
              v-else 
              style="width: 50px; height: 50px; display: flex; align-items: center; justify-content: center; background: #f5f7fa; color: #909399; border-radius: 4px;"
            >
              <el-icon><Picture /></el-icon>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getPackageList, deletePackage, getPackageProducts } from '../../api/package';
import { Picture } from '@element-plus/icons-vue';

const router = useRouter();

// 防止重复请求的锁
const isRequestLocked = ref(false);

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

// 记录页面是否已经初始化
const isInitialized = ref(false);

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '暂无数据';
  
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return dateString; // 如果转换失败，返回原始字符串
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  } catch (error) {
    console.error('日期格式化错误:', error);
    return dateString; // 发生错误时返回原始字符串
  }
};

// 格式化时长
const formatDuration = (minutes) => {
  if (!minutes || isNaN(minutes)) return '0 分钟';
  
  // 将分钟数转换为数字类型
  const mins = typeof minutes === 'string' ? parseFloat(minutes) : minutes;
  
  if (mins < 60) {
    // 不足1小时，显示分钟
    return `${mins} 分钟`;
  } else {
    // 超过1小时，显示小时和分钟
    const hours = Math.floor(mins / 60);
    const remainingMins = mins % 60;
    
    if (remainingMins === 0) {
      // 整数小时
      return `${hours} 小时`;
    } else {
      // 小时+分钟
      return `${hours} 小时 ${remainingMins} 分钟`;
    }
  }
};

// 获取套餐列表
const getList = async () => {
  // 如果已经在加载中或请求被锁定，则跳过
  if (loading.value || isRequestLocked.value) {
    console.log('请求被锁定或正在加载中，跳过此次请求');
    return;
  }
  
  // 锁定请求，防止短时间内重复调用
  isRequestLocked.value = true;
  loading.value = true;
  
  try {
    const res = await getPackageList(queryParams);
    console.log('套餐列表原始数据:', res.data);
    
    // 处理返回的数据
    if (res.data && res.data.list && Array.isArray(res.data.list)) {
      packageList.value = res.data.list.map(item => {
        // 处理有效时长
        let durationMins = 0;
        if (item.duration_minutes !== undefined) {
          durationMins = item.duration_minutes;
        } else if (item.durationMinutes !== undefined) {
          durationMins = item.durationMinutes;
        } else if (item.duration_hours !== undefined) {
          durationMins = item.duration_hours * 60;
        }
        
        return {
          ...item,
          // 确保isActive字段存在（兼容is_active和isActive两种命名）
          isActive: item.isActive !== undefined ? item.isActive : (item.is_active !== undefined ? item.is_active : true),
          // 确保价格是数字类型
          price: typeof item.price === 'string' ? parseFloat(item.price) : item.price,
          // 确保创建时间和更新时间字段存在
          created_at: item.created_at || item.createdAt || '',
          updated_at: item.updated_at || item.updatedAt || '',
          // 设置处理后的有效时长
          durationMinutes: durationMins
        };
      });
      total.value = res.data.total || packageList.value.length;
    } else {
      packageList.value = [];
      total.value = 0;
    }
  } catch (error) {
    console.error('获取套餐列表失败:', error);
    ElMessage.error('获取套餐列表失败');
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
    console.log('获取套餐商品，ID:', row.id);
    const res = await getPackageProducts(row.id);
    console.log('套餐商品原始数据:', res.data);
    
    // 处理返回的数据结构
    if (res.data.products && Array.isArray(res.data.products)) {
      // 标准返回格式
      packageProducts.value = res.data.products.map(product => ({
        id: product.id,
        name: product.name,
        price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
        imageUrl: product.image_url || product.imageUrl || product.image || ''
      }));
    } else if (res.data.list && Array.isArray(res.data.list)) {
      // 列表格式
      packageProducts.value = res.data.list.map(product => ({
        id: product.id,
        name: product.name,
        price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
        imageUrl: product.image_url || product.imageUrl || product.image || ''
      }));
    } else if (Array.isArray(res.data)) {
      // 直接返回数组
      packageProducts.value = res.data.map(product => ({
        id: product.id,
        name: product.name,
        price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
        imageUrl: product.image_url || product.imageUrl || product.image || ''
      }));
    } else {
      // 尝试从对象中提取数据
      const extractedProducts = [];
      for (const key in res.data) {
        if (res.data[key] && typeof res.data[key] === 'object' && !Array.isArray(res.data[key])) {
          const product = res.data[key];
          if (product.id && product.name) {
            extractedProducts.push({
              id: product.id,
              name: product.name,
              price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
              imageUrl: product.image_url || product.imageUrl || product.image || ''
            });
          }
        }
      }
      
      if (extractedProducts.length > 0) {
        packageProducts.value = extractedProducts;
      } else {
        packageProducts.value = [];
        console.warn('未找到有效的商品数据');
      }
    }
    
    console.log('处理后的商品数据:', packageProducts.value);
    
    if (packageProducts.value.length === 0) {
      ElMessage.info('该套餐暂无关联商品');
    }
  } catch (error) {
    console.error('获取套餐商品失败:', error);
    ElMessage.error('获取套餐商品失败');
    packageProducts.value = [];
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
  console.log('套餐列表页面挂载');
  isInitialized.value = true;
  getList();
});

// 当页面被激活时（从缓存中恢复）重新加载数据
onActivated(() => {
  console.log('套餐列表页面激活');
  if (isInitialized.value) {
    getList();
  }
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