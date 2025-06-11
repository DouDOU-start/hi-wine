import { get, post, put, del } from './request';

/**
 * 创建标准CRUD API函数集
 * @param {string} baseUrl - API基础路径
 * @param {Object} options - 配置选项
 * @returns {Object} 包含标准CRUD操作的对象
 */
export function createApi(baseUrl, options = {}) {
  const { 
    idField = 'id',
    listMethod = 'get',
    detailMethod = 'get',
    createMethod = 'post',
    updateMethod = 'put',
    deleteMethod = 'delete',
    statusUpdateUrl,
    batchUrl
  } = options;

  const api = {
    // 获取列表
    getList(params) {
      return get(baseUrl, params);
    },

    // 获取详情
    getDetail(id) {
      return get(`${baseUrl}/${id}`);
    },

    // 创建
    create(data) {
      return post(baseUrl, data);
    },

    // 更新
    update(id, data) {
      return put(`${baseUrl}/${id}`, data);
    },

    // 删除
    delete(id) {
      return del(`${baseUrl}/${id}`);
    }
  };

  // 添加状态更新方法（如果配置了）
  if (statusUpdateUrl) {
    api.updateStatus = (id, status) => {
      const url = typeof statusUpdateUrl === 'function' 
        ? statusUpdateUrl(id) 
        : `${baseUrl}/${id}${statusUpdateUrl}`;
      return put(url, { status });
    };
  }

  // 添加批量操作方法（如果配置了）
  if (batchUrl) {
    api.batchUpdate = (ids, data) => {
      const url = typeof batchUrl === 'function'
        ? batchUrl()
        : `${baseUrl}${batchUrl}`;
      return put(url, { ids, ...data });
    };
  }

  return api;
} 