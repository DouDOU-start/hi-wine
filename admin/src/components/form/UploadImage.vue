<template>
  <div class="upload-image">
    <el-upload
      :action="action"
      :headers="headers"
      :multiple="multiple"
      :limit="limit"
      :file-list="fileList"
      :list-type="listType"
      :disabled="disabled"
      :accept="accept"
      :before-upload="beforeUpload"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-exceed="handleExceed"
      :on-remove="handleRemove"
      :on-preview="handlePreview"
      :auto-upload="autoUpload"
      :drag="drag"
    >
      <template v-if="drag">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">拖拽文件到此处或 <em>点击上传</em></div>
      </template>
      <template v-else>
        <el-button :type="buttonType" :icon="buttonIcon">{{ buttonText }}</el-button>
      </template>
      <template #tip>
        <div v-if="tip" class="el-upload__tip">{{ tip }}</div>
      </template>
    </el-upload>
    
    <!-- 图片预览对话框 -->
    <el-dialog v-model="dialogVisible" title="图片预览" width="50%">
      <img :src="previewUrl" alt="Preview" class="preview-image" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { UploadFilled, Plus } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

const props = defineProps({
  // 上传地址
  action: {
    type: String,
    required: true
  },
  // 请求头
  headers: {
    type: Object,
    default: () => ({})
  },
  // 是否支持多选
  multiple: {
    type: Boolean,
    default: false
  },
  // 最大上传数量
  limit: {
    type: Number,
    default: 5
  },
  // 文件列表
  modelValue: {
    type: Array,
    default: () => []
  },
  // 列表类型
  listType: {
    type: String,
    default: 'picture-card'
  },
  // 是否禁用
  disabled: {
    type: Boolean,
    default: false
  },
  // 接受的文件类型
  accept: {
    type: String,
    default: 'image/*'
  },
  // 是否自动上传
  autoUpload: {
    type: Boolean,
    default: true
  },
  // 是否启用拖拽上传
  drag: {
    type: Boolean,
    default: false
  },
  // 按钮类型
  buttonType: {
    type: String,
    default: 'primary'
  },
  // 按钮文本
  buttonText: {
    type: String,
    default: '上传图片'
  },
  // 按钮图标
  buttonIcon: {
    type: String,
    default: Plus
  },
  // 提示文本
  tip: {
    type: String,
    default: '只能上传jpg/png文件，且不超过2MB'
  },
  // 最大文件大小（MB）
  maxSize: {
    type: Number,
    default: 2
  }
});

const emit = defineEmits(['update:modelValue', 'success', 'error', 'exceed', 'remove']);

// 文件列表
const fileList = ref([]);

// 预览相关
const dialogVisible = ref(false);
const previewUrl = ref('');

// 根据传入的值更新文件列表
watch(() => props.modelValue, (newValue) => {
  if (newValue && newValue.length > 0) {
    fileList.value = newValue.map(item => {
      if (typeof item === 'string') {
        return {
          name: item.substring(item.lastIndexOf('/') + 1),
          url: item
        };
      }
      return item;
    });
  } else {
    fileList.value = [];
  }
}, { immediate: true, deep: true });

// 上传前验证
const beforeUpload = (file) => {
  // 检查文件类型
  const isImage = file.type.startsWith('image/');
  if (!isImage) {
    ElMessage.error('只能上传图片文件!');
    return false;
  }
  
  // 检查文件大小
  const isLimitSize = file.size / 1024 / 1024 < props.maxSize;
  if (!isLimitSize) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB!`);
    return false;
  }
  
  return true;
};

// 上传成功处理
const handleSuccess = (response, file, fileList) => {
  if (response.code === 200) {
    // 更新文件列表
    updateModelValue(fileList);
    
    // 发出成功事件
    emit('success', response, file, fileList);
    
    ElMessage.success('上传成功');
  } else {
    ElMessage.error(response.message || '上传失败');
    handleRemove(file, fileList);
  }
};

// 上传失败处理
const handleError = (error, file, fileList) => {
  ElMessage.error('上传失败: ' + (error.message || '未知错误'));
  emit('error', error, file, fileList);
};

// 超出上传数量限制
const handleExceed = (files, fileList) => {
  ElMessage.warning(`最多只能上传 ${props.limit} 个文件`);
  emit('exceed', files, fileList);
};

// 移除文件
const handleRemove = (file, fileList) => {
  updateModelValue(fileList);
  emit('remove', file, fileList);
};

// 预览文件
const handlePreview = (file) => {
  previewUrl.value = file.url;
  dialogVisible.value = true;
};

// 更新绑定值
const updateModelValue = (files) => {
  const value = files.map(file => {
    if (file.response && file.response.data) {
      return file.response.data;
    }
    return file.url || file;
  });
  emit('update:modelValue', value);
};
</script>

<style scoped>
.upload-image {
  width: 100%;
}

.preview-image {
  width: 100%;
  object-fit: contain;
}

:deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
  line-height: 100px;
}

:deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 100px;
  height: 100px;
}
</style> 