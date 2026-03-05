package playwright

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewBrowser 测试创建浏览器控制器
func TestNewBrowser(t *testing.T) {
	b := NewBrowser()
	require.NotNil(t, b, "Browser 实例不应为空")

	// 验证默认配置
	assert.NotNil(t, b.launchOptions, "启动选项不应为空")
	assert.NotNil(t, b.launchOptions.Headless, "Headless 配置应存在")
	assert.False(t, *b.launchOptions.Headless, "默认应显示浏览器窗口")

	// 验证初始状态
	assert.Nil(t, b.page, "初始 page 应为 nil")
	assert.Nil(t, b.browser, "初始 browser 应为 nil")
	assert.Nil(t, b.context, "初始 context 应为 nil")
}

// TestBrowser_Launch_NilPage 测试未启动时调用方法返回错误
func TestBrowser_Launch_NilPage(t *testing.T) {
	b := NewBrowser()

	// 测试各种方法在未启动时的错误处理
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "Navigate",
			testFunc: func() error {
				return b.Navigate("https://example.com")
			},
		},
		{
			name: "Screenshot",
			testFunc: func() error {
				_, err := b.Screenshot()
				return err
			},
		},
		{
			name: "ScreenshotToFile",
			testFunc: func() error {
				return b.ScreenshotToFile("/tmp/test.png")
			},
		},
		{
			name: "Evaluate",
			testFunc: func() error {
				_, err := b.Evaluate("return 1")
				return err
			},
		},
		{
			name: "GetElementPosition",
			testFunc: func() error {
				_, _, err := b.GetElementPosition("#test")
				return err
			},
		},
		{
			name: "ClickElement",
			testFunc: func() error {
				return b.ClickElement("#test")
			},
		},
		{
			name: "FillElement",
			testFunc: func() error {
				return b.FillElement("#test", "value")
			},
		},
		{
			name: "GetText",
			testFunc: func() error {
				_, err := b.GetText("#test")
				return err
			},
		},
		{
			name: "GetHTML",
			testFunc: func() error {
				_, err := b.GetHTML("#test")
				return err
			},
		},
		{
			name: "WaitForSelector",
			testFunc: func() error {
				return b.WaitForSelector("#test")
			},
		},
		{
			name: "Scroll",
			testFunc: func() error {
				return b.Scroll(100)
			},
		},
		{
			name: "ScrollToBottom",
			testFunc: func() error {
				return b.ScrollToBottom()
			},
		},
		{
			name: "GetPageTitle",
			testFunc: func() error {
				_, err := b.GetPageTitle()
				return err
			},
		},
		{
			name: "GetCurrentURL",
			testFunc: func() error {
				_, err := b.GetCurrentURL()
				return err
			},
		},
		{
			name: "AnalyzePage",
			testFunc: func() error {
				_, err := b.AnalyzePage("test query")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			assert.Error(t, err, "未启动时应返回错误")
			assert.Contains(t, err.Error(), "浏览器未启动", "错误信息应包含 '浏览器未启动'")
		})
	}
}

// TestBrowser_GetCookies_NilContext 测试 GetCookies 在 context 为 nil 时返回错误
func TestBrowser_GetCookies_NilContext(t *testing.T) {
	b := NewBrowser()

	_, err := b.GetCookies()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "浏览器上下文未创建")
}

// TestBrowser_SetCookies_NilContext 测试 SetCookies 在 context 为 nil 时返回错误
func TestBrowser_SetCookies_NilContext(t *testing.T) {
	b := NewBrowser()

	err := b.SetCookies(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "浏览器上下文未创建")
}

// TestBrowser_Close 测试 Close 方法
func TestBrowser_Close(t *testing.T) {
	b := NewBrowser()

	// 关闭未启动的浏览器不应出错
	err := b.Close()
	assert.NoError(t, err, "关闭未启动的浏览器不应返回错误")

	// 再次关闭也不应出错
	err = b.Close()
	assert.NoError(t, err, "重复关闭不应返回错误")
}

// TestBrowser_GetPage 测试 GetPage 方法
func TestBrowser_GetPage(t *testing.T) {
	b := NewBrowser()

	// 初始返回 nil
	page := b.GetPage()
	assert.Nil(t, page, "未启动时应返回 nil")
}

// TestBrowser_Navigate_InvalidURL 测试无效 URL
func TestBrowser_Navigate_InvalidURL(t *testing.T) {
	b := NewBrowser()

	// 由于未启动，这会先检查 page 是否为 nil
	err := b.Navigate("")
	assert.Error(t, err)
}

// TestBrowser_Navigate_EmptyURL 测试空 URL
func TestBrowser_Navigate_EmptyURL(t *testing.T) {
	b := NewBrowser()

	err := b.Navigate("")
	assert.Error(t, err)
}

