// 选课推荐 API 封装
window.CourseRecommendationApi = {
    // 查询课程推荐
    async query(keyword) {
        return await window.request.get('/api/v1/course-recommendation/query', {
            params: { keyword }
        });
    },

    // 提交课程推荐
    async recommend(data) {
        return await window.request.post('/api/v1/course-recommendation/recommend', data);
    },

    // 格式化时间戳
    formatTime(timestamp) {
        const date = new Date(timestamp * 1000);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hour = String(date.getHours()).padStart(2, '0');
        const minute = String(date.getMinutes()).padStart(2, '0');
        return `${year}-${month}-${day} ${hour}:${minute}`;
    },

    // 高亮关键词
    highlightKeyword(text, keyword) {
        if (!text || !keyword) return text || '';
        const escaped = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        const regex = new RegExp(`(${escaped})`, 'gi');
        return text.replace(regex, '<mark class="bg-warning/30 text-warning px-0.5 rounded">$1</mark>');
    }
};
