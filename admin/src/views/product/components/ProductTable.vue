<template>
  <div class="product-table">
    <!-- 表格工具栏 -->
    <table-toolbar 
      :columns="columns" 
      :loading="loading"
      @refresh="$emit('refresh')"
      @size-change="handleSizeChange"
      @column-change="handleColumnChange"
    >
      <template #left>
        <el-button 
          type="primary" 
          @click="$emit('add')"
          v-if="hasPermission('product:add')"
        >
          <el-icon><Plus /></el-icon>添加商品
        </el-button>
        <el-button 
          type="success" 
          @click="handleBatchAction('上架')"
          v-if="hasPermission('product:update') && selectedRows.length > 0"
        >
          <el-icon><Top /></el-icon>批量上架
        </el-button>
        <el-button 
          type="info" 
          @click="handleBatchAction('下架')"
          v-if="hasPermission('product:update') && selectedRows.length > 0"
        >
          <el-icon><Bottom /></el-icon>批量下架
        </el-button>
        <el-button 
          type="danger" 
          @click="handleBatchAction('删除')"
          v-if="hasPermission('product:delete') && selectedRows.length > 0"
        >
          <el-icon><Delete /></el-icon>批量删除
        </el-button>
      </template>
      <template #right>
        <el-tooltip content="导出数据" placement="top">
          <el-button 
            circle 
            plain 
            @click="$emit('export')" 
            v-if="hasPermission('product:export')"
          >
            <el-icon><Download /></el-icon>
          </el-button>
        </el-tooltip>
      </template>
    </table-toolbar>
    
    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="tableData"
      :size="tableSize"
      border
      stripe
      highlight-current-row
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
    >
      <el-table-column type="selection" width="50" fixed="left" />
      
      <el-table-column 
        v-for="col in visibleColumns" 
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :width="col.width"
        :min-width="col.minWidth"
        :sortable="col.sortable"
        :fixed="col.fixed"
        :show-overflow-tooltip="col.showOverflowTooltip"
        :align="col.align || 'left'"
      >
        <template #default="scope">
          <!-- 图片列 -->
          <template v-if="col.prop === 'image'">
            <el-image 
              :src="scope.row.imageUrl || scope.row.image_url || scope.row.image" 
              :preview-src-list="[scope.row.imageUrl || scope.row.image_url || scope.row.image]"
              fit="cover"
              class="product-image"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
          
          <!-- 价格列 -->
          <template v-else-if="col.prop === 'price'">
            <span class="price">{{ formatPrice(scope.row.price) }}</span>
          </template>
          
          <!-- 状态列 -->
          <template v-else-if="col.prop === 'status'">
            <el-tag 
              :type="scope.row.status === 1 ? 'success' : 'info'" 
              effect="light"
            >
              {{ scope.row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
          
          <!-- 库存列 -->
          <template v-else-if="col.prop === 'stock'">
            <el-tag 
              :type="getStockTagType(scope.row.stock)" 
              effect="light"
            >
              {{ scope.row.stock }}
            </el-tag>
          </template>
          
          <!-- 创建时间列 -->
          <template v-else-if="col.prop === 'createdAt'">
            {{ formatDate(scope.row.createdAt) }}
          </template>
          
          <!-- 操作列 -->
          <template v-else-if="col.prop === 'actions'">
            <el-button 
              type="primary" 
              link 
              @click="$emit('edit', scope.row)"
              v-if="hasPermission('product:update')"
            >
              编辑
            </el-button>
            <el-button 
              :type="scope.row.status === 1 ? 'info' : 'success'" 
              link 
              @click="$emit('toggle-status', scope.row)"
              v-if="hasPermission('product:update')"
            >
              {{ scope.row.status === 1 ? '下架' : '上架' }}
            </el-button>
            <el-button 
              type="danger" 
              link 
              @click="handleDelete(scope.row)"
              v-if="hasPermission('product:delete')"
            >
              删除
            </el-button>
          </template>
          
          <!-- 默认列 -->
          <template v-else>
            {{ scope.row[col.prop] }}
          </template>
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handlePageSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, inject } from 'vue';
import { ElMessageBox } from 'element-plus';
import { Plus, Top, Bottom, Delete, Download, Picture } from '@element-plus/icons-vue';
import { formatPrice, formatDate } from '@/utils/format';

const props = defineProps({
  data: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  }
});

