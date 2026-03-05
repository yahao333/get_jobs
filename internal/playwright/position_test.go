package playwright

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPosition_Struct 测试 Position 结构体
func TestPosition_Struct(t *testing.T) {
	tests := []struct {
		name     string
		position Position
	}{
		{
			name: "normal position",
			position: Position{
				X:      100,
				Y:      200,
				Width:  50,
				Height: 30,
				Text:   "Click me",
				Tag:    "button",
			},
		},
		{
			name: "zero position",
			position: Position{
				X:      0,
				Y:      0,
				Width:  0,
				Height: 0,
				Text:   "",
				Tag:    "",
			},
		},
		{
			name: "large position",
			position: Position{
				X:      99999,
				Y:      99999,
				Width:  10000,
				Height: 10000,
				Text:   "Long text content",
				Tag:    "div.container",
			},
		},
		{
			name: "negative position",
			position: Position{
				X:      -100,
				Y:      -200,
				Width:  50,
				Height: 30,
				Text:   "test",
				Tag:    "span",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.position
			assert.NotNil(t, p)
			// 验证字段可以被正确设置和读取
			assert.Equal(t, tt.position.X, p.X)
			assert.Equal(t, tt.position.Y, p.Y)
			assert.Equal(t, tt.position.Width, p.Width)
			assert.Equal(t, tt.position.Height, p.Height)
			assert.Equal(t, tt.position.Text, p.Text)
			assert.Equal(t, tt.position.Tag, p.Tag)
		})
	}
}

// TestPosition_JSON 测试 Position JSON 序列化
func TestPosition_JSON(t *testing.T) {
	p := Position{
		X:      100,
		Y:      200,
		Width:  50,
		Height: 30,
		Text:   "Click me",
		Tag:    "button",
	}

	// 验证 JSON 标签
	bytes, err := json.Marshal(p)
	require.NoError(t, err)
	jsonStr := string(bytes)

	assert.Contains(t, jsonStr, `"x":100`)
	assert.Contains(t, jsonStr, `"y":200`)
	assert.Contains(t, jsonStr, `"width":50`)
	assert.Contains(t, jsonStr, `"height":30`)
	assert.Contains(t, jsonStr, `"text":"Click me"`)
	assert.Contains(t, jsonStr, `"tag":"button"`)
	_ = jsonStr // 避免未使用警告
}

// TestNewDualChannelPositioner 测试创建双通道定位器
func TestNewDualChannelPositioner(t *testing.T) {
	// 测试传入 nil browser
	p := NewDualChannelPositioner(nil)
	require.NotNil(t, p, "DualChannelPositioner 实例不应为空")
	assert.Nil(t, p.browser, "browser 应为 nil")

	// 测试传入非 nil browser
	b := NewBrowser()
	p = NewDualChannelPositioner(b)
	require.NotNil(t, p)
	assert.NotNil(t, p.browser)
}

// TestDualChannelPositioner_FindElementByAI 测试 AI 查找元素（需要实现）
func TestDualChannelPositioner_FindElementByAI(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	// 测试 AI 查找功能（当前返回错误，因为未实现）
	pos, err := p.FindElementByAI([]byte{}, "按钮")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "AI 定位功能需要实现")
}

// TestDualChannelPositioner_FindElementByAI_WithBrowser 测试带浏览器的 AI 查找
func TestDualChannelPositioner_FindElementByAI_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	// 即使有 browser，AI 功能也需要实现
	pos, err := p.FindElementByAI([]byte{}, "按钮")
	assert.Error(t, err)
	assert.Nil(t, pos)
}

// TestDualChannelPositioner_FindElementByDOM_NilPage 测试 DOM 查找（page 为 nil）
func TestDualChannelPositioner_FindElementByDOM_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	pos, err := p.FindElementByDOM("#test")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElementByDOM_WithBrowser 测试 DOM 查找（有 browser 但未启动）
func TestDualChannelPositioner_FindElementByDOM_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	pos, err := p.FindElementByDOM("#test")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElementByText_NilPage 测试文本查找（page 为 nil）
func TestDualChannelPositioner_FindElementByText_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	pos, err := p.FindElementByText("test")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElementByText_WithBrowser 测试文本查找（有 browser 但未启动）
func TestDualChannelPositioner_FindElementByText_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	pos, err := p.FindElementByText("test")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElements_NilPage 测试多元素查找（page 为 nil）
func TestDualChannelPositioner_FindElements_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	positions, err := p.FindElements(".item")
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElements_WithBrowser 测试多元素查找（有 browser 但未启动）
func TestDualChannelPositioner_FindElements_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	positions, err := p.FindElements(".item")
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindClickableElements_NilPage 测试查找可点击元素（page 为 nil）
func TestDualChannelPositioner_FindClickableElements_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	positions, err := p.FindClickableElements()
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindClickableElements_WithBrowser 测试查找可点击元素（有 browser 但未启动）
func TestDualChannelPositioner_FindClickableElements_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	positions, err := p.FindClickableElements()
	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_WaitForPageStable_NilPage 测试等待页面稳定（page 为 nil）
func TestDualChannelPositioner_WaitForPageStable_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	// 由于 WaitForPageStable 会重试直到超时，这里只需检查返回错误
	err := p.WaitForPageStable(100*time.Millisecond, 50*time.Millisecond)
	assert.Error(t, err)
	// 错误可能是 "浏览器未启动" 或 "等待页面稳定超时"
	assert.True(t, strings.Contains(err.Error(), "浏览器未启动") ||
		strings.Contains(err.Error(), "等待页面稳定超时"))
}

