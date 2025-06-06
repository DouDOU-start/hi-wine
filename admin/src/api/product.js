import request from '../utils/request';

// 获取商品列表
export function getProductList(params) {
  return request({
    url: '/admin/products',
    method: 'get',
    params
  });
}

// 获取商品详情
export function getProductDetail(id) {
  return request({
    url: `/admin/products/${id}`,
    method: 'get'
  });
}

// 添加商品
export function addProduct(data) {
  return request({
    url: '/admin/products',
    method: 'post',
    data
  });
}

// 更新商品
export function updateProduct(id, data) {
  return request({
    url: `/admin/products/${id}`,
    method: 'put',
    data
  });
}

// 删除商品
export function deleteProduct(id) {
  return request({
    url: `/admin/products/${id}`,
    method: 'delete'
  });
}

// 上下架商品
export function updateProductStatus(id, status) {
  console.log('API调用 - 更新商品状态:', id, status);
  // 确保status是布尔值
  const isActive = status === 1;
  
  // 先获取商品详情，然后只修改状态
  return request({
    url: `/admin/products/${id}`,
    method: 'get'
  }).then(response => {
    const product = response.data;
    // 修改状态字段
    product.is_active = isActive;
    
    // 发送更新请求
    return request({
      url: `/admin/products/${id}`,
      method: 'put',
      data: product
    });
  });
}

// 批量更新商品状态
export function batchUpdateProductStatus(ids, status) {
  // 确保status是布尔值
  const isActive = status === 1;
  return request({
    url: '/admin/products/batch',
    method: 'put',
    data: { ids, is_active: isActive }
  });
}

// 上传商品图片
export function uploadProductImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  return request({
    url: '/admin/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
} 