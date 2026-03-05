package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	// 测试日志初始化 - 无文件输出
	t.Run("init without file", func(t *testing.T) {
		err := InitLogger("info", "")
		require.NoError(t, err)
		assert.NotNil(t, GetLogger())

		// 测试日志输出
		Info("测试信息日志")
		Debug("测试调试日志")
		Warn("测试警告日志")
		Error("测试错误日志")
	})

	// 测试带文件路径的日志初始化
	t.Run("init with file", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "test.log")

		err := InitLogger("debug", logFile)
		require.NoError(t, err)

		// 写入一些日志
		Info("测试日志写入文件")

		// 验证文件存在
		info, err := os.Stat(logFile)
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0), "日志文件应该有内容")
	})

	// 测试不同日志级别
	t.Run("different log levels", func(t *testing.T) {
		levels := []string{"debug", "info", "warn", "error", "invalid"}
		for _, level := range levels {
			err := InitLogger(level, "")
			require.NoError(t, err, "级别 %s 应该能正常初始化", level)
		}
	})

	// 测试带路径的日志文件创建目录
	t.Run("init with nested path", func(t *testing.T) {
		tmpDir := t.TempDir()
		logFile := filepath.Join(tmpDir, "subdir", "nested", "test.log")

		err := InitLogger("info", logFile)
		require.NoError(t, err)

		// 写入日志后同步，确保数据写入文件
		Info("测试嵌套路径日志")
		GetLogger().Sync()

		// 验证文件存在
		info, err := os.Stat(logFile)
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	})
}

func TestLoggerMethods(t *testing.T) {
	// 先初始化日志
	err := InitLogger("debug", "")
	require.NoError(t, err)

	logger := GetLogger()
	require.NotNil(t, logger)

	// 测试各种日志方法
	t.Run("debug methods", func(t *testing.T) {
		logger.Debug("debug message")
		logger.Debugf("debug %s", "formatted")
	})

	t.Run("info methods", func(t *testing.T) {
		logger.Info("info message")
		logger.Infof("info %s", "formatted")
	})

	t.Run("warn methods", func(t *testing.T) {
		logger.Warn("warn message")
		logger.Warnf("warn %s", "formatted")
	})

	t.Run("error methods", func(t *testing.T) {
		logger.Error("error message")
		logger.Errorf("error %s", "formatted")
	})

	t.Run("with fields", func(t *testing.T) {
		sugar := logger.With("key", "value", "number", 42)
		sugar.Info("message with fields")
	})
}

func TestLoggerSync(t *testing.T) {
	// 测试 Sync 方法
	err := InitLogger("info", "")
	require.NoError(t, err)

	logger := GetLogger()
	require.NotNil(t, logger)

	// 写入一些日志后同步
	Info("测试同步")
	logger.Sync()
}

func TestGlobalLoggerFunctions(t *testing.T) {
	// 测试全局便捷函数
	err := InitLogger("debug", "")
	require.NoError(t, err)

	// 测试全局函数
	Debug("debug全局函数")
	Debugf("debugf全局函数: %s", "test")
	Info("info全局函数")
	Infof("infof全局函数: %s", "test")
	Warn("warn全局函数")
	Warnf("warnf全局函数: %s", "test")
	Error("error全局函数")
	Errorf("errorf全局函数: %s", "test")
	Sync()
}

func TestGetLogger(t *testing.T) {
	// 测试获取未初始化时的默认日志
	// 由于 init() 已经初始化，这里测试 GetLogger 返回非 nil
	logger := GetLogger()
	assert.NotNil(t, logger)
}
