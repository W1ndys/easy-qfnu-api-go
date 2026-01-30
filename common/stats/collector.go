package stats

import (
	"log"
	"sync"
	"time"
)

// RequestLog 请求日志结构
type RequestLog struct {
	Path       string
	Method     string
	StatusCode int
	LatencyMs  int64
	ClientIP   string
	UserAgent  string
	CreatedAt  int64 // Unix 时间戳
}

var (
	logChan     chan RequestLog
	collectorOnce sync.Once
)

const (
	chanBufferSize = 1000 // Channel 缓冲大小
	batchSize      = 50   // 批量写入条数
	flushInterval  = 5    // 强制刷新间隔（秒）
)

// InitCollector 初始化异步收集器
func InitCollector() {
	collectorOnce.Do(func() {
		logChan = make(chan RequestLog, chanBufferSize)
		go runCollector()
	})
}

// Collect 收集请求日志（非阻塞）
func Collect(log RequestLog) {
	select {
	case logChan <- log:
		// 成功入队
	default:
		// 队列满时丢弃，不阻塞请求
	}
}

func runCollector() {
	batch := make([]RequestLog, 0, batchSize)
	ticker := time.NewTicker(time.Duration(flushInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case logItem := <-logChan:
			batch = append(batch, logItem)
			if len(batch) >= batchSize {
				flushBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

func flushBatch(batch []RequestLog) {
	if len(batch) == 0 {
		return
	}

	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务失败: %v", err)
		return
	}

	stmt, err := tx.Prepare(`
		INSERT INTO api_request_logs (path, method, status_code, latency_ms, client_ip, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Printf("准备语句失败: %v", err)
		return
	}
	defer stmt.Close()

	for _, item := range batch {
		_, err := stmt.Exec(item.Path, item.Method, item.StatusCode, item.LatencyMs, item.ClientIP, item.UserAgent, item.CreatedAt)
		if err != nil {
			log.Printf("插入日志失败: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("提交事务失败: %v", err)
	}
}
