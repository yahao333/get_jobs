package boss

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/loks666/get_jobs/internal/config"
	"github.com/loks666/get_jobs/internal/storage"
	"github.com/playwright-community/playwright-go"
)

// BossClient Boss直聘客户端
type BossClient struct {
	browser     *playwright.Browser
	page        *playwright.Page
	cookieFile  string
	isLoggedIn  bool
}

// NewBossClient 创建 Boss 直聘客户端
func NewBossClient(cookieFile string) *BossClient {
	return &BossClient{
		cookieFile: cookieFile,
		isLoggedIn: false,
	}
}

// LoginStatus 登录状态
type LoginStatus struct {
	IsLoggedIn bool   `json:"is_logged_in"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	ExpireAt   string `json:"expire_at"`
}

// CheckLoginStatus 检查登录状态
func (b *BossClient) CheckLoginStatus() (*LoginStatus, error) {
	// 尝试加载 Cookie
	cookies, err := b.loadCookies()
	if err != nil {
		config.Warn("加载 Cookie 失败: ", err)
	}

	if len(cookies) > 0 {
		// 添加 Cookie 到浏览器
		if b.page != nil {
			optionalCookies := make([]playwright.OptionalCookie, 0)
			for _, c := range cookies {
				domain := c.Domain
				path := c.Path
				expires := c.Expires
				httpOnly := c.HttpOnly
				secure := c.Secure
				sameSite := c.SameSite
				optionalCookies = append(optionalCookies, playwright.OptionalCookie{
					Name:     c.Name,
					Value:    c.Value,
					Domain:   &domain,
					Path:     &path,
					Expires:  &expires,
					HttpOnly: &httpOnly,
					Secure:   &secure,
					SameSite: sameSite,
				})
			}
			(*b.page).Context().AddCookies(optionalCookies)
		}

		// 访问Boss直聘首页检查是否已登录
		if err := b.navigateToMainPage(); err != nil {
			return &LoginStatus{IsLoggedIn: false}, nil
		}

		// 检查是否显示用户名
		isLoggedIn, userName := b.checkLoginElement()
		if isLoggedIn {
			b.isLoggedIn = true
			return &LoginStatus{
				IsLoggedIn: true,
				UserName:   userName,
			}, nil
		}
	}

	return &LoginStatus{IsLoggedIn: false}, nil
}

// navigateToMainPage 导航到首页
func (b *BossClient) navigateToMainPage() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}
	_, err := (*b.page).Goto("https://www.zhipin.com/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	return err
}

// checkLoginElement 检查登录元素
func (b *BossClient) checkLoginElement() (bool, string) {
	// 尝试查找用户头像或用户名元素
	// Boss直聘登录后通常会显示用户名
	selectors := []string{
		".user-name",
		".header-user-name",
		"[class*='user']",
		".nick-name",
	}

	for _, selector := range selectors {
		locator := (*b.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			text, _ := locator.First().TextContent()
			if text != "" {
				return true, text
			}
		}
	}

	// 检查是否存在登录按钮（未登录）
	loginSelectors := []string{
		".login-btn",
		"[class*='login']",
	}

	for _, selector := range loginSelectors {
		locator := (*b.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return false, ""
		}
	}

	return false, ""
}

// WaitForLogin 等待用户扫码登录
func (b *BossClient) WaitForLogin(timeout time.Duration) error {
	config.Info("等待用户扫码登录...")
	config.Info("请在浏览器窗口中扫描二维码登录 Boss直聘")

	// 导航到登录页面
	if err := b.navigateToLoginPage(); err != nil {
		return fmt.Errorf("导航到登录页面失败: %w", err)
	}

	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("登录超时")
		}

		// 检查是否登录成功
		isLoggedIn, userName := b.checkLoginElement()
		if isLoggedIn {
			config.Info("登录成功! 用户: ", userName)
			b.isLoggedIn = true

			// 保存 Cookie
			if err := b.saveCookies(); err != nil {
				config.Warn("保存 Cookie 失败: ", err)
			}
			return nil
		}

		time.Sleep(2 * time.Second)
	}
}

// navigateToLoginPage 导航到登录页面
func (b *BossClient) navigateToLoginPage() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}
	_, err := (*b.page).Goto("https://www.zhipin.com/user/login.html", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	return err
}

// loadCookies 从文件加载 Cookie
func (b *BossClient) loadCookies() ([]storage.Cookie, error) {
	if b.cookieFile == "" {
		return nil, fmt.Errorf("Cookie 文件路径未设置")
	}

	data, err := os.ReadFile(b.cookieFile)
	if err != nil {
		return nil, fmt.Errorf("读取 Cookie 文件失败: %w", err)
	}

	var cookies []storage.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("解析 Cookie 失败: %w", err)
	}

	return cookies, nil
}

// saveCookies 保存 Cookie 到文件
func (b *BossClient) saveCookies() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}

	cookies, err := (*b.page).Context().Cookies()
	if err != nil {
		return fmt.Errorf("获取 Cookie 失败: %w", err)
	}

	// 转换为存储格式
	storageCookies := make([]storage.Cookie, len(cookies))
	for i, c := range cookies {
		storageCookies[i] = storage.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  c.Expires,
			HttpOnly: c.HttpOnly,
			Secure:   c.Secure,
			SameSite: c.SameSite,
		}
	}

	data, err := json.MarshalIndent(storageCookies, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 Cookie 失败: %w", err)
	}

	// 确保目录存在
	dir := b.cookieFile
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

	if err := os.WriteFile(b.cookieFile, data, 0644); err != nil {
		return fmt.Errorf("写入 Cookie 文件失败: %w", err)
	}

	config.Info("Cookie 已保存到: ", b.cookieFile)
	return nil
}

// SetBrowser 设置浏览器实例
func (b *BossClient) SetBrowser(browser *playwright.Browser, page *playwright.Page) {
	b.browser = browser
	b.page = page
}

// IsLoggedIn 检查是否已登录
func (b *BossClient) IsLoggedIn() bool {
	return b.isLoggedIn
}
