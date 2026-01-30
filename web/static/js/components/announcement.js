// 公告组件
function announcementBanner() {
    return {
        announcements: [],
        async init() {
            try {
                const res = await axios.get('/api/v1/site/announcements');
                this.announcements = res.data.data || [];
            } catch (e) {
                console.error('获取公告失败', e);
            }
        },
        getTypeClass(type) {
            switch (type) {
                case 'warning': return 'bg-yellow-50 border-yellow-200 text-yellow-800';
                case 'error': return 'bg-red-50 border-red-200 text-red-800';
                default: return 'bg-blue-50 border-blue-200 text-blue-800';
            }
        }
    }
}
