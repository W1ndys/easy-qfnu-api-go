/**
 * 培养方案 API
 * 提供培养方案相关接口
 */
window.CoursePlanApi = {
  /**
   * 获取培养方案
   * @returns {Promise<Object>}
   */
  async getCoursePlan() {
    return await window.request.get('/api/v1/zhjw/course-plan');
  },

  /**
   * 根据完成情况返回颜色类名
   * @param {string} status - 完成情况
   * @returns {string} Tailwind 颜色类
   */
  getStatusColorClass(status) {
    if (!status) return 'text-[#8E8E93]';
    if (status.includes('优') || status.includes('良')) return 'text-success';
    if (status.includes('中') || status.includes('及格')) return 'text-info';
    if (status.includes('已修')) return 'text-primary';
    if (status.includes('未修')) return 'text-[#8E8E93]';
    return 'text-[#1C1C1E]';
  },

  /**
   * 计算总学分统计
   * @param {Array} groups - 课程组列表
   * @returns {Object}
   */
  calculateTotalStats(groups) {
    let totalRequired = 0;
    let totalEarned = 0;
    let courseCount = 0;
    let completedCount = 0;
    let groupCount = groups.length;
    let completedGroupCount = 0;

    groups.forEach(group => {
      totalRequired += group.required_credits || 0;
      totalEarned += group.earned_credits || 0;
      // 检查模块是否完成（已修学分 >= 应修学分）
      if ((group.earned_credits || 0) >= (group.required_credits || 0)) {
        completedGroupCount++;
      }
      if (group.courses) {
        courseCount += group.courses.length;
        completedCount += group.courses.filter(c => c.status && c.status.includes('已修')).length;
      }
    });

    return {
      totalRequired: totalRequired.toFixed(1),
      totalEarned: totalEarned.toFixed(1),
      progress: totalRequired > 0 ? Math.round((totalEarned / totalRequired) * 100) : 0,
      courseCount,
      completedCount,
      groupCount,
      completedGroupCount
    };
  }
};
