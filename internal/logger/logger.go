package logger

import (
	"fmt"
	"os"
	"time"
)

// Level 日志级别
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

var levelColors = map[Level]string{
	DEBUG: "\033[36m", // cyan
	INFO:  "\033[32m", // green
	WARN:  "\033[33m", // yellow
	ERROR: "\033[31m", // red
}

const resetColor = "\033[0m"

// Logger 日志记录器
type Logger struct {
	module string
	level  Level
}

// 全局默认日志级别
var globalLevel = INFO

// SetGlobalLevel 设置全局日志级别
func SetGlobalLevel(level Level) {
	globalLevel = level
}

// New 创建新的日志记录器
func New(module string) *Logger {
	return &Logger{
		module: module,
		level:  globalLevel,
	}
}

// log 内部日志方法
func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(format, args...)

	color := levelColors[level]
	levelName := levelNames[level]

	fmt.Fprintf(os.Stderr, "%s%s%s [%s] %s: %s\n",
		color, levelName, resetColor,
		timestamp, l.module, msg)
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

// WithError 带错误的日志
func (l *Logger) WithError(err error) *Logger {
	if err != nil {
		l.Error("error: %v", err)
	}
	return l
}
