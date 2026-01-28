/**
 * 选课结果 API
 * 提供选课结果相关接口和辅助方法
 */
window.SelectionApi = {
  /**
   * 生成学期选项列表
   * @returns {Array<{value: string, label: string}>}
   */
  generateTermOptions() {
    const options = [];
    const currentYear = new Date().getFullYear();

    for (let year = currentYear + 1; year >= 2020; year--) {
      options.push({ value: `${year - 1}-${year}-3`, label: `${year - 1}-${year} 学年 第三学期` });
      options.push({ value: `${year - 1}-${year}-2`, label: `${year - 1}-${year} 学年 第二学期` });
      options.push({ value: `${year - 1}-${year}-1`, label: `${year - 1}-${year} 学年 第一学期` });
    }

    return options;
  },

  /**
   * 查询选课结果
   * @param {string} [term] - 学期
   * @returns {Promise<Object>}
   */
  async getSelectionResults(term) {
    const queryParams = new URLSearchParams();
    if (term) queryParams.append('term', term);

    const queryString = queryParams.toString();
    const url = `/api/v1/zhjw/selection-results${queryString ? '?' + queryString : ''}`;

    return await window.request.get(url);
  },

  /**
   * 计算总学分
   * @param {Array} results - 选课结果列表
   * @returns {number}
   */
  calculateTotalCredits(results) {
    if (!results || results.length === 0) return 0;
    return results.reduce((sum, r) => sum + (parseFloat(r.credit) || 0), 0).toFixed(1);
  },

  /**
   * 计算总学时
   * @param {Array} results - 选课结果列表
   * @returns {number}
   */
  calculateTotalHours(results) {
    if (!results || results.length === 0) return 0;
    return results.reduce((sum, r) => sum + (parseInt(r.hours) || 0), 0);
  }
};
