<template>
  <div class="package-edit-container">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑套餐' : '添加套餐' }}</h2>
      <el-button @click="goBack">返回列表</el-button>
    </div>
    
    <el-card class="form-container">
      <el-form
        ref="packageFormRef"
        :model="packageForm"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="套餐名称" prop="name">
          <el-input v-model="packageForm.name" placeholder="请输入套餐名称" />
        </el-form-item>
        
        <el-form-item label="套餐价格" prop="price">
          <el-input-number 
            v-model="packageForm.price"
            :precision="2"
            :step="0.01"
            :min="0"
            controls-position="right"
            style="width: 200px;"
          />
        </el-form-item>
        
        <el-form-item label="有效时长" prop="durationMinutes">
          <el-input-number
            v-model="packageForm.durationMinutes"
            :min="1"
            :step="30"
            controls-position="right"
            style="width: 200px;"
          />
          <span class="unit-text">分钟</span>
        </el-form-item>
        
        <el-form-item label="套餐描述" prop="description">
          <el-input
            v-model="packageForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入套餐描述"
          />
        </el-form-item>
        
        <el-form-item label="状态" prop="isActive">
          <el-switch
            v-model="packageForm.isActive"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        
        <el-divider content-position="left">关联商品</el-divider>
        
        <el-form-item label="选择商品">
          <el-button type="primary" @click="showProductSelector">选择商品</el-button>
          <div class="selected-count" v-if="selectedProducts.length > 0">
            已选择 {{ selectedProducts.length }} 个商品
          </div>
        </el-form-item>
        
        <el-table
          v-if="selectedProducts.length > 0"
          :data="selectedProducts"
          border
          style="width: 100%; margin-bottom: 20px;"
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
          <el-table-column label="操作" width="100">
            <template #default="scope">
              <el-button
                type="danger"
                size="small"
                @click="removeProduct(scope.$index)"
              >移除</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建套餐' }}
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 商品选择对话框 -->
    <el-dialog
      v-model="productSelectorVisible"
      title="选择商品"
      width="800px"
    >
      <div class="product-filter">
        <el-form :inline="true" :model="productQueryParams">
          <el-form-item label="商品名称">
            <el-input v-model="productQueryParams.name" placeholder="请输入商品名称" clearable />
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="productQueryParams.categoryId" placeholder="全部分类" clearable>
              <el-option
                v-for="item in categories"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="searchProducts">查询</el-button>
            <el-button @click="resetProductSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <el-table
        v-loading="productsLoading"
        :data="products"
        border
        style="width: 100%"
        @selection-change="handleSelectionChange"
        ref="productTableRef"
      >
        <el-table-column type="selection" width="55" />
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
      
      <div class="product-pagination">
        <el-pagination
          background
          layout="total, prev, pager, next"
          :total="productTotal"
          :page-size="productQueryParams.pageSize"
          :current-page="productQueryParams.pageNum"
          @current-change="handleProductPageChange"
        />
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="productSelectorVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmSelectProducts">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onActivated, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { createPackage, getPackageDetail, updatePackage, getPackageProducts, associatePackageProducts, removeProductFromPackage } from '../../api/package';
import { getProductList } from '../../api/product';
import { getCategoryList } from '../../api/category';

const router = useRouter();
const route = useRoute();
const packageId = computed(() => route.params.id);
const isEdit = computed(() => !!packageId.value);

// 防止重复请求的锁
const isDetailRequestLocked = ref(false);
const isProductsRequestLocked = ref(false);
const isCategoriesRequestLocked = ref(false);

// 记录页面是否已经初始化
const isInitialized = ref(false);

// 表单引用
const packageFormRef = ref(null);

// 表单数据
const packageForm = reactive({
  name: '',
  price: 0,
  durationMinutes: 120,
  description: '',
  isActive: true
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入套餐名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入套餐价格', trigger: 'blur' }
  ],
  durationMinutes: [
    { required: true, message: '请输入有效时长', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入套餐描述', trigger: 'blur' }
  ]
};

// 提交状态
const submitting = ref(false);

// 选中的商品
const selectedProducts = ref([]);
const selectedProductIds = ref([]);

// 商品选择对话框
const productSelectorVisible = ref(false);
const productsLoading = ref(false);
const products = ref([]);
const productTotal = ref(0);
const tempSelectedProducts = ref([]);
const categories = ref([]);
const productTableRef = ref(null);

// 商品查询参数
const productQueryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  name: '',
  categoryId: null
});

// 返回列表页
const goBack = () => {
  router.push('/package/list');
};

