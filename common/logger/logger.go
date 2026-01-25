package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化全局日志配置
// logDir: 日志存放目录 (如 "./logs")
// logName: 日志文件名 (如 "app.log")
// level: 日志级别 (debug, info, warn, error)
func InitLogger(logDir, logName, level string) {
	// 1. 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("无法创建日志目录: " + err.Error())
	}

	// 2. 配置 Lumberjack (日志切割)
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, logName),
		MaxSize:    10,   // 每个日志文件最大 10MB
		MaxBackups: 5,    // 保留最近 5 个文件
		MaxAge:     30,   // 保留最近 30 天
		Compress:   true, // 是否压缩旧文件
	}

	// 3. 配置多重输出 (同时输出到 控制台 和 文件)
	// io.MultiWriter 可以把流分发给多个目的地
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// 4. 设置日志级别
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

	// 5. 创建 Handler (使用 JSON 格式)
	opts := &slog.HandlerOptions{
		Level: logLevel,
		// AddSource: true, // 开发环境可以开启，会显示文件名和行号，生产环境建议关闭以提升性能
	}
	handler := slog.NewJSONHandler(multiWriter, opts)

	// 6. 设置为全局默认 Logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("日志系统初始化完成", "dir", logDir, "level", level)
}
