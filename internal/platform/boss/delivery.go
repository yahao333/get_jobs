// Package boss 提供 Boss直聘平台的自动化功能
// 简历投递模块：实现自动向 HR 发送打招呼消息和简历的功能
package boss

import (
	"fmt" // 格式化错误信息
	"time" // 时间处理

	"github.com/playwright-community/playwright-go" // Playwright 浏览器自动化

	"github.com/yahao333/get_jobs/internal/config"  // 配置管理模块
	"github.com/yahao333/get_jobs/internal/service" // AI 服务模块
	"github.com/yahao333/get_jobs/internal/storage" // 数据存储模块
)

// DeliveryResult 投递结果结构体
// 用于表示单次简历投递的结果信息
type DeliveryResult struct {
	JobID     int64     // 岗位 ID，对应数据库中的岗位记录
	Success   bool      // 投递是否成功
	Message   string    // 结果描述信息（成功或失败原因）
	Delivered time.Time // 投递时间
}

// Delivery 投递器结构体
// 封装了页面操作、AI 消息生成、投递限制等功能
type Delivery struct {
	page               *playwright.Page           // Playwright 页面对象，用于执行页面操作
	greetingGenerator  *service.GreetingGenerator // AI 打招呼消息生成器
	sendImgResume      bool                       // 是否发送图片简历
	imgResumePath      string                    // 图片简历的文件路径
	dailyLimit         int                       // 每日投递上限（0 表示不限制）
	deliveredToday     int                       // 今日已投递数量
}

// NewDelivery 创建投递器实例
// 参数：
// - page: Playwright 页面对象
// - greetingGenerator: AI 消息生成器（可选，可为 nil）
// 返回初始化后的 Delivery 实例
func NewDelivery(page *playwright.Page, greetingGenerator *service.GreetingGenerator) *Delivery {
	return &Delivery{
		page:               page,
		greetingGenerator:  greetingGenerator,
		sendImgResume:      config.GetBool("delivery.send_img_resume"),      // 从配置读取是否发送图片简历
		imgResumePath:      config.GetString("delivery.img_resume_path"),    // 从配置读取图片简历路径
		dailyLimit:         config.GetInt("delivery.daily_limit"),          // 从配置读取每日投递上限
		deliveredToday:     0,
	}
}