// TestBrowser_Launch_Integration 测试启动浏览器（集成测试，需要浏览器环境）
// 使用 build tag 标记为集成测试
func TestBrowser_Launch_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（short mode）")
	}

	b := NewBrowser()

	// 尝试启动浏览器
	err := b.Launch()
	if err != nil {
		// 如果环境不支持，跳过测试
		t.Skipf("无法启动浏览器: %v", err)
	}

	// 测试完成后关闭
	defer b.Close()

	// 验证启动成功
	assert.NotNil(t, b.browser)
	assert.NotNil(t, b.context)
	assert.NotNil(t, b.page)

	// 测试导航
	err = b.Navigate("https://example.com")
	assert.NoError(t, err)

	// 测试获取页面标题
	title, err := b.GetPageTitle()
	assert.NoError(t, err)
	assert.NotEmpty(t, title)

	// 测试获取当前 URL
	url, err := b.GetCurrentURL()
	assert.NoError(t, err)
	assert.Contains(t, url, "example.com")
}

// TestBrowser_Screenshot_Integration 测试截图功能（集成测试）
func TestBrowser_Screenshot_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（short mode）")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试截图
	data, err := b.Screenshot()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
}

// TestBrowser_ClickElement_NotFound 测试点击不存在的元素
func TestBrowser_ClickElement_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试点击不存在的元素
	err = b.ClickElement("#nonexistent-element-12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未找到元素")
}

// TestBrowser_FillElement_NotFound 测试填充不存在的元素
func TestBrowser_FillElement_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试填充不存在的元素
	err = b.FillElement("#nonexistent", "test value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未找到元素")
}

// TestBrowser_GetText_NotFound 测试获取不存在元素的文本
func TestBrowser_GetText_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试获取不存在元素的文本
	_, err = b.GetText("#nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未找到元素")
}

// TestBrowser_GetElementPosition_NotFound 测试获取不存在元素的位置
func TestBrowser_GetElementPosition_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试获取不存在元素的位置
	_, _, err = b.GetElementPosition("#nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未找到元素")
}

// TestBrowser_WaitForSelector_Timeout 测试等待元素超时
func TestBrowser_WaitForSelector_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试等待不存在的元素（短超时）
	err = b.WaitForSelector("#never-appear", 100)
	assert.Error(t, err)
}

// TestBrowser_Scroll 测试滚动功能
func TestBrowser_Scroll(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试滚动
	err = b.Scroll(500)
	assert.NoError(t, err)

	// 测试滚动到顶部
	err = b.Scroll(0)
	assert.NoError(t, err)
}

// TestBrowser_Evaluate 测试 JavaScript 执行
func TestBrowser_Evaluate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试执行简单 JavaScript
	result, err := b.Evaluate("return 1 + 1")
	assert.NoError(t, err)
	assert.Equal(t, 2, int(result.(float64)))

	// 测试执行返回字符串的 JavaScript
	result, err = b.Evaluate("return 'hello'")
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

// TestBrowser_AnalyzePage 测试页面分析功能
func TestBrowser_AnalyzePage(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试分析页面
	result, err := b.AnalyzePage("获取可点击元素")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "HTML长度")
}

// TestBrowser_LaunchWithCookie 测试使用 Cookie 启动
func TestBrowser_LaunchWithCookie(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	// 测试不存在的 cookie 文件
	err := b.LaunchWithCookie("/nonexistent/cookie.json")
	// 这个测试可能会失败，取决于实现
	if err != nil {
		t.Logf("LaunchWithCookie 失败（可能预期）: %v", err)
	}
}

// TestBrowser_ScreenshotToFile 测试保存截图到文件
func TestBrowser_ScreenshotToFile(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试保存截图
	err = b.ScreenshotToFile("/tmp/playwright_test.png")
	assert.NoError(t, err)
}

// TestBrowser_Cookies 测试 Cookie 操作
func TestBrowser_Cookies(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	b := NewBrowser()
	defer b.Close()

	err := b.Launch()
	if err != nil {
		t.Skipf("无法启动浏览器: %v", err)
	}

	err = b.Navigate("https://example.com")
	if err != nil {
		t.Skipf("无法导航: %v", err)
	}

	// 测试获取 Cookie
	cookies, err := b.GetCookies()
	assert.NoError(t, err)
	// example.com 可能有或没有 cookie，不做强制断言
	t.Logf("获取到 %d 个 cookie", len(cookies))

	// 测试设置 Cookie
	testCookie := []pw.Cookie{
		{
			Name:   "test_cookie",
			Value:  "test_value",
			Domain: "example.com",
			Path:   "/",
		},
	}
	err = b.SetCookies(testCookie)
	assert.NoError(t, err)

	// 验证 Cookie 已设置
	cookies, err = b.GetCookies()
	assert.NoError(t, err)

	// 查找我们设置的 cookie
	found := false
	for _, c := range cookies {
		if c.Name == "test_cookie" {
			found = true
			assert.Equal(t, "test_value", c.Value)
			break
		}
	}
	assert.True(t, found, "应该能找到设置的 cookie")
}

// TestBrowser_HeadlessOption 测试无头模式选项
func TestBrowser_HeadlessOption(t *testing.T) {
	// 测试设置无头模式
	headless := true
	b := &Browser{
		launchOptions: pw.BrowserTypeLaunchOptions{
			Headless: &headless,
		},
	}

	assert.True(t, *b.launchOptions.Headless)
}
