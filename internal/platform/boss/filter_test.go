package boss

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	appConfig "github.com/yahao333/get_jobs/internal/config"
)

// TestParseSalaryKUnit 测试解析 K 为单位的薪资字符串
func TestParseSalaryKUnit(t *testing.T) {
	tests := []struct {
		name    string
		salary  string
		wantMin int
		wantMax int
	}{
		{
			name:    "标准 K 单位薪资",
			salary:  "15K-30K",
			wantMin: 15000,
			wantMax: 30000,
		},
		{
			name:    "小数值 K 单位薪资",
			salary:  "10K-25K",
			wantMin: 10000,
			wantMax: 25000,
		},
		{
			name:    "大数值 K 单位薪资",
			salary:  "50K-80K",
			wantMin: 50000,
			wantMax: 80000,
		},
		{
			name:    "带空格的 K 单位薪资",
			salary:  "15K - 30K",
			wantMin: 15000,
			wantMax: 30000,
		},
		{
			name:    "小写 k",
			salary:  "15k-30k",
			wantMin: 15000,
			wantMax: 30000,
		},
		{
			name:    "单边 K 值",
			salary:  "20K",
			wantMin: 0,
			wantMax: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := parseSalary(tt.salary)
			assert.Equal(t, tt.wantMin, min, "最小薪资不匹配")
			assert.Equal(t, tt.wantMax, max, "最大薪资不匹配")
		})
	}
}

// TestParseSalaryWanUnit 测试解析万为单位的薪资字符串
func TestParseSalaryWanUnit(t *testing.T) {
	tests := []struct {
		name    string
		salary  string
		wantMin int
		wantMax int
	}{
		{
			name:    "标准万单位薪资",
			salary:  "1.5万-3万",
			wantMin: 15000,
			wantMax: 30000,
		},
		{
			name:    "带小数的万单位",
			salary:  "2万-4.5万",
			wantMin: 20000,
			wantMax: 45000,
		},
		{
			name:    "大数值万单位",
			salary:  "5万-10万",
			wantMin: 50000,
			wantMax: 100000,
		},
		{
			name:    "带空格的万单位",
			salary:  "1.5万 - 3万",
			wantMin: 15000,
			wantMax: 30000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := parseSalary(tt.salary)
			assert.Equal(t, tt.wantMin, min, "最小薪资不匹配")
			assert.Equal(t, tt.wantMax, max, "最大薪资不匹配")
		})
	}
}

// TestParseSalaryInvalid 测试无效的薪资字符串
func TestParseSalaryInvalid(t *testing.T) {
	tests := []struct {
		name    string
		salary  string
		wantMin int
		wantMax int
	}{
		{
			name:    "空字符串",
			salary:  "",
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "无单位的数字",
			salary:  "15000-30000",
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "无效格式",
			salary:  "abc-def",
			wantMin: 0,
			wantMax: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := parseSalary(tt.salary)
			assert.Equal(t, tt.wantMin, min, "最小薪资不匹配")
			assert.Equal(t, tt.wantMax, max, "最大薪资不匹配")
		})
	}
}

// TestParseSalaryToK 测试将薪资转换为 K 单位
func TestParseSalaryToK(t *testing.T) {
	tests := []struct {
		name   string
		salary string
		wantK  int
	}{
		{
			name:   "标准 K 单位",
			salary: "15K",
			wantK:  15,
		},
		{
			name:   "万单位转换为 K",
			salary: "2万",
			wantK:  20,
		},
		{
			name:   "小数万单位转换为 K",
			salary: "1.5万",
			wantK:  15,
		},
		{
			name:   "空字符串",
			salary: "",
			wantK:  0,
		},
		{
			name:   "带空格",
			salary: "15 K",
			wantK:  15,
		},
		{
			name:   "小写 k",
			salary: "20k",
			wantK:  20,
		},
		{
			name:   "无效格式",
			salary: "abc",
			wantK:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSalaryToK(tt.salary)
			assert.Equal(t, tt.wantK, result, "K 单位转换结果不匹配")
		})
	}
}

