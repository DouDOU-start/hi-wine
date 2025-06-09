<template>
  <div class="page-container">
    <!-- 页面头部 - 使用更现代的设计 -->
    <div class="page-header">
      <div class="page-title-container">
        <el-icon class="page-icon"><Goods /></el-icon>
        <div class="page-title-content">
          <h2 class="page-title">{{ isEdit ? '编辑商品' : '添加商品' }}</h2>
          <div class="page-subtitle">{{ isEdit ? '修改现有商品信息' : '创建新的商品' }}</div>
        </div>
      </div>
      <el-button type="primary" plain @click="goBack">
        <el-icon><Back /></el-icon>返回列表
      </el-button>
    </div>
    
    <!-- 步骤指示器 - 添加视觉引导 -->
    <el-steps :active="1" finish-status="success" simple class="steps-bar">
      <el-step title="填写信息" />
      <el-step title="预览确认" />
      <el-step title="提交保存" />
    </el-steps>
    
    <!-- 主表单 - 使用卡片组布局 -->
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="product-form"
      :disabled="submitting"
    >
      <el-row :gutter="24">
        <!-- 左侧表单 -->
        <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
          <!-- 基本信息卡片 -->
          <el-card shadow="hover" class="form-card">
            <template #header>
              <div class="card-header">
                <div class="card-title">
                  <el-icon class="card-icon"><InfoFilled /></el-icon>
                  <span>基本信息</span>
                </div>
                <el-tag size="small" effect="plain" type="info">必填</el-tag>
              </div>
            </template>
            
            <el-form-item label="商品名称" prop="name">
              <el-input 
                v-model="form.name" 
                placeholder="请输入商品名称" 
                maxlength="50" 
                show-word-limit
                prefix-icon="Document"
              />
            </el-form-item>
            
            <el-form-item label="商品分类" prop="categoryId">
              <el-select 
                v-model="form.categoryId" 
                placeholder="请选择商品分类" 
                style="width: 100%"
                filterable
              >
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                  <div class="category-option">
                    <el-icon><Folder /></el-icon>
                    <span>{{ item.name }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="商品价格" prop="price">
                  <el-input-number
                    v-model="form.price"
                    :precision="2"
                    :step="0.1"
                    :min="0"
                    style="width: 100%;"
                    controls-position="right"
                  >
                    <template #prefix>¥</template>
                  </el-input-number>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="商品库存" prop="stock">
                  <el-input-number
                    v-model="form.stock"
                    :min="0"
                    :precision="0"
                    style="width: 100%;"
                    controls-position="right"
                  >
                    <template #prefix>
                      <el-icon><Goods /></el-icon>
                    </template>
                  </el-input-number>
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="商品状态" prop="status">
              <el-switch
                v-model="form.status"
                :active-value="1"
                :inactive-value="0"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949;"
                inline-prompt
                active-text="上架"
                inactive-text="下架"
                class="status-switch"
              />
              <span class="status-text">{{ form.status === 1 ? '商品将在保存后立即上架销售' : '商品将保存为下架状态' }}</span>
            </el-form-item>
          </el-card>
          
          <!-- 商品描述卡片 -->
          <el-card shadow="hover" class="form-card">
            <template #header>
              <div class="card-header">
                <div class="card-title">
                  <el-icon class="card-icon"><Document /></el-icon>
                  <span>商品描述</span>
                </div>
                <el-tooltip content="详细的商品描述有助于提高销售转化率" placement="top">
                  <el-icon class="help-icon"><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            
            <el-form-item label="商品描述" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="8"
                placeholder="请输入商品描述，包括特点、口感、产地等信息"
                maxlength="500"
                show-word-limit
                resize="none"
                class="description-textarea"
              />
            </el-form-item>
          </el-card>
        </el-col>
        
        <!-- 右侧图片上传 -->
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-card shadow="hover" class="form-card image-card">
            <template #header>
              <div class="card-header">
                <div class="card-title">
                  <el-icon class="card-icon"><Picture /></el-icon>
                  <span>商品图片</span>
                </div>
                <el-tag size="small" effect="plain" type="info">必填</el-tag>
              </div>
            </template>
            
            <el-form-item prop="image" class="image-form-item">
              <div class="upload-container">
                <el-upload
                  class="product-image-uploader"
                  action="/api/admin/upload"
                  :show-file-list="false"
                  :on-success="handleUploadSuccess"
                  :before-upload="beforeUpload"
                  :headers="uploadHeaders"
                  :disabled="submitting"
                  :on-error="handleUploadError"
                  drag
                >
                  <div class="upload-area" :class="{'has-image': form.image}">
                    <template v-if="form.image">
                      <img :src="getImageUrl(form.image)" class="product-image" />
                      <div class="image-overlay">
                        <el-icon class="upload-icon"><RefreshRight /></el-icon>
                        <div class="upload-text">点击更换图片</div>
                      </div>
                    </template>
                    <template v-else>
                      <div class="upload-placeholder">
                        <el-icon class="upload-icon"><Upload /></el-icon>
                        <div class="upload-text">点击或拖拽图片上传</div>
                      </div>
                    </template>
                  </div>
                </el-upload>
              </div>
              
              <div class="upload-tips">
                <el-alert
                  title="图片上传要求"
                  type="info"
                  :closable="false"
                  class="upload-alert"
                >
                  <template #default>
                    <ul class="upload-requirements">
                      <li><el-icon><Check /></el-icon> 建议上传正方形图片，大小不超过5MB</li>
                      <li><el-icon><Check /></el-icon> 支持JPG、PNG、GIF格式</li>
                      <li><el-icon><Check /></el-icon> 图片清晰度不低于300x300像素</li>
                    </ul>
                  </template>
                </el-alert>
              </div>
            </el-form-item>
            
            <!-- 商品预览卡片 -->
            <div class="preview-section" v-if="form.name || form.image">
              <div class="preview-title">
                <el-icon><View /></el-icon>
                <span>商品预览</span>
              </div>
              <div class="product-preview-card">
                <div class="preview-image-container">
                  <img v-if="form.image" :src="getImageUrl(form.image)" class="preview-image" />
                  <div v-else class="no-image">暂无图片</div>
                </div>
                <div class="preview-info">
                  <div class="preview-name">{{ form.name || '商品名称' }}</div>
                  <div class="preview-price">¥ {{ form.price.toFixed(2) }}</div>
                  <div class="preview-category">
                    {{ getCategoryName(form.categoryId) || '未分类' }}
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- 表单操作按钮 - 更现代的设计 -->
      <div class="form-actions">
        <div class="action-left">
          <el-button plain @click="goBack">
            <el-icon><ArrowLeft /></el-icon>取消
          </el-button>
          <el-button @click="resetForm">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
        </div>
        <div class="action-right">
          <el-button type="primary" @click="submitForm" :loading="submitting" size="large">
            <el-icon><Check /></el-icon>{{ isEdit ? '更新商品' : '保存商品' }}
          </el-button>
        </div>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, onActivated } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  Plus, Check, Refresh, Close, Back, ArrowLeft, InfoFilled,
  Goods, Document, Picture, Upload, View, RefreshRight,
  Folder, QuestionFilled
} from '@element-plus/icons-vue';
import { addProduct, updateProduct, getProductDetail } from '../../api/product';
import { getCategoryList } from '../../api/category';
import { getToken } from '../../utils/auth';

