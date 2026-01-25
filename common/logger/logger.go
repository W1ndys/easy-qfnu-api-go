package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logDir, logName, level string) {
	// 1. 确保日志目录存在
	_ = os.MkdirAll(logDir, 0755)

	// 2. 解析日志级别
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

	// ==========================================
	// Handler 1: 控制台输出 (漂亮、彩色、文本格式)
	// ==========================================
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.TimeOnly, // 控制台只显示时间 "15:04:05"，不需要日期，清爽
		NoColor:    false,         // 强制开启颜色
	})

	// ==========================================
	// Handler 2: 文件输出 (JSON 格式，保留所有细节)
	// ==========================================
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, logName),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	fileHandler := slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
		Level: logLevel,
	})

	// ==========================================
	// 3. 混合器: 同时分发给上面两个 Handler
	// ==========================================
	// Fanout 会把一条日志同时发给 Console 和 File，而且格式互不影响
	multiHandler := slogmulti.Fanout(consoleHandler, fileHandler)

	// 4. 设置全局 Logger
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	slog.Info("日志系统已升级", "mode", "Console(Color)+File(JSON)")
}