// 获取套餐详情
const getDetail = async () => {
  if (!isEdit.value) return;
  
  // 如果请求已被锁定，则跳过
  if (isDetailRequestLocked.value) {
    console.log('获取套餐详情请求被锁定，跳过此次请求');
    return;
  }
  
  // 锁定请求
  isDetailRequestLocked.value = true;
  
  try {
    const res = await getPackageDetail(packageId.value);
    console.log('套餐详情原始数据:', res.data);
    
    // 处理后端返回的数据
    const packageData = res.data;
    
    // 填充表单 - 适配后端字段名称
    packageForm.name = packageData.name || packageData.package?.name || '';
    packageForm.price = packageData.price || packageData.package?.price || 0;
    packageForm.durationMinutes = packageData.duration_minutes || packageData.durationMinutes || packageData.package?.duration_minutes || 120;
    packageForm.description = packageData.description || packageData.package?.description || '';
    packageForm.isActive = packageData.is_active !== undefined ? packageData.is_active : 
                           packageData.isActive !== undefined ? packageData.isActive : 
                           packageData.package?.is_active !== undefined ? packageData.package.is_active : true;
    
    // 直接从with-products接口获取商品信息
    try {
      let productsData = [];
      
      // 从返回数据中提取商品信息
      if (packageData.products && Array.isArray(packageData.products)) {
        // 新接口返回的数据格式
        productsData = packageData.products;
      } else if (packageData.products_list && Array.isArray(packageData.products_list)) {
        productsData = packageData.products_list;
      }
      
      console.log('从详情中提取的商品数据:', productsData);
      
      if (productsData.length > 0) {
        // 处理商品数据
        selectedProducts.value = productsData.map(product => {
          return {
            id: product.id,
            name: product.name,
            price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
            imageUrl: product.image_url || product.imageUrl || product.image || ''
          };
        });
        
        selectedProductIds.value = selectedProducts.value.map(item => item.id);
      } else {
        // 如果没有找到商品数据，清空选中的商品
        selectedProducts.value = [];
        selectedProductIds.value = [];
      }
    } catch (productError) {
      console.error('处理套餐商品失败:', productError);
      ElMessage.warning('处理套餐商品失败，请手动选择商品');
      selectedProducts.value = [];
      selectedProductIds.value = [];
    }
  } catch (error) {
    console.error('获取套餐详情失败:', error);
    ElMessage.error('获取套餐详情失败');
  } finally {
    // 延迟解锁，防止短时间内重复请求
    setTimeout(() => {
      isDetailRequestLocked.value = false;
    }, 300);
  }
};

// 提交表单
const submitForm = () => {
  packageFormRef.value.validate(async (valid) => {
    if (!valid) return;
    
    if (selectedProducts.value.length === 0) {
      ElMessage.warning('请至少选择一个商品');
      return;
    }
    
    submitting.value = true;
    
    try {
      // 适配后端API字段名称
      const packageData = {
        name: packageForm.name,
        price: packageForm.price,
        duration_minutes: packageForm.durationMinutes, // 使用下划线格式
        description: packageForm.description,
        is_active: packageForm.isActive // 使用下划线格式
      };
      
      console.log('提交的套餐数据:', packageData);
      
      let packageRes;
      
      if (isEdit.value) {
        // 更新套餐
        packageRes = await updatePackage(packageId.value, packageData);
      } else {
        // 创建套餐
        packageRes = await createPackage(packageData);
        
        // 新增模式下，需要关联商品
        const newPackageId = packageRes.data.id;
        await associatePackageProducts(newPackageId, selectedProductIds.value);
      }
      
      ElMessage.success(isEdit.value ? '更新套餐成功' : '创建套餐成功');
      router.push('/package/list');
    } catch (error) {
      console.error(isEdit.value ? '更新套餐失败:' : '创建套餐失败:', error);
      ElMessage.error(isEdit.value ? '更新套餐失败' : '创建套餐失败');
    } finally {
      submitting.value = false;
    }
  });
};

// 重置表单
const resetForm = () => {
  packageFormRef.value.resetFields();
  if (!isEdit.value) {
    selectedProducts.value = [];
    selectedProductIds.value = [];
  } else {
    getDetail();
  }
};

// 显示商品选择器
const showProductSelector = () => {
  productSelectorVisible.value = true;
  tempSelectedProducts.value = [...selectedProducts.value];
  getProducts().then(() => {
    // 在商品加载完成后，设置默认选中的行
    nextTick(() => {
      if (productTableRef.value && selectedProductIds.value.length > 0) {
        // 遍历商品列表，选中已有的商品
        products.value.forEach(product => {
          if (selectedProductIds.value.includes(product.id)) {
            productTableRef.value.toggleRowSelection(product, true);
          }
        });
      }
    });
  });
};

// 获取商品列表
const getProducts = async () => {
  // 如果请求已被锁定，则跳过
  if (isProductsRequestLocked.value) {
    console.log('获取商品列表请求被锁定，跳过此次请求');
    return Promise.resolve();
  }
  
  // 锁定请求
  isProductsRequestLocked.value = true;
  productsLoading.value = true;
  
  try {
    const res = await getProductList(productQueryParams);
    console.log('商品列表原始数据:', res.data);
    
    // 处理商品列表数据
    if (res.data && res.data.list && Array.isArray(res.data.list)) {
      products.value = res.data.list.map(product => {
        return {
          id: product.id,
          name: product.name,
          price: typeof product.price === 'number' ? product.price : parseFloat(product.price || '0'),
          imageUrl: product.image_url || product.imageUrl || product.image || '',
          // 检查是否已选中
          selected: selectedProductIds.value.includes(product.id)
        };
      });
      productTotal.value = res.data.total || products.value.length;
    } else {
      products.value = [];
      productTotal.value = 0;
    }
    return Promise.resolve();
  } catch (error) {
    console.error('获取商品列表失败:', error);
    ElMessage.error('获取商品列表失败');
    return Promise.reject(error);
  } finally {
    productsLoading.value = false;
    // 延迟解锁，防止短时间内重复请求
    setTimeout(() => {
      isProductsRequestLocked.value = false;
    }, 300);
  }
};

