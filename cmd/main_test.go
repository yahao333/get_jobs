package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/storage"
)

// TestMainInitialization 集成测试 - 主程序初始化流程
func TestMainInitialization(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()

	// 创建临时配置文件
	configContent := `
app:
  name: "get_jobs_test"
  version: "1.0.0"
  debug: true
  log_level: "debug"
  log_file: ""

database:
  type: "sqlite"
  path: "` + filepath.Join(tmpDir, "test.db") + `"

web:
  host: "127.0.0.1"
  port: 9999
`
	configFile := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 切换到临时目录
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 测试配置加载
	t.Run("load config", func(t *testing.T) {
		err := config.LoadConfig("config.yaml")
		require.NoError(t, err, "配置加载应该成功")

		// 验证配置值
		assert.Equal(t, "get_jobs_test", config.GetString("app.name"))
		assert.Equal(t, "1.0.0", config.GetString("app.version"))
		assert.Equal(t, true, config.GetBool("app.debug"))
	})

	// 测试数据库初始化
	t.Run("init database", func(t *testing.T) {
		err := storage.InitDB()
		require.NoError(t, err, "数据库初始化应该成功")
		assert.NotNil(t, storage.GetDB())
	})

	// 测试数据库连接
	t.Run("database connection", func(t *testing.T) {
		db := storage.GetDB()
		err := db.Raw("SELECT 1").Error
		require.NoError(t, err, "数据库连接应该正常")
	})
}

// TestMainInitializationWithInvalidConfig 测试无效配置
func TestMainInitializationWithInvalidConfig(t *testing.T) {
	// 测试不存在的配置文件
	err := config.LoadConfig("nonexistent.yaml")
	assert.Error(t, err, "加载不存在的配置文件应该失败")
}

// TestMainInitializationWithNestedPath 测试嵌套路径配置
func TestMainInitializationWithNestedPath(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建嵌套目录
	nestedDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(nestedDir, 0755)

	// 创建配置文件
	configContent := `
app:
  name: "nested_test"
  version: "1.0.0"
  debug: false
  log_level: "info"
  log_file: ""

database:
  type: "sqlite"
  path: "` + filepath.Join(tmpDir, "data", "test.db") + `"

web:
  host: "0.0.0.0"
  port: 8888
`
	configFile := filepath.Join(nestedDir, "config.yaml")
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 切换到嵌套目录
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(nestedDir)
	require.NoError(t, err)

	// 加载配置
	err = config.LoadConfig("config.yaml")
	require.NoError(t, err)

	// 初始化数据库
	err = storage.InitDB()
	require.NoError(t, err)
}
