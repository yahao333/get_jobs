// Package boss 提供 Boss直聘平台的自动化功能
// 岗位搜索模块：实现岗位搜索、页面解析、数据保存等功能
package boss

import (
	"fmt"  // 格式化错误信息
	"time" // 时间处理

	"github.com/playwright-community/playwright-go" // Playwright 浏览器自动化

	"github.com/yahao333/get_jobs/internal/config"  // 配置管理模块
	"github.com/yahao333/get_jobs/internal/storage" // 数据存储模块
)

// JobCard 岗位卡片信息结构体
// 用于存储从 Boss直聘网站解析出的岗位基本信息
type JobCard struct {
	EncryptID      string // 岗位加密 ID，用于构建岗位详情页 URL
	EncryptUserID  string // HR 加密 ID，用于标识 HR 身份
	CompanyName    string // 公司名称
	JobName        string // 职位名称
	Salary         string // 薪资范围（如 "15K-30K"）
	Location       string // 工作地点
	Experience     string // 经验要求
	Degree         string // 学历要求
	HRName         string // HR 姓名
	HRPosition     string // HR 职位
	HRActiveStatus string // HR 活跃状态（如 "今日活跃"、"3天前在线"）
	JobDescription string // 职位描述（详细内容）
	JobURL         string // 岗位详情页 URL
}

// SearchJobs 搜索岗位
// 根据提供的搜索条件访问 Boss直聘搜索页面，获取符合条件的岗位列表
// 搜索流程：
// 1. 构建搜索 URL
// 2. 导航到搜索页面
// 3. 滚动页面触发懒加载，获取更多岗位
// 4. 解析页面中的岗位列表
// 参数：
// - cityCode: 城市编码（如 "101010100" 表示北京）
// - keyword: 搜索关键词（如 "Go后端"）
// - jobType: 职位类型
// - salary: 薪资筛选（可选）
// - experience: 经验要求（可选）
// - degree: 学历要求（可选）
// 返回值：
// - []JobCard: 岗位卡片列表
// - error: 搜索过程中的错误
func (b *BossClient) SearchJobs(cityCode, keyword, jobType, salary, experience, degree string) ([]JobCard, error) {
	// 检查页面是否初始化
	if b.page == nil {
		return nil, fmt.Errorf("浏览器页面未初始化")
	}

	// 步骤1: 构建搜索 URL
	// Boss直聘搜索页面 URL 格式
	searchURL := fmt.Sprintf(
		"https://www.zhipin.com/web/geek/job?query=%s&cityCode=%s&jobType=%s",
		keyword, cityCode, jobType,
	)

	config.Info("搜索岗位: ", searchURL)

	// 步骤2: 导航到搜索页面
	_, err := (*b.page).Goto(searchURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return nil, fmt.Errorf("导航到搜索页面失败: %w", err)
	}

	// 等待页面加载完成
	time.Sleep(2 * time.Second)

	// 步骤3: 滚动页面触发懒加载
	// Boss直聘使用无限滚动加载更多岗位，需要滚动页面才能获取完整列表
	if err := b.scrollToLoadAll(); err != nil {
		config.Warn("滚动加载失败: ", err)
	}

	// 步骤4: 解析岗位列表
	jobCards, err := b.parseJobList()
	if err != nil {
		return nil, fmt.Errorf("解析岗位列表失败: %w", err)
	}

	config.Info("找到 ", len(jobCards), " 个岗位")
	return jobCards, nil
}

// scrollToLoadAll 滚动加载所有岗位
// 通过模拟用户滚动行为触发懒加载，持续滚动直到没有新岗位加载为止
// 算法说明：
// - 最多滚动 10 次
// - 每次滚动后等待 1 秒让新岗位加载
// - 如果连续两次滚动后岗位数量不变，认为已加载完毕，停止滚动
// 返回值：
// - error: 滚动过程中的错误
func (b *BossClient) scrollToLoadAll() error {
	// 检查页面是否初始化
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}

	lastCount := 0
	// 最多滚动 10 次
	for i := 0; i < 10; i++ {
		// 滚动到页面底部（每次滚动 800 像素）
		_, err := (*b.page).Evaluate("window.scrollBy(0, 800)")
		if err != nil {
			return err
		}

		// 等待新岗位加载
		time.Sleep(1 * time.Second)

		// 检查是否还有新岗位
		newCount, _ := (*b.page).Locator(".job-card").Count()
		// 连续两次数量相同且大于 0，表示已加载完毕
		if newCount == lastCount && newCount > 0 {
			break
		}
		lastCount = newCount
	}

	return nil
}

