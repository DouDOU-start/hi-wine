import request from '../utils/request';

// 获取商品列表
export function getProductList(params) {
  return request({
    url: '/api/product/list',
    method: 'get',
    params
  });
}

// 添加商品
export function addProduct(data) {
  return request({
    url: '/api/product/add',
    method: 'post',
    data
  });
}

// 更新商品
export function updateProduct(data) {
  return request({
    url: '/api/product/update',
    method: 'post',
    data
  });
}

// 删除商品
export function deleteProduct(id) {
  return request({
    url: '/api/product/delete',
    method: 'post',
    data: { id }
  });
}

// 获取商品详情
export function getProductDetail(id) {
  return request({
    url: '/api/product/detail',
    method: 'get',
    params: { id }
  });
}

// 上下架商品
export function updateProductStatus(id, status) {
  return request({
    url: '/api/product/status',
    method: 'post',
    data: { id, status }
  });
}

// 批量更新商品状态
export function batchUpdateProductStatus(ids, status) {
  return request({
    url: '/api/product/batch-status',
    method: 'post',
    data: { ids, status }
  });
}

// 上传商品图片
export function uploadProductImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  return request({
    url: '/api/admin/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
} 