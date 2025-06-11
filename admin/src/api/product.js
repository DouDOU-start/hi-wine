import { createApi } from '../utils/apiFactory';
import { post } from '../utils/request';

// 创建基础API
const productApi = createApi('/admin/products', {
  batchUrl: '/batch'
});

// 导出基础API方法
export const {
  getList: getProductList,
  getDetail: getProductDetail,
  create: addProduct,
  update: updateProduct,
  delete: deleteProduct,
  batchUpdate: batchUpdateProductStatus
} = productApi;

// 上下架商品
export function updateProductStatus(id, status) {
  // 确保status是布尔值
  const isActive = status === 1;
  
  return productApi.update(id, { 
    is_active: isActive, 
    status: status 
  });
}

// 上传商品图片
export function uploadProductImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  return post('/admin/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
} 