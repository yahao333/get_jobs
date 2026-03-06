// Package boss 提供 Boss直聘平台的自动化功能
// 包括登录认证、岗位搜索、简历投递等核心功能
package boss

import (
	"encoding/json" // JSON 序列化与反序列化
	"fmt"          // 格式化错误信息
	"os"           // 文件系统操作
	"time"         // 时间处理

	"github.com/playwright-community/playwright-go"          // Playwright 浏览器自动化
	"github.com/yahao333/get_jobs/internal/config"           // 配置管理模块
	"github.com/yahao333/get_jobs/internal/storage"          // 数据存储模块
)

// BossClient Boss直聘客户端结构体
// 封装了浏览器实例、页面对象和登录状态，用于管理与 Boss直聘网站的交互
type BossClient struct {
	browser    *playwright.Browser // Playwright 浏览器实例，用于控制浏览器
	page       *playwright.Page   // 当前页面对象，用于执行页面操作
	cookieFile string             // Cookie 存储文件路径，用于持久化登录状态
	isLoggedIn bool               // 登录状态标志，true 表示已登录
}

// NewBossClient 创建 Boss 直聘客户端实例
// 参数 cookieFile 指定 Cookie 持久化存储的文件路径
// 返回初始化后的 BossClient 实例
func NewBossClient(cookieFile string) *BossClient {
	return &BossClient{
		cookieFile: cookieFile,
		isLoggedIn: false,
	}
}

// LoginStatus 登录状态结构体
// 用于表示用户的登录状态和基本信息
type LoginStatus struct {
	IsLoggedIn bool   `json:"is_logged_in"` // 是否已登录
	UserID     string `json:"user_id"`     // 用户 ID
	UserName   string `json:"user_name"`   // 用户名称
	ExpireAt   string `json:"expire_at"`   // 过期时间
}