const router = useRouter();
const route = useRoute();

// 表单引用
const formRef = ref(null);

// 是否为编辑模式
const isEdit = computed(() => !!route.params.id);

// 提交状态
const submitting = ref(false);

// 分类选项
const categoryOptions = ref([]);

// 获取分类名称
const getCategoryName = (categoryId) => {
  if (!categoryId) return '';
  const category = categoryOptions.value.find(item => String(item.id) === String(categoryId));
  return category ? category.name : '';
};

// 上传请求头
const uploadHeaders = computed(() => {
  const token = getToken();
  console.log('上传图片使用的token:', token);
  
  // 检查token是否存在
  if (!token) {
    ElMessage.warning('未检测到登录凭证，请重新登录');
    setTimeout(() => {
      router.push('/login');
    }, 2000);
    return {};
  }
  
  return {
    Authorization: `Bearer ${token}`
  };
});

// 表单数据
const form = reactive({
  id: '',
  name: '',
  categoryId: null,
  price: 0,
  stock: 0,
  image: '',
  status: 1,
  description: ''
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入商品价格', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入商品库存', trigger: 'blur' }
  ],
  image: [
    { required: true, message: '请上传商品图片', trigger: 'change' }
  ]
};

// 获取分类列表
const fetchCategoryList = async () => {
  try {
    const response = await getCategoryList();
    categoryOptions.value = response.data.list || [];
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.error('获取分类列表失败');
  }
};

