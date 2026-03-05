package config

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 日志管理器
type Logger struct {
	sugar  *zap.SugaredLogger
	logger *zap.Logger
}

// 全局日志实例
var globalLogger *Logger

// InitLogger 初始化日志系统
func InitLogger(level string, filePath string) error {
	// 创建日志目录
	if filePath != "" {
		dir := filePath
		for i := len(dir) - 1; i >= 0; i-- {
			if dir[i] == '/' || dir[i] == '\\' {
				dir = dir[:i]
				break
			}
		}
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("创建日志目录失败: %w", err)
			}
		}
	}

	// 解析日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	var cores []zapcore.Core
	if filePath != "" {
		// 文件写入器
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %w", err)
		}
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		writer := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(fileEncoder, writer, zapLevel))
	}

	// 控制台写入器
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLevel))

	// 创建 logger
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))

	globalLogger = &Logger{
		sugar:  logger.Sugar(),
		logger: logger,
	}

	log.SetFlags(0)
	log.SetOutput(zapcore.AddSync(os.Stdout))

	return nil
}

// GetLogger 获取全局日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		// 如果未初始化，创建默认日志
		InitLogger("info", "")
	}
	return globalLogger
}

// Debug 调试日志
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debugf 格式化调试日志
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Info 信息日志
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warn 警告日志
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Error 错误日志
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Sync 同步日志缓冲
func (l *Logger) Sync() {
	l.logger.Sync()
}

// With 创建带有额外字段的日志
func (l *Logger) With(fields ...interface{}) *zap.SugaredLogger {
	return l.sugar.With(fields...)
}

// 便捷函数
var (
	Debug  func(args ...interface{})
	Debugf func(template string, args ...interface{})
	Info   func(args ...interface{})
	Infof  func(template string, args ...interface{})
	Warn   func(args ...interface{})
	Warnf  func(template string, args ...interface{})
	Error  func(args ...interface{})
	Errorf func(template string, args ...interface{})
	Fatal  func(args ...interface{})
	Fatalf func(template string, args ...interface{})
	Sync   func()
	With   func(fields ...interface{}) *zap.SugaredLogger
)

func init() {
	// 设置便捷函数
	Debug = func(args ...interface{}) { GetLogger().Debug(args...) }
	Debugf = func(template string, args ...interface{}) { GetLogger().Debugf(template, args...) }
	Info = func(args ...interface{}) { GetLogger().Info(args...) }
	Infof = func(template string, args ...interface{}) { GetLogger().Infof(template, args...) }
	Warn = func(args ...interface{}) { GetLogger().Warn(args...) }
	Warnf = func(template string, args ...interface{}) { GetLogger().Warnf(template, args...) }
	Error = func(args ...interface{}) { GetLogger().Error(args...) }
	Errorf = func(template string, args ...interface{}) { GetLogger().Errorf(template, args...) }
	Fatal = func(args ...interface{}) { GetLogger().Fatal(args...) }
	Fatalf = func(template string, args ...interface{}) { GetLogger().Fatalf(template, args...) }
	Sync = func() { GetLogger().Sync() }
	With = func(fields ...interface{}) *zap.SugaredLogger { return GetLogger().With(fields...) }
}
