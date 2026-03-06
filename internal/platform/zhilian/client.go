// Package zhilian 提供智联招聘平台的自动化功能
// 注意：当前版本为占位实现，功能待完善
package zhilian

import (
	"fmt"

	"github.com/playwright-community/playwright-go"

	"github.com/yahao333/get_jobs/internal/config"
	"github.com/yahao333/get_jobs/internal/platform"
	"github.com/yahao333/get_jobs/internal/service"
)

// ZhilianClient 智联招聘客户端
// 实现 platform.PlatformClient 接口
type ZhilianClient struct {
	browser           *playwright.Browser
	page              *playwright.Page
	cookieFile        string
	isLoggedIn        bool
	deliveredToday    int
	greetingGenerator *service.GreetingGenerator
}

// 初始化时注册平台
func init() {
	platform.RegisterPlatform(platform.PlatformZhilian, func() platform.PlatformClient {
		return NewZhilianClient("")
	})
}

// NewZhilianClient 创建智联招聘客户端
func NewZhilianClient(cookieFile string) *ZhilianClient {
	return &ZhilianClient{
		cookieFile: cookieFile,
		isLoggedIn: false,
	}
}

// Init 初始化客户端
func (c *ZhilianClient) Init() error {
	config.Info("初始化智联招聘客户端...")
	// TODO: 启动浏览器
	return nil
}

// Login 登录
func (c *ZhilianClient) Login() error {
	config.Info("智联招聘登录功能待实现")
	return fmt.Errorf("智联招聘登录功能待实现")
}

// CheckLogin 检查登录状态
func (c *ZhilianClient) CheckLogin() (bool, error) {
	return c.isLoggedIn, nil
}

// Search 搜索岗位
func (c *ZhilianClient) Search(cityCode, keyword string, criteria *platform.FilterCriteria) ([]platform.JobInfo, error) {
	config.Info("智联招聘搜索功能待实现")
	return nil, fmt.Errorf("智联招聘搜索功能待实现")
}

// Deliver 投递简历
func (c *ZhilianClient) Deliver(job *platform.JobInfo, message string) (bool, string, error) {
	config.Info("智联招聘投递功能待实现")
	return false, "功能待实现", fmt.Errorf("智联招聘投递功能待实现")
}

// GetDeliveryCount 获取今日投递数量
func (c *ZhilianClient) GetDeliveryCount() int {
	return c.deliveredToday
}

// Close 关闭客户端
func (c *ZhilianClient) Close() error {
	if c.browser != nil {
		c.browser.Close()
	}
	return nil
}

// SetBrowser 设置浏览器
func (c *ZhilianClient) SetBrowser(browser *playwright.Browser, page *playwright.Page) {
	c.browser = browser
	c.page = page
}

// ZhilianJobCard 智联招聘岗位卡片
type ZhilianJobCard struct {
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
func (j *ZhilianJobCard) SaveToJobInfo() *platform.JobInfo {
	return &platform.JobInfo{
		Platform:       platform.PlatformZhilian,
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
	"北京": "530",
	"上海": "538",
	"广州": "763",
	"深圳": "765",
	"杭州": "653",
	"南京": "635",
	"苏州": "639",
	"成都": "828",
	"武汉": "736",
	"西安": "854",
}
