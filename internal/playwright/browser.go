package playwright

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/yahao333/get_jobs/internal/config"
)

// Browser 浏览器控制器
type Browser struct {
	page          *playwright.Page
	browser       playwright.Browser
	context       playwright.BrowserContext
	launchOptions playwright.BrowserTypeLaunchOptions
}

// NewBrowser 创建浏览器控制器
func NewBrowser() *Browser {
	return &Browser{
		launchOptions: playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false), // 默认显示浏览器窗口
		},
	}
}

// Launch 启动浏览器
func (b *Browser) Launch() error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("启动 Playwright 失败: %w", err)
	}

	// 启动 Chromium
	chromium := pw.Chromium
	browser, err := chromium.Launch(b.launchOptions)
	if err != nil {
		return fmt.Errorf("启动浏览器失败: %w", err)
	}

	// 创建上下文
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	if err != nil {
		return fmt.Errorf("创建浏览器上下文失败: %w", err)
	}

	// 创建页面
	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("创建页面失败: %w", err)
	}

	b.page = &page
	b.browser = browser
	b.context = context

	config.Info("浏览器启动成功")
	return nil
}

// LaunchWithCookie 使用 Cookie 启动浏览器
func (b *Browser) LaunchWithCookie(cookieFile string) error {
	if err := b.Launch(); err != nil {
		return err
	}

	// TODO: 从文件加载 Cookie
	// cookies, err := loadCookies(cookieFile)
	// if err != nil {
	// 	return err
	// }
	// (*b.context).AddCookies(cookies)

	return nil
}

// Navigate 导航到 URL
func (b *Browser) Navigate(url string) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}
	_, err := (*b.page).Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("导航失败: %w", err)
	}
	config.Debug("已导航到: ", url)
	return nil
}

// GetPage 获取页面
func (b *Browser) GetPage() *playwright.Page {
	return b.page
}

// Screenshot 截图
func (b *Browser) Screenshot(options ...playwright.PageScreenshotOptions) ([]byte, error) {
	if b.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}
	return (*b.page).Screenshot(options...)
}

// ScreenshotToFile 截图并保存到文件
func (b *Browser) ScreenshotToFile(filePath string) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}
	_, err := (*b.page).Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(filePath),
	})
	return err
}

// Evaluate 执行 JavaScript
func (b *Browser) Evaluate(script string) (interface{}, error) {
	if b.page == nil {
		return nil, fmt.Errorf("浏览器未启动")
	}
	return (*b.page).Evaluate(script)
}

// GetElementPosition 获取元素位置
func (b *Browser) GetElementPosition(selector string) (x, y int, err error) {
	if b.page == nil {
		return 0, 0, fmt.Errorf("浏览器未启动")
	}

	locator := (*b.page).Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return 0, 0, err
	}

	if count == 0 {
		return 0, 0, fmt.Errorf("未找到元素: %s", selector)
	}

	box, err := locator.First().BoundingBox()
	if err != nil {
		return 0, 0, err
	}

	x = int(box.X + box.Width/2)
	y = int(box.Y + box.Height/2)
	return x, y, nil
}

// ClickElement 点击元素
func (b *Browser) ClickElement(selector string) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}

	locator := (*b.page).Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("未找到元素: %s", selector)
	}

	return locator.First().Click()
}

// FillElement 填充输入框
func (b *Browser) FillElement(selector, value string) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}

	locator := (*b.page).Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("未找到元素: %s", selector)
	}

	return locator.First().Fill(value)
}

// GetText 获取元素文本
func (b *Browser) GetText(selector string) (string, error) {
	if b.page == nil {
		return "", fmt.Errorf("浏览器未启动")
	}

	locator := (*b.page).Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("未找到元素: %s", selector)
	}

	return locator.First().TextContent()
}

