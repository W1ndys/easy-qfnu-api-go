package stats

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbPath = "./data/stats.db"
	mu     sync.Mutex
)

// GetDB 获取统计数据库连接
// 如果数据库文件被删除，会自动重新创建
func GetDB() *sql.DB {
	mu.Lock()
	defer mu.Unlock()

	// 检查数据库文件是否存在，如果不存在则重新初始化
	if db != nil {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			// 文件被删除，关闭旧连接并重新初始化
			db.Close()
			db = nil
		}
	}

	if db == nil {
		initDB()
	}

	return db
}

func initDB() {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建数据目录失败: %v", err)
		return
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("打开统计数据库失败: %v", err)
		return
	}

	// 设置连接池参数
	db.SetMaxOpenConns(1) // SQLite 单写
	db.SetMaxIdleConns(1)

	// 创建表
	createTables()
}

func createTables() {
	if db == nil {
		return
	}

	// API 请求日志表 - 使用 Unix 时间戳
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS api_request_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			method TEXT NOT NULL,
			status_code INTEGER NOT NULL,
			latency_ms INTEGER NOT NULL,
			client_ip TEXT NOT NULL,
			user_agent TEXT,
			created_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Printf("创建 api_request_logs 表失败: %v", err)
	}

	// 创建索引以加速查询
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_created_at ON api_request_logs(created_at)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_path ON api_request_logs(path)`)

	// 搜索热词统计表 - 使用 Unix 时间戳
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS search_keywords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			keyword TEXT NOT NULL UNIQUE,
			search_count INTEGER DEFAULT 1,
			last_searched_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Printf("创建 search_keywords 表失败: %v", err)
	}

	// 系统信息表 - 使用 Unix 时间戳
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS system_info (
			id INTEGER PRIMARY KEY,
			start_time INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Printf("创建 system_info 表失败: %v", err)
	}
}

// RecordStartTime 记录系统启动时间（Unix 时间戳）
func RecordStartTime() {
	db := GetDB()
	if db == nil {
		return
	}
	now := time.Now().Unix()
	_, err := db.Exec(`
		INSERT OR REPLACE INTO system_info (id, start_time)
		VALUES (1, ?)
	`, now)
	if err != nil {
		log.Printf("记录启动时间失败: %v", err)
	}
}

// Close 关闭数据库连接
func Close() {
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		db.Close()
		db = nil
	}
}
