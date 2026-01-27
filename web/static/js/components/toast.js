/**
 * Toast 通知组件
 * 符合 UI.md 毛玻璃效果规范
 */
window.Toast = {
  container: null,

  /**
   * 初始化 Toast 容器
   */
  init() {
    if (this.container) return;

    this.container = document.createElement('div');
    this.container.id = 'toast-container';
    document.body.appendChild(this.container);
  },

  /**
   * 显示 Toast 消息
   * @param {string} message - 消息内容
   * @param {string} type - 类型: success, error, warning, info
   * @param {number} duration - 显示时长(毫秒), 默认 2500
   */
  show(message, type = 'info', duration = 2500) {
    this.init();

    const toast = document.createElement('div');
    toast.className = 'toast toast-enter';

    // 根据类型添加图标
    const icons = {
      success: '<span class="text-[#34C759]">✓</span>',
      error: '<span class="text-[#FF3B30]">✕</span>',
      warning: '<span class="text-[#FF9500]">!</span>',
      info: '<span class="text-[#007AFF]">i</span>'
    };

    toast.innerHTML = `${icons[type] || ''}${message}`;
    this.container.appendChild(toast);

    // 触发进入动画
    requestAnimationFrame(() => {
      toast.classList.remove('toast-enter');
      toast.classList.add('toast-enter-active');
    });

    // 自动移除
    setTimeout(() => {
      toast.classList.remove('toast-enter-active');
      toast.classList.add('toast-leave');
      requestAnimationFrame(() => {
        toast.classList.remove('toast-leave');
        toast.classList.add('toast-leave-active');
      });

      setTimeout(() => {
        if (toast.parentNode) {
          toast.parentNode.removeChild(toast);
        }
      }, 200);
    }, duration);
  },

  /**
   * 快捷方法
   */
  success(message, duration) {
    this.show(message, 'success', duration);
  },

  error(message, duration) {
    this.show(message, 'error', duration);
  },

  warning(message, duration) {
    this.show(message, 'warning', duration);
  },

  info(message, duration) {
    this.show(message, 'info', duration);
  }
};
