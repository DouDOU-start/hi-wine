<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">商品管理</div>
      <div class="data-summary">
        <el-tag type="info">总商品: {{ total }}</el-tag>
        <el-tag v-if="lowStockCount > 0" type="warning">库存不足: {{ lowStockCount }}</el-tag>
        <el-tag v-if="outOfStockCount > 0" type="danger">缺货: {{ outOfStockCount }}</el-tag>
      </div>
    </div>
    
    <el-card shadow="hover" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form" size="default">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
            <el-form-item label="商品名称">
              <el-input v-model="searchForm.name" placeholder="请输入商品名称" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
            <el-form-item label="分类">
              <el-select 
                v-model="searchForm.categoryId" 
                placeholder="请选择分类" 
                clearable
                filterable
                popper-class="category-select-dropdown"
                :popper-options="{ boundariesPadding: 0, gpuAcceleration: false }"
              >
                <template #prefix>
                  <el-icon class="category-icon"><Grid /></el-icon>
                </template>
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                  <div class="category-option">
                    <el-icon class="category-child-icon"><Document /></el-icon>
                    <span>{{ item.name }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
            <el-form-item label="状态">
              <el-select 
                v-model="searchForm.status" 
                placeholder="请选择状态" 
                clearable
                popper-class="status-select-dropdown"
              >
                <template #prefix>
                  <el-icon class="status-icon"><SetUp /></el-icon>
                </template>
                <el-option-group label="商品状态">
                  <el-option label="上架" :value="1">
                    <div class="status-option">
                      <el-tag size="small" type="success" effect="dark" class="status-tag-select">上架</el-tag>
                      <span class="status-desc">商品可见且可购买</span>
                    </div>
                  </el-option>
                  <el-option label="下架" :value="0">
                    <div class="status-option">
                      <el-tag size="small" type="info" effect="dark" class="status-tag-select">下架</el-tag>
                      <span class="status-desc">商品不可见且不可购买</span>
                    </div>
                  </el-option>
                </el-option-group>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="4" class="search-buttons">
            <el-form-item>
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>查询
              </el-button>
              <el-button @click="resetSearch">
                <el-icon><Refresh /></el-icon>重置
              </el-button>
              <el-button type="text" @click="toggleAdvanced">
                {{ showAdvanced ? '收起' : '高级筛选' }}
                <el-icon>
                  <component :is="showAdvanced ? 'ArrowUp' : 'ArrowDown'" />
                </el-icon>
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 高级搜索区域 -->
        <el-collapse-transition>
          <div v-show="showAdvanced">
            <el-row :gutter="20">
              <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
                <el-form-item label="价格区间">
                  <el-input-number v-model="searchForm.minPrice" placeholder="最低价" :min="0" :precision="2" style="width: 120px;" />
                  <span class="price-separator">至</span>
                  <el-input-number v-model="searchForm.maxPrice" placeholder="最高价" :min="0" :precision="2" style="width: 120px;" />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
                <el-form-item label="库存">
                  <el-select 
                    v-model="searchForm.stockStatus" 
                    placeholder="库存状态" 
                    clearable
                    popper-class="stock-select-dropdown"
                  >
                    <template #prefix>
                      <el-icon class="stock-icon"><Goods /></el-icon>
                    </template>
                    <el-option-group label="库存状态">
                      <el-option label="充足" :value="'normal'">
                        <div class="stock-option">
                          <div class="stock-indicator">
                            <el-icon class="stock-normal-icon"><CircleCheckFilled /></el-icon>
                            <el-tag size="small" type="success" effect="plain" class="stock-tag-select">充足</el-tag>
                          </div>
                          <span class="stock-desc">库存 > 10</span>
                        </div>
                      </el-option>
                      <el-option label="不足" :value="'low'">
                        <div class="stock-option">
                          <div class="stock-indicator">
                            <el-icon class="stock-low-icon"><WarningFilled /></el-icon>
                            <el-tag size="small" type="warning" effect="plain" class="stock-tag-select">不足</el-tag>
                          </div>
                          <span class="stock-desc">0 < 库存 ≤ 10</span>
                        </div>
                      </el-option>
                      <el-option label="缺货" :value="'out'">
                        <div class="stock-option">
                          <div class="stock-indicator">
                            <el-icon class="stock-out-icon"><CircleCloseFilled /></el-icon>
                            <el-tag size="small" type="danger" effect="plain" class="stock-tag-select">缺货</el-tag>
                          </div>
                          <span class="stock-desc">库存 = 0</span>
                        </div>
                      </el-option>
                    </el-option-group>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8" :md="6" :lg="6" :xl="4">
                <el-form-item label="创建时间">
                  <el-date-picker
                    v-model="searchForm.dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="YYYY-MM-DD"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </el-collapse-transition>
      </el-form>
      
      <!-- 快速筛选标签 -->
      <div class="quick-filter">
        <span class="quick-filter-label">快速筛选:</span>
        <div class="filter-button-group">
          <el-button 
            :class="['filter-button', { 'active': activeFilter === 'all' }]" 
            size="small" 
            @click="applyQuickFilter('all')"
          >
            <el-icon><Tickets /></el-icon>全部商品
          </el-button>
          <el-button 
            :class="['filter-button', { 'active': activeFilter === 'lowStock' }]" 
            size="small" 
            @click="applyQuickFilter('lowStock')"
            type="warning"
            plain
          >
            <el-icon><WarningFilled /></el-icon>库存不足
            <el-badge v-if="lowStockCount > 0" :value="lowStockCount" class="filter-badge" />
          </el-button>
          <el-button 
            :class="['filter-button', { 'active': activeFilter === 'outOfStock' }]" 
            size="small" 
            @click="applyQuickFilter('outOfStock')"
            type="danger"
            plain
          >
            <el-icon><RemoveFilled /></el-icon>缺货商品
            <el-badge v-if="outOfStockCount > 0" :value="outOfStockCount" class="filter-badge" />
          </el-button>
          <el-button 
            :class="['filter-button', { 'active': activeFilter === 'onSale' }]" 
            size="small" 
            @click="applyQuickFilter('onSale')"
            type="success"
            plain
          >
            <el-icon><Select /></el-icon>已上架
          </el-button>
          <el-button 
            :class="['filter-button', { 'active': activeFilter === 'offSale' }]" 
            size="small" 
            @click="applyQuickFilter('offSale')"
            type="info"
            plain
          >
            <el-icon><TurnOff /></el-icon>已下架
          </el-button>
        </div>
      </div>
    </el-card>
    
    <el-card shadow="hover" class="table-card">
      <div class="table-header">
        <div class="left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>添加商品
          </el-button>
          <el-dropdown v-if="selectedIds.length > 0" @command="handleBatchCommand">
            <el-button type="primary" plain>
              批量操作<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="delete">批量删除</el-dropdown-item>
                <el-dropdown-item command="onSale">批量上架</el-dropdown-item>
                <el-dropdown-item command="offSale">批量下架</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button v-else disabled>批量操作</el-button>
        </div>
        <div class="right">
          <el-switch
            v-model="tableConfig.showImage"
            active-text="显示图片"
            inactive-text="隐藏图片"
          />
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="productList"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
        :header-cell-style="{ background: '#f5f7fa' }"
      >
        <el-table-column type="selection" width="55" fixed="left" />
        <el-table-column prop="id" label="ID" width="80" sortable />
        <el-table-column prop="image" label="图片" width="100" v-if="tableConfig.showImage">
          <template #default="scope">
            <el-image 
              v-if="scope.row.image_url || scope.row.image" 
              :src="getImageUrl(scope.row.image_url || scope.row.image)" 
              style="width: 60px; height: 60px; object-fit: cover; border-radius: 6px;"
              fit="cover"
              :preview-src-list="[getImageUrl(scope.row.image_url || scope.row.image)]"
              :initial-index="0"
              lazy
              @error="handleImageError(scope.row)"
              preview-teleported
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div v-else class="no-image" style="width: 60px; height: 60px; display: flex; align-items: center; justify-content: center; background: #f5f7fa; color: #909399; border-radius: 6px;">
              <el-icon><PictureFilled /></el-icon>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" min-width="200" show-overflow-tooltip sortable>
          <template #default="scope">
            <div class="product-name">
              {{ scope.row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="categoryName" label="分类" width="120" />
        <el-table-column prop="price" label="价格" width="100" sortable>
          <template #default="scope">
            <span class="price">￥{{ typeof scope.row.price === 'number' ? scope.row.price.toFixed(2) : scope.row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="140" sortable>
          <template #default="scope">
            <div class="stock-cell">
              <div class="stock-info">
                <span :class="{ 'stock-warning': scope.row.stock <= 10 && scope.row.stock > 0, 'stock-danger': scope.row.stock === 0 }">
                  {{ scope.row.stock }}
                </span>
                <el-progress 
                  :percentage="Math.min(100, (scope.row.stock / 100) * 100)" 
                  :stroke-width="6" 
                  :status="scope.row.stock === 0 ? 'exception' : scope.row.stock <= 10 ? 'warning' : 'success'"
                  class="stock-progress"
                />
              </div>
              <el-tag 
                v-if="scope.row.stock <= 10 && scope.row.stock > 0" 
                size="small" 
                type="warning"
                effect="plain"
                class="stock-tag"
              >
                不足
              </el-tag>
              <el-tag 
                v-else-if="scope.row.stock === 0" 
                size="small" 
                type="danger"
                effect="plain"
                class="stock-tag"
              >
                缺货
              </el-tag>
              <el-tag 
                v-else 
                size="small" 
                type="success"
                effect="plain"
                class="stock-tag"
              >
                充足
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="sales" label="销量" width="100" sortable />
        <el-table-column prop="status" label="状态" width="140">
          <template #default="scope">
            <div class="status-cell">
              <el-switch
                v-model="scope.row.status"
                :active-value="1"
                :inactive-value="0"
                @change="() => { scope.row.userTriggered = true; handleStatusChange(scope.row); }"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #909399;"
                inline-prompt
                active-text="上"
                inactive-text="下"
                class="status-switch"
              />
              <el-tag 
                :type="scope.row.status === 1 ? 'success' : 'info'" 
                effect="plain"
                size="small"
                class="status-tag"
              >
                <el-icon v-if="scope.row.status === 1"><CircleCheckFilled /></el-icon>
                <el-icon v-else><CircleCloseFilled /></el-icon>
                {{ scope.row.status === 1 ? '已上架' : '已下架' }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              @click="handleEdit(scope.row)"
            >
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(scope.row)"
            >
              <el-icon><Delete /></el-icon>删除
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
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, onBeforeMount, onActivated } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getProductList, deleteProduct, updateProductStatus, batchUpdateProductStatus } from '../../api/product';
import { getCategoryList } from '../../api/category';
import { 
  Search, Refresh, ArrowDown, ArrowUp, Plus, Edit, Delete, 
  Picture, PictureFilled, Grid, Folder, SetUp, Goods, 
  Tickets, WarningFilled, RemoveFilled, Select, TurnOff, 
  Document, CircleCheckFilled, CircleCloseFilled 
} from '@element-plus/icons-vue';

const router = useRouter();

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '暂无数据';
  
  try {
    const date = new Date(dateTimeStr);
    if (isNaN(date.getTime())) return dateTimeStr; // 如果转换失败，返回原始字符串
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  } catch (error) {
    console.error('日期格式化错误:', error);
    return dateTimeStr; // 发生错误时返回原始字符串
  }
};

// 加载状态
const loading = ref(false);
// 防止重复请求的锁
const isRequestLocked = ref(false);

// 商品列表数据
const productList = ref([]);

// 分类选项
const categoryOptions = ref([]);

// 分类分组
const categoryGroups = ref([]);

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 高级搜索显示状态
const showAdvanced = ref(false);

// 当前激活的快速筛选
const activeFilter = ref('all');

// 表格配置
const tableConfig = reactive({
  showImage: true
});

// 搜索表单
const searchForm = reactive({
  name: '',
  categoryId: '',
  status: '',
  minPrice: null,
  maxPrice: null,
  stockStatus: '',
  dateRange: []
});

// 选中的商品ID
const selectedIds = ref([]);

// 库存不足和缺货商品数量
const lowStockCount = computed(() => {
  return productList.value.filter(item => item.stock <= 10 && item.stock > 0).length;
});

const outOfStockCount = computed(() => {
  return productList.value.filter(item => item.stock === 0).length;
});

// 切换高级搜索显示状态
const toggleAdvanced = () => {
  showAdvanced.value = !showAdvanced.value;
};

// 应用快速筛选
const applyQuickFilter = (filter) => {
  // 设置当前激活的筛选
  activeFilter.value = filter;
  
  // 重置搜索表单
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = key === 'dateRange' ? [] : '';
  });
  searchForm.minPrice = null;
  searchForm.maxPrice = null;
  
  // 应用对应的筛选条件
  switch (filter) {
    case 'all':
      // 不需要设置任何条件，使用默认值
      break;
    case 'lowStock':
      searchForm.stockStatus = 'low';
      break;
    case 'outOfStock':
      searchForm.stockStatus = 'out';
      break;
    case 'onSale':
      searchForm.status = 1;
      break;
    case 'offSale':
      searchForm.status = 0;
      break;
  }
  
  // 执行搜索
  currentPage.value = 1;
  fetchProductList();
};

// 批量操作命令处理
const handleBatchCommand = (command) => {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择商品');
    return;
  }
  
  switch (command) {
    case 'delete':
      handleBatchDelete();
      break;
    case 'onSale':
      handleBatchStatusChange(1);
      break;
    case 'offSale':
      handleBatchStatusChange(0);
      break;
  }
};