// GetHTML 获取元素 HTML
func (b *Browser) GetHTML(selector string) (string, error) {
	if b.page == nil {
		return "", fmt.Errorf("浏览器未启动")
	}

	locator := (*b.page).Locator(selector)
	count, err := locator.Count()
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("未找到元素: %s", selector)
	}

	return locator.First().InnerHTML()
}

// WaitForSelector 等待元素出现
func (b *Browser) WaitForSelector(selector string, timeout ...time.Duration) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}

	t := 30 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	_, err := (*b.page).WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(t.Milliseconds())),
	})
	return err
}

// Scroll 滚动页面
func (b *Browser) Scroll(amount int) error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}

	_, err := (*b.page).Evaluate(fmt.Sprintf("window.scrollBy(0, %d)", amount))
	return err
}

// ScrollToBottom 滚动到页面底部
func (b *Browser) ScrollToBottom() error {
	if b.page == nil {
		return fmt.Errorf("浏览器未启动")
	}

	_, err := (*b.page).Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	return err
}

// GetCookies 获取 Cookie
func (b *Browser) GetCookies() ([]playwright.Cookie, error) {
	if b.context == nil {
		return nil, fmt.Errorf("浏览器上下文未创建")
	}
	return b.context.Cookies()
}

// SetCookies 设置 Cookie
func (b *Browser) SetCookies(cookies []playwright.Cookie) error {
	if b.context == nil {
		return fmt.Errorf("浏览器上下文未创建")
	}
	var optionalCookies []playwright.OptionalCookie
	for _, c := range cookies {
		name := c.Name
		value := c.Value
		domain := c.Domain
		path := c.Path
		expires := c.Expires
		httpOnly := c.HttpOnly
		secure := c.Secure
		sameSite := c.SameSite

		oc := playwright.OptionalCookie{
			Name:     name,
			Value:    value,
			Domain:   &domain,
			Path:     &path,
			Expires:  &expires,
			HttpOnly: &httpOnly,
			Secure:   &secure,
			SameSite: sameSite,
		}
		optionalCookies = append(optionalCookies, oc)
	}
	return b.context.AddCookies(optionalCookies)
}

// Close 关闭浏览器
func (b *Browser) Close() error {
	if b.browser != nil {
		b.browser.Close()
		config.Info("浏览器已关闭")
	}
	return nil
}

// GetPageTitle 获取页面标题
func (b *Browser) GetPageTitle() (string, error) {
	if b.page == nil {
		return "", fmt.Errorf("浏览器未启动")
	}
	return (*b.page).Title()
}

// GetCurrentURL 获取当前 URL
func (b *Browser) GetCurrentURL() (string, error) {
	if b.page == nil {
		return "", fmt.Errorf("浏览器未启动")
	}
	return (*b.page).URL(), nil
}

// AnalyzePage 分析页面，返回 DOM 信息
func (b *Browser) AnalyzePage(query string) (string, error) {
	if b.page == nil {
		return "", fmt.Errorf("浏览器未启动")
	}

	// 获取页面 HTML
	html, err := (*b.page).Content()
	if err != nil {
		return "", err
	}

	// 获取所有可点击元素的简化信息
	script := `
		(function() {
			var elements = [];
			var clickable = document.querySelectorAll('a, button, [role="button"], input[type="submit"], .btn, .button');
			for (var i = 0; i < Math.min(clickable.length, 50); i++) {
				var el = clickable[i];
				var rect = el.getBoundingClientRect();
				if (rect.width > 0 && rect.height > 0) {
					elements.push({
						tag: el.tagName.toLowerCase(),
						text: el.innerText ? el.innerText.substring(0, 50) : '',
						classes: el.className || '',
						x: rect.x + rect.width / 2,
						y: rect.y + rect.height / 2
					});
				}
			}
			return JSON.stringify(elements);
		})()
	`

	result, err := (*b.page).Evaluate(script)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("HTML长度: %d 字节\n可点击元素数量: %s\n查询: %s", len(html), result, query), nil
}
