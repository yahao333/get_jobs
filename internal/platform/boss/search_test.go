package boss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJobCardStructure 测试 JobCard 结构体
func TestJobCardStructure(t *testing.T) {
	job := JobCard{
		EncryptID:      "enc123",
		EncryptUserID:  "encuser456",
		CompanyName:    "字节跳动",
		JobName:        "后端开发工程师",
		Salary:         "25K-50K",
		Location:       "北京",
		Experience:     "3-5年",
		Degree:         "本科及以上",
		HRName:         "李四",
		HRPosition:     "HRBP",
		HRActiveStatus: "今日活跃",
		JobDescription: "负责业务后端开发",
		JobURL:         "https://www.zhipin.com/job_detail/abc.html",
	}

	assert.Equal(t, "enc123", job.EncryptID)
	assert.Equal(t, "encuser456", job.EncryptUserID)
	assert.Equal(t, "字节跳动", job.CompanyName)
	assert.Equal(t, "后端开发工程师", job.JobName)
	assert.Equal(t, "25K-50K", job.Salary)
	assert.Equal(t, "北京", job.Location)
	assert.Equal(t, "3-5年", job.Experience)
	assert.Equal(t, "本科及以上", job.Degree)
	assert.Equal(t, "李四", job.HRName)
	assert.Equal(t, "HRBP", job.HRPosition)
	assert.Equal(t, "今日活跃", job.HRActiveStatus)
	assert.Equal(t, "负责业务后端开发", job.JobDescription)
	assert.Equal(t, "https://www.zhipin.com/job_detail/abc.html", job.JobURL)
}

// TestJobCardEmptyFields 测试空字段的 JobCard
func TestJobCardEmptyFields(t *testing.T) {
	job := JobCard{}

	assert.Equal(t, "", job.EncryptID)
	assert.Equal(t, "", job.EncryptUserID)
	assert.Equal(t, "", job.CompanyName)
	assert.Equal(t, "", job.JobName)
	assert.Equal(t, "", job.Salary)
	assert.Equal(t, "", job.Location)
	assert.Equal(t, "", job.Experience)
	assert.Equal(t, "", job.Degree)
	assert.Equal(t, "", job.HRName)
	assert.Equal(t, "", job.HRPosition)
	assert.Equal(t, "", job.HRActiveStatus)
	assert.Equal(t, "", job.JobDescription)
	assert.Equal(t, "", job.JobURL)
}

// TestJobCardURLFormat 测试 JobURL 格式
func TestJobCardURLFormat(t *testing.T) {
	tests := []struct {
		name     string
		jobURL   string
		expected string
	}{
		{
			name:     "完整 URL",
			jobURL:   "https://www.zhipin.com/job_detail/abc123.html",
			expected: "https://www.zhipin.com/job_detail/abc123.html",
		},
		{
			name:     "空 URL",
			jobURL:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := JobCard{
				JobURL: tt.jobURL,
			}
			assert.Equal(t, tt.expected, job.JobURL)
		})
	}
}

// TestJobCardSalaryFormats 测试不同薪资格式
func TestJobCardSalaryFormats(t *testing.T) {
	jobs := []JobCard{
		{Salary: "15K-30K"},
		{Salary: "1.5万-3万"},
		{Salary: "8K-15K"},
		{Salary: "面议"},
		{Salary: ""},
	}

	assert.Equal(t, "15K-30K", jobs[0].Salary)
	assert.Equal(t, "1.5万-3万", jobs[1].Salary)
	assert.Equal(t, "8K-15K", jobs[2].Salary)
	assert.Equal(t, "面议", jobs[3].Salary)
	assert.Equal(t, "", jobs[4].Salary)
}

// TestJobCardHRStatus 测试 HR 活跃状态
func TestJobCardHRStatus(t *testing.T) {
	statuses := []struct {
		status   string
		expected string
	}{
		{"今日活跃", "今日活跃"},
		{"3日内活跃", "3日内活跃"},
		{"1周内活跃", "1周内活跃"},
		{"1年前在线", "1年前在线"},
		{"3年前在线", "3年前在线"},
		{"", ""},
	}

	for _, s := range statuses {
		job := JobCard{
			HRActiveStatus: s.status,
		}
		assert.Equal(t, s.expected, job.HRActiveStatus)
	}
}

// TestMultipleJobCards 测试多个岗位卡片
func TestMultipleJobCards(t *testing.T) {
	jobs := []JobCard{
		{
			CompanyName: "阿里巴巴",
			JobName:    "Go后端开发",
			Salary:     "20K-40K",
			Location:   "杭州",
		},
		{
			CompanyName: "腾讯",
			JobName:    "后端开发",
			Salary:     "25K-50K",
			Location:   "深圳",
		},
		{
			CompanyName: "字节跳动",
			JobName:    "服务端开发",
			Salary:     "30K-60K",
			Location:   "北京",
		},
	}

	assert.Equal(t, 3, len(jobs))

	// 验证每个岗位
	assert.Equal(t, "阿里巴巴", jobs[0].CompanyName)
	assert.Equal(t, "腾讯", jobs[1].CompanyName)
	assert.Equal(t, "字节跳动", jobs[2].CompanyName)

	// 验证薪资排序
	assert.True(t, jobs[0].Salary < jobs[1].Salary)
	assert.True(t, jobs[1].Salary < jobs[2].Salary)
}
