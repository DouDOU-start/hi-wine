import { get, post, put, del, download } from '../utils/request';

// 获取桌号列表
export function getTableList(params) {
  return get('/admin/table-qrcodes', params);
}

// 获取桌号详情
export function getTableDetail(id) {
  return get(`/admin/table-qrcodes/${id}`);
}

// 创建桌号
export function createTable(data) {
  return post('/admin/table-qrcodes', data);
}

// 更新桌号
export function updateTable(id, data) {
  return put(`/admin/table-qrcodes/${id}`, data);
}

// 删除桌号
export function deleteTable(id) {
  return del(`/admin/table-qrcodes/${id}`);
}

// 重新生成桌号二维码
export function regenerateQrcode(id) {
  return post(`/admin/table-qrcodes/${id}/regenerate`);
}

// 下载桌号二维码
export function downloadQrcode(id) {
  return download(`/admin/table-qrcodes/${id}/download`, {}, `table-qrcode-${id}.png`);
} 