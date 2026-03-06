// Package job51 提供前程无忧(51Job)平台的自动化功能
// 注意：当前版本为占位实现，功能待完善
package job51

import (
	"fmt"

	"github.com/playwright-community/playwright-go"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/platform"
	"github.com/yahao333/get_jobs/internal/service"
)

// Job51Client 51Job客户端
// 实现 platform.PlatformClient 接口
type Job51Client struct {
	browser           *playwright.Browser
	page              *playwright.Page
	cookieFile        string
	isLoggedIn        bool
	deliveredToday    int
	greetingGenerator *service.GreetingGenerator
}

// 初始化时注册平台
func init() {
	platform.RegisterPlatform(platform.Platform51Job, func() platform.PlatformClient {
		return NewJob51Client("")
	})
}

// NewJob51Client 创建51Job客户端
func NewJob51Client(cookieFile string) *Job51Client {
	return &Job51Client{
		cookieFile: cookieFile,
		isLoggedIn: false,
	}
}

// Init 初始化客户端
func (c *Job51Client) Init() error {
	config.Info("初始化51Job客户端...")
	// TODO: 启动浏览器
	return nil
}

// Login 登录
func (c *Job51Client) Login() error {
	config.Info("51Job登录功能待实现")
	return fmt.Errorf("51Job登录功能待实现")
}

// CheckLogin 检查登录状态
func (c *Job51Client) CheckLogin() (bool, error) {
	return c.isLoggedIn, nil
}

// Search 搜索岗位
func (c *Job51Client) Search(cityCode, keyword string, criteria *platform.FilterCriteria) ([]platform.JobInfo, error) {
	config.Info("51Job搜索功能待实现")
	return nil, fmt.Errorf("51Job搜索功能待实现")
}

// Deliver 投递简历
func (c *Job51Client) Deliver(job *platform.JobInfo, message string) (bool, string, error) {
	config.Info("51Job投递功能待实现")
	return false, "功能待实现", fmt.Errorf("51Job投递功能待实现")
}

// GetDeliveryCount 获取今日投递数量
func (c *Job51Client) GetDeliveryCount() int {
	return c.deliveredToday
}

// Close 关闭客户端
func (c *Job51Client) Close() error {
	if c.browser != nil {
		c.browser.Close()
	}
	return nil
}

// SetBrowser 设置浏览器
func (c *Job51Client) SetBrowser(browser *playwright.Browser, page *playwright.Page) {
	c.browser = browser
	c.page = page
}

// Job51JobCard 51Job岗位卡片
type Job51JobCard struct {
	JobID          string
	CompanyName    string
	JobName        string
	Salary         string
	Location       string
	Experience     string
	Degree         string
	HRName         string
	JobDescription string
	JobURL         string
}

// SaveToJobInfo 转换为通用岗位信息
func (j *Job51JobCard) SaveToJobInfo() *platform.JobInfo {
	return &platform.JobInfo{
		Platform:       platform.Platform51Job,
		JobID:          j.JobID,
		CompanyName:    j.CompanyName,
		JobName:        j.JobName,
		Salary:         j.Salary,
		Location:       j.Location,
		Experience:     j.Experience,
		Degree:         j.Degree,
		HRName:         j.HRName,
		JobDescription: j.JobDescription,
		JobURL:         j.JobURL,
	}
}

// 预定义的热门城市编码
var CityCodes = map[string]string{
	"北京": "010000",
	"上海": "020000",
	"广州": "030200",
	"深圳": "030090",
	"杭州": "060200",
	"南京": "060300",
	"苏州": "060500",
	"成都": "070200",
	"武汉": "170000",
	"西安": "270000",
}
