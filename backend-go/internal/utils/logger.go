package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xiaocainiao633/Code_agent/backend-go/internal/config"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// SimpleLogger 简单日志实现
type SimpleLogger struct {
	level  string
	output *os.File
}

// NewLogger 创建新的日志实例
func NewLogger(cfg *config.LoggingConfig) (Logger, error) {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.Output)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &SimpleLogger{
		level:  cfg.Level,
		output: file,
	}, nil
}

// Debug 调试日志
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.log("DEBUG", msg, args...)
}

// Info 信息日志
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.log("INFO", msg, args...)
}

// Warn 警告日志
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.log("WARN", msg, args...)
}

// Error 错误日志
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.log("ERROR", msg, args...)
}

// Fatal 致命错误日志
func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.log("FATAL", msg, args...)
	os.Exit(1)
}

// log 内部日志方法
func (l *SimpleLogger) log(level, msg string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var formattedMsg string
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	} else {
		formattedMsg = msg
	}

	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, formattedMsg)
	
	// 同时输出到文件和控制台
	l.output.WriteString(logEntry)
	log.Print(logEntry) // 也输出到控制台
}

// shouldLog 判断是否应该记录日志
func (l *SimpleLogger) shouldLog(level string) bool {
	levelPriority := map[string]int{
		"DEBUG": 0,
		"INFO":  1,
		"WARN":  2,
		"ERROR": 3,
		"FATAL": 4,
	}

	currentLevel := levelPriority[l.level]
	msgLevel := levelPriority[level]
	
	return msgLevel >= currentLevel
}

// Close 关闭日志文件
func (l *SimpleLogger) Close() error {
	if l.output != nil {
		return l.output.Close()
	}
	return nil
}

// GlobalLogger 全局日志实例
var GlobalLogger Logger

// InitLogger 初始化全局日志
func InitLogger(cfg *config.LoggingConfig) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}
	GlobalLogger = logger
	return nil
}

// Debug 全局调试日志
func Debug(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Debug(msg, args...)
	}
}

// Info 全局信息日志
func Info(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Info(msg, args...)
	}
}

// Warn 全局警告日志
func Warn(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Warn(msg, args...)
	}
}

// Error 全局错误日志
func Error(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Error(msg, args...)
	}
}

// Fatal 全局致命错误日志
func Fatal(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Fatal(msg, args...)
	}
}