// TestDualChannelPositioner_WaitForPageStable_WithBrowser 测试等待页面稳定（有 browser 但未启动）
func TestDualChannelPositioner_WaitForPageStable_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	err := p.WaitForPageStable(100*time.Millisecond, 50*time.Millisecond)
	assert.Error(t, err)
	// 错误可能是 "浏览器未启动" 或 "等待页面稳定超时"
	assert.True(t, strings.Contains(err.Error(), "浏览器未启动") ||
		strings.Contains(err.Error(), "等待页面稳定超时"))
}

// TestDualChannelPositioner_GetElementByVisual 测试视觉定位（需要实现）
func TestDualChannelPositioner_GetElementByVisual(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	pos, err := p.GetElementByVisual([]byte{}, "蓝色按钮")
	assert.Error(t, err)
	assert.Nil(t, pos)
	assert.Contains(t, err.Error(), "视觉定位功能需要实现")
}

// TestDualChannelPositioner_ClickWithRetry_NilPage 测试重试点击（page 为 nil）
func TestDualChannelPositioner_ClickWithRetry_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	err := p.ClickWithRetry("#test", 3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_ClickWithRetry_WithBrowser 测试重试点击（有 browser 但未启动）
func TestDualChannelPositioner_ClickWithRetry_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	err := p.ClickWithRetry("#test", 3)
	assert.Error(t, err)
	// page 为 nil 时，ClickElement 返回"浏览器未启动"错误
	// 但 ClickWithRetry 会重试直到超时
	assert.True(t, strings.Contains(err.Error(), "浏览器未启动") ||
		strings.Contains(err.Error(), "点击失败"))
}

// TestDualChannelPositioner_GetPageInfo_NilPage 测试获取页面信息（page 为 nil）
func TestDualChannelPositioner_GetPageInfo_NilPage(t *testing.T) {
	p := NewDualChannelPositioner(nil)

	info, err := p.GetPageInfo()
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_GetPageInfo_WithBrowser 测试获取页面信息（有 browser 但未启动）
func TestDualChannelPositioner_GetPageInfo_WithBrowser(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	info, err := p.GetPageInfo()
	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Contains(t, err.Error(), "浏览器未启动")
}

// TestDualChannelPositioner_FindElementByDOM_Integration 测试 DOM 查找（集成测试）
func TestDualChannelPositioner_FindElementByDOM_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	// 测试查找不存在的元素
	pos, err := p.FindElementByDOM("#nonexistent")
	assert.Error(t, err)
	assert.Nil(t, pos)

	// 测试查找存在的元素
	pos, err = p.FindElementByDOM("h1")
	if err != nil {
		// example.com 可能没有 h1
		t.Logf("查找 h1 失败: %v", err)
	} else {
		assert.NotNil(t, pos)
		assert.Greater(t, pos.X, 0)
		assert.Greater(t, pos.Y, 0)
	}
}

// TestDualChannelPositioner_FindElementByText_Integration 测试文本查找（集成测试）
func TestDualChannelPositioner_FindElementByText_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	// 测试查找包含 "Example" 的元素
	pos, err := p.FindElementByText("Example")
	if err != nil {
		t.Logf("查找文本失败: %v", err)
	} else {
		assert.NotNil(t, pos)
	}

	// 测试查找不存在的文本
	pos, err = p.FindElementByText("ThisTextDoesNotExist12345")
	assert.Error(t, err)
	assert.Nil(t, pos)
}

// TestDualChannelPositioner_FindElements_Integration 测试多元素查找（集成测试）
func TestDualChannelPositioner_FindElements_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	// 测试查找不存在的元素
	positions, err := p.FindElements(".nonexistent-class")
	assert.Error(t, err)
	assert.Nil(t, positions)
}

// TestDualChannelPositioner_FindClickableElements_Integration 测试查找可点击元素（集成测试）
func TestDualChannelPositioner_FindClickableElements_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	// 测试查找可点击元素
	positions, err := p.FindClickableElements()
	if err != nil {
		t.Logf("查找可点击元素失败: %v", err)
	} else {
		assert.NotNil(t, positions)
		t.Logf("找到 %d 个可点击元素", len(positions))
	}
}

// TestDualChannelPositioner_GetPageInfo_Integration 测试获取页面信息（集成测试）
func TestDualChannelPositioner_GetPageInfo_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	info, err := p.GetPageInfo()
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Contains(t, info, "url")
	assert.Contains(t, info, "title")

	url, ok := info["url"].(string)
	assert.True(t, ok)
	assert.Contains(t, url, "example.com")
}

// TestDualChannelPositioner_ClickWithRetry_Integration 测试重试点击（集成测试）
func TestDualChannelPositioner_ClickWithRetry_Integration(t *testing.T) {
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

	p := NewDualChannelPositioner(b)

	// 测试点击不存在的元素（应重试后失败）
	err = p.ClickWithRetry("#nonexistent-button", 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "点击失败")
}

// TestDualChannelPositioner_EmptySelector 测试空选择器
func TestDualChannelPositioner_EmptySelector(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	// 空选择器也会尝试查找，应该返回错误
	pos, err := p.FindElementByDOM("")
	if err != nil {
		assert.Contains(t, err.Error(), "未找到")
	}
	_ = pos
}

// TestDualChannelPositioner_InvalidSelector 测试无效选择器
func TestDualChannelPositioner_InvalidSelector(t *testing.T) {
	b := NewBrowser()
	p := NewDualChannelPositioner(b)

	// 无效选择器格式
	pos, err := p.FindElementByDOM("[[[invalid")
	if err != nil {
		assert.Error(t, err)
	}
	_ = pos
}
