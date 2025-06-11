import { createApi } from '../utils/apiFactory';
import { post, download } from '../utils/request';

// 创建基础API
const tableApi = createApi('/admin/table-qrcodes');

// 导出基础API方法
export const {
  getList: getTableList,
  getDetail: getTableDetail,
  create: createTable,
  update: updateTable,
  delete: deleteTable
} = tableApi;

// 重新生成桌号二维码
export function regenerateQrcode(id) {
  return post(`/admin/table-qrcodes/${id}/regenerate`);
}

// 下载桌号二维码
export function downloadQrcode(id) {
  return download(`/admin/table-qrcodes/${id}/download`, {}, `table-qrcode-${id}.png`);
} 