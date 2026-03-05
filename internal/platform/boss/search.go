package boss

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/loks666/get_jobs/internal/config"
	"github.com/loks666/get_jobs/internal/storage"
)

// JobCard 岗位卡片信息
type JobCard struct {
	EncryptID      string // 岗位加密ID
	EncryptUserID  string // HR加密ID
	CompanyName    string // 公司名称
	JobName        string // 职位名称
	Salary         string // 薪资范围
	Location       string // 工作地点
	Experience     string // 经验要求
	Degree         string // 学历要求
	HRName         string // HR姓名
	HRPosition     string // HR职位
	HRActiveStatus string // HR活跃状态
	JobDescription string // 职位描述
	JobURL         string // 岗位URL
}

// SearchJobs 搜索岗位
func (b *BossClient) SearchJobs(cityCode, keyword, jobType, salary, experience, degree string) ([]JobCard, error) {
	if b.page == nil {
		return nil, fmt.Errorf("浏览器页面未初始化")
	}

	// 构建搜索 URL
	searchURL := fmt.Sprintf(
		"https://www.zhipin.com/web/geek/job?query=%s&cityCode=%s&jobType=%s",
		keyword, cityCode, jobType,
	)

	config.Info("搜索岗位: ", searchURL)

	// 导航到搜索页面
	_, err := (*b.page).Goto(searchURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return nil, fmt.Errorf("导航到搜索页面失败: %w", err)
	}

	// 等待页面加载
	time.Sleep(2 * time.Second)

	// 滚动页面触发懒加载
	if err := b.scrollToLoadAll(); err != nil {
		config.Warn("滚动加载失败: ", err)
	}

	// 解析岗位列表
	jobCards, err := b.parseJobList()
	if err != nil {
		return nil, fmt.Errorf("解析岗位列表失败: %w", err)
	}

	config.Info("找到 ", len(jobCards), " 个岗位")
	return jobCards, nil
}

// scrollToLoadAll 滚动加载所有岗位
func (b *BossClient) scrollToLoadAll() error {
	if b.page == nil {
		return fmt.Errorf("浏览器页面未初始化")
	}

	lastCount := 0
	for i := 0; i < 10; i++ {
		// 获取当前岗位数量
		count, _ := (*b.page).Locator(".job-card").Count()

		// 滚动到页面底部
		_, err := (*b.page).Evaluate("window.scrollBy(0, 800)")
		if err != nil {
			return err
		}

		// 等待加载
		time.Sleep(1 * time.Second)

		// 检查是否还有新岗位
		newCount, _ := (*b.page).Locator(".job-card").Count()
		if newCount == lastCount && newCount > 0 {
			// 连续两次相同，停止滚动
			break
		}
		lastCount = newCount
	}

	return nil
}

// parseJobList 解析岗位列表
func (b *BossClient) parseJobList() ([]JobCard, error) {
	if b.page == nil {
		return nil, fmt.Errorf("浏览器页面未初始化")
	}

	// 获取所有岗位卡片
	locator := (*b.page).Locator(".job-card")
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		// 尝试其他选择器
		locator = (*b.page).Locator(".job-list-box .job-card")
		count, _ = locator.Count()
		if count == 0 {
			return nil, fmt.Errorf("未找到岗位卡片")
		}
	}

	cards := make([]JobCard, 0, count)
	for i := 0; i < count; i++ {
		card, err := b.parseJobCard(locator.Nth(i))
		if err != nil {
			continue
		}
		cards = append(cards, *card)
	}

	return cards, nil
}

// parseJobCard 解析单个岗位卡片
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

	// 获取经验学历要求
	tagEl := locator.Locator(".tag-list")
	tagText, _ := tagEl.TextContent()
	card.Experience = tagText

	// 获取 HR 信息
	hrNameEl := locator.Locator(".boss-name")
	hrName, _ := hrNameEl.TextContent()
	card.HRName = hrName

	// 获取链接
	linkEl := locator.Locator("a")
	href, _ := linkEl.Attribute("href")
	if href != "" {
		card.JobURL = "https://www.zhipin.com" + href
	}

	return card, nil
}

// GetJobDetail 获取岗位详情
func (b *BossClient) GetJobDetail(jobURL string) (*JobCard, error) {
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
func (b *BossClient) SaveJobs(jobs []JobCard) error {
	for _, job := range jobs {
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
			DeliveryStatus: "pending",
		}

		if err := storage.Create(&bossData); err != nil {
			config.Warn("保存岗位失败: ", err)
		}
	}

	config.Info("已保存 ", len(jobs), " 个岗位到数据库")
	return nil
}