// TestCheckCompanyBlacklist 测试公司黑名单检查
func TestCheckCompanyBlacklist(t *testing.T) {
	config := &FilterConfig{
		CompanyBlacklist: []string{"培训机构", "保险公司", "中介", "骗子公司"},
	}

	tests := []struct {
		name        string
		companyName string
		wantBlocked bool
	}{
		{
			name:        "公司名匹配黑名单",
			companyName: "某培训机构",
			wantBlocked: true,
		},
		{
			name:        "公司名包含黑名单关键词",
			companyName: "某某保险公司分公司",
			wantBlocked: true,
		},
		{
			name:        "公司名不在黑名单",
			companyName: "阿里巴巴",
			wantBlocked: false,
		},
		{
			name:        "空公司名",
			companyName: "",
			wantBlocked: false,
		},
		{
			name:        "部分匹配",
			companyName: "华为技术有限公司",
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.checkCompanyBlacklist(tt.companyName)
			assert.Equal(t, tt.wantBlocked, result)
		})
	}
}

// TestCheckHRBlacklist 测试 HR 黑名单检查
func TestCheckHRBlacklist(t *testing.T) {
	config := &FilterConfig{
		HRBlacklist: []string{"招聘专员", "招聘经理", "猎头"},
	}

	tests := []struct {
		name        string
		hrName      string
		hrPosition  string
		wantBlocked bool
	}{
		{
			name:        "HR职位匹配黑名单",
			hrName:      "张三",
			hrPosition:  "招聘专员",
			wantBlocked: true,
		},
		{
			name:        "HR职位包含黑名单关键词",
			hrName:      "李四",
			hrPosition:  "高级猎头顾问",
			wantBlocked: true,
		},
		{
			name:        "HR职位不在黑名单",
			hrName:      "王五",
			hrPosition:  "技术经理",
			wantBlocked: false,
		},
		{
			name:        "空HR职位",
			hrName:      "赵六",
			hrPosition:  "",
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.checkHRBlacklist(tt.hrName, tt.hrPosition)
			assert.Equal(t, tt.wantBlocked, result)
		})
	}
}

// TestCheckJobBlacklist 测试职位黑名单检查
func TestCheckJobBlacklist(t *testing.T) {
	config := &FilterConfig{
		JobBlacklist: []string{"实习", "兼职", "代理", "招商"},
	}

	tests := []struct {
		name        string
		jobName     string
		wantBlocked bool
	}{
		{
			name:        "职位名匹配黑名单",
			jobName:     "后端开发实习",
			wantBlocked: true,
		},
		{
			name:        "职位名包含黑名单关键词",
			jobName:     "诚聘兼职",
			wantBlocked: true,
		},
		{
			name:        "职位名不在黑名单",
			jobName:     "Go后端开发工程师",
			wantBlocked: false,
		},
		{
			name:        "空职位名",
			jobName:     "",
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.checkJobBlacklist(tt.jobName)
			assert.Equal(t, tt.wantBlocked, result)
		})
	}
}

// TestIsDeadHR 测试判断 HR 是否不活跃
func TestIsDeadHR(t *testing.T) {
	config := &FilterConfig{
		FilterDeadHR: true,
	}

	tests := []struct {
		name         string
		activeStatus string
		wantIsDead   bool
	}{
		{
			name:         "1年前在线",
			activeStatus: "1年前在线",
			wantIsDead:   true,
		},
		{
			name:         "3年前在线",
			activeStatus: "3年前在线",
			wantIsDead:   true,
		},
		{
			name:         "7年前在线",
			activeStatus: "7年前在线",
			wantIsDead:   true,
		},
		{
			name:         "今日活跃",
			activeStatus: "今日活跃",
			wantIsDead:   false,
		},
		{
			name:         "刚刚活跃",
			activeStatus: "刚刚活跃",
			wantIsDead:   false,
		},
		{
			name:         "昨日活跃",
			activeStatus: "昨日活跃",
			wantIsDead:   false,
		},
		{
			name:         "3日内活跃",
			activeStatus: "3日内活跃",
			wantIsDead:   false,
		},
		{
			name:         "空状态",
			activeStatus: "",
			wantIsDead:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.isDeadHR(tt.activeStatus)
			assert.Equal(t, tt.wantIsDead, result)
		})
	}
}