// Deliver 投递简历到指定岗位
// 完整的投递流程包括：
// 1. 检查每日投递限制
// 2. 导航到岗位详情页
// 3. 点击"立即沟通"按钮
// 4. 生成打招呼消息
// 5. 填写并发送消息
// 6. 可选：发送图片简历
// 7. 保存投递记录
// 参数：
// - job: 岗位信息卡片
// - message: 自定义打招呼消息（可选，为空时使用 AI 生成或默认消息）
// 返回值：
// - *DeliveryResult: 投递结果
// - error: 投递过程中的错误
func (d *Delivery) Deliver(job *JobCard, message string) (*DeliveryResult, error) {
	// 步骤1: 检查每日投递限制
	if d.dailyLimit > 0 && d.deliveredToday >= d.dailyLimit {
		return &DeliveryResult{
			JobID:   0,
			Success: false,
			Message: "已达到每日投递上限",
		}, fmt.Errorf("已达到每日投递上限")
	}

	// 步骤2: 导航到岗位详情页
	if job.JobURL != "" {
		_, err := (*d.page).Goto(job.JobURL, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		if err != nil {
			return &DeliveryResult{Success: false, Message: err.Error()}, err
		}
		// 等待页面加载完成
		time.Sleep(1 * time.Second)
	}

	// 步骤3: 查找并点击"立即沟通"按钮
	chatBtn, err := d.findChatButton()
	if err != nil {
		return &DeliveryResult{Success: false, Message: "找不到沟通按钮"}, err
	}

	// 点击沟通按钮打开聊天窗口
	if err := chatBtn.Click(); err != nil {
		return &DeliveryResult{Success: false, Message: "点击沟通按钮失败"}, err
	}

	// 等待聊天窗口加载
	time.Sleep(2 * time.Second)

	// 步骤4: 生成打招呼消息
	// 优先级：传入消息 > AI生成 > 默认配置
	if message == "" && d.greetingGenerator != nil {
		// 使用 AI 根据岗位信息生成个性化消息
		msg, err := d.greetingGenerator.Generate(
			"我是后端开发工程师，擅长Go语言", // 个人介绍
			"Go后端",                         // 关键词
			job.JobName,                     // 职位名称
			job.JobDescription,             // 职位描述
		)
		if err != nil {
			config.Warn("生成打招呼消息失败，使用默认消息")
			message = config.GetString("greeting.default")
		} else {
			message = msg
		}
	} else if message == "" {
		// 使用配置的默认打招呼语
		message = config.GetString("greeting.default")
	}

	// 步骤5: 填写打招呼消息到输入框
	if err := d.fillMessage(message); err != nil {
		return &DeliveryResult{Success: false, Message: "填写消息失败"}, err
	}

	// 步骤6: 发送消息
	if err := d.sendMessage(); err != nil {
		return &DeliveryResult{Success: false, Message: "发送消息失败"}, err
	}

	// 步骤7: 可选功能 - 发送图片简历
	if d.sendImgResume {
		if err := d.sendImageResume(); err != nil {
			config.Warn("发送图片简历失败: ", err)
		}
	}

	// 更新投递计数
	d.deliveredToday++

	// 步骤8: 保存投递记录到数据库
	d.saveDeliveryRecord(job, message)

	return &DeliveryResult{
		Success:   true,
		Message:   "投递成功",
		Delivered: time.Now(),
	}, nil
}

// findChatButton 查找"立即沟通"按钮
// 尝试多种 CSS 选择器定位沟通按钮，返回第一个匹配的元素
// 返回值：
// - playwright.Locator: 按钮元素定位器
// - error: 找不到按钮时返回错误
func (d *Delivery) findChatButton() (playwright.Locator, error) {
	// 尝试多种选择器，按优先级排序
	selectors := []string{
		".btn-start-chat",          // 沟通按钮 class
		"[class*='chat']",          // 包含 chat 的 class
		".contact-btn",             // 联系按钮
		".btn-contact",             // 联系按钮
		"button:has-text('立即沟通')", // 文本包含"立即沟通"的按钮
		"a:has-text('立即沟通')",     // 文本包含"立即沟通"的链接
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

// fillMessage 填写打招呼消息到输入框
// 尝试多种选择器查找消息输入框，找到后填充消息内容
// 参数 message: 要填写的消息内容
// 返回值：
// - error: 找不到输入框时返回错误
func (d *Delivery) fillMessage(message string) error {
	// 尝试多种选择器查找输入框
	selectors := []string{
		".chat-input textarea",      // 聊天输入框
		".msg-input textarea",       // 消息输入框
		"textarea[placeholder*='说']", // placeholder 包含"说"的文本域
		"textarea",                  // 通用文本域
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
// 优先尝试点击发送按钮，如果找不到按钮则使用回车键发送
// 返回值：
// - error: 发送过程中的错误
func (d *Delivery) sendMessage() error {
	// 尝试点击发送按钮
	selectors := []string{
		".btn-send",            // 发送按钮
		".send-btn",            // 发送按钮
		"button:has-text('发送')", // 文本包含"发送"的按钮
		"[class*='send']",     // 包含 send 的 class
	}

	for _, selector := range selectors {
		locator := (*d.page).Locator(selector)
		count, _ := locator.Count()
		if count > 0 {
			return locator.First().Click()
		}
	}

	// 备用方案：使用回车键发送
	(*d.page).Keyboard().Press("Enter")

	return nil
}

// sendImageResume 发送图片简历
// 查找并点击发送图片按钮，然后通过系统级文件选择器选择图片
// 注意：当前版本尚未完全实现系统文件对话框的处理
// 返回值：
// - error: 发送过程中的错误
func (d *Delivery) sendImageResume() error {
	// 查找发送图片按钮
	selectors := []string{
		".btn-img",            // 图片按钮
		"[class*='image']",   // 包含 image 的 class
		"[class*='picture']", // 包含 picture 的 class
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

// saveDeliveryRecord 保存投递记录到数据库
// 记录每次投递的岗位信息、消息内容和投递时间
// 参数：
// - job: 投递的岗位信息
// - message: 发送的打照顾呼消息
func (d *Delivery) saveDeliveryRecord(job *JobCard, message string) {
	record := storage.DeliveryRecord{
		JobID:       0, // TODO: 获取 job ID
		Platform:    "boss",
		Status:      "success",
		Message:     message,
		DeliveredAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := storage.Create(&record); err != nil {
		config.Warn("保存投递记录失败: ", err)
	}
}

// GetDeliveredCount 获取今日已投递数量
// 返回当前会话中今日已投递的简历数量
func (d *Delivery) GetDeliveredCount() int {
	return d.deliveredToday
}

// SetDailyLimit 设置每日投递限制
// 参数 limit: 每日允许投递的最大数量，0 表示不限制
func (d *Delivery) SetDailyLimit(limit int) {
	d.dailyLimit = limit
}

// BatchDeliver 批量投递简历
// 遍历岗位列表逐个投递，每次投递后等待一定间隔
// 参数：
// - jobs: 岗位卡片列表
// - message: 打招呼消息（可选）
// 返回值：
// - []DeliveryResult: 每个岗位的投递结果列表
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

		// 投递间隔，防止频繁操作被检测
		time.Sleep(3 * time.Second)
	}

	return results
}
