import { createApi } from '../utils/apiFactory';
import { put } from '../utils/request';

// 创建基础API
const categoryApi = createApi('/admin/categories');

// 导出基础API方法
export const {
  getList: getCategoryList,
  getDetail: getCategoryDetail,
  create: addCategory,
  update: updateCategory,
  delete: deleteCategory
} = categoryApi;

// 更新分类状态
export function updateCategoryStatus(id, status) {
  return put(`/admin/categories/${id}`, { isActive: status });
} 