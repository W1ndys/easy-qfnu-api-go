/**
 * 考试安排 API
 * 提供考试安排相关接口和辅助方法
 */
window.ExamApi = {
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
   * 查询考试安排
   * @param {string} [term] - 学期
   * @returns {Promise<Object>}
   */
  async getExamSchedules(term) {
    const queryParams = new URLSearchParams();
    if (term) queryParams.append('term', term);

    const queryString = queryParams.toString();
    const url = `/api/v1/zhjw/exam-schedules${queryString ? '?' + queryString : ''}`;

    return await window.request.get(url);
  },

  /**
   * 解析考试时间，判断是否已过期
   * @param {string} examTime - 考试时间字符串
   * @returns {Object} { isPast: boolean, isToday: boolean, isSoon: boolean }
   */
  parseExamTimeStatus(examTime) {
    if (!examTime) return { isPast: false, isToday: false, isSoon: false };

    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

    // 尝试从考试时间中提取日期 (格式可能是 "2026-01-15 09:00-11:00" 或类似格式)
    const dateMatch = examTime.match(/(\d{4})-(\d{2})-(\d{2})/);
    if (!dateMatch) return { isPast: false, isToday: false, isSoon: false };

    const examDate = new Date(parseInt(dateMatch[1]), parseInt(dateMatch[2]) - 1, parseInt(dateMatch[3]));
    const diffDays = Math.ceil((examDate - today) / (1000 * 60 * 60 * 24));

    return {
      isPast: diffDays < 0,
      isToday: diffDays === 0,
      isSoon: diffDays > 0 && diffDays <= 7
    };
  },

  /**
   * 获取考试状态的颜色类
   * @param {string} examTime - 考试时间
   * @returns {string}
   */
  getExamStatusClass(examTime) {
    const status = this.parseExamTimeStatus(examTime);
    if (status.isPast) return 'text-[#8E8E93]';
    if (status.isToday) return 'text-danger';
    if (status.isSoon) return 'text-warning';
    return 'text-primary';
  }
};
