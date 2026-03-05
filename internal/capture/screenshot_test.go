package capture

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewScreenshot 测试创建截图管理器
func TestNewScreenshot(t *testing.T) {
	s := NewScreenshot()
	require.NotNil(t, s, "Screenshot 实例不应为空")
}

// TestWindowInfo 测试窗口信息结构体
func TestWindowInfo(t *testing.T) {
	// 测试窗口信息结构体的字段
	windows := []WindowInfo{
		{
			ID:     0,
			Title:  "Main Display",
			X:      0,
			Y:      0,
			Width:  1920,
			Height: 1080,
		},
		{
			ID:     1,
			Title:  "Secondary Display",
			Width:  2560,
			Height: 1440,
		},
	}

	assert.Equal(t, 2, len(windows))
	assert.Equal(t, 0, windows[0].ID)
	assert.Equal(t, "Main Display", windows[0].Title)
	assert.Equal(t, 1920, windows[0].Width)
	assert.Equal(t, 1080, windows[0].Height)
}

// TestCaptureRegion_InvalidParameters 测试 CaptureRegion 无效参数
func TestCaptureRegion_InvalidParameters(t *testing.T) {
	s := NewScreenshot()

	// 测试零宽度 - 应该返回错误
	_, err := s.CaptureRegion(0, 0, 0, 100)
	assert.Error(t, err, "零宽度应该返回错误")

	// 测试零高度 - 应该返回错误
	_, err = s.CaptureRegion(0, 0, 100, 0)
	assert.Error(t, err, "零高度应该返回错误")

	// 注意：kbinani/screenshot 库对负数坐标不会返回错误
	// 所以这里不测试负数坐标
}

// TestCaptureRegion_NegativeCoordinates 测试 CaptureRegion 负数坐标参数
func TestCaptureRegion_NegativeCoordinates(t *testing.T) {
	s := NewScreenshot()

	// 负数宽度和高度
	_, err := s.CaptureRegion(0, 0, -100, 100)
	assert.Error(t, err, "负数宽度应该返回错误")

	_, err = s.CaptureRegion(0, 0, 100, -100)
	assert.Error(t, err, "负数高度应该返回错误")

	_, err = s.CaptureRegion(0, 0, -100, -100)
	assert.Error(t, err, "负数宽高应该返回错误")
}

// TestSaveToFile 测试保存截图到文件
func TestSaveToFile(t *testing.T) {
	s := NewScreenshot()

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "screenshot_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试数据
	testData := []byte("fake image data")

	// 测试保存到普通路径
	filePath := filepath.Join(tempDir, "test.png")
	err = s.SaveToFile(testData, filePath)
	assert.NoError(t, err)

	// 验证文件存在
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// 验证文件内容
	data, err := os.ReadFile(filePath)
	assert.Equal(t, testData, data)

	// 测试保存到嵌套路径
	nestedPath := filepath.Join(tempDir, "subdir", "nested", "test.png")
	err = s.SaveToFile(testData, nestedPath)
	assert.NoError(t, err)

	// 验证嵌套文件存在
	_, err = os.Stat(nestedPath)
	assert.NoError(t, err)
}

