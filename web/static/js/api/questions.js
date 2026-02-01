/**
 * 智慧曲园入学题库搜索 API
 * 提供试题搜索相关接口
 */
window.QuestionsApi = {
  /**
   * 搜索试题
   * @param {string} keyword - 搜索关键词
   * @returns {Promise<Object>}
   */
  async search(keyword) {
    const queryParams = new URLSearchParams();
    if (keyword) queryParams.append('keyword', keyword);

    const queryString = queryParams.toString();
    const url = `/api/v1/questions/search${queryString ? '?' + queryString : ''}`;

    return await window.request.get(url);
  },

  /**
   * 题目类型映射
   */
  typeMap: {
    '单选题': 'primary',
    '多选题': 'info',
    '判断题': 'success',
    '填空题': 'warning'
  },

  /**
   * 获取题目类型的颜色类
   * @param {string} type - 题目类型
   * @returns {string}
   */
  getTypeColorClass(type) {
    const color = this.typeMap[type] || 'primary';
    return `bg-${color}/10 text-${color}`;
  },

  /**
   * 高亮关键词
   * @param {string} text - 原文
   * @param {string} keyword - 关键词
   * @returns {string}
   */
  highlightKeyword(text, keyword) {
    if (!text || !keyword) return text || '';
    const regex = new RegExp(`(${keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
    return text.replace(regex, '<mark class="bg-warning/30 px-0.5 rounded">$1</mark>');
  }
};
