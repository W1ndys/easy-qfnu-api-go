/**
 * 课表查询 API
 * 提供课表相关接口和辅助方法
 */
window.ScheduleApi = {
  /**
   * 查询课表
   * @param {string} [date] - 日期 (e.g., 2026-01-01)
   * @returns {Promise<Object>}
   */
  async getSchedule(date) {
    const queryParams = new URLSearchParams();
    if (date) queryParams.append('date', date);

    const queryString = queryParams.toString();
    const url = `/api/v1/zhjw/schedule${queryString ? '?' + queryString : ''}`;

    return await window.request.get(url);
  },

  /**
   * 获取今天的日期字符串
   * @returns {string} 格式 YYYY-MM-DD
   */
  getTodayString() {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0');
    const day = String(today.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  },

  /**
   * 星期几映射
   */
  dayOfWeekMap: {
    1: '星期一',
    2: '星期二',
    3: '星期三',
    4: '星期四',
    5: '星期五',
    6: '星期六',
    7: '星期日'
  },

  /**
   * 获取星期几的文本
   * @param {number} day - 1-7
   * @returns {string}
   */
  getDayOfWeekText(day) {
    return this.dayOfWeekMap[day] || `星期${day}`;
  },

  /**
   * 节次时间映射
   */
  periodTimeMap: {
    1: '08:00-08:45',
    2: '08:55-09:40',
    3: '10:00-10:45',
    4: '10:55-11:40',
    5: '14:00-14:45',
    6: '14:55-15:40',
    7: '16:00-16:45',
    8: '16:55-17:40',
    9: '19:00-19:45',
    10: '19:55-20:40',
    11: '20:50-21:35'
  },

  /**
   * 格式化节次显示
   * @param {Array<number>} periods - 节次数组
   * @returns {string}
   */
  formatPeriods(periods) {
    if (!periods || periods.length === 0) return '--';
    if (periods.length === 1) return `第${periods[0]}节`;
    return `第${periods[0]}-${periods[periods.length - 1]}节`;
  },

  /**
   * 获取节次时间范围
   * @param {Array<number>} periods - 节次数组
   * @returns {string}
   */
  getPeriodTimeRange(periods) {
    if (!periods || periods.length === 0) return '';
    const start = this.periodTimeMap[periods[0]]?.split('-')[0] || '';
    const end = this.periodTimeMap[periods[periods.length - 1]]?.split('-')[1] || '';
    return start && end ? `${start}-${end}` : '';
  },

  /**
   * 按星期分组课程
   * @param {Array} courses - 课程列表
   * @returns {Object} 按星期分组的课程 { 1: [...], 2: [...], ... }
   */
  groupByDayOfWeek(courses) {
    const grouped = {};
    courses.forEach(course => {
      const day = course.timeParsed?.dayOfWeek;
      if (day) {
        if (!grouped[day]) grouped[day] = [];
        grouped[day].push(course);
      }
    });
    // 按节次排序
    Object.keys(grouped).forEach(day => {
      grouped[day].sort((a, b) => {
        const aPeriod = a.timeParsed?.periodArray?.[0] || 0;
        const bPeriod = b.timeParsed?.periodArray?.[0] || 0;
        return aPeriod - bPeriod;
      });
    });
    return grouped;
  }
};
