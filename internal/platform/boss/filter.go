// Package boss 提供 Boss直聘平台的自动化功能
// 岗位过滤模块：实现多维度岗位过滤，包括公司黑名单、HR黑名单、薪资范围等
package boss

import (
	"fmt"     // 格式化输出
	"strings" // 字符串处理

	"github.com/yahao333/get_jobs/internal/config"  // 配置管理模块
	"github.com/yahao333/get_jobs/internal/storage" // 数据存储模块
)

// FilterConfig 过滤配置结构体
// 定义了各种过滤规则，用于筛选符合要求的岗位
type FilterConfig struct {
	FilterDeadHR     bool     // 是否过滤不活跃 HR（如长时间未在线的 HR）
	CompanyBlacklist []string // 公司黑名单关键词列表
	HRBlacklist      []string // HR 职位黑名单关键词列表
	JobBlacklist     []string // 职位名称黑名单关键词列表
	MinSalary        string   // 最低薪资要求（如 "15K" 或 "1.5万"）
	MaxSalary        string   // 最高薪资要求（如 "30K" 或 "3万"）
}

// DefaultFilterConfig 获取默认过滤配置
// 从配置文件中读取过滤相关设置，返回配置好的 FilterConfig 实例
// 配置项说明：
// - filter.filter_dead_hr: 是否启用不活跃 HR 过滤
// - filter.company_blacklist: 公司黑名单关键词列表
// - filter.hr_blacklist: HR 职位黑名单关键词列表
// - filter.job_blacklist: 职位名称黑名单关键词列表
func DefaultFilterConfig() *FilterConfig {
	return &FilterConfig{
		FilterDeadHR:     config.GetBool("filter.filter_dead_hr"),
		CompanyBlacklist: config.GetStringSlice("filter.company_blacklist"),
		HRBlacklist:      config.GetStringSlice("filter.hr_blacklist"),
		JobBlacklist:     config.GetStringSlice("filter.job_blacklist"),
	}
}

// FilterJob 过滤单个岗位
// 根据配置的过滤规则对岗位进行全面检查，返回是否通过及过滤原因
// 过滤顺序：
// 1. 公司黑名单检查 - 检查公司名称是否包含黑名单关键词
// 2. HR 黑名单检查 - 检查 HR 职位名称是否包含黑名单关键词
// 3. 职位黑名单检查 - 检查职位名称是否包含黑名单关键词
// 4. 不活跃 HR 检查 - 检查 HR 是否长时间未在线
// 5. 薪资范围检查 - 检查薪资是否符合配置的范围要求
// 参数：
// - job: 岗位卡片信息
// 返回值：
// - bool: 岗位是否通过所有过滤条件
// - string: 如果未通过，返回具体的过滤原因
func (f *FilterConfig) FilterJob(job *JobCard) (bool, string) {
	// 检查1: 公司黑名单
	if f.checkCompanyBlacklist(job.CompanyName) {
		return false, "公司黑名单"
	}

	// 检查2: HR 黑名单
	if f.checkHRBlacklist(job.HRName, job.HRPosition) {
		return false, "HR黑名单"
	}

	// 检查3: 职位黑名单
	if f.checkJobBlacklist(job.JobName) {
		return false, "职位黑名单"
	}

	// 检查4: 不活跃 HR
	if f.FilterDeadHR && f.isDeadHR(job.HRActiveStatus) {
		return false, "不活跃HR"
	}

	// 检查5: 薪资范围
	if !f.checkSalary(job.Salary) {
		return false, "薪资不符合要求"
	}

	return true, ""
}

