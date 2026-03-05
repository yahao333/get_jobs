package boss

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/loks666/get_jobs/internal/config"
	"github.com/loks666/get_jobs/internal/service"
	"github.com/loks666/get_jobs/internal/storage"
)

// DeliveryResult 投递结果
type DeliveryResult struct {
	JobID     int64
	Success   bool
	Message   string
	Delivered time.Time
}

// Delivery 投递器
type Delivery struct {
	page              *playwright.Page
	greetingGenerator *service.GreetingGenerator
	sendImgResume    bool
	imgResumePath    string
	dailyLimit       int
	deliveredToday   int
}

// NewDelivery 创建投递器
func NewDelivery(page *playwright.Page, greetingGenerator *service.GreetingGenerator) *Delivery {
	return &Delivery{
		page:              page,
		greetingGenerator: greetingGenerator,
		sendImgResume:    config.GetBool("delivery.send_img_resume"),
		imgResumePath:    config.GetString("delivery.img_resume_path"),
		dailyLimit:       config.GetInt("delivery.daily_limit"),
		deliveredToday:   0,
	}
}

// Deliver 投递简历
func (d *Delivery) Deliver(job *JobCard, message string) (*DeliveryResult, error) {
	// 检查每日投递限制
	if d.dailyLimit > 0 && d.deliveredToday >= d.dailyLimit {
		return &DeliveryResult{
			JobID:   0,
			Success: false,
			Message: "已达到每日投递上限",
		}, fmt.Errorf("已达到每日投递上限")
	}

	// 导航到岗位详情页
	if job.JobURL != "" {
		_, err := (*d.page).Goto(job.JobURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		if err != nil {
			return &DeliveryResult{Success: false, Message: err.Error()}, err
		}
		time.Sleep(1 * time.Second)
	}

	// 查找"立即沟通"按钮
	chatBtn, err := d.findChatButton()
	if err != nil {
		return &DeliveryResult{Success: false, Message: "找不到沟通按钮"}, err
	}

	// 点击沟通按钮
	if err := chatBtn.Click(); err != nil {
		return &DeliveryResult{Success: false, Message: "点击沟通按钮失败"}, err
	}

	time.Sleep(2 * time.Second)

	// 生成打招呼消息
	if message == "" && d.greetingGenerator != nil {
		msg, err := d.greetingGenerator.Generate(
			"我是后端开发工程师，擅长Go语言",
			"Go后端",
			job.JobName,
			job.JobDescription,
		)
		if err != nil {
			config.Warn("生成打招呼消息失败，使用默认消息")
			message = config.GetString("greeting.default")
		} else {
			message = msg
		}
	} else if message == "" {
		message = config.GetString("greeting.default")
	}

	// 填写打招呼消息
	if err := d.fillMessage(message); err != nil {
		return &DeliveryResult{Success: false, Message: "填写消息失败"}, err
	}

	// 发送消息
	if err := d.sendMessage(); err != nil {
		return &DeliveryResult{Success: false, Message: "发送消息失败"}, err
	}

	// 可选：发送图片简历
	if d.sendImgResume {
		if err := d.sendImageResume(); err != nil {
			config.Warn("发送图片简历失败: ", err)
		}
	}

	// 更新投递计数
	d.deliveredToday++

	// 保存投递记录
	d.saveDeliveryRecord(job, message)

	return &DeliveryResult{
		Success:   true,
		Message:   "投递成功",
		Delivered: time.Now(),
	}, nil
}

// findChatButton 查找沟通按钮
func (d *Delivery) findChatButton() (playwright.Locator, error) {
	// 尝试多种选择器
	selectors := []string{
		".btn-start-chat",
		"[class*='chat']",
		".contact-btn",
		".btn-contact",
		"button:has-text('立即沟通')",
		"a:has-text('立即沟通')",
	}

	for _, selector := range selectors {
		locator := (*d.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First(), nil
		}
	}

	// 使用 AI 视觉分析查找
	config.Warn("使用 DOM 选择器未找到沟通按钮，尝试 AI 视觉分析")
	return nil, fmt.Errorf("找不到沟通按钮")
}

// fillMessage 填写消息
func (d *Delivery) fillMessage(message string) error {
	// 尝试多种选择器
	selectors := []string{
		".chat-input textarea",
		".msg-input textarea",
		"textarea[placeholder*='说']",
		"textarea",
	}

	for _, selector := range selectors {
		locator := (*d.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First().Fill(message)
		}
	}

	return fmt.Errorf("找不到消息输入框")
}

// sendMessage 发送消息
func (d *Delivery) sendMessage() error {
	// 尝试点击发送按钮
	selectors := []string{
		".btn-send",
		".send-btn",
		"button:has-text('发送')",
		"[class*='send']",
	}

	for _, selector := range selectors {
		locator := (*d.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First().Click()
		}
	}

	// 尝试使用回车发送
	(*d.page).Keyboard().Press("Enter")

	return nil
}

// sendImageResume 发送图片简历
func (d *Delivery) sendImageResume() error {
	// 查找发送图片按钮
	selectors := []string{
		".btn-img",
		"[class*='image']",
		"[class*='picture']",
	}

	var imgBtn playwright.Locator
	for _, selector := range selectors {
		locator := (*d.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			imgBtn = locator.First()
			break
		}
	}

	if imgBtn == nil {
		return fmt.Errorf("找不到发送图片按钮")
	}

	// 点击发送图片按钮
	if err := imgBtn.Click(); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	// 文件选择对话框通常需要使用操作系统级别的文件选择
	// 由于 Playwright 无法直接处理系统文件对话框
	// 这里需要使用 robotgo 来处理
	// 暂时返回错误，需要后续集成
	config.Warn("图片简历发送需要集成系统级文件选择")

	return nil
}

// saveDeliveryRecord 保存投递记录
func (d *Delivery) saveDeliveryRecord(job *JobCard, message string) {
	record := storage.DeliveryRecord{
		JobID:       0, // TODO: 获取 job ID
		Platform:    "boss",
		Status:      "success",
		Message:     message,
		DeliveredAt: time.Now().Format("2006-01-02 15:04:05"),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := storage.Create(&record); err != nil {
		config.Warn("保存投递记录失败: ", err)
	}
}

// GetDeliveredCount 获取今日投递数量
func (d *Delivery) GetDeliveredCount() int {
	return d.deliveredToday
}

// SetDailyLimit 设置每日投递限制
func (d *Delivery) SetDailyLimit(limit int) {
	d.dailyLimit = limit
}

// BatchDeliver 批量投递
func (d *Delivery) BatchDeliver(jobs []JobCard, message string) []DeliveryResult {
	results := make([]DeliveryResult, 0, len(jobs))

	for i, job := range jobs {
		config.Info("投递岗位: ", i+1, "/", len(jobs), " - ", job.JobName)

		result, err := d.Deliver(&job, message)
		result.JobID = int64(i)
		if err != nil {
			config.Error("投递失败: ", err)
		}

		results = append(results, *result)

		// 投递间隔
		time.Sleep(3 * time.Second)
	}

	return results
}
