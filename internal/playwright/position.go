package playwright

import (
	"fmt"
	"time"

	"github.com/yahao333/get_jobs/internal/config"
)

// Position 元素位置
type Position struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Text   string `json:"text"`
	Tag    string `json:"tag"`
}

// DualChannelPositioner 双通道定位器
// 结合 AI 视觉分析和 Playwright DOM 分析
type DualChannelPositioner struct {
	browser *Browser
}

// NewDualChannelPositioner 创建双通道定位器
func NewDualChannelPositioner(browser *Browser) *DualChannelPositioner {
	return &DualChannelPositioner{
		browser: browser,
	}
}

// FindElementByAI 使用 AI 查找元素位置
// 注意：这里返回模拟数据，实际需要调用 AI 服务
func (d *DualChannelPositioner) FindElementByAI(screenshot []byte, description string) (*Position, error) {
	config.Debug("使用 AI 分析查找元素: ", description)

	// TODO: 调用 AI 服务分析截图
	// 实际实现需要：
	// 1. 将截图发送给 AI 模型
	// 2. AI 返回元素位置和描述

	// 返回示例数据
	return nil, fmt.Errorf("AI 定位功能需要实现 AI 服务集成")
}

// FindElementByDOM 使用 Playwright DOM 查找元素
func (d *DualChannelPositioner) FindElementByDOM(selector string) (*Position, error) {
	if d.browser == nil || d.browser.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}

	page := *d.browser.page
	locator := page.Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("未找到元素: %s", selector)
	}

	first := locator.First()
	box, err := first.BoundingBox()
	if err != nil {
		return nil, err
	}

	// 获取元素文本
	text, err := first.TextContent()
	if err != nil {
		return nil, fmt.Errorf("获取元素文本失败: %w", err)
	}

	return &Position{
		X:      int(box.X + box.Width/2),
		Y:      int(box.Y + box.Height/2),
		Width:  int(box.Width),
		Height: int(box.Height),
		Text:   text,
		Tag:    selector,
	}, nil
}

// FindElementByText 使用文本内容查找元素
func (d *DualChannelPositioner) FindElementByText(text string) (*Position, error) {
	if d.browser == nil || d.browser.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}

	page := *d.browser.page

	// 尝试多种选择器方式查找包含文本的元素
	selectors := []string{
		fmt.Sprintf(`button:has-text("%s")`, text),
		fmt.Sprintf(`a:has-text("%s")`, text),
		fmt.Sprintf(`text=%s`, text),
		fmt.Sprintf(`//button[contains(text(), "%s")]`, text),
		fmt.Sprintf(`//a[contains(text(), "%s")]`, text),
	}

	for _, selector := range selectors {
		locator := page.Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			box, err := locator.First().BoundingBox()
			if err != nil {
				continue
			}

			return &Position{
				X:      int(box.X + box.Width/2),
				Y:      int(box.Y + box.Height/2),
				Width:  int(box.Width),
				Height: int(box.Height),
				Text:   text,
				Tag:    selector,
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到包含文本的元素: %s", text)
}

// FindElements 查找多个元素
func (d *DualChannelPositioner) FindElements(selector string) ([]*Position, error) {
	if d.browser == nil || d.browser.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}

	page := *d.browser.page
	locator := page.Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("未找到元素: %s", selector)
	}

	positions := make([]*Position, 0, count)
	for i := 0; i < count; i++ {
		box, err := locator.Nth(i).BoundingBox()
		if err != nil {
			continue
		}

		text, err := locator.Nth(i).TextContent()
		if err != nil {
			text = ""
		}

		positions = append(positions, &Position{
			X:      int(box.X + box.Width/2),
			Y:      int(box.Y + box.Height/2),
			Width:  int(box.Width),
			Height: int(box.Height),
			Text:   text,
			Tag:    selector,
		})
	}

	return positions, nil
}