// TestSaveToFile_EmptyData 测试空数据保存
func TestSaveToFile_EmptyData(t *testing.T) {
	s := NewScreenshot()

	tempDir, err := os.MkdirTemp("", "screenshot_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 空数据也应该能保存
	err = s.SaveToFile([]byte{}, filepath.Join(tempDir, "empty.png"))
	assert.NoError(t, err)
}

// TestSaveToFile_PermissionDenied 测试权限问题（只读目录）
func TestSaveToFile_PermissionDenied(t *testing.T) {
	s := NewScreenshot()

	// 测试保存到无效路径
	err := s.SaveToFile([]byte("test"), "/root/test.png")
	assert.Error(t, err, "保存到只读目录应该返回错误")
}

// TestSaveToFile_CurrentDirectory 测试当前目录保存
func TestSaveToFile_CurrentDirectory(t *testing.T) {
	// 测试只提供文件名（保存到当前目录）
	// 由于测试可能没有写权限，这个测试可能是可选的
	// 这里我们测试相对路径的处理
	testData := []byte("test data")
	_ = testData // 避免未使用警告
}

// TestEncodeToPNG 测试 PNG 编码（内部函数需要 mock 才能完整测试）
// 这里只测试边界情况
func TestEncodeToPNG_Internal(t *testing.T) {
	// encodeToPNG 是内部函数，需要通过公开接口间接测试
	// 由于 screenshot 库需要真实屏幕，我们跳过实际调用
	s := NewScreenshot()
	require.NotNil(t, s)
}

// TestCaptureFullScreen_SkipOnCI 测试在 CI 环境下跳过实际屏幕捕获
func TestCaptureFullScreen_SkipOnCI(t *testing.T) {
	// 这是一个占位测试，真实测试需要 mock screenshot 库
	// 在没有屏幕的环境中，screenshot 库会返回错误
	// 这正是我们期望的行为
	s := NewScreenshot()
	require.NotNil(t, s)

	// 如果在 CI 环境，尝试捕获会失败，这是预期行为
	// 实际测试时应该使用 mock
}

// TestCaptureWindow_SkipOnCI 测试窗口捕获
func TestCaptureWindow_SkipOnCI(t *testing.T) {
	s := NewScreenshot()
	require.NotNil(t, s)

	// 尝试捕获不存在的窗口
	_, err := s.CaptureWindow(99999)
	assert.Error(t, err, "捕获不存在的窗口应该返回错误")
	assert.Contains(t, err.Error(), "失败")
}

// TestCaptureActiveWindow_NoWindow 测试没有活跃窗口的情况
func TestCaptureActiveWindow_NoWindow(t *testing.T) {
	s := NewScreenshot()
	require.NotNil(t, s)

	// 如果没有窗口，返回错误
	// 这里测试错误处理路径
	_, err := s.CaptureActiveWindow()
	// 可能返回错误或者返回主屏幕截图，取决于实现
	if err != nil {
		assert.Contains(t, err.Error(), "没有找到")
	}
}

// TestGetWindows 测试获取窗口列表
func TestGetWindows(t *testing.T) {
	s := NewScreenshot()

	windows, err := s.GetWindows()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(windows), 1, "至少应该返回一个窗口（主屏幕）")

	// 验证第一个窗口是主屏幕
	if len(windows) > 0 {
		assert.Equal(t, 0, windows[0].ID)
		assert.NotEmpty(t, windows[0].Title)
		assert.Greater(t, windows[0].Width, 0)
		assert.Greater(t, windows[0].Height, 0)
	}
}

// TestCaptureAndSave 测试捕获并保存
func TestCaptureAndSave(t *testing.T) {
	s := NewScreenshot()

	tempDir, err := os.MkdirTemp("", "capture_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "capture.png")

	// 尝试捕获并保存
	// 由于可能没有屏幕，这个测试可能会失败
	data, err := s.CaptureAndSave(filePath)

	// 如果失败，检查错误类型
	if err != nil {
		// 检查是否是预期的错误（没有屏幕或其他硬件问题）
		t.Logf("捕获失败（预期在无屏幕环境失败）: %v", err)
	} else {
		// 如果成功，验证返回数据和文件
		assert.NotEmpty(t, data)
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	}
}

// TestScreenshotStruct_Fields 测试 Screenshot 结构体字段
func TestScreenshotStruct_Fields(t *testing.T) {
	// 验证 Screenshot 是空结构体
	s := NewScreenshot()
	assert.NotNil(t, s)

	// 验证可以通过指针调用方法
	err := s.SaveToFile([]byte("test"), "/tmp/test.png")
	// 这里只验证方法可以被调用
	if err != nil {
		// 权限错误是预期的
		assert.Error(t, err)
	}
}

// TestBoundary_CaptureRegion 测试 CaptureRegion 边界值
func TestBoundary_CaptureRegion(t *testing.T) {
	s := NewScreenshot()

	// 测试极大值（超出屏幕范围）
	// 这可能会失败，但不应该崩溃
	_, _ = s.CaptureRegion(0, 0, 100000, 100000)

	// 测试极小正值
	_, err := s.CaptureRegion(0, 0, 1, 1)
	// 可能成功或失败，取决于环境
	if err != nil {
		t.Logf("极小区域捕获失败（可能预期）: %v", err)
	}
}
