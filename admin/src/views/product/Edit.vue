<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">{{ isEdit ? '编辑商品' : '添加商品' }}</div>
      <el-button @click="goBack">
        <el-icon><Back /></el-icon>返回列表
      </el-button>
    </div>
    
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="product-form"
      :disabled="submitting"
    >
      <el-row :gutter="20">
        <!-- 左侧表单 -->
        <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
          <!-- 基本信息 -->
          <el-card shadow="hover" class="form-card">
            <template #header>
              <div class="card-header">
                <span>基本信息</span>
              </div>
            </template>
            
            <el-form-item label="商品名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入商品名称" maxlength="50" show-word-limit />
            </el-form-item>
            
            <el-form-item label="商品分类" prop="categoryId">
              <el-select v-model="form.categoryId" placeholder="请选择商品分类" style="width: 100%">
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
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
                  />
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
                  />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="商品状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :label="1">上架</el-radio>
                <el-radio :label="0">下架</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-card>
          
          <!-- 商品描述 -->
          <el-card shadow="hover" class="form-card">
            <template #header>
              <div class="card-header">
                <span>商品描述</span>
              </div>
            </template>
            
            <el-form-item label="商品描述" prop="description">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="6"
                placeholder="请输入商品描述"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
          </el-card>
        </el-col>
        
        <!-- 右侧图片上传 -->
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <el-card shadow="hover" class="form-card">
            <template #header>
              <div class="card-header">
                <span>商品图片</span>
              </div>
            </template>
            
            <el-form-item prop="image" class="image-form-item">
              <el-upload
                class="product-image-uploader"
                action="/api/admin/upload/image"
                :show-file-list="false"
                :on-success="handleUploadSuccess"
                :before-upload="beforeUpload"
                :headers="uploadHeaders"
                :disabled="submitting"
              >
                <div class="upload-area">
                  <img v-if="form.image" :src="form.image" class="product-image" />
                  <div v-else class="upload-placeholder">
                    <el-icon class="upload-icon"><Plus /></el-icon>
                    <div class="upload-text">点击上传图片</div>
                  </div>
                </div>
              </el-upload>
              
              <div class="upload-tips">
                <el-alert
                  title="图片上传提示"
                  type="info"
                  :closable="false"
                >
                  <template #default>
                    <p>- 建议上传正方形图片，大小不超过5MB</p>
                    <p>- 支持JPG、PNG、GIF格式</p>
                    <p>- 图片清晰度不低于300x300像素</p>
                  </template>
                </el-alert>
              </div>
            </el-form-item>
            
            <div class="image-preview" v-if="form.image">
              <div class="preview-title">图片预览</div>
              <div class="preview-container">
                <img :src="form.image" class="preview-image" />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- 表单操作按钮 -->
      <div class="form-actions">
        <el-button type="primary" @click="submitForm" :loading="submitting">
          <el-icon><Check /></el-icon>{{ isEdit ? '更新商品' : '保存商品' }}
        </el-button>
        <el-button @click="resetForm">
          <el-icon><Refresh /></el-icon>重置
        </el-button>
        <el-button @click="goBack">
          <el-icon><Close /></el-icon>取消
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Plus, Check, Refresh, Close, Back } from '@element-plus/icons-vue';
import { addProduct, updateProduct, getProductDetail } from '../../api/product';
import { getCategoryList } from '../../api/category';

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

// 上传请求头
const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${localStorage.getItem('token')}`
  };
});

// 表单数据
const form = reactive({
  id: '',
  name: '',
  categoryId: '',
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
    if (response.data && response.data.product) {
      const product = response.data.product;
      Object.keys(form).forEach(key => {
        if (product[key] !== undefined) {
          form[key] = product[key];
        }
      });
      ElMessage.success('商品信息加载成功');
    } else {
      ElMessage.error('商品信息不存在');
      goBack();
    }
  } catch (error) {
    console.error('获取商品详情失败:', error);
    ElMessage.error('获取商品详情失败');
    goBack();
  }
};

// 上传前验证
const beforeUpload = (file) => {
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
  if (res.code === 0 && res.data && res.data.url) {
    form.image = res.data.url;
    
    // 测试图片是否可以加载
    const img = new Image();
    img.onload = () => {
      ElMessage.success('图片上传成功');
    };
    img.onerror = () => {
      ElMessage.warning('图片上传成功，但可能无法正常显示，请检查网络或服务器配置');
    };
    img.src = res.data.url;
  } else {
    ElMessage.error(res.message || '上传失败');
  }
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
      if (isEdit.value) {
        await updateProduct(form);
        ElMessage.success('商品更新成功');
      } else {
        await addProduct(form);
        ElMessage.success('商品添加成功');
        resetForm();
      }
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败: ' + (error.message || '未知错误'));
    } finally {
      submitting.value = false;
    }
  });
};

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  
  if (!isEdit.value) {
    form.id = '';
    form.name = '';
    form.categoryId = '';
    form.price = 0;
    form.stock = 0;
    form.image = '';
    form.status = 1;
    form.description = '';
  } else {
    // 如果是编辑模式，重新获取商品详情
    fetchProductDetail(route.params.id);
  }
};

// 返回列表
const goBack = () => {
  router.push('/product/list');
};

// 页面加载时获取数据
onMounted(async () => {
  await fetchCategoryList();
  
  if (isEdit.value) {
    await fetchProductDetail(route.params.id);
  }
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

.product-form {
  margin-bottom: 20px;
}

.form-card {
  margin-bottom: 20px;
  transition: all 0.3s;
}

.form-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
}

.image-form-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.product-image-uploader {
  width: 100%;
  cursor: pointer;
}

.upload-area {
  width: 100%;
  height: 200px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  transition: all 0.3s;
}

.upload-area:hover {
  border-color: #409EFF;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #8c939d;
}

.upload-icon {
  font-size: 28px;
  margin-bottom: 8px;
}

.upload-text {
  font-size: 14px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.upload-tips {
  margin-top: 10px;
  width: 100%;
}

.image-preview {
  margin-top: 20px;
}

.preview-title {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.preview-container {
  width: 100%;
  height: 150px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  background-color: #f5f7fa;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 30px;
}

/* 响应式调整 */
@media screen and (max-width: 768px) {
  .form-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .form-actions .el-button {
    width: 100%;
    margin-left: 0;
    margin-bottom: 10px;
  }
}
</style> 