// checkCompanyBlacklist 检查公司黑名单
// 使用字符串包含匹配，检查公司名称是否包含黑名单中的任意关键词
// 参数：
// - companyName: 公司名称
// 返回值：
// - bool: 公司是否在黑名单中
func (f *FilterConfig) checkCompanyBlacklist(companyName string) bool {
	// 空公司名直接放行
	if companyName == "" {
		return false
	}
	// 遍历黑名单关键词，使用字符串包含匹配
	for _, keyword := range f.CompanyBlacklist {
		if strings.Contains(companyName, keyword) {
			config.Debug("过滤公司（黑名单）: ", companyName, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// checkHRBlacklist 检查 HR 职位黑名单
// 使用字符串包含匹配，检查 HR 职位名称是否包含黑名单中的任意关键词
// 注意：此处检查的是 HR 的职位名称（如"招聘专员"），而非 HR 姓名
// 参数：
// - hrName: HR 姓名（当前未使用，保留扩展）
// - hrPosition: HR 职位名称
// 返回值：
// - bool: HR 职位是否在黑名单中
func (f *FilterConfig) checkHRBlacklist(hrName, hrPosition string) bool {
	// 空职位名直接放行
	if hrPosition == "" {
		return false
	}
	// 遍历黑名单关键词
	for _, keyword := range f.HRBlacklist {
		if strings.Contains(hrPosition, keyword) {
			config.Debug("过滤HR（黑名单）: ", hrPosition, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// checkJobBlacklist 检查职位名称黑名单
// 使用字符串包含匹配，检查职位名称是否包含黑名单中的任意关键词
// 参数：
// - jobName: 职位名称
// 返回值：
// - bool: 职位是否在黑名单中
func (f *FilterConfig) checkJobBlacklist(jobName string) bool {
	// 空职位名直接放行
	if jobName == "" {
		return false
	}
	// 遍历黑名单关键词
	for _, keyword := range f.JobBlacklist {
		if strings.Contains(jobName, keyword) {
			config.Debug("过滤职位（黑名单）: ", jobName, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// isDeadHR 判断 HR 是否不活跃
// 通过检查 HR 活跃状态中是否包含"年前"来判断
// 例如："7年前在线"表示该 HR 已经 7 年没有登录，属于不活跃状态
// 参数：
// - activeStatus: HR 活跃状态字符串
// 返回值：
// - bool: HR 是否不活跃
func (f *FilterConfig) isDeadHR(activeStatus string) bool {
	// 空状态直接放行
	if activeStatus == "" {
		return false
	}
	// "年前在线" 表示不活跃
	if strings.Contains(activeStatus, "年前") {
		config.Debug("过滤不活跃HR: ", activeStatus)
		return true
	}
	return false
}

// checkSalary 检查薪资范围是否符合要求
// 支持 K（千）和万两种单位，支持薪资范围字符串解析
// 算法说明：
// - 如果未配置薪资范围（MinSalary 和 MaxSalary 都为空），直接返回 true
// - 解析岗位薪资字符串，提取最低和最高薪资
// - 薪资范围匹配逻辑：
//   - 如果配置了最低薪资，岗位最高薪资必须 >= 最低薪资
//   - 如果配置了最高薪资，岗位最低薪资必须 <= 最高薪资
//   - 两个条件都满足时返回 true
//
// 参数：
// - salary: 岗位薪资字符串，格式如 "15K-30K" 或 "1.5万-3万"
// 返回值：
// - bool: 薪资是否符合要求
func (f *FilterConfig) checkSalary(salary string) bool {
	// 如果薪资为空或未配置范围要求，直接通过
	if salary == "" || (f.MinSalary == "" && f.MaxSalary == "") {
		return true
	}

	// 解析薪资字符串，如 "15K-30K"
	salaryMin, salaryMax := parseSalary(salary)
	// 如果无法解析薪资（如格式不正确），默认通过
	if salaryMin == 0 && salaryMax == 0 {
		return true
	}

	// 解析配置的薪资范围（转换为 K 单位进行对比）
	// 例如："15K" 转换为 15，"1.5万" 转换为 15
	minK := parseSalaryToK(f.MinSalary)
	maxK := parseSalaryToK(f.MaxSalary)

	// 将 K 单位转换为元
	minSalaryReq := minK * 1000
	maxSalaryReq := maxK * 1000

	// 检查是否符合范围
	// 最低要求：如果配置的最低薪资 > 0，岗位最低薪资必须 >= 最低薪资
	if minSalaryReq > 0 && salaryMin < minSalaryReq {
		return false
	}
	// 最高要求：如果配置的最高薪资 > 0，岗位最高薪资必须 <= 最高薪资
	if maxSalaryReq > 0 && salaryMax > maxSalaryReq {
		return false
	}

	return true
}

// parseSalary 解析薪资字符串
// 支持两种格式：
// - K 单位：如 "15K-30K"，返回 15000-30000（单位：元）
// - 万单位：如 "1.5万-3万"，返回 15000-30000（单位：元）
// 参数：
// - salary: 薪资字符串
// 返回值：
// - min: 最低薪资（单位：元）
// - max: 最高薪资（单位：元）
func parseSalary(salary string) (min, max int) {
	// 标准化处理：转大写，去空格
	salary = strings.ToUpper(salary)
	salary = strings.ReplaceAll(salary, " ", "")

	// 处理 K 单位（如 15K-30K）
	if strings.Contains(salary, "K") {
		var minStr, maxStr string
		parts := strings.Split(salary, "-")
		if len(parts) == 2 {
			minStr = strings.ReplaceAll(parts[0], "K", "")
			maxStr = strings.ReplaceAll(parts[1], "K", "")
		}
		if minStr != "" {
			fmt.Sscanf(minStr, "%d", &min)
		}
		if maxStr != "" {
			fmt.Sscanf(maxStr, "%d", &max)
		}
		// K 转换为元：K * 1000
		return min * 1000, max * 1000
	}

	// 处理万单位（如 1.5万-3万）
	if strings.Contains(salary, "万") {
		var minF, maxF float64
		parts := strings.Split(salary, "-")
		if len(parts) == 2 {
			minStr := strings.ReplaceAll(parts[0], "万", "")
			maxStr := strings.ReplaceAll(parts[1], "万", "")
			fmt.Sscanf(minStr, "%f", &minF)
			fmt.Sscanf(maxStr, "%f", &maxF)
		}
		// 万转换为元：万 * 10000
		return int(minF * 10000), int(maxF * 10000)
	}

	return 0, 0
}

// parseSalaryToK 将薪资字符串转换为 K 单位
// 用于与配置中的薪资范围进行对比
// 例如："15K" 转换为 15，"1.5万" 转换为 15
// 参数：
// - salary: 薪资字符串（如 "15K" 或 "1.5万"）
// 返回值：
// - int: K 单位的薪资值（如 15 表示 15K）
func parseSalaryToK(salary string) int {
	if salary == "" {
		return 0
	}
	salary = strings.ToUpper(salary)
	salary = strings.ReplaceAll(salary, " ", "")

	var value int
	if strings.Contains(salary, "K") {
		// K 单位：直接取数值
		str := strings.ReplaceAll(salary, "K", "")
		fmt.Sscanf(str, "%d", &value)
	} else if strings.Contains(salary, "万") {
		// 万单位：万 * 10 = K
		str := strings.ReplaceAll(salary, "万", "")
		var valueF float64
		fmt.Sscanf(str, "%f", &valueF)
		value = int(valueF * 10)
	}
	return value
}

// FilterJobs 批量过滤岗位
// 对岗位列表进行逐个过滤，返回通过过滤的岗位列表和过滤原因列表
// 参数：
// - jobs: 岗位卡片列表
// 返回值：
// - []JobCard: 通过过滤的岗位列表
// - []string: 被过滤的岗位及原因列表
func (f *FilterConfig) FilterJobs(jobs []JobCard) ([]JobCard, []string) {
	filtered := make([]JobCard, 0)
	reasons := make([]string, 0)

	for _, job := range jobs {
		passed, reason := f.FilterJob(&job)
		if passed {
			// 通过过滤，加入结果列表
			filtered = append(filtered, job)
		} else {
			// 未通过过滤，记录原因
			reasons = append(reasons, job.JobName+" - "+reason)
		}
	}

	config.Info("过滤前: ", len(jobs), " 个岗位, 过滤后: ", len(filtered), " 个岗位")
	return filtered, reasons
}

// AddCompanyToBlacklist 添加公司到黑名单
// 将公司名称添加到黑名单并持久化存储
// 参数：
// - companyName: 公司名称
// 返回值：
// - error: 添加失败时返回错误
func (f *FilterConfig) AddCompanyToBlacklist(companyName string) error {
	blacklist := storage.Blacklist{
		Keyword: companyName,
		Type:    "company",
		Source:  "auto",
	}
	return storage.Create(&blacklist)
}

// AddHRToBlacklist 添加 HR 到黑名单
// 将 HR 姓名添加到黑名单并持久化存储
// 参数：
// - hrName: HR 姓名
// 返回值：
// - error: 添加失败时返回错误
func (f *FilterConfig) AddHRToBlacklist(hrName string) error {
	blacklist := storage.Blacklist{
		Keyword: hrName,
		Type:    "hr",
		Source:  "auto",
	}
	return storage.Create(&blacklist)
}
