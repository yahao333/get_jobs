// Package platform 定义招聘平台的通用接口
// 支持 Boss直聘、猎聘、前程无忧、智联招聘等主流招聘平台
package platform

import "time"

// Platform 招聘平台类型
type Platform string

const (
	PlatformBoss   Platform = "boss"   // Boss直聘
	PlatformLiepin Platform = "liepin" // 猎聘
	Platform51Job  Platform = "51job"  // 前程无忧
	PlatformZhilian Platform = "zhilian" // 智联招聘
)

// PlatformInfo 平台信息
type PlatformInfo struct {
	Name        string    // 平台名称
	Code        Platform  // 平台代码
	URL         string    // 官网 URL
	LoginURL    string    // 登录页面 URL
	SearchURL   string    // 搜索页面 URL
	Description string    // 平台描述
}

// AllPlatforms 所有支持的平台
var AllPlatforms = []PlatformInfo{
	{
		Name:        "Boss直聘",
		Code:        PlatformBoss,
		URL:         "https://www.zhipin.com",
		LoginURL:    "https://www.zhipin.com/user/login.html",
		SearchURL:   "https://www.zhipin.com/web/geek/job",
		Description: "BOSS直聘，求职者可以直接和boss沟通的招聘平台",
	},
	{
		Name:        "猎聘",
		Code:        PlatformLiepin,
		URL:         "https://www.liepin.com",
		LoginURL:    "https://www.liepin.com/account/login/",
		SearchURL:   "https://www.liepin.com/zhaopin/",
		Description: "猎聘，高端人才招聘平台",
	},
	{
		Name:        "前程无忧",
		Code:        Platform51Job,
		URL:         "https://www.51job.com",
		LoginURL:    "https://login.51job.com/login.php",
		SearchURL:   "https://search.51job.com",
		Description: "前程无忧，综合性招聘平台",
	},
	{
		Name:        "智联招聘",
		Code:        PlatformZhilian,
		URL:         "https://www.zhaopin.com",
		LoginURL:    "https://www.zhaopin.com/login/",
		SearchURL:   "https://www.zhaopin.com/sou/",
		Description: "智联招聘，老牌招聘平台",
	},
}

// GetPlatformInfo 获取平台信息
func GetPlatformInfo(p Platform) *PlatformInfo {
	for _, info := range AllPlatforms {
		if info.Code == p {
			return &info
		}
	}
	return nil
}

// JobInfo 通用岗位信息
// 所有平台解析的岗位信息都转换为这个统一格式
type JobInfo struct {
	Platform     Platform // 平台来源
	JobID        string   // 平台岗位ID
	CompanyName  string   // 公司名称
	JobName      string   // 职位名称
	Salary       string   // 薪资范围
	Location     string   // 工作地点
	Experience   string   // 经验要求
	Degree       string   // 学历要求
	HRName       string   // HR姓名
	HRPosition   string   // HR职位
	HRActiveStatus string // HR活跃状态
	JobDescription string // 职位描述
	JobURL       string   // 岗位详情页URL
	PublishTime  string   // 发布时间
}

// FilterCriteria 筛选条件
type FilterCriteria struct {
	MinSalary        string   // 最低薪资
	MaxSalary        string   // 最高薪资
	Experience       string   // 经验要求
	Degree           string   // 学历要求
	CompanyBlacklist []string // 公司黑名单
	HRBlacklist      []string // HR黑名单
	JobBlacklist     []string // 职位黑名单
	FilterDeadHR     bool     // 过滤不活跃HR
}

// PlatformClient 平台客户端接口
// 所有平台实现必须实现这个接口
type PlatformClient interface {
	// Init 初始化客户端
	Init() error

	// Login 登录（扫码或cookie）
	Login() error

	// CheckLogin 检查登录状态
	CheckLogin() (bool, error)

	// Search 搜索岗位
	Search(cityCode, keyword string, criteria *FilterCriteria) ([]JobInfo, error)

	// Deliver 投递简历
	Deliver(job *JobInfo, message string) (bool, string, error)

	// GetDeliveryCount 获取今日投递数量
	GetDeliveryCount() int

	// Close 关闭客户端
	Close() error
}

// PlatformFactory 平台工厂
// 根据平台类型创建对应的客户端实例
var PlatformFactory = make(map[Platform]func() PlatformClient)

// RegisterPlatform 注册平台
func RegisterPlatform(p Platform, creator func() PlatformClient) {
	PlatformFactory[p] = creator
}

// CreatePlatform 创建平台客户端
func CreatePlatform(p Platform) (PlatformClient, error) {
	creator, ok := PlatformFactory[p]
	if !ok {
		return nil, ErrPlatformNotSupported
	}
	return creator(), nil
}

// ErrPlatformNotSupported 平台不支持错误
var ErrPlatformNotSupported = &PlatformError{
	Code:    "PLATFORM_NOT_SUPPORTED",
	Message: "不支持的平台类型",
}

// PlatformError 平台错误
type PlatformError struct {
	Code    string
	Message string
}

func (e *PlatformError) Error() string {
	return e.Message
}

// SearchOptions 搜索选项
type SearchOptions struct {
	CityCode   string        // 城市编码
	Keyword    string        // 搜索关键词
	JobType    string        // 职位类型
	PageSize   int           // 每页数量
	MaxPages   int           // 最大页数
	Timeout    time.Duration // 超时时间
}

// DefaultSearchOptions 获取默认搜索选项
func DefaultSearchOptions() *SearchOptions {
	return &SearchOptions{
		PageSize: 15,
		MaxPages: 5,
		Timeout:  5 * time.Minute,
	}
}
