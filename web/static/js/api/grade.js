/**
 * 成绩查询 API
 * 提供成绩相关接口和辅助方法
 */
window.GradeApi = {
  /**
   * 课程性质选项
   */
  courseTypeOptions: [
    { value: '01', label: '公共课' },
    { value: '02', label: '公共基础课' },
    { value: '03', label: '专业基础课' },
    { value: '04', label: '专业课' },
    { value: '05', label: '专业选修课' },
    { value: '06', label: '公共选修课' },
    { value: '07', label: '专业任选课' },
    { value: '08', label: '实践教学环节' },
    { value: '09', label: '公共任选课' },
    { value: '10', label: '教师教育基础课程（必修）' },
    { value: '11', label: '专业必修课' },
    { value: '12', label: '学科基础必修课' },
    { value: '13', label: '专业方向限选课' },
    { value: '14', label: '考试报名虚拟课程' },
    { value: '15', label: '教师教育选修课程' },
    { value: '16', label: '公共必修课' }
  ],

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
   * 查询成绩
   * @param {Object} params - 查询参数
   * @param {string} [params.term] - 学期
   * @param {string} [params.courseType] - 课程性质
   * @param {string} [params.courseName] - 课程名称
   * @param {string} [params.displayType] - 显示方式 (all/max)
   * @returns {Promise<Object>}
   */
  async getGrades(params = {}) {
    const queryParams = new URLSearchParams();

    if (params.term) queryParams.append('term', params.term);
    if (params.courseType) queryParams.append('course_type', params.courseType);
    if (params.courseName) queryParams.append('course_name', params.courseName);
    if (params.displayType) queryParams.append('display_type', params.displayType);

    const queryString = queryParams.toString();
    const url = `/api/v1/zhjw/grades${queryString ? '?' + queryString : ''}`;

    return await window.request.get(url);
  },

  /**
   * 根据成绩返回颜色类名
   * @param {string} score - 成绩
   * @returns {string} Tailwind 颜色类
   */
  getScoreColorClass(score) {
    const num = parseFloat(score);

    if (isNaN(num)) {
      // 非数字成绩（优/良/中/及格/不及格）
      if (['优', '良'].some(s => score?.includes(s))) return 'text-success';
      if (['中', '及格'].some(s => score?.includes(s))) return 'text-info';
      if (score?.includes('不及格')) return 'text-danger';
      return 'text-[#1C1C1E]';
    }

    // 数字成绩
    if (num >= 90) return 'text-success';
    if (num >= 80) return 'text-info';
    if (num >= 70) return 'text-[#00C7BE]'; // cyan
    if (num >= 60) return 'text-warning';
    return 'text-danger';
  },

  /**
   * 计算自定义统计
   * @param {Array} grades - 成绩列表
   * @param {Array} selectedIndices - 选中的索引数组
   * @returns {Object} 统计结果
   */
  calculateCustomStats(grades, selectedIndices) {
    const selected = grades.filter((_, i) => selectedIndices[i]);

    let totalCredits = 0;
    let weightedGPASum = 0;
    let weightedScoreSum = 0;
    let validGPACredits = 0;
    let validScoreCredits = 0;

    selected.forEach(g => {
      const credit = parseFloat(g.credit) || 0;
      const gpa = parseFloat(g.gpa);
      const score = parseFloat(g.score);

      totalCredits += credit;

      if (!isNaN(gpa) && gpa >= 0 && credit > 0) {
        weightedGPASum += gpa * credit;
        validGPACredits += credit;
      }

      if (!isNaN(score) && credit > 0) {
        weightedScoreSum += score * credit;
        validScoreCredits += credit;
      }
    });

    return {
      courseCount: selected.length,
      totalCredits: totalCredits.toFixed(1),
      weightedGPA: validGPACredits > 0 ? (weightedGPASum / validGPACredits).toFixed(2) : '--',
      avgScore: validScoreCredits > 0 ? (weightedScoreSum / validScoreCredits).toFixed(2) : '--'
    };
  },

  /**
   * 计算通过率
   * @param {Array} grades - 成绩列表
   * @returns {number} 通过率百分比 (0-100)
   */
  calculatePassRate(grades) {
    if (!grades || grades.length === 0) return 0;

    const passed = grades.filter(g => {
      const score = parseFloat(g.score);
      if (!isNaN(score)) return score >= 60;
      return ['优', '良', '中', '及格'].some(s => g.score?.includes(s)) && !g.score?.includes('不及格');
    }).length;

    return Math.round((passed / grades.length) * 100);
  }
};
