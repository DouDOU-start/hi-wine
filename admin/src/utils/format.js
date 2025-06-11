/**
 * 格式化日期时间
 * @param {string|number|Date} timestamp - 时间戳或日期对象
 * @param {string} format - 格式化模式，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDate(timestamp, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!timestamp) return '-';
  
  const date = timestamp instanceof Date ? timestamp : new Date(timestamp);
  
  if (isNaN(date.getTime())) return '-';
  
  const formatOptions = {
    YYYY: date.getFullYear(),
    MM: String(date.getMonth() + 1).padStart(2, '0'),
    DD: String(date.getDate()).padStart(2, '0'),
    HH: String(date.getHours()).padStart(2, '0'),
    mm: String(date.getMinutes()).padStart(2, '0'),
    ss: String(date.getSeconds()).padStart(2, '0')
  };
  
  return Object.entries(formatOptions).reduce(
    (result, [key, value]) => result.replace(key, value),
    format
  );
}

/**
 * 格式化价格
 * @param {number|string} price - 价格数值
 * @param {number} decimals - 小数位数，默认为2
 * @param {string} currency - 货币符号，默认为¥
 * @returns {string} 格式化后的价格字符串
 */
export function formatPrice(price, decimals = 2, currency = '¥') {
  const number = Number(price);
  return `${currency}${(isNaN(number) ? 0 : number).toFixed(decimals)}`;
}

/**
 * 格式化文件大小
 * @param {number} bytes - 文件大小（字节）
 * @param {number} decimals - 小数位数，默认为2
 * @returns {string} 格式化后的文件大小字符串
 */
export function formatFileSize(bytes, decimals = 2) {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(decimals))} ${sizes[i]}`;
}

/**
 * 将对象的属性名从下划线转换为驼峰
 * @param {Object} obj - 要转换的对象
 * @returns {Object} 转换后的对象
 */
export function toCamelCase(obj) {
  if (typeof obj !== 'object' || obj === null) return obj;
  
  if (Array.isArray(obj)) {
    return obj.map(item => toCamelCase(item));
  }
  
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [
      key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase()),
      toCamelCase(value)
    ])
  );
}

/**
 * 将对象的属性名从驼峰转换为下划线
 * @param {Object} obj - 要转换的对象
 * @returns {Object} 转换后的对象
 */
export function toSnakeCase(obj) {
  if (typeof obj !== 'object' || obj === null) return obj;
  
  if (Array.isArray(obj)) {
    return obj.map(item => toSnakeCase(item));
  }
  
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [
      key.replace(/([A-Z])/g, letter => `_${letter.toLowerCase()}`),
      toSnakeCase(value)
    ])
  );
} 