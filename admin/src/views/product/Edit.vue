<template>
  <div class="page-container">
    <div class="page-header">
      <div class="page-title">{{ isEdit ? '编辑商品' : '添加商品' }}</div>
      <el-button @click="goBack">返回列表</el-button>
    </div>
    
    <el-card shadow="never" class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="product-form"
      >
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        
        <el-form-item label="商品分类" prop="categoryId">
          <el-select v-model="form.categoryId" placeholder="请选择商品分类">
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="商品价格" prop="price">
          <el-input-number
            v-model="form.price"
            :precision="2"
            :step="0.1"
            :min="0"
            style="width: 200px;"
          />
        </el-form-item>
        
        <el-form-item label="商品库存" prop="stock">
          <el-input-number
            v-model="form.stock"
            :min="0"
            :precision="0"
            style="width: 200px;"
          />
        </el-form-item>
        
        <el-form-item label="商品图片" prop="image">
          <el-upload
            class="product-image-uploader"
            action="/api/admin/upload/image"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
            :before-upload="beforeUpload"
            :headers="uploadHeaders"
          >
            <img v-if="form.image" :src="form.image" class="product-image" />
            <el-icon v-else class="product-image-uploader-icon"><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">建议上传正方形图片，大小不超过5MB</div>
        </el-form-item>
        
        <el-form-item label="商品状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="商品描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入商品描述"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">{{ isEdit ? '更新' : '保存' }}</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
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
    const response = await getCategoryList({ size: 100 });
    categoryOptions.value = response.data.list || [];
  } catch (error) {
    console.error('获取分类列表失败:', error);
  }
};

// 获取商品详情
const fetchProductDetail = async (id) => {
  try {
    const response = await getProductDetail(id);
    Object.assign(form, response.data.product);
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
  console.log('上传响应:', res);
  if (res.code === 0 && res.data && res.data.url) {
    console.log('图片URL:', res.data.url);
    form.image = res.data.url;
    
    // 测试图片是否可以加载
    const img = new Image();
    img.onload = () => {
      console.log('图片加载成功');
      ElMessage.success('上传成功');
    };
    img.onerror = (error) => {
      console.error('图片加载失败:', error);
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
    if (!valid) return;
    
    submitting.value = true;
    try {
      if (isEdit.value) {
        await updateProduct(form);
        ElMessage.success('更新成功');
      } else {
        await addProduct(form);
        ElMessage.success('添加成功');
        resetForm();
      }
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
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

.form-card {
  margin-bottom: 20px;
}

.product-form {
  max-width: 800px;
}

.product-image-uploader {
  width: 200px;
  height: 200px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  display: inline-block;
  vertical-align: top;
}

.product-image-uploader:hover {
  border-color: #409EFF;
}

.product-image-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 200px;
  height: 200px;
  line-height: 200px;
  text-align: center;
}

.product-image {
  width: 200px;
  height: 200px;
  display: block;
  object-fit: cover;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style> 