// CheckLoginStatus 检查当前登录状态
// 流程：
// 1. 从文件加载已保存的 Cookie
// 2. 将 Cookie 添加到浏览器上下文
// 3. 访问首页检查是否显示用户名
// 4. 返回登录状态和用户信息
// 返回值：
// - LoginStatus: 包含登录状态和用户信息
// - error: 检查过程中的错误信息
func (b *BossClient) CheckLoginStatus() (*LoginStatus, error) {
	// 步骤1: 尝试从文件加载 Cookie
	cookies, err := b.loadCookies()
	if err != nil {
		config.Warn("加载 Cookie 失败: ", err)
	}

	// 步骤2: 如果存在 Cookie，添加到浏览器并验证有效性
	if len(cookies) > 0 {
		if b.page == nil {
			return &LoginStatus{IsLoggedIn: false}, nil
		}

		// 将存储的 Cookie 转换为 Playwright 需要的格式并添加
		optionalCookies := make([]playwright.OptionalCookie, 0)
		for _, c := range cookies {
			domain := c.Domain
			path := c.Path
			// 将 time.Time 转换为 Unix 时间戳（float64）
			expires := float64(c.Expires.Unix())
			httpOnly := c.HttpOnly
			secure := c.Secure
			// 将 SameSite 字符串转换为 Playwright 枚举类型
			sameSite := playwright.SameSiteAttribute(c.SameSite)
			optionalCookies = append(optionalCookies, playwright.OptionalCookie{
				Name:     c.Name,
				Value:    c.Value,
				Domain:   &domain,
				Path:     &path,
				Expires:  &expires,
				HttpOnly: &httpOnly,
				Secure:   &secure,
				SameSite: &sameSite,
			})
		}
		// 添加 Cookie 到浏览器上下文
		(*b.page).Context().AddCookies(optionalCookies)

		// 步骤3: 访问 Boss直聘首页检查登录状态
		if err := b.navigateToMainPage(); err != nil {
			return &LoginStatus{IsLoggedIn: false}, nil
		}

		// 步骤4: 检查页面是否显示用户名（已登录的标志）
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

// navigateToMainPage 导航到 Boss直聘首页
// 使用 Playwright 访问主页面，等待网络空闲后完成
// 返回错误信息（如果有）
func (b *BossClient) navigateToMainPage() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}
	_, err := (*b.page).Goto("https://www.zhipin.com/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	return err
}

// checkLoginElement 检查页面上的登录元素
// 通过查找用户名元素或登录按钮来判断当前登录状态
// 返回值：
// - bool: 是否已登录
// - string: 用户名（如果已登录）
func (b *BossClient) checkLoginElement() (bool, string) {
	// 尝试多种选择器查找已登录用户的用户名元素
	// Boss直聘登录后通常在页面头部显示用户名
	selectors := []string{
		".user-name",       // 用户名class
		".header-user-name", // 头部用户名
		"[class*='user']",  // 包含user的class
		".nick-name",       // 昵称
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

	// 检查是否存在登录按钮（表示未登录）
	loginSelectors := []string{
		".login-btn",      // 登录按钮
		"[class*='login']", // 包含login的class
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

// WaitForLogin 等待用户通过二维码扫码登录
// 此方法会阻塞直到用户成功登录或超时
// 参数 timeout 指定等待超时时间
// 登录成功后会自动保存 Cookie 到文件
func (b *BossClient) WaitForLogin(timeout time.Duration) error {
	config.Info("等待用户扫码登录...")
	config.Info("请在浏览器窗口中扫描二维码登录 Boss直聘")

	// 导航到登录页面显示二维码
	if err := b.navigateToLoginPage(); err != nil {
		return fmt.Errorf("导航到登录页面失败: %w", err)
	}

	startTime := time.Now()
	// 轮询检查登录状态
	for {
		// 检查是否超时
		if time.Since(startTime) > timeout {
			return fmt.Errorf("登录超时")
		}

		// 检查是否登录成功
		isLoggedIn, userName := b.checkLoginElement()
		if isLoggedIn {
			config.Info("登录成功! 用户: ", userName)
			b.isLoggedIn = true

			// 登录成功后保存 Cookie 以便下次免登录
			if err := b.saveCookies(); err != nil {
				config.Warn("保存 Cookie 失败: ", err)
			}
			return nil
		}

		// 每2秒检查一次
		time.Sleep(2 * time.Second)
	}
}

// navigateToLoginPage 导航到 Boss直聘登录页面
// 返回错误信息（如果有）
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
// 读取 JSON 格式的 Cookie 数据并反序列化为 Cookie 结构体切片
// 返回值：
// - []storage.Cookie: Cookie 列表
// - error: 读取或解析过程中的错误
func (b *BossClient) loadCookies() ([]storage.Cookie, error) {
	// 检查文件路径是否设置
	if b.cookieFile == "" {
		return nil, fmt.Errorf("Cookie 文件路径未设置")
	}

	// 读取文件内容
	data, err := os.ReadFile(b.cookieFile)
	if err != nil {
		return nil, fmt.Errorf("读取 Cookie 文件失败: %w", err)
	}

	// 解析 JSON 数据
	var cookies []storage.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("解析 Cookie 失败: %w", err)
	}

	return cookies, nil
}

// saveCookies 保存当前浏览器上下文中的 Cookie 到文件
// 将 Cookie 序列化为 JSON 格式并写入文件，实现登录状态持久化
// 注意：会创建必要的目录结构
func (b *BossClient) saveCookies() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}

	// 获取浏览器上下文中的所有 Cookie
	cookies, err := (*b.page).Context().Cookies()
	if err != nil {
		return fmt.Errorf("获取 Cookie 失败: %w", err)
	}

	// 转换为存储格式（从 Playwright Cookie 转换为我们自定义的 storage.Cookie）
	storageCookies := make([]storage.Cookie, len(cookies))
	for i, c := range cookies {
		var sameSite string
		// 处理 SameSite 可能为 nil 的情况
		if c.SameSite != nil {
			sameSite = string(*c.SameSite)
		}
		// 将 Unix 时间戳转换回 time.Time
		storageCookies[i] = storage.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  time.Unix(int64(c.Expires), 0),
			HttpOnly: c.HttpOnly,
			Secure:   c.Secure,
			SameSite: sameSite,
		}
	}

	// 序列化为格式化的 JSON
	data, err := json.MarshalIndent(storageCookies, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 Cookie 失败: %w", err)
	}

	// 确保目录存在（提取文件所在目录并创建）
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

	// 写入文件
	if err := os.WriteFile(b.cookieFile, data, 0644); err != nil {
		return fmt.Errorf("写入 Cookie 文件失败: %w", err)
	}

	config.Info("Cookie 已保存到: ", b.cookieFile)
	return nil
}

// SetBrowser 设置浏览器实例和页面对象
// 在创建 BossClient 后需要调用此方法绑定浏览器和页面
// 参数：
// - browser: Playwright 浏览器实例
// - page: Playwright 页面对象
func (b *BossClient) SetBrowser(browser *playwright.Browser, page *playwright.Page) {
	b.browser = browser
	b.page = page
}

// IsLoggedIn 检查当前客户端是否已登录
// 返回 true 表示已登录，false 表示未登录
func (b *BossClient) IsLoggedIn() bool {
	return b.isLoggedIn
}
