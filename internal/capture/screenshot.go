package capture

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"

	"github.com/yahao333/get_jobs/internal/config"
)

// Screenshot 截图管理器
type Screenshot struct{}

// NewScreenshot 创建截图管理器
func NewScreenshot() *Screenshot {
	return &Screenshot{}
}

// CaptureFullScreen 捕获整个屏幕
func (s *Screenshot) CaptureFullScreen() ([]byte, error) {
	// 获取屏幕数量
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return nil, fmt.Errorf("没有可用的屏幕")
	}

	// 捕获第一个屏幕（主屏幕）
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		return nil, fmt.Errorf("捕获屏幕失败: %w", err)
	}

	return s.encodeToPNG(img)
}

// CaptureWindow 捕获指定窗口
func (s *Screenshot) CaptureWindow(windowID int) ([]byte, error) {
	// 如果 windowID 为 0，通常表示主屏幕（GetWindows 返回的 ID）
	if windowID == 0 {
		return s.CaptureFullScreen()
	}

	// 获取窗口边界
	x, y, w, h := robotgo.GetBounds(windowID)
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("捕获窗口失败: 无效的窗口 ID 或窗口不可见 (ID: %d)", windowID)
	}

	// 捕获窗口区域
	img, err := robotgo.CaptureImg(x, y, w, h)
	if err != nil {
		return nil, fmt.Errorf("捕获窗口失败: %w", err)
	}

	// 转换为 *image.RGBA
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return nil, fmt.Errorf("捕获窗口失败: 图像格式不支持")
	}

	return s.encodeToPNG(rgba)
}

// CaptureRegion 捕获指定区域
func (s *Screenshot) CaptureRegion(x, y, width, height int) ([]byte, error) {
	// 获取屏幕数量
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return nil, fmt.Errorf("没有可用的屏幕")
	}

	// 捕获指定区域
	img, err := screenshot.Capture(x, y, width, height)
	if err != nil {
		return nil, fmt.Errorf("捕获区域失败: %w", err)
	}

	return s.encodeToPNG(img)
}

// CaptureActiveWindow 捕获当前活跃窗口
func (s *Screenshot) CaptureActiveWindow() ([]byte, error) {
	// 获取当前活跃窗口 PID
	pid := robotgo.GetPid()
	if pid <= 0 {
		return nil, fmt.Errorf("没有找到活跃窗口")
	}

	return s.CaptureWindow(int(pid))
}

// GetWindows 获取所有窗口信息
func (s *Screenshot) GetWindows() ([]WindowInfo, error) {
	// 注意：kbinani/screenshot 库不直接提供窗口列表功能
	// 需要结合其他方式获取窗口
	// 这里返回一个简化的实现

	// 获取屏幕区域
	display := screenshot.GetDisplayBounds(0)
	/*
		if err != nil {
			return nil, err
		}
	*/

	windows := []WindowInfo{
		{
			ID:     0,
			Title:  "Main Display",
			X:      display.Min.X,
			Y:      display.Min.Y,
			Width:  display.Dx(),
			Height: display.Dy(),
		},
	}

	return windows, nil
}

// WindowInfo 窗口信息
type WindowInfo struct {
	ID     int
	Title  string
	X, Y   int
	Width  int
	Height int
}

// encodeToPNG 将图像编码为 PNG
func (s *Screenshot) encodeToPNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("编码 PNG 失败: %w", err)
	}
	return buf.Bytes(), nil
}

// SaveToFile 保存截图到文件
func (s *Screenshot) SaveToFile(data []byte, filePath string) error {
	// 确保目录存在
	dir := filePath
	for i := len(dir) - 1; i >= 0; i-- {
		if dir[i] == '/' || dir[i] == '\\' {
			dir = dir[:i]
			break
		}
	}
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	return os.WriteFile(filePath, data, 0644)
}

// CaptureAndSave 捕获并保存截图
func (s *Screenshot) CaptureAndSave(filePath string) ([]byte, error) {
	data, err := s.CaptureFullScreen()
	if err != nil {
		return nil, err
	}

	if err := s.SaveToFile(data, filePath); err != nil {
		return nil, err
	}

	config.Info("截图已保存: ", filePath)
	return data, nil
}
