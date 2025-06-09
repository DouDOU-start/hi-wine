<template>
  <div class="table-edit-container">
    <div class="page-header">
      <h2>{{ isEdit ? '编辑桌号' : '添加桌号' }}</h2>
      <el-button @click="goBack">返回列表</el-button>
    </div>
    
    <el-card class="form-container">
      <el-form
        ref="tableFormRef"
        :model="tableForm"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="桌号" prop="tableNumber">
          <el-input v-model="tableForm.tableNumber" placeholder="请输入桌号" />
        </el-form-item>
        
        <el-form-item label="状态" prop="isActive">
          <el-switch
            v-model="tableForm.isActive"
            active-text="激活"
            inactive-text="禁用"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '保存修改' : '创建桌号' }}
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 二维码预览 -->
    <el-card v-if="isEdit && tableDetail.qrcodeUrl" class="qrcode-container">
      <template #header>
        <div class="card-header">
          <span class="header-title">桌号二维码</span>
        </div>
      </template>
      <div class="qrcode-preview">
        <el-image
          :src="tableDetail.qrcodeUrl"
          :preview-src-list="[tableDetail.qrcodeUrl]"
          fit="contain"
          style="width: 200px; height: 200px;"
        />
      </div>
      <div class="qrcode-actions">
        <el-button type="success" @click="handleRegenerateQrcode">重新生成二维码</el-button>
        <el-button type="primary" @click="handleDownloadQrcode">下载二维码</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getTableDetail, createTable, updateTable, regenerateQrcode, downloadQrcode } from '../../api/table';

const router = useRouter();
const route = useRoute();
const tableId = computed(() => route.params.id);
const isEdit = computed(() => !!tableId.value);

// 表单引用
const tableFormRef = ref(null);

// 表单数据
const tableForm = reactive({
  tableNumber: '',
  isActive: true
});

// 表单验证规则
const rules = {
  tableNumber: [
    { required: true, message: '请输入桌号', trigger: 'blur' },
    { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
  ]
};

// 提交状态
const submitting = ref(false);

// 桌号详情
const tableDetail = ref({});

// 返回列表页
const goBack = () => {
  router.push('/table/list');
};

// 获取桌号详情
const getDetail = async () => {
  if (!isEdit.value) return;
  
  try {
    const response = await getTableDetail(tableId.value);
    tableDetail.value = response.data || {};
    
    // 填充表单
    tableForm.tableNumber = tableDetail.value.tableNumber || '';
    tableForm.isActive = tableDetail.value.isActive !== false;
  } catch (error) {
    ElMessage.error('获取桌号详情失败');
    console.error(error);
  }
};

// 提交表单
const submitForm = () => {
  tableFormRef.value.validate(async (valid) => {
    if (!valid) return;
    
    submitting.value = true;
    
    try {
      const tableData = {
        tableNumber: tableForm.tableNumber,
        isActive: tableForm.isActive
      };
      
      if (isEdit.value) {
        // 更新桌号
        await updateTable(tableId.value, tableData);
        ElMessage.success('更新桌号成功');
      } else {
        // 创建桌号
        await createTable(tableData);
        ElMessage.success('创建桌号成功');
      }
      
      router.push('/table/list');
    } catch (error) {
      ElMessage.error(isEdit.value ? '更新桌号失败' : '创建桌号失败');
      console.error(error);
    } finally {
      submitting.value = false;
    }
  });
};

// 重置表单
const resetForm = () => {
  tableFormRef.value.resetFields();
  if (isEdit.value) {
    getDetail();
  }
};

// 重新生成二维码
const handleRegenerateQrcode = async () => {
  ElMessageBox.confirm('确认重新生成二维码吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await regenerateQrcode(tableId.value);
      ElMessage.success('二维码重新生成成功');
      getDetail(); // 刷新详情
    } catch (error) {
      ElMessage.error('重新生成二维码失败');
      console.error(error);
    }
  }).catch(() => {});
};

// 下载二维码
const handleDownloadQrcode = async () => {
  try {
    await downloadQrcode(tableId.value);
    ElMessage.success('二维码下载成功');
  } catch (error) {
    ElMessage.error('下载二维码失败');
    console.error(error);
  }
};

// 初始化
onMounted(() => {
  getDetail();
});
</script>

<style scoped>
.table-edit-container {
  padding: 0 20px 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  font-weight: 500;
  color: #303133;
  margin: 0;
}

.form-container {
  margin-bottom: 20px;
}

.qrcode-container {
  margin-bottom: 20px;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.qrcode-preview {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.qrcode-actions {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 10px;
}
</style> 