// 获取商品详情
const fetchProductDetail = async (id) => {
  try {
    const response = await getProductDetail(id);
    // 检查响应数据结构
    console.log('商品详情响应:', response);
    
    // 根据截图中的响应格式调整数据处理逻辑
    if (response.data) {
      const product = response.data;
      console.log('原始商品数据:', product);
      console.log('分类ID类型:', typeof product.category_id, '值:', product.category_id);
      
      // 映射后端字段到表单字段
      form.id = product.id;
      form.name = product.name || '';
      form.price = typeof product.price === 'string' ? parseFloat(product.price) : product.price || 0;
      form.stock = typeof product.stock === 'string' ? parseInt(product.stock) : product.stock || 0;
      form.image = product.image_url || product.image || '';
      form.description = product.description || '';
      // 状态字段可能是is_active(布尔值)或status(数字)
      form.status = product.status === 1 || product.is_active === true ? 1 : 0;
      // 分类ID处理，保持为数字类型
      form.categoryId = product.category_id || null;
      
      console.log('表单数据设置后:', form);
      console.log('表单分类ID类型:', typeof form.categoryId, '值:', form.categoryId);
      
      ElMessage.success('商品信息加载成功');
      // 赋值后刷新表单校验
      if (formRef.value) {
        formRef.value.clearValidate();
      }
    } else {
      ElMessage.error('商品信息不存在');
      goBack();
    }
  } catch (error) {
    console.error('获取商品详情失败:', error);
    ElMessageBox.confirm('商品信息加载失败，是否重试？', '提示', {
      confirmButtonText: '重试',
      cancelButtonText: '返回列表',
      type: 'warning'
    }).then(() => {
      fetchProductDetail(id);
    }).catch(() => {
      goBack();
    });
  }
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

// 上传前验证
const beforeUpload = (file) => {
  // 检查token是否正确设置
  const token = getToken();
  console.log('上传前检查token:', token);
  console.log('上传请求头:', uploadHeaders.value);
  
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png' || file.type === 'image/gif';
  const isLt5M = file.size / 1024 / 1024 < 5;

  if (!isJPG) {
    ElMessage.error('上传图片只能是 JPG/PNG/GIF 格式!');
  }
  if (!isLt5M) {
    ElMessage.error('上传图片大小不能超过 5MB!');
  }
  return isJPG && isLt5M;
};

// 上传成功回调
const handleUploadSuccess = (res) => {
  console.log('上传响应:', res);
  
  // 检查响应格式
  if (res.code === 200 && res.data && res.data.url) {
    form.image = res.data.url;
    
    // 测试图片是否可以加载
    const img = new Image();
    img.onload = () => {
      ElMessage.success('图片上传成功');
    };
    img.onerror = () => {
      console.error('图片加载失败:', res.data.url);
      ElMessage.warning('图片上传成功，但可能无法正常显示，请检查网络或服务器配置');
    };
    img.src = getImageUrl(res.data.url);
  } else if (res.data && res.data.url) {
    // 兼容其他响应格式
    form.image = res.data.url;
    ElMessage.success('图片上传成功');
  } else {
    console.error('上传响应格式异常:', res);
    ElMessage.error(res.message || '上传失败，响应格式异常');
  }
};

// 上传错误回调
const handleUploadError = (error) => {
  console.error('上传错误:', error);
  ElMessage.error('上传失败，请检查网络或服务器配置');
};

// 提交表单
const submitForm = () => {
  formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning('请完善表单信息');
      return;
    }
    
    submitting.value = true;
    try {
      // 准备提交的数据，确保格式与后端API一致
      const submitData = {
        name: form.name,
        category_id: form.categoryId,
        price: parseFloat(form.price),
        stock: parseInt(form.stock),
        image_url: form.image,
        description: form.description,
        is_active: parseInt(form.status)
      };
      
      if (isEdit.value) {
        await updateProduct(form.id, submitData);
        ElMessage.success('商品更新成功');
        // 返回列表页面，列表页面的onActivated钩子会重新加载数据
        goBack();
      } else {
        await addProduct(submitData);
        ElMessage.success('商品添加成功');
        // 返回列表页面，列表页面的onActivated钩子会重新加载数据
        goBack();
      }
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败: ' + (error.message || '未知错误'));
    } finally {
      submitting.value = false;
    }
  });
};

// 重置表单数据
const resetFormData = () => {
  console.log('重置表单数据');
  form.id = '';
  form.name = '';
  form.categoryId = null;
  form.price = 0;
  form.stock = 0;
  form.image = '';
  form.status = 1;
  form.description = '';
  
  // 延迟执行表单验证重置，确保DOM更新后执行
  setTimeout(() => {
    if (formRef.value) {
      formRef.value.clearValidate();
    }
  }, 0);
};

// 重置表单
const resetForm = () => {
  resetFormData();
};

// 返回列表
const goBack = () => {
  router.push('/product/list');
};

// 页面加载时获取数据
onMounted(async () => {
  console.log('商品编辑页面已挂载', route.path);
  await fetchCategoryList();
  
  if (isEdit.value) {
    // 编辑模式，获取商品详情
    await fetchProductDetail(route.params.id);
  } else {
    // 添加模式，始终重置表单数据
    resetFormData();
    // 如果是从列表页点击"添加商品"按钮过来的，清除标志
    if (localStorage.getItem('product_form_reset') === 'true') {
      localStorage.removeItem('product_form_reset');
    }
  }
});