// 获取分类列表
const getCategories = async () => {
  // 如果请求已被锁定，则跳过
  if (isCategoriesRequestLocked.value) {
    console.log('获取分类列表请求被锁定，跳过此次请求');
    return;
  }
  
  // 锁定请求
  isCategoriesRequestLocked.value = true;
  
  try {
    const res = await getCategoryList({ pageSize: 100 });
    console.log('分类列表原始数据:', res.data);
    
    // 处理分类数据
    if (res.data && res.data.list && Array.isArray(res.data.list)) {
      categories.value = res.data.list.map(category => {
        return {
          id: category.id,
          name: category.name
        };
      });
    } else if (res.data && typeof res.data === 'object') {
      // 尝试从对象中提取分类数据
      const extractedCategories = [];
      for (const key in res.data) {
        if (!isNaN(parseInt(key)) || key === 'list') {
          const item = res.data[key];
          if (item && typeof item === 'object' && item.name) {
            extractedCategories.push({
              id: item.id || parseInt(key),
              name: item.name
            });
          }
        }
      }
      
      if (extractedCategories.length > 0) {
        categories.value = extractedCategories;
      }
    } else {
      categories.value = [];
    }
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.warning('获取分类列表失败');
  } finally {
    // 延迟解锁，防止短时间内重复请求
    setTimeout(() => {
      isCategoriesRequestLocked.value = false;
    }, 300);
  }
};

// 搜索商品
const searchProducts = () => {
  productQueryParams.pageNum = 1;
  getProducts();
};

// 重置商品搜索
const resetProductSearch = () => {
  productQueryParams.name = '';
  productQueryParams.categoryId = null;
  searchProducts();
};

// 处理商品分页变化
const handleProductPageChange = (page) => {
  productQueryParams.pageNum = page;
  getProducts();
};

// 处理商品选择变化
const handleSelectionChange = (selection) => {
  tempSelectedProducts.value = selection;
};

// 确认选择商品
const confirmSelectProducts = async () => {
  const newSelectedProducts = [...tempSelectedProducts.value];
  const newSelectedProductIds = newSelectedProducts.map(item => item.id);
  
  if (isEdit.value) {
    // 编辑模式下，通过API关联商品
    try {
      await associatePackageProducts(packageId.value, newSelectedProductIds);
      ElMessage.success('商品关联更新成功');
      // 更新本地数据
      selectedProducts.value = newSelectedProducts;
      selectedProductIds.value = newSelectedProductIds;
    } catch (error) {
      console.error('更新关联商品失败:', error);
      ElMessage.error('更新关联商品失败');
    }
  } else {
    // 新增模式下，仅更新本地数据
    selectedProducts.value = newSelectedProducts;
    selectedProductIds.value = newSelectedProductIds;
  }
  
  productSelectorVisible.value = false;
};

// 移除商品
const removeProduct = (index) => {
  const product = selectedProducts.value[index];
  
  if (isEdit.value) {
    // 编辑模式下，通过API移除商品
    ElMessageBox.confirm(`确定要将商品"${product.name}"从套餐中移除吗?`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await removeProductFromPackage(packageId.value, product.id);
        ElMessage.success('商品移除成功');
        // 更新本地数据
        selectedProducts.value.splice(index, 1);
        selectedProductIds.value = selectedProducts.value.map(item => item.id);
      } catch (error) {
        console.error('移除商品失败:', error);
        ElMessage.error('移除商品失败');
      }
    }).catch(() => {});
  } else {
    // 新增模式下，直接从本地数据中移除
    selectedProducts.value.splice(index, 1);
    selectedProductIds.value = selectedProducts.value.map(item => item.id);
  }
};

// 初始化
onMounted(() => {
  console.log('套餐编辑页面挂载');
  isInitialized.value = true;
  getDetail();
  getCategories();
});

// 当页面被激活时（从缓存中恢复）重新获取数据
onActivated(() => {
  console.log('套餐编辑页面激活');
  if (isInitialized.value) {
    if (isEdit.value) {
      getDetail();
    }
    getCategories();
  }
});
</script>

<style scoped>
.package-edit-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-container {
  margin-bottom: 20px;
}

.unit-text {
  margin-left: 10px;
  color: #606266;
}

.selected-count {
  display: inline-block;
  margin-left: 20px;
  color: #409EFF;
}

.product-filter {
  margin-bottom: 20px;
}

.product-pagination {
  margin-top: 20px;
  text-align: center;
}
</style> 