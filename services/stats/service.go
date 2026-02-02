package stats

import (
	"github.com/W1ndys/easy-qfnu-api-go/common/stats"
)

// DashboardData 大屏数据结构
type DashboardData struct {
	TotalRequests      int64              `json:"totalRequests"`
	TodayRequests      int64              `json:"todayRequests"`
	UniqueIPs          int64              `json:"uniqueIPs"`
	TodayUniqueIPs     int64              `json:"todayUniqueIPs"`
	AvgLatencyMs       float64            `json:"avgLatencyMs"`
	StartTime          int64              `json:"startTime"`
	APIStats           []APIStat          `json:"apiStats"`
	StatusCodeStats    []StatusCodeStat   `json:"statusCodeStats"`
	TopKeywords        []KeywordStat      `json:"topKeywords"`
}

// APIStat 各接口统计
type APIStat struct {
	Path        string  `json:"path"`
	Count       int64   `json:"count"`
	AvgLatency  float64 `json:"avgLatency"`
}

// StatusCodeStat 状态码分布
type StatusCodeStat struct {
	StatusCode int   `json:"statusCode"`
	Count      int64 `json:"count"`
}

// KeywordStat 搜索热词
type KeywordStat struct {
	Keyword       string `json:"keyword"`
	SearchCount   int64  `json:"searchCount"`
	LastSearched  int64  `json:"lastSearched"`
}

// TrendData 趋势数据
type TrendData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetDashboardData 获取大屏数据
func GetDashboardData() (*DashboardData, error) {
	db := stats.GetDB()
	data := &DashboardData{}

	// 1. 总请求数
	db.QueryRow(`SELECT COUNT(*) FROM api_request_logs`).Scan(&data.TotalRequests)

	// 2. 今日请求数（东八区）
	db.QueryRow(`
		SELECT COUNT(*) FROM api_request_logs
		WHERE datetime(created_at, 'unixepoch', '+8 hours') >= datetime('now', '+8 hours', 'start of day')
	`).Scan(&data.TodayRequests)

	// 3. 独立 IP 数
	db.QueryRow(`SELECT COUNT(DISTINCT client_ip) FROM api_request_logs`).Scan(&data.UniqueIPs)

	// 4. 今日独立 IP 数（东八区）
	db.QueryRow(`
		SELECT COUNT(DISTINCT client_ip) FROM api_request_logs
		WHERE datetime(created_at, 'unixepoch', '+8 hours') >= datetime('now', '+8 hours', 'start of day')
	`).Scan(&data.TodayUniqueIPs)

	// 5. 平均响应时间
	db.QueryRow(`SELECT COALESCE(AVG(latency_ms), 0) FROM api_request_logs`).Scan(&data.AvgLatencyMs)

	// 6. 系统启动时间
	db.QueryRow(`SELECT COALESCE(start_time, 0) FROM system_info WHERE id = 1`).Scan(&data.StartTime)

	// 7. 各接口统计（Top 10）
	rows, err := db.Query(`
		SELECT path, COUNT(*) as cnt, AVG(latency_ms) as avg_lat
		FROM api_request_logs
		GROUP BY path
		ORDER BY cnt DESC
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat APIStat
			rows.Scan(&stat.Path, &stat.Count, &stat.AvgLatency)
			data.APIStats = append(data.APIStats, stat)
		}
	}

	// 8. 状态码分布
	rows, err = db.Query(`
		SELECT status_code, COUNT(*) as cnt
		FROM api_request_logs
		GROUP BY status_code
		ORDER BY cnt DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat StatusCodeStat
			rows.Scan(&stat.StatusCode, &stat.Count)
			data.StatusCodeStats = append(data.StatusCodeStats, stat)
		}
	}

	// 9. 搜索热词 Top 10
	rows, err = db.Query(`
		SELECT keyword, search_count, last_searched_at
		FROM search_keywords
		ORDER BY search_count DESC
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat KeywordStat
			rows.Scan(&stat.Keyword, &stat.SearchCount, &stat.LastSearched)
			// 对关键词进行脱敏处理
			stat.Keyword = maskKeyword(stat.Keyword)
			data.TopKeywords = append(data.TopKeywords, stat)
		}
	}

	return data, nil
}

// maskKeyword 对关键词进行脱敏处理
// 规则：保留首尾各一个字符，中间用星号替换
func maskKeyword(keyword string) string {
	runes := []rune(keyword)
	length := len(runes)

	if length <= 1 {
		return keyword
	}
	if length == 2 {
		return string(runes[0]) + "*"
	}

	// 保留首尾各一个字符，中间全部用星号替换
	masked := string(runes[0])
	for i := 1; i < length-1; i++ {
		masked += "*"
	}
	masked += string(runes[length-1])

	return masked
}

// GetTrendData 获取调用趋势数据
func GetTrendData(days int) ([]TrendData, error) {
	if days <= 0 {
		days = 7
	}
	if days > 30 {
		days = 30
	}

	db := stats.GetDB()
	rows, err := db.Query(`
		SELECT date(created_at, 'unixepoch', '+8 hours') as day, COUNT(*) as cnt
		FROM api_request_logs
		WHERE datetime(created_at, 'unixepoch', '+8 hours') >= datetime('now', '+8 hours', 'start of day', '-' || ? || ' days')
		GROUP BY day
		ORDER BY day ASC
	`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []TrendData
	for rows.Next() {
		var t TrendData
		rows.Scan(&t.Date, &t.Count)
		trends = append(trends, t)
	}

	return trends, nil
}