// parseJobList 解析岗位列表
// 从当前页面中提取所有岗位卡片信息
// 返回值：
// - []JobCard: 岗位卡片列表
// - error: 解析过程中的错误
func (b *BossClient) parseJobList() ([]JobCard, error) {
	// 检查页面是否初始化
	if b.page == nil {
		return nil, fmt.Errorf("浏览器页面未初始化")
	}

	// 获取所有岗位卡片
	locator := (*b.page).Locator(".job-card")
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	// 如果找不到，尝试备用选择器
	if count == 0 {
		locator = (*b.page).Locator(".job-list-box .job-card")
		count, _ = locator.Count()
		if count == 0 {
			return nil, fmt.Errorf("未找到岗位卡片")
		}
	}

	// 逐个解析岗位卡片
	cards := make([]JobCard, 0, count)
	for i := 0; i < count; i++ {
		card, err := b.parseJobCard(locator.Nth(i))
		if err != nil {
			// 单个卡片解析失败继续解析下一个
			continue
		}
		cards = append(cards, *card)
	}

	return cards, nil
}

// parseJobCard 解析单个岗位卡片
// 从页面元素的 DOM 中提取岗位的各个字段信息
// 提取的字段包括：职位名称、公司名称、薪资、工作地点、经验学历要求、HR 姓名、岗位链接
// 参数：
// - locator: 岗位卡片的 Playwright 定位器
// 返回值：
// - *JobCard: 解析出的岗位信息
// - error: 解析过程中的错误
func (b *BossClient) parseJobCard(locator playwright.Locator) (*JobCard, error) {
	card := &JobCard{}

	// 获取职位名称
	jobNameEl := locator.Locator(".job-name")
	jobName, _ := jobNameEl.TextContent()
	card.JobName = jobName

	// 获取公司名称
	companyEl := locator.Locator(".company-name")
	companyName, _ := companyEl.TextContent()
	card.CompanyName = companyName

	// 获取薪资
	salaryEl := locator.Locator(".salary")
	salary, _ := salaryEl.TextContent()
	card.Salary = salary

	// 获取工作地点
	locationEl := locator.Locator(".job-area")
	location, _ := locationEl.TextContent()
	card.Location = location

	// 获取经验学历要求（通常在同一个标签列表中）
	tagEl := locator.Locator(".tag-list")
	tagText, _ := tagEl.TextContent()
	card.Experience = tagText

	// 获取 HR 姓名
	hrNameEl := locator.Locator(".boss-name")
	hrName, _ := hrNameEl.TextContent()
	card.HRName = hrName

	// 获取岗位详情页链接
	linkEl := locator.Locator("a")
	href, _ := linkEl.GetAttribute("href")
	if href != "" {
		// 拼接完整的 URL
		card.JobURL = "https://www.zhipin.com" + href
	}

	return card, nil
}

// GetJobDetail 获取岗位详情
// 访问岗位详情页，获取更详细的岗位信息，包括职位描述和 HR 活跃状态
// 参数：
// - jobURL: 岗位详情页 URL
// 返回值：
// - *JobCard: 包含详情的岗位信息
// - error: 获取过程中的错误
func (b *BossClient) GetJobDetail(jobURL string) (*JobCard, error) {
	// 检查页面是否初始化
	if b.page == nil {
		return nil, fmt.Errorf("浏览器页面未初始化")
	}

	// 导航到岗位详情页
	_, err := (*b.page).Goto(jobURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return nil, fmt.Errorf("导航到岗位详情页失败: %w", err)
	}

	time.Sleep(1 * time.Second)

	card := &JobCard{
		JobURL: jobURL,
	}

	// 获取职位描述
	descEl := (*b.page).Locator(".job-detail .detail-section .text")
	desc, _ := descEl.TextContent()
	card.JobDescription = desc

	// 获取 HR 活跃状态
	activeEl := (*b.page).Locator(".boss-active-status")
	activeStatus, _ := activeEl.TextContent()
	card.HRActiveStatus = activeStatus

	return card, nil
}

// SaveJobs 保存岗位到数据库
// 将解析出的岗位列表批量保存到数据库
// 参数：
// - jobs: 岗位卡片列表
// 返回值：
// - error: 保存过程中的错误
func (b *BossClient) SaveJobs(jobs []JobCard) error {
	for _, job := range jobs {
		// 转换为数据库存储模型
		bossData := storage.BossData{
			EncryptID:      job.EncryptID,
			EncryptUserID:  job.EncryptUserID,
			CompanyName:    job.CompanyName,
			JobName:        job.JobName,
			Salary:         job.Salary,
			Location:       job.Location,
			Experience:     job.Experience,
			Degree:         job.Degree,
			HRName:         job.HRName,
			HRPosition:     job.HRPosition,
			HRActiveStatus: job.HRActiveStatus,
			JobDescription: job.JobDescription,
			JobURL:         job.JobURL,
			DeliveryStatus: "pending", // 初始投递状态为待投递
		}

		if err := storage.Create(&bossData); err != nil {
			config.Warn("保存岗位失败: ", err)
		}
	}

	config.Info("已保存 ", len(jobs), " 个岗位到数据库")
	return nil
}
