/**
 * Axios 请求实例
 * 配置 baseURL、超时、拦截器
 */
(function() {
  // 创建 Axios 实例
  const request = axios.create({
    baseURL: `${window.location.protocol}//${window.location.hostname}:8080`,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json'
    }
  });

  // 请求拦截器
  request.interceptors.request.use(
    (config) => {
      // 从 Alpine store 或全局获取 cookie
      const authCookie = window.authCookie || '';
      if (authCookie) {
        config.headers['Authorization'] = authCookie;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // 响应拦截器
  request.interceptors.response.use(
    (response) => {
      // 成功: 返回 response.data
      return response.data;
    },
    (error) => {
      // 获取错误信息
      let message = '请求失败';

      if (error.response) {
        const { status, data } = error.response;

        // 提取错误消息
        message = data?.error || data?.msg || data?.message || `请求失败 (${status})`;

        // 401 未授权: 清除 cookie，但不跳转
        if (status === 401) {
          window.authCookie = '';
          if (window.Storage) {
            window.Storage.removeCookie('auth_cookie');
          }
          message = 'Cookie 无效或已过期，请重新输入';
        }
      } else if (error.request) {
        message = '网络错误，请检查网络连接';
      } else {
        message = error.message || '请求配置错误';
      }

      // Toast 提示错误（如果 Toast 组件可用）
      if (window.Toast) {
        window.Toast.error(message);
      }

      // 返回带有错误信息的 rejected Promise
      return Promise.reject({ message, originalError: error });
    }
  );

  // 挂载到全局
  window.request = request;
})();
