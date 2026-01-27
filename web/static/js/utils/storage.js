/**
 * 本地存储工具
 * 提供 Cookie 和 localStorage 操作
 */
window.Storage = {
  /**
   * 获取 Cookie 值
   * @param {string} name - Cookie 名称
   * @returns {string|null}
   */
  getCookie(name) {
    const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    return match ? decodeURIComponent(match[2]) : null;
  },

  /**
   * 设置 Cookie
   * @param {string} name - Cookie 名称
   * @param {string} value - Cookie 值
   * @param {number} days - 过期天数，默认 7 天
   * @param {string} path - 路径，默认 /
   */
  setCookie(name, value, days = 7, path = '/') {
    const expires = new Date(Date.now() + days * 864e5).toUTCString();
    document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=${path}`;
  },

  /**
   * 删除 Cookie
   * @param {string} name - Cookie 名称
   * @param {string} path - 路径，默认 /
   */
  removeCookie(name, path = '/') {
    document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=${path}`;
  },

  /**
   * 获取 localStorage 中的 JSON 对象
   * @param {string} key - 键名
   * @param {*} defaultValue - 默认值
   * @returns {*}
   */
  getJSON(key, defaultValue = null) {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : defaultValue;
    } catch (e) {
      console.error('Storage.getJSON error:', e);
      return defaultValue;
    }
  },

  /**
   * 将 JSON 对象存储到 localStorage
   * @param {string} key - 键名
   * @param {*} value - 值（会被 JSON 序列化）
   */
  setJSON(key, value) {
    try {
      localStorage.setItem(key, JSON.stringify(value));
    } catch (e) {
      console.error('Storage.setJSON error:', e);
    }
  },

  /**
   * 删除 localStorage 项
   * @param {string} key - 键名
   */
  remove(key) {
    localStorage.removeItem(key);
  },

  /**
   * 清空 localStorage
   */
  clear() {
    localStorage.clear();
  }
};
