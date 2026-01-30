package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	statsDB   *sql.DB
	appDB     *sql.DB
	statsPath = "./data/stats.db"
	appPath   = "./data/app.db"
	mu        sync.Mutex
)

// GetStatsDB 获取统计数据库连接
func GetStatsDB() *sql.DB {
	mu.Lock()
	defer mu.Unlock()

	if statsDB != nil {
		if _, err := os.Stat(statsPath); os.IsNotExist(err) {
			statsDB.Close()
			statsDB = nil
		}
	}

	if statsDB == nil {
		statsDB = openDB(statsPath)
		initStatsTables()
	}

	return statsDB
}

// GetAppDB 获取应用数据库连接
func GetAppDB() *sql.DB {
	mu.Lock()
	defer mu.Unlock()

	if appDB != nil {
		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			appDB.Close()
			appDB = nil
		}
	}

	if appDB == nil {
		appDB = openDB(appPath)
		initAppTables()
	}

	return appDB
}

// openDB 打开数据库连接
func openDB(path string) *sql.DB {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建数据目录失败: %v", err)
		return nil
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("打开数据库失败 %s: %v", path, err)
		return nil
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db
}

// Close 关闭所有数据库连接
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if statsDB != nil {
		statsDB.Close()
		statsDB = nil
	}
	if appDB != nil {
		appDB.Close()
		appDB = nil
	}
}

// initStatsTables 初始化统计数据库表
func initStatsTables() {
	if statsDB == nil {
		return
	}

	// API 请求日志表
	statsDB.Exec(`
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

	statsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_created_at ON api_request_logs(created_at)`)
	statsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_logs_path ON api_request_logs(path)`)

	// 搜索热词统计表
	statsDB.Exec(`
		CREATE TABLE IF NOT EXISTS search_keywords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			keyword TEXT NOT NULL UNIQUE,
			search_count INTEGER DEFAULT 1,
			last_searched_at INTEGER NOT NULL
		)
	`)

	// 系统信息表
	statsDB.Exec(`
		CREATE TABLE IF NOT EXISTS system_info (
			id INTEGER PRIMARY KEY,
			start_time INTEGER NOT NULL
		)
	`)
}

// initAppTables 初始化应用数据库表
func initAppTables() {
	if appDB == nil {
		return
	}

	// 系统配置表
	appDB.Exec(`
		CREATE TABLE IF NOT EXISTS system_config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at INTEGER NOT NULL
		)
	`)

	// 公告表
	appDB.Exec(`
		CREATE TABLE IF NOT EXISTS announcements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			type TEXT DEFAULT 'info',
			is_active INTEGER DEFAULT 1,
			sort_order INTEGER DEFAULT 0,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)
	`)
}