// TestCheckSalary 测试薪资范围检查
func TestCheckSalary(t *testing.T) {
	tests := []struct {
		name      string
		salary    string
		minSalary string
		maxSalary string
		wantPass  bool
	}{
		{
			name:      "薪资在范围内 - K单位",
			salary:    "20K-40K",
			minSalary: "15K",
			maxSalary: "50K",
			wantPass:  true,
		},
		{
			name:      "薪资低于最小值",
			salary:    "10K-20K",
			minSalary: "15K",
			maxSalary: "50K",
			wantPass:  false,
		},
		{
			name:      "薪资高于最大值",
			salary:    "60K-80K",
			minSalary: "15K",
			maxSalary: "50K",
			wantPass:  false,
		},
		{
			name:      "只有最小值限制 - 通过",
			salary:    "20K-40K",
			minSalary: "15K",
			maxSalary: "",
			wantPass:  true,
		},
		{
			name:      "只有最小值限制 - 不通过",
			salary:    "10K-20K",
			minSalary: "15K",
			maxSalary: "",
			wantPass:  false,
		},
		{
			name:      "只有最大值限制 - 通过",
			salary:    "20K-40K",
			minSalary: "",
			maxSalary: "50K",
			wantPass:  true,
		},
		{
			name:      "只有最大值限制 - 不通过",
			salary:    "60K-80K",
			minSalary: "",
			maxSalary: "50K",
			wantPass:  false,
		},
		{
			name:      "无限制 - 空薪资",
			salary:    "",
			minSalary: "15K",
			maxSalary: "50K",
			wantPass:  true,
		},
		{
			name:      "无限制 - 空配置",
			salary:    "20K-40K",
			minSalary: "",
			maxSalary: "",
			wantPass:  true,
		},
		{
			name:      "万单位薪资范围检查",
			salary:    "2万-4万",
			minSalary: "1.5万",
			maxSalary: "5万",
			wantPass:  true,
		},
		{
			name:      "万单位薪资低于最小值",
			salary:    "1万-2万",
			minSalary: "1.5万",
			maxSalary: "5万",
			wantPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &FilterConfig{
				MinSalary: tt.minSalary,
				MaxSalary: tt.maxSalary,
			}
			result := config.checkSalary(tt.salary)
			assert.Equal(t, tt.wantPass, result)
		})
	}
}

// TestFilterJob 测试完整的岗位过滤功能
func TestFilterJob(t *testing.T) {
	config := &FilterConfig{
		FilterDeadHR:     true,
		CompanyBlacklist: []string{"培训机构", "保险公司"},
		HRBlacklist:      []string{"猎头"},
		JobBlacklist:     []string{"实习", "兼职"},
		MinSalary:        "15K",
		MaxSalary:        "50K",
	}

	tests := []struct {
		name       string
		job        *JobCard
		wantPass   bool
		wantReason string
	}{
		{
			name: "岗位通过所有过滤",
			job: &JobCard{
				CompanyName:    "阿里巴巴",
				JobName:        "Go后端开发",
				Salary:         "20K-40K",
				HRName:         "张三",
				HRPosition:     "技术经理",
				HRActiveStatus: "今日活跃",
			},
			wantPass:   true,
			wantReason: "",
		},
		{
			name: "公司黑名单过滤",
			job: &JobCard{
				CompanyName: "某培训机构",
				JobName:     "Go后端开发",
				Salary:      "20K-40K",
			},
			wantPass:   false,
			wantReason: "公司黑名单",
		},
		{
			name: "HR黑名单过滤",
			job: &JobCard{
				CompanyName: "阿里巴巴",
				JobName:     "Go后端开发",
				Salary:      "20K-40K",
				HRPosition:  "猎头顾问",
			},
			wantPass:   false,
			wantReason: "HR黑名单",
		},
		{
			name: "职位黑名单过滤",
			job: &JobCard{
				CompanyName: "阿里巴巴",
				JobName:     "后端开发实习",
				Salary:      "20K-40K",
			},
			wantPass:   false,
			wantReason: "职位黑名单",
		},
		{
			name: "不活跃HR过滤",
			job: &JobCard{
				CompanyName:    "阿里巴巴",
				JobName:        "Go后端开发",
				Salary:         "20K-40K",
				HRActiveStatus: "3年前在线",
			},
			wantPass:   false,
			wantReason: "不活跃HR",
		},
		{
			name: "薪资低于最小值过滤",
			job: &JobCard{
				CompanyName: "阿里巴巴",
				JobName:     "Go后端开发",
				Salary:      "10K-15K",
			},
			wantPass:   false,
			wantReason: "薪资不符合要求",
		},
		{
			name: "薪资高于最大值过滤",
			job: &JobCard{
				CompanyName: "阿里巴巴",
				JobName:     "Go后端开发",
				Salary:      "60K-80K",
			},
			wantPass:   false,
			wantReason: "薪资不符合要求",
		},
		{
			name: "禁用不活跃HR过滤时应该通过",
			job: &JobCard{
				CompanyName:    "阿里巴巴",
				JobName:        "Go后端开发",
				Salary:         "20K-40K",
				HRActiveStatus: "3年前在线",
			},
			wantPass:   true,
			wantReason: "",
		},
	}

	// 测试禁用不活跃HR过滤的场景
	configNoDeadHR := &FilterConfig{
		FilterDeadHR:     false,
		CompanyBlacklist: []string{"培训机构", "保险公司"},
		HRBlacklist:      []string{"猎头"},
		JobBlacklist:     []string{"实习", "兼职"},
		MinSalary:        "15K",
		MaxSalary:        "50K",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个新的 config 引用，以便修改 FilterDeadHR
			useConfig := config
			if tt.name == "禁用不活跃HR过滤时应该通过" {
				useConfig = configNoDeadHR
			}

			pass, reason := useConfig.FilterJob(tt.job)
			assert.Equal(t, tt.wantPass, pass, "过滤结果不匹配")
			if !tt.wantPass {
				assert.Equal(t, tt.wantReason, reason, "过滤原因不匹配")
			}
		})
	}
}