// FindClickableElements 查找所有可点击元素
func (d *DualChannelPositioner) FindClickableElements() ([]*Position, error) {
	if d.browser == nil || d.browser.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}

	page := *d.browser.page

	// 获取所有可点击元素
	script := `
		(function() {
			var elements = [];
			var clickable = document.querySelectorAll('a, button, [role="button"], input[type="submit"], input[type="button"]');
			for (var i = 0; i < clickable.length; i++) {
				var el = clickable[i];
				var rect = el.getBoundingClientRect();
				if (rect.width > 5 && rect.height > 5 && rect.top >= 0) {
					elements.push({
						tag: el.tagName.toLowerCase(),
						text: el.innerText ? el.innerText.trim().substring(0, 100) : (el.value || ''),
						x: rect.x + rect.width / 2,
						y: rect.y + rect.height / 2,
						width: rect.width,
						height: rect.height
					});
				}
			}
			return elements;
		})()
	`

	result, err := page.Evaluate(script)
	if err != nil {
		return nil, err
	}

	// 解析结果
	elements, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("解析元素失败")
	}

	positions := make([]*Position, 0, len(elements))
	for _, e := range elements {
		el, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		pos := &Position{
			X:      int(el["x"].(float64)),
			Y:      int(el["y"].(float64)),
			Width:  int(el["width"].(float64)),
			Height: int(el["height"].(float64)),
		}

		if text, ok := el["text"].(string); ok {
			pos.Text = text
		}
		if tag, ok := el["tag"].(string); ok {
			pos.Tag = tag
		}

		positions = append(positions, pos)
	}

	return positions, nil
}

// WaitForPageStable 等待页面稳定（懒加载完成）
func (d *DualChannelPositioner) WaitForPageStable(maxWait time.Duration, checkInterval time.Duration) error {
	config.Debug("等待页面稳定...")

	elapsed := time.Duration(0)
	var lastCount int

	for elapsed < maxWait {
		// 获取当前可点击元素数量
		positions, err := d.FindClickableElements()
		if err != nil {
			time.Sleep(checkInterval)
			elapsed += checkInterval
			continue
		}

		currentCount := len(positions)

		if currentCount == lastCount && currentCount > 0 {
			// 连续两次相同，认为页面稳定
			config.Debug("页面已稳定，元素数量: ", currentCount)
			return nil
		}

		lastCount = currentCount
		time.Sleep(checkInterval)
		elapsed += checkInterval
	}

	return fmt.Errorf("等待页面稳定超时")
}

// GetElementByVisual 使用视觉特征查找元素
// 结合 AI 视觉分析和 DOM 信息
func (d *DualChannelPositioner) GetElementByVisual(screenshot []byte, visualDescription string) (*Position, error) {
	config.Debug("使用视觉特征查找元素: ", visualDescription)

	// 双通道策略：
	// 1. 使用 AI 分析截图，识别目标元素
	// 2. 使用 Playwright 验证和精确定位

	// TODO: 实现 AI 视觉分析

	// 暂时返回错误，需要先实现 AI 服务
	return nil, fmt.Errorf("视觉定位功能需要实现 AI 服务")
}

// ClickWithRetry 带重试的点击
func (d *DualChannelPositioner) ClickWithRetry(selector string, maxRetries int) error {
	if d.browser == nil {
		return fmt.Errorf("浏览器未启动")
	}

	for i := 0; i < maxRetries; i++ {
		err := d.browser.ClickElement(selector)
		if err == nil {
			config.Debug("点击成功: ", selector)
			return nil
		}

		config.Warn("点击失败，尝试重试: ", i+1, "/", maxRetries)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("点击失败，已重试 %d 次", maxRetries)
}

// GetPageInfo 获取页面信息
func (d *DualChannelPositioner) GetPageInfo() (map[string]interface{}, error) {
	if d.browser == nil || d.browser.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}

	page := *d.browser.page

	url := page.URL()
	title, _ := page.Title()

	return map[string]interface{}{
		"url":   url,
		"title": title,
	}, nil
}
