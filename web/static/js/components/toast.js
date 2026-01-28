/**
 * Toast 通知组件 (Floating Cards Style)
 * 动感卡片式设计，支持多状态 Emoji 和颜色区分
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
   * @param {number} duration - 显示时长(毫秒), 默认 3000
   */
  show(message, type = 'info', duration = 3000) {
    this.init();

    const toast = document.createElement('div');

    // 配置样式映射
    const config = {
      success: {
        emoji: '🎉',
        iconBg: 'bg-[#34C759]/15',
        titleColor: 'text-[#1C1C1E]',
        msgColor: 'text-[#3C3C43]/60'
      },
      error: {
        emoji: '💣',
        iconBg: 'bg-[#FF3B30]/15',
        titleColor: 'text-[#1C1C1E]',
        msgColor: 'text-[#3C3C43]/60'
      },
      warning: {
        emoji: '🔔',
        iconBg: 'bg-[#FF9500]/15',
        titleColor: 'text-[#1C1C1E]',
        msgColor: 'text-[#3C3C43]/60'
      },
      info: {
        emoji: '🦄',
        iconBg: 'bg-[#007AFF]/15',
        titleColor: 'text-[#1C1C1E]',
        msgColor: 'text-[#3C3C43]/60'
      }
    };

    const style = config[type] || config.info;

    // 应用基础类和动画类
    toast.className = 'toast-card toast-enter group cursor-default select-none';

    // 构建内容结构
    toast.innerHTML = `
      <div class="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center text-lg ${style.iconBg} backdrop-blur-sm group-hover:scale-110 transition-transform duration-300">
        ${style.emoji}
      </div>
      <div class="flex flex-col min-w-0">
        <span class="text-[15px] font-semibold ${style.titleColor} leading-snug">${message}</span>
      </div>
    `;

    this.container.appendChild(toast);

    // 触发进入动画
    requestAnimationFrame(() => {
      toast.classList.remove('toast-enter');
      toast.classList.add('toast-enter-active');
    });

    // 悬停暂停计时逻辑
    let timer;
    const startTimer = () => {
      timer = setTimeout(() => {
        removeToast();
      }, duration);
    };

    const removeToast = () => {
      toast.classList.remove('toast-enter-active');
      toast.classList.add('toast-leave-active');

      // 等待动画结束后移除 DOM
      setTimeout(() => {
        if (toast.parentNode) {
          toast.parentNode.removeChild(toast);
        }
      }, 300);
    };

    // 绑定鼠标事件
    toast.addEventListener('mouseenter', () => clearTimeout(timer));
    toast.addEventListener('mouseleave', startTimer);

    // 开始计时
    startTimer();
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