const emit = defineEmits([
  'refresh', 'add', 'edit', 'delete', 'toggle-status', 
  'export', 'batch-action', 'page-change', 'size-change',
  'sort-change', 'update:currentPage', 'update:pageSize'
]);

// 从父组件注入分类选项
const categoryOptions = inject('categoryOptions', ref([]));

// 根据分类ID获取分类名称
const getCategoryNameById = (categoryId) => {
  if (!categoryId) return '未分类';
  
  // 确保categoryOptions是ref对象且其值是数组
  if (!categoryOptions || !categoryOptions.value || !Array.isArray(categoryOptions.value)) {
    console.error('分类选项不是有效数组:', categoryOptions);
    return '未分类';
  }
  
  const category = categoryOptions.value.find(cat => Number(cat.value) === Number(categoryId));
  return category ? category.label : '未分类';
};

// 表格数据
const tableData = computed(() => {
  console.log('表格数据:', props.data);
  
  // 检查数据格式，确保关键字段存在
  if (props.data && props.data.length > 0) {
    return props.data.map(item => {
      // 确保所有必要的字段都存在
      return {
        id: item.id,
        name: item.name || '未命名商品',
        image: item.image || '',
        imageUrl: item.imageUrl || item.image_url || '',
        price: item.price || 0,
        stock: item.stock || 0,
        status: typeof item.status !== 'undefined' ? item.status : 0,
        categoryId: item.categoryId || item.category_id || 0,
        categoryName: item.categoryName || getCategoryNameById(item.categoryId || item.category_id) || '未分类',
        createdAt: item.createdAt || item.created_at || item.createTime || new Date().toISOString(),
        ...item // 保留其他字段
      };
    });
  }
  return [];
});

// 表格尺寸
const tableSize = ref('default');

// 选中的行
const selectedRows = ref([]);

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 80, sortable: true },
  { prop: 'image', label: '商品图片', width: 100 },
  { prop: 'name', label: '商品名称', minWidth: 150, showOverflowTooltip: true },
  { prop: 'categoryName', label: '分类', width: 120 },
  { prop: 'price', label: '价格', width: 100, sortable: true },
  { prop: 'stock', label: '库存', width: 100, sortable: true },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'createdAt', label: '创建时间', width: 180, sortable: true },
  { prop: 'actions', label: '操作', width: 200, fixed: 'right', align: 'center' }
];

// 可见列
const visibleColumns = ref([...columns]);

// 处理表格尺寸变化
const handleSizeChange = (size) => {
  tableSize.value = size;
};

// 处理列变化
const handleColumnChange = (cols) => {
  visibleColumns.value = cols;
};

// 处理选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection;
};

// 处理排序变化
const handleSortChange = (sort) => {
  emit('sort-change', sort);
};

// 处理页码变化
const handleCurrentChange = (page) => {
  emit('update:currentPage', page);
  emit('page-change', page, props.pageSize);
};

// 处理每页条数变化
const handlePageSizeChange = (size) => {
  emit('update:pageSize', size);
  emit('size-change', props.currentPage, size);
};

// 处理删除
const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除商品 "${row.name}" 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    emit('delete', row);
  }).catch(() => {});
};

// 处理批量操作
const handleBatchAction = (action) => {
  if (selectedRows.value.length === 0) {
    return;
  }
  
  const ids = selectedRows.value.map(row => row.id);
  let message = '';
  let actionType = '';
  
  switch (action) {
    case '上架':
      message = `确定要上架选中的 ${ids.length} 个商品吗?`;
      actionType = 'online';
      break;
    case '下架':
      message = `确定要下架选中的 ${ids.length} 个商品吗?`;
      actionType = 'offline';
      break;
    case '删除':
      message = `确定要删除选中的 ${ids.length} 个商品吗?`;
      actionType = 'delete';
      break;
  }
  
  ElMessageBox.confirm(message, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    emit('batch-action', actionType, ids);
  }).catch(() => {});
};

// 获取库存标签类型
const getStockTagType = (stock) => {
  if (stock <= 0) return 'danger';
  if (stock <= 10) return 'warning';
  return 'success';
};

// 检查权限
const hasPermission = (permission) => {
  // 这里可以实现实际的权限检查逻辑
  return true;
};
</script>

<style scoped>
.product-table {
  margin-top: 20px;
}

.product-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
}

.image-placeholder {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.image-placeholder .el-icon {
  font-size: 24px;
  color: #909399;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>