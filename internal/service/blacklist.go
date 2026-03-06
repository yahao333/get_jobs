// Package service 提供业务逻辑服务
// 黑名单自动更新模块：根据聊天记录自动识别并添加黑名单
package service

import (
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/storage"
)

// BlacklistUpdater 黑名单自动更新器
type BlacklistUpdater struct {
	page     *playwright.Page
	keywords []string
}

// NewBlacklistUpdater 创建黑名单更新器
func NewBlacklistUpdater(page *playwright.Page) *BlacklistUpdater {
	keywords := config.GetStringSlice("blacklist.keywords")
	if len(keywords) == 0 {
		keywords = []string{"不", "感谢", "但", "遗憾", "抱歉", "不合适"}
	}

	return &BlacklistUpdater{
		page:     page,
		keywords: keywords,
	}
}

// UpdateFromChatHistory 从聊天记录更新黑名单
// 访问聊天记录页面，分析历史消息，将匹配关键词的公司或HR加入黑名单
func (b *BlacklistUpdater) UpdateFromChatHistory() (int, error) {
	if b.page == nil {
		return 0, nil
	}

	config.Info("开始从聊天记录更新黑名单...")

	// 导航到聊天记录页面
	// Boss直聘的聊天记录页面
	chatURL := "https://www.zhipin.com/web/geek/chat/"
	_, err := (*b.page).Goto(chatURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return 0, err
	}

	time.Sleep(2 * time.Second)

	// 获取聊天列表
	count, _ := (*b.page).Locator(".chat-item").Count()
	if count == 0 {
		config.Info("没有找到聊天记录")
		return 0, nil
	}

	addedCount := 0

	// 遍历聊天记录
	for i := 0; i < min(count, 50); i++ {
		// 点击聊天项
		chatItem := (*b.page).Locator(".chat-item").Nth(i)
		if err := chatItem.Click(); err != nil {
			continue
		}

		time.Sleep(1 * time.Second)

		// 获取聊天内容
		messages, _ := (*b.page).Locator(".message-content").AllTextContents()

		// 检查是否包含拒绝关键词
		isRejected := b.checkRejectKeywords(messages)

		if isRejected {
			// 获取公司名称
			companyName, _ := (*b.page).Locator(".company-name").First().TextContent()

			// 添加到黑名单
			if companyName != "" {
				blacklist := storage.Blacklist{
					Keyword:   companyName,
					Type:      "company",
					Source:    "auto",
					CreatedAt: time.Now(),
				}

				if err := storage.Create(&blacklist); err == nil {
					addedCount++
					config.Info("已添加公司到黑名单: ", companyName)
				}
			}
		}

		// 返回列表
		(*b.page).Keyboard().Press("Escape")
		time.Sleep(500 * time.Millisecond)
	}

	config.Info("黑名单更新完成，新增 ", addedCount, " 条记录")
	return addedCount, nil
}

// checkRejectKeywords 检查消息是否包含拒绝关键词
func (b *BlacklistUpdater) checkRejectKeywords(messages []string) bool {
	fullText := strings.Join(messages, "")

	for _, keyword := range b.keywords {
		if strings.Contains(fullText, keyword) {
			return true
		}
	}
	return false
}

// ManualAdd 手动添加黑名单
func (b *BlacklistUpdater) ManualAdd(keyword, btype string) error {
	blacklist := storage.Blacklist{
		Keyword:   keyword,
		Type:      btype,
		Source:    "manual",
		CreatedAt: time.Now(),
	}
	return storage.Create(&blacklist)
}

// GetBlacklist 获取黑名单列表
func (b *BlacklistUpdater) GetBlacklist(btype string) ([]storage.Blacklist, error) {
	var blacklists []storage.Blacklist
	err := storage.Where(&blacklists, "type = ?", btype)
	return blacklists, err
}

// DeleteBlacklist 删除黑名单
func (b *BlacklistUpdater) DeleteBlacklist(id int64) error {
	return storage.Delete(&storage.Blacklist{}, "id = ?", id)
}

// min 返回两个整数中较小的值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
