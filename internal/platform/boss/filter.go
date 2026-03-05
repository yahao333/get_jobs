package boss

import (
	"fmt"
	"strings"

	"github.com/loks666/get_jobs/internal/config"
	"github.com/loks666/get_jobs/internal/storage"
)

// FilterConfig 过滤配置
type FilterConfig struct {
	FilterDeadHR      bool
	CompanyBlacklist []string
	HRBlacklist      []string
	JobBlacklist     []string
	MinSalary        string
	MaxSalary        string
}

// DefaultFilterConfig 获取默认过滤配置
func DefaultFilterConfig() *FilterConfig {
	return &FilterConfig{
		FilterDeadHR:      config.GetBool("filter.filter_dead_hr"),
		CompanyBlacklist: config.GetStringSlice("filter.company_blacklist"),
		HRBlacklist:      config.GetStringSlice("filter.hr_blacklist"),
		JobBlacklist:     config.GetStringSlice("filter.job_blacklist"),
	}
}

// FilterJob 过滤岗位
func (f *FilterConfig) FilterJob(job *JobCard) (bool, string) {
	// 检查公司黑名单
	if f.checkCompanyBlacklist(job.CompanyName) {
		return false, "公司黑名单"
	}

	// 检查 HR 黑名单
	if f.checkHRBlacklist(job.HRName, job.HRPosition) {
		return false, "HR黑名单"
	}

	// 检查职位黑名单
	if f.checkJobBlacklist(job.JobName) {
		return false, "职位黑名单"
	}

	// 检查不活跃 HR
	if f.FilterDeadHR && f.isDeadHR(job.HRActiveStatus) {
		return false, "不活跃HR"
	}

	// 检查薪资范围
	if !f.checkSalary(job.Salary) {
		return false, "薪资不符合要求"
	}

	return true, ""
}

// checkCompanyBlacklist 检查公司黑名单
func (f *FilterConfig) checkCompanyBlacklist(companyName string) bool {
	if companyName == "" {
		return false
	}
	for _, keyword := range f.CompanyBlacklist {
		if strings.Contains(companyName, keyword) {
			config.Debug("过滤公司（黑名单）: ", companyName, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// checkHRBlacklist 检查 HR 黑名单
func (f *FilterConfig) checkHRBlacklist(hrName, hrPosition string) bool {
	if hrPosition == "" {
		return false
	}
	for _, keyword := range f.HRBlacklist {
		if strings.Contains(hrPosition, keyword) {
			config.Debug("过滤HR（黑名单）: ", hrPosition, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// checkJobBlacklist 检查职位黑名单
func (f *FilterConfig) checkJobBlacklist(jobName string) bool {
	if jobName == "" {
		return false
	}
	for _, keyword := range f.JobBlacklist {
		if strings.Contains(jobName, keyword) {
			config.Debug("过滤职位（黑名单）: ", jobName, " 匹配 ", keyword)
			return true
		}
	}
	return false
}

// isDeadHR 判断 HR 是否不活跃
func (f *FilterConfig) isDeadHR(activeStatus string) bool {
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

// checkSalary 检查薪资范围
func (f *FilterConfig) checkSalary(salary string) bool {
	if salary == "" || (f.MinSalary == "" && f.MaxSalary == "") {
		return true
	}

	// 解析薪资字符串，如 "15K-30K"
	salaryMin, salaryMax := parseSalary(salary)
	if salaryMin == 0 && salaryMax == 0 {
		return true
	}

	// 解析配置的薪资范围
	minK := parseSalaryToK(f.MinSalary)
	maxK := parseSalaryToK(f.MaxSalary)

	// 检查是否符合范围
	if minK > 0 && salaryMax < minK {
		return false
	}
	if maxK > 0 && salaryMin > maxK {
		return false
	}

	return true
}

// parseSalary 解析薪资字符串
func parseSalary(salary string) (min, max int) {
	// 示例: "15K-30K" 或 "1.5万-3万"
	salary = strings.ToUpper(salary)
	salary = strings.ReplaceAll(salary, " ", "")

	// 处理 "K" 单位
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
		return min * 1000, max * 1000
	}

	// 处理 "万" 单位
	if strings.Contains(salary, "万") {
		var minF, maxF float64
		parts := strings.Split(salary, "-")
		if len(parts) == 2 {
			minStr := strings.ReplaceAll(parts[0], "万", "")
			maxStr := strings.ReplaceAll(parts[1], "万", "")
			fmt.Sscanf(minStr, "%f", &minF)
			fmt.Sscanf(maxStr, "%f", &maxF)
		}
		return int(minF * 10000), int(maxF * 10000)
	}

	return 0, 0
}

// parseSalaryToK 将薪资字符串转换为 K 单位
func parseSalaryToK(salary string) int {
	if salary == "" {
		return 0
	}
	salary = strings.ToUpper(salary)
	salary = strings.ReplaceAll(salary, " ", "")

	var value int
	if strings.Contains(salary, "K") {
		str := strings.ReplaceAll(salary, "K", "")
		fmt.Sscanf(str, "%d", &value)
	} else if strings.Contains(salary, "万") {
		str := strings.ReplaceAll(salary, "万", "")
		var valueF float64
		fmt.Sscanf(str, "%f", &valueF)
		value = int(valueF * 10)
	}
	return value
}

// FilterJobs 批量过滤岗位
func (f *FilterConfig) FilterJobs(jobs []JobCard) ([]JobCard, []string) {
	filtered := make([]JobCard, 0)
	reasons := make([]string, 0)

	for _, job := range jobs {
		passed, reason := f.FilterJob(&job)
		if passed {
			filtered = append(filtered, job)
		} else {
			reasons = append(reasons, job.JobName+" - "+reason)
		}
	}

	config.Info("过滤前: ", len(jobs), " 个岗位, 过滤后: ", len(filtered), " 个岗位")
	return filtered, reasons
}

// AddCompanyToBlacklist 添加公司到黑名单
func (f *FilterConfig) AddCompanyToBlacklist(companyName string) error {
	blacklist := storage.Blacklist{
		Keyword:   companyName,
		Type:      "company",
		Source:    "auto",
	}
	return storage.Create(&blacklist)
}

// AddHRToBlacklist 添加 HR 到黑名单
func (f *FilterConfig) AddHRToBlacklist(hrName string) error {
	blacklist := storage.Blacklist{
		Keyword:   hrName,
		Type:      "hr",
		Source:    "auto",
	}
	return storage.Create(&blacklist)
}
