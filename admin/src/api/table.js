import request from '../utils/request';

// 获取桌号列表
export function getTableList(params) {
  return request({
    url: '/admin/table-qrcodes',
    method: 'get',
    params
  });
}

// 获取桌号详情
export function getTableDetail(id) {
  return request({
    url: `/admin/table-qrcodes/${id}`,
    method: 'get'
  });
}

// 创建桌号
export function createTable(data) {
  return request({
    url: '/admin/table-qrcodes',
    method: 'post',
    data
  });
}

// 更新桌号
export function updateTable(id, data) {
  return request({
    url: `/admin/table-qrcodes/${id}`,
    method: 'put',
    data
  });
}

// 删除桌号
export function deleteTable(id) {
  return request({
    url: `/admin/table-qrcodes/${id}`,
    method: 'delete'
  });
}

// 重新生成桌号二维码
export function regenerateQrcode(id) {
  return request({
    url: `/admin/table-qrcodes/${id}/regenerate`,
    method: 'post'
  });
}

// 下载桌号二维码
export function downloadQrcode(id) {
  return request({
    url: `/admin/table-qrcodes/${id}/download`,
    method: 'get',
    responseType: 'blob'
  });
} 