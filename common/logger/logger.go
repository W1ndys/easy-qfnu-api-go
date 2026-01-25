package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化全局日志配置
// logName 建议只传前缀，例如 "qfnu-api" (不要带 .log)
func InitLogger(logDir, logNamePrefix, level string) {
	// 1. 确保日志目录存在
	_ = os.MkdirAll(logDir, 0755)

	// 2. 生成带时间戳的文件名
	// Go 的时间格式化必须用固定的参考时间: 2006-01-02 15:04:05
	// Windows文件名不支持冒号，所以用 15-04-05
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// 最终文件名: qfnu-api-2026-01-26_12-57-05.log
	fullFileName := fmt.Sprintf("%s-%s.log", logNamePrefix, timestamp)
	logPath := filepath.Join(logDir, fullFileName)

	// 3. 解析日志级别
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// 4. 配置 Handler

	// A. 控制台输出
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
	})

	// B. 文件输出 (Lumberjack)
	// 注意: 因为每次启动文件名都不同，Lumberjack 的 MaxBackups 对"历史运行"的清理可能失效
	// 它只能负责切割"本次运行"产生的大文件。
	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // 10MB 切割
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	fileHandler := slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
		Level: logLevel,
	})

	// 5. 组合并设置默认
	multiHandler := slogmulti.Fanout(consoleHandler, fileHandler)
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	// 打印一下日志文件的位置，方便你找
	slog.Info("日志系统初始化完成", "file", logPath)
}
