// Package service 提供业务逻辑服务
// 图片简历发送模块：实现自动发送图片简历的功能
package service

import (
	"fmt"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/playwright-community/playwright-go"

	"github.com/loks666/get_jobs/internal/config"
)

// ImageResumeSender 图片简历发送器
type ImageResumeSender struct {
	page         *playwright.Page
	resumePath   string
	isCompatible bool // 当前平台是否支持
}

// NewImageResumeSender 创建图片简历发送器
func NewImageResumeSender(page *playwright.Page) *ImageResumeSender {
	resumePath := config.GetString("delivery.img_resume_path")
	if resumePath == "" {
		resumePath = "./resources/resume.jpg"
	}

	return &ImageResumeSender{
		page:         page,
		resumePath:   resumePath,
		isCompatible: true, // macOS 和 Windows 都支持
	}
}

// Send 发送图片简历
// 流程：
// 1. 查找发送图片按钮并点击
// 2. 等待文件选择对话框出现
// 3. 使用系统级 API 输入文件路径
// 4. 确认发送
func (s *ImageResumeSender) Send() error {
	if s.page == nil {
		return fmt.Errorf("页面未初始化")
	}

	// 检查图片是否存在
	if _, err := os.Stat(s.resumePath); err != nil {
		return fmt.Errorf("图片简历文件不存在: %w", err)
	}

	config.Info("开始发送图片简历: ", s.resumePath)

	// 步骤1: 查找并点击发送图片按钮
	if err := s.clickImageButton(); err != nil {
		return fmt.Errorf("找不到发送图片按钮: %w", err)
	}

	// 等待文件对话框出现
	time.Sleep(1 * time.Second)

	// 步骤2: 使用系统级 API 输入文件路径
	if err := s.selectFile(); err != nil {
		return fmt.Errorf("选择文件失败: %w", err)
	}

	// 等待文件上传
	time.Sleep(2 * time.Second)

	// 步骤3: 点击发送按钮
	if err := s.clickSendButton(); err != nil {
		return fmt.Errorf("点击发送按钮失败: %w", err)
	}

	config.Info("图片简历发送成功")
	return nil
}

// clickImageButton 查找并点击发送图片按钮
func (s *ImageResumeSender) clickImageButton() error {
	// 尝试多种选择器
	selectors := []string{
		".btn-img",              // 图片按钮
		"[class*='image']",     // 包含 image 的 class
		"[class*='picture']",   // 包含 picture 的 class
		"[class*='resume']",    // 包含 resume 的 class
		"button[title*='图片']", // title 包含图片
		"button[title*='简历']", // title 包含简历
	}

	for _, selector := range selectors {
		locator := (*s.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First().Click()
		}
	}

	return fmt.Errorf("找不到发送图片按钮")
}

// selectFile 选择文件
// 由于 Playwright 无法处理系统文件对话框
// 使用 robotgo 来模拟键盘输入文件路径
func (s *ImageResumeSender) selectFile() error {
	// macOS 和 Windows 处理方式不同
	// macOS: 使用 Cmd+Shift+G 打开文件路径输入框
	// Windows: 直接粘贴文件路径

	// 获取绝对路径
	absPath, err := getAbsolutePath(s.resumePath)
	if err != nil {
		return err
	}

	// 使用 robotgo 输入文件路径
	// 注意：这种方法可能在不同系统上表现不同

	// 方法1: 直接粘贴路径 (适用于 macOS)
	robotgo.TypeStr(absPath)
	time.Sleep(500 * time.Millisecond)

	// 按回车确认
	robotgo.KeyTap("return")
	time.Sleep(1 * time.Second)

	return nil
}

// clickSendButton 点击发送按钮
func (s *ImageResumeSender) clickSendButton() error {
	// 查找发送按钮
	selectors := []string{
		".btn-send",            // 发送按钮
		".send-btn",            // 发送按钮
		"button:has-text('发送')", // 文本包含发送
		"[class*='send']",     // 包含 send
	}

	for _, selector := range selectors {
		locator := (*s.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First().Click()
		}
	}

	// 备用方案：使用回车发送
	robotgo.KeyTap("return")

	return nil
}

// CheckResumeExists 检查图片简历是否存在
func (s *ImageResumeSender) CheckResumeExists() (bool, error) {
	_, err := os.Stat(s.resumePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SetResumePath 设置图片简历路径
func (s *ImageResumeSender) SetResumePath(path string) {
	s.resumePath = path
}

// getAbsolutePath 获取绝对路径
func getAbsolutePath(path string) (string, error) {
	// 如果已经是绝对路径，直接返回
	if path[0] == '/' || (len(path) > 2 && path[1] == ':') {
		return path, nil
	}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir + "/" + path, nil
}