// 批量修改商品状态
const handleBatchStatusChange = (status) => {
  const statusText = status === 1 ? '上架' : '下架';
  
  ElMessageBox.confirm(`确定要批量${statusText}选中的${selectedIds.value.length}个商品吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      // 调用批量更新状态API
      await batchUpdateProductStatus(selectedIds.value, status);
      ElMessage.success(`批量${statusText}成功`);
      fetchProductList();
    } catch (error) {
      console.error(`批量${statusText}失败:`, error);
      ElMessage.error(`批量${statusText}失败`);
    }
  }).catch(() => {});
};

// 获取商品列表
const fetchProductList = async () => {
  // 如果已经在加载中或请求被锁定，则跳过
  if (loading.value || isRequestLocked.value) {
    return;
  }
  
  // 锁定请求，防止短时间内重复调用
  isRequestLocked.value = true;
  loading.value = true;
  
  try {
    // 构建查询参数
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      name: searchForm.name,
      categoryId: searchForm.categoryId,
      status: searchForm.status,
      minPrice: searchForm.minPrice,
      maxPrice: searchForm.maxPrice,
      stockStatus: searchForm.stockStatus
    };
    
    // 如果有日期范围，添加到查询参数
    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.startDate = searchForm.dateRange[0];
      params.endDate = searchForm.dateRange[1];
    }
    
    console.log('获取商品列表，参数:', params);
    const response = await getProductList(params);
    console.log('商品列表响应:', response);
    
    // 处理返回的数据，确保每个商品有分类名称
    if (response.data && response.data.list) {
      productList.value = response.data.list.map(item => {
        // 查找对应的分类名称
        const category = categoryOptions.value.find(cat => cat.id === item.category_id);
        
        // 调试日志
        console.log('处理商品数据:', item);
        
        return {
          ...item,
          categoryName: category ? category.name : '未分类',
          // 处理状态字段，后端返回的可能是status或is_active
          status: item.status !== undefined ? item.status : (item.is_active === true ? 1 : 0),
          // 确保价格是数字类型
          price: typeof item.price === 'string' ? parseFloat(item.price) : item.price,
          // 确保图片字段存在
          image: item.image_url || item.image,
          // 确保创建时间和更新时间字段存在
          created_at: item.created_at || '暂无数据',
          updated_at: item.updated_at || '暂无数据'
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
    // 延迟解锁，防止短时间内重复请求
    setTimeout(() => {
      isRequestLocked.value = false;
    }, 300);
  }
};

// 获取分类列表
const fetchCategoryList = async () => {
  try {
    console.log('获取分类列表...');
    const response = await getCategoryList();
    console.log('分类列表原始响应:', response);
    
    // 直接输出原始数据结构
    console.log('分类数据结构:', JSON.stringify(response, null, 2));
    
    // 根据截图显示的数据结构进行处理
    if (response.data) {
      // 从截图可以看出，返回的是包含list数组的对象
      if (response.data.list && Array.isArray(response.data.list)) {
        // 处理标准的list数组结构
        categoryOptions.value = response.data.list.map(item => ({
          id: item.id,
          name: item.name
        }));
      } else {
        // 从截图看，返回的是一个对象，其中包含数字索引的分类对象
        const extractedCategories = [];
        
        // 检查是否存在直接的list对象（非数组）
        if (response.data.list && typeof response.data.list === 'object') {
          console.log('发现list对象结构');
          const listData = response.data.list;
          
          // 遍历对象中的每个键
          for (const key in listData) {
            if (!isNaN(parseInt(key))) {
              const item = listData[key];
              if (item && typeof item === 'object') {
                extractedCategories.push({
                  id: item.id || parseInt(key),
                  name: item.name || `分类${key}`
                });
              }
            }
          }
        } else {
          // 尝试直接从data中提取
          for (const key in response.data) {
            if (!isNaN(parseInt(key))) {
              const item = response.data[key];
              if (item && typeof item === 'object') {
                extractedCategories.push({
                  id: item.id || parseInt(key),
                  name: item.name || `分类${key}`
                });
              }
            }
          }
        }
        
        if (extractedCategories.length > 0) {
          categoryOptions.value = extractedCategories;
          console.log('提取的分类数据:', categoryOptions.value);
        } else {
          // 如果没有提取到数据，尝试使用调试数据
          console.log('未能提取到分类数据，使用硬编码数据');
          categoryOptions.value = [
            { id: 1, name: '短饮' },
            { id: 2, name: '长饮' }
          ];
        }
      }
      
      console.log('最终分类选项:', categoryOptions.value);
    }
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.error('获取分类列表失败');
    
    // 出错时使用默认分类数据
    categoryOptions.value = [
      { id: 1, name: '短饮' },
      { id: 2, name: '长饮' }
    ];
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
    searchForm[key] = key === 'dateRange' ? [] : '';
  });
  searchForm.minPrice = null;
  searchForm.maxPrice = null;
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
  // 在跳转前设置一个标志，表示需要清空编辑页面的缓存数据
  localStorage.setItem('product_form_reset', 'true');
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
  const originalStatus = row.status;
  try {
    console.log('更新商品状态开始:', row.id, '新状态:', row.status);
    // 静默处理，不显示成功消息
    const response = await updateProductStatus(row.id, row.status);
    console.log('状态更新成功，响应:', response);
    
    // 只有用户手动操作时才显示消息
    if (row.userTriggered) {
      ElMessage.success(`商品已${row.status === 1 ? '上架' : '下架'}`);
      row.userTriggered = false;
    }
  } catch (error) {
    console.error('更新商品状态失败:', error);
    // 只有用户手动操作时才显示错误消息
    if (row.userTriggered) {
      ElMessage.error(`更新商品状态失败: ${error.message || '未知错误'}`);
      row.userTriggered = false;
    }
    // 恢复原状态
    row.status = originalStatus === 1 ? 0 : 1;
  }
};

// 图片错误处理
const handleImageError = (row) => {
  console.error(`图片加载失败: ${row.name}`, row.image_url || row.image);
  // 避免多次提示，不显示消息提示
};

// 获取完整图片URL
const getImageUrl = (url) => {
  if (!url) return '';
  
  // 如果是相对路径，添加基础URL
  if (url.startsWith('/')) {
    // 使用当前域名作为基础URL，而不是硬编码
    // 这样可以适应不同的部署环境
    return url;
  }
  
  // 如果已经是完整URL，直接返回
  return url;
};

// 初始化
onMounted(() => {
  console.log('商品列表页面已挂载');
  // 获取分类数据
  fetchCategoryList();
  // 获取商品列表
  fetchProductList();
});

// 添加onActivated生命周期钩子，当页面从缓存中激活时重新加载数据
onActivated(() => {
  console.log('商品列表页面已激活');
  // 重新加载数据
  fetchProductList();
});
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.data-summary {
  display: flex;
  gap: 10px;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin-bottom: 10px;
}

.search-buttons {
  display: flex;
  justify-content: flex-start;
}

.quick-filter {
  margin-top: 15px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-filter-label {
  font-size: 14px;
  color: #606266;
  margin-right: 8px;
}

.price-separator {
  margin: 0 5px;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.product-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.price {
  font-weight: bold;
  color: #ff6b6b;
}

.stock-warning {
  color: #e6a23c;
  font-weight: bold;
}

.stock-danger {
  color: #f56c6c;
  font-weight: bold;
}

.image-error, .no-image {
  width: 60px;
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  color: #909399;
  font-size: 24px;
}

/* 分类和状态选择器增强样式 */
.category-option,
.status-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.category-option {
  padding: 6px 0;
}

.category-icon,
.status-icon,
.stock-icon {
  color: #409EFF;
}

.stock-normal-icon {
  color: #67C23A;
}

.stock-low-icon {
  color: #E6A23C;
}

.stock-out-icon {
  color: #F56C6C;
}

.category-child-icon {
  color: #909399;
  font-size: 14px;
}

.group-label {
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #606266;
  background-color: #f5f7fa;
  display: flex;
  align-items: center;
  gap: 6px;
  border-radius: 4px;
  margin: 4px 0;
}

.item-count {
  margin-left: auto;
  font-size: 11px;
  border-radius: 10px;
  padding: 0 6px;
}

.status-option,
.stock-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
}

.stock-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-desc,
.stock-desc {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}

.status-tag-select,
.stock-tag-select {
  width: 48px;
  text-align: center;
  padding: 0 6px;
}

:deep(.category-select-dropdown) {
  max-height: 300px;
}

:deep(.status-select-dropdown),
:deep(.stock-select-dropdown) {
  min-width: 200px !important;
}

:deep(.el-select-dropdown__item) {
  padding: 0 12px;
  height: auto;
}

:deep(.el-select-group__wrap:not(:last-of-type)) {
  margin-bottom: 8px;
}

:deep(.el-select-group__title) {
  padding: 8px 12px;
  font-size: 14px;
  font-weight: bold;
  color: #606266;
  background-color: #f5f7fa;
  margin-bottom: 4px;
}

:deep(.el-select .el-input__prefix) {
  display: flex;
  align-items: center;
  color: #409EFF;
}

/* 表格中状态和库存列样式 */
.stock-cell,
.status-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stock-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.stock-progress {
  margin-top: 4px;
  width: 100%;
}

.stock-cell {
  justify-content: space-between;
}

.status-cell {
  justify-content: flex-start;
}

.stock-tag,
.status-tag {
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-switch {
  margin-right: 8px;
}

/* 快速筛选按钮样式 */
.filter-button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.filter-button {
  position: relative;
  transition: all 0.3s;
}

.filter-button.active {
  transform: translateY(-2px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.filter-button .el-icon {
  margin-right: 4px;
}

.filter-badge {
  margin-left: 4px;
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .search-buttons {
    justify-content: flex-start;
    margin-top: 10px;
  }
  
  .quick-filter {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .quick-filter-label {
    margin-bottom: 8px;
  }
}
</style> 