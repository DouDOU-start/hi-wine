<template>
  <div>
    <el-button type="primary" @click="showAddDialog = true">新增桌台</el-button>
    <el-table :data="tableList" style="margin-top: 20px;">
      <el-table-column prop="id" label="ID" width="80"/>
      <el-table-column prop="name" label="桌台名称"/>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
            {{ scope.row.status === 1 ? '可用' : '不可用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="scope">
          <el-button size="mini" @click="editTable(scope.row)">编辑</el-button>
          <el-button size="mini" type="danger" @click="deleteTable(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="showAddDialog" title="桌台">
      <el-form :model="form">
        <el-form-item label="桌台名称">
          <el-input v-model="form.name"/>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status">
            <el-option label="可用" :value="1"/>
            <el-option label="不可用" :value="0"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTable">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTableList, createTable, updateTable, deleteTable as delTableApi } from '@/api/table'

const tableList = ref([])
const showAddDialog = ref(false)
const form = ref({ id: null, name: '', status: 1 })

const fetchTableList = async () => {
  const res = await getTableList()
  tableList.value = res.data || []
}
const saveTable = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入桌台名称')
    return
  }
  if (form.value.id) {
    await updateTable(form.value)
    ElMessage.success('更新成功')
  } else {
    await createTable(form.value)
    ElMessage.success('创建成功')
  }
  showAddDialog.value = false
  fetchTableList()
}
const editTable = (row) => {
  form.value = { ...row }
  showAddDialog.value = true
}
const deleteTable = (id) => {
  ElMessageBox.confirm('确定删除该桌台吗？', '提示').then(async () => {
    await delTableApi({ id })
    ElMessage.success('删除成功')
    fetchTableList()
  })
}
onMounted(fetchTableList)
</script> 