package stats

import (
	"log"
	"time"
)

// RecordKeyword 记录搜索关键词
func RecordKeyword(keyword string) {
	if keyword == "" {
		return
	}

	db := GetDB()
	now := time.Now().Unix()

	// 使用 UPSERT 语法：存在则更新计数，不存在则插入
	_, err := db.Exec(`
		INSERT INTO search_keywords (keyword, search_count, last_searched_at)
		VALUES (?, 1, ?)
		ON CONFLICT(keyword) DO UPDATE SET
			search_count = search_count + 1,
			last_searched_at = ?
	`, keyword, now, now)

	if err != nil {
		log.Printf("记录搜索关键词失败: %v", err)
	}
}