// TestFilterJobs 测试批量岗位过滤
func TestFilterJobs(t *testing.T) {
	config := &FilterConfig{
		FilterDeadHR:     true,
		CompanyBlacklist: []string{"培训机构"},
		HRBlacklist:      []string{},
		JobBlacklist:     []string{"实习"},
		MinSalary:        "10K",
		MaxSalary:        "50K",
	}

	jobs := []JobCard{
		{
			CompanyName: "阿里巴巴",
			JobName:     "Go后端开发",
			Salary:      "20K-40K",
		},
		{
			CompanyName: "某培训机构",
			JobName:     "Go后端开发",
			Salary:      "20K-40K",
		},
		{
			CompanyName: "腾讯",
			JobName:     "后端开发实习",
			Salary:      "20K-40K",
		},
		{
			CompanyName: "字节跳动",
			JobName:     "Go后端开发",
			Salary:      "5K-10K", // 低于最小值
		},
		{
			CompanyName: "美团",
			JobName:     "Go后端开发",
			Salary:      "30K-60K", // 高于最大值
		},
		{
			CompanyName: "百度",
			JobName:     "Go后端开发",
			Salary:      "25K-45K",
		},
	}

	filtered, reasons := config.FilterJobs(jobs)

	// 应该过滤掉 4 个岗位: 培训机构、实习、薪资低、薪资高
	assert.Equal(t, 2, len(filtered), "过滤后的岗位数量不匹配")
	assert.Equal(t, 4, len(reasons), "过滤原因数量不匹配")

	// 验证保留下来的岗位
	assert.Equal(t, "阿里巴巴", filtered[0].CompanyName)
	assert.Equal(t, "百度", filtered[1].CompanyName)
}

// TestFilterConfigDefault 测试默认配置
func TestFilterConfigDefault(t *testing.T) {
	// 初始化配置，避免 panic
	appConfig.Config = viper.New()

	// 注意: 这个测试依赖配置文件，如果配置文件不存在或有不同值，测试可能失败
	// 我们只测试结构体是否正确创建
	config := DefaultFilterConfig()
	assert.NotNil(t, config)
	// FilterConfig 应该有这些字段
	_ = config.FilterDeadHR
	_ = config.CompanyBlacklist
	_ = config.HRBlacklist
	_ = config.JobBlacklist
	_ = config.MinSalary
	_ = config.MaxSalary
}
