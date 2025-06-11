<template>
  <div class="page-container">
    <div class="page-title">分类管理</div>
    
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <div class="left">
          <el-button type="primary" @click="handleAdd">添加分类</el-button>
        </div>
      </div>
      
      <el-table
        v-loading="loading"
        :data="categoryList"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="sortOrder" label="排序" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isActive"
              :active-value="true"
              :inactive-value="false"
              @change="handleStatusChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column prop="updatedAt" label="更新时间" width="180" />
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
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 添加/编辑分类对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑分类' : '添加分类'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="form.sortOrder" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.isActive"
            :active-value="true"
            :inactive-value="false"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getCategoryList, addCategory, updateCategory, deleteCategory, updateCategoryStatus } from '../../api/category';

// 防止重复请求的锁
const isRequestLocked = ref(false);
// 记录页面是否已经初始化
const isInitialized = ref(false);

// 加载状态
const loading = ref(false);

// 分类列表数据
const categoryList = ref([]);

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 表单引用
const formRef = ref(null);

// 对话框控制
const dialogVisible = ref(false);
const isEdit = ref(false);

// 表单数据
const form = reactive({
  id: '',
  name: '',
  sortOrder: 0,
  isActive: true
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  sortOrder: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ]
};

// 获取分类列表
const fetchCategoryList = async () => {
  if (isRequestLocked.value) return;
  isRequestLocked.value = true;

  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value
    };
    
    const response = await getCategoryList(params);
    categoryList.value = response.data.list || [];
    total.value = response.data.total || 0;
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.error('获取分类列表失败');
  } finally {
    loading.value = false;
    isRequestLocked.value = false;
  }
};

// 分页大小变化
const handleSizeChange = (size) => {
  pageSize.value = size;
  fetchCategoryList();
};

// 页码变化
const handleCurrentChange = (page) => {
  currentPage.value = page;
  fetchCategoryList();
};

// 添加分类
const handleAdd = () => {
  isEdit.value = false;
  resetForm();
  dialogVisible.value = true;
};

// 编辑分类
const handleEdit = (row) => {
  isEdit.value = true;
  resetForm();
  Object.assign(form, {
    id: row.id,
    name: row.name,
    sortOrder: row.sortOrder,
    isActive: row.isActive
  });
  dialogVisible.value = true;
};

// 删除分类
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除分类"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await deleteCategory(row.id);
      if (response && response.code === 200) {
        ElMessage.success('删除成功');
        fetchCategoryList();
      } else {
        ElMessage.error(response?.message || '删除分类失败');
      }
    } catch (error) {
      console.error('删除分类失败:', error);
      ElMessage.error('删除分类失败');
    }
  }).catch(() => {});
};

// 更新分类状态
const handleStatusChange = async (row) => {
  try {
    const response = await updateCategoryStatus(row.id, row.isActive);
    if (response && response.code === 200) {
      ElMessage.success(`已${row.isActive ? '启用' : '禁用'}分类"${row.name}"`);
    } else {
      // 状态更新失败，恢复原状态
      row.isActive = !row.isActive;
      ElMessage.error(response?.message || '更新分类状态失败');
    }
  } catch (error) {
    console.error('更新分类状态失败:', error);
    ElMessage.error('更新分类状态失败');
    // 恢复原状态
    row.isActive = !row.isActive;
  }
};

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  form.id = '';
  form.name = '';
  form.sortOrder = 0;
  form.isActive = true;
};

// 提交表单
const submitForm = () => {
  formRef.value.validate(async (valid) => {
    if (!valid) return;
    
    try {
      const formData = {
        name: form.name,
        sortOrder: form.sortOrder,
        isActive: form.isActive
      };
      
      let response;
      if (isEdit.value) {
        response = await updateCategory(form.id, formData);
        if (response && response.code === 200) {
          ElMessage.success('更新成功');
          dialogVisible.value = false;
          fetchCategoryList();
        } else {
          ElMessage.error(response?.message || '更新分类失败');
        }
      } else {
        response = await addCategory(formData);
        if (response && response.code === 200) {
          ElMessage.success('添加成功');
          dialogVisible.value = false;
          fetchCategoryList();
        } else {
          ElMessage.error(response?.message || '添加分类失败');
        }
      }
    } catch (error) {
      console.error('操作失败:', error);
      ElMessage.error('操作失败');
    }
  });
};

// 页面加载时获取数据
onMounted(() => {
  console.log('分类列表页面已挂载');
  if (!isInitialized.value) {
    fetchCategoryList();
    isInitialized.value = true;
  }
});

// 当页面从缓存中激活时触发（切换tab时）
onActivated(() => {
  console.log('分类列表页面已激活');
  // 避免重复请求数据
  if (!isRequestLocked.value) {
    fetchCategoryList();
  }
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