// 当页面从缓存中激活时触发
onActivated(() => {
  console.log('商品编辑页面已激活', route.path);
  // 如果当前是添加商品路径，则清空表单
  if (route.path === '/product/add') {
    console.log('检测到添加商品路径，清空表单');
    resetFormData();
  }
});

// 监听路由变化
watch(
  () => route.path,
  (newPath) => {
    console.log('路由路径变化:', newPath);
    if (newPath === '/product/add') {
      console.log('路由变为添加商品路径，清空表单');
      resetFormData();
    }
  }
);

// 新增：监听路由参数变化，动态获取详情
watch(
  () => route.params.id,
  async (newId, oldId) => {
    if (isEdit.value && newId && newId !== oldId) {
      await fetchProductDetail(newId);
    }
  }
);
</script>

<style scoped>
.page-container {
  padding: 24px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 84px);
}

/* 页面头部样式 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background-color: #fff;
  padding: 16px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.page-title-container {
  display: flex;
  align-items: center;
}

.page-icon {
  font-size: 24px;
  color: #409EFF;
  margin-right: 12px;
}

.page-title-content {
  display: flex;
  flex-direction: column;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.page-subtitle {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

/* 步骤条样式 */
.steps-bar {
  margin-bottom: 24px;
  padding: 16px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

/* 表单样式 */
.product-form {
  margin-bottom: 20px;
}

.form-card {
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.form-card:hover {
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.card-icon {
  margin-right: 8px;
  color: #409EFF;
  font-size: 18px;
}

.help-icon {
  color: #909399;
  cursor: pointer;
  font-size: 16px;
}

.help-icon:hover {
  color: #409EFF;
}

/* 状态开关样式 */
.status-switch {
  margin-right: 12px;
}

.status-text {
  font-size: 14px;
  color: #909399;
}

/* 分类选项样式 */
.category-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 描述文本域样式 */
.description-textarea :deep(.el-textarea__inner) {
  font-family: Arial, sans-serif;
  line-height: 1.6;
  padding: 12px;
}

/* 图片上传区域样式 */
.image-card {
  height: 100%;
}

.upload-container {
  width: 100%;
}

.upload-area {
  width: 100%;
  height: 240px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  transition: all 0.3s;
  position: relative;
  background-color: #fafafa;
}

.upload-area:hover {
  border-color: #409EFF;
  background-color: #f5f7fa;
}

.upload-area.has-image:hover .image-overlay {
  opacity: 1;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #909399;
}

.upload-icon {
  font-size: 36px;
  margin-bottom: 12px;
  color: #c0c4cc;
}

.upload-text {
  font-size: 14px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
  color: #fff;
}

.image-overlay .upload-icon {
  color: #fff;
  font-size: 32px;
}

.upload-tips {
  margin-top: 16px;
  width: 100%;
}

.upload-alert {
  border-radius: 8px;
}

.upload-requirements {
  padding-left: 0;
  list-style: none;
  margin: 8px 0 0 0;
}

.upload-requirements li {
  display: flex;
  align-items: center;
  margin-bottom: 6px;
  font-size: 13px;
}

.upload-requirements li .el-icon {
  margin-right: 6px;
  color: #67c23a;
}

/* 商品预览样式 */
.preview-section {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px dashed #e4e7ed;
}

.preview-title {
  display: flex;
  align-items: center;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.preview-title .el-icon {
  margin-right: 6px;
  color: #409EFF;
}

.product-preview-card {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.preview-image-container {
  width: 100%;
  height: 180px;
  overflow: hidden;
  background-color: #f5f7fa;
  display: flex;
  justify-content: center;
  align-items: center;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  color: #909399;
  font-size: 14px;
}

.preview-info {
  padding: 12px;
}

.preview-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-price {
  font-size: 18px;
  color: #f56c6c;
  font-weight: 600;
  margin-bottom: 8px;
}

.preview-category {
  font-size: 13px;
  color: #909399;
  background-color: #f5f7fa;
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
}

/* 表单操作按钮样式 */
.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 32px;
  background-color: #fff;
  padding: 16px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.action-left {
  display: flex;
  gap: 12px;
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .page-container {
    padding: 16px;
  }
  
  .form-actions {
    flex-direction: column-reverse;
    gap: 16px;
  }
  
  .action-left, .action-right {
    width: 100%;
    display: flex;
    justify-content: center;
  }
  
  .action-left {
    gap: 12px;
  }
  
  .action-right .el-button {
    width: 100%;
  }
}
</style> 