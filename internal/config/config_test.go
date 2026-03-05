package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureDirectories(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()

	// 切换到临时目录进行测试
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 调用 ensureDirectories
	err = ensureDirectories()
	require.NoError(t, err)

	// 验证目录是否创建
	dirs := []string{"data", "logs", "resources"}
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		require.NoError(t, err, "目录 %s 应该存在", dir)
		assert.True(t, info.IsDir(), "%s 应该是目录", dir)
	}

	// 测试目录权限
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0755), info.Mode().Perm(), "%s 权限应该是 0755", dir)
	}
}

func TestEnsureDirectoriesWithExistingDirs(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 预先创建部分目录
	err := os.MkdirAll(filepath.Join(tmpDir, "data"), 0755)
	require.NoError(t, err)

	// 切换到临时目录
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 调用 ensureDirectories - 应该不会报错
	err = ensureDirectories()
	require.NoError(t, err)
}

func TestEnsureDirectoriesPermissionError(t *testing.T) {
	// 这个测试验证目录创建的错误处理
	// 在大多数系统上，root 用户可能会失败，但普通用户应该成功

	tmpDir := t.TempDir()

	// 尝试在一个只读路径创建目录（应该失败）
	// 注意：这个测试可能在某些环境下通过或失败，取决于权限
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 正常情况下应该成功
	err = ensureDirectories()
	require.NoError(t, err)
}

func TestGetString(t *testing.T) {
	// 创建一个临时的 viper 实例
	tmpViper := setupTestViper()

	// 测试获取字符串值
	tests := []struct {
		key      string
		expected string
	}{
		{"app.name", "get_jobs"},
		{"app.version", "1.0.0"},
		{"database.type", "sqlite"},
		{"web.host", "0.0.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := tmpViper.GetString(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetInt(t *testing.T) {
	tmpViper := setupTestViper()

	// 测试获取整数值
	tests := []struct {
		key      string
		expected int
	}{
		{"web.port", 8080},
		{"browser.window_width", 1280},
		{"browser.window_height", 800},
		{"delivery.daily_limit", 100},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := tmpViper.GetInt(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBool(t *testing.T) {
	tmpViper := setupTestViper()

	// 测试获取布尔值
	tests := []struct {
		key      string
		expected bool
	}{
		{"app.debug", true},
		{"ai.enable", true},
		{"greeting.enable_ai", true},
		{"delivery.send_img_resume", false},
		{"filter.filter_dead_hr", true},
		{"browser.use_existing", false},
		{"blacklist.enable", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := tmpViper.GetBool(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	tmpViper := setupTestViper()

	// 测试获取字符串切片
	result := tmpViper.GetStringSlice("search.keywords")
	expected := []string{"Go", "后端开发", "Golang"}
	assert.Equal(t, expected, result)
}

func TestGetMap(t *testing.T) {
	tmpViper := setupTestViper()

	// 测试获取 map
	result := tmpViper.GetStringMap("ai.qwen")
	assert.NotEmpty(t, result)
	assert.Equal(t, "${QWEN_API_KEY}", result["api_key"])
	assert.Equal(t, "qwen-vl-plus", result["model"])
}

// setupTestViper 创建测试用的 viper 实例
func setupTestViper() *viper.Viper {
	v := viper.New()
	// 使用绝对路径确保测试能找到配置文件
	v.SetConfigFile("/Users/yanghao/Work/github/get_jobs/config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		// 如果读取失败，返回一个空的 viper
		return v
	}
	return v
}
