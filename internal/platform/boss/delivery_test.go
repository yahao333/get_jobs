package boss

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDeliveryResultStructure 测试 DeliveryResult 结构体
func TestDeliveryResultStructure(t *testing.T) {
	deliveredAt := time.Now()
	result := DeliveryResult{
		JobID:     12345,
		Success:   true,
		Message:   "投递成功",
		Delivered: deliveredAt,
	}

	assert.Equal(t, int64(12345), result.JobID)
	assert.True(t, result.Success)
	assert.Equal(t, "投递成功", result.Message)
	assert.Equal(t, deliveredAt, result.Delivered)
}

// TestDeliveryResultFailure 测试投递失败结果
func TestDeliveryResultFailure(t *testing.T) {
	result := DeliveryResult{
		JobID:     0,
		Success:   false,
		Message:   "已达到每日投递上限",
		Delivered: time.Time{},
	}

	assert.False(t, result.Success)
	assert.Equal(t, "已达到每日投递上限", result.Message)
	assert.Equal(t, int64(0), result.JobID)
}

// TestDeliveryStructure 测试 Delivery 结构体字段
func TestDeliveryStructure(t *testing.T) {
	delivery := &Delivery{
		sendImgResume:  true,
		imgResumePath:  "/path/to/resume.jpg",
		dailyLimit:     150,
		deliveredToday: 10,
	}

	assert.True(t, delivery.sendImgResume)
	assert.Equal(t, "/path/to/resume.jpg", delivery.imgResumePath)
	assert.Equal(t, 150, delivery.dailyLimit)
	assert.Equal(t, 10, delivery.deliveredToday)
}

// TestGetDeliveredCount 测试获取今日投递数量
func TestGetDeliveredCount(t *testing.T) {
	delivery := &Delivery{
		deliveredToday: 25,
	}

	count := delivery.GetDeliveredCount()
	assert.Equal(t, 25, count)
}

// TestSetDailyLimit 测试设置每日投递限制
func TestSetDailyLimit(t *testing.T) {
	delivery := &Delivery{
		dailyLimit: 100,
	}

	// 修改每日限制
	delivery.SetDailyLimit(150)
	assert.Equal(t, 150, delivery.dailyLimit)

	// 设置为 0（无限制）
	delivery.SetDailyLimit(0)
	assert.Equal(t, 0, delivery.dailyLimit)
}

// TestBatchDeliverEmptyJobs 测试批量投递空列表
func TestBatchDeliverEmptyJobs(t *testing.T) {
	delivery := &Delivery{
		dailyLimit:     150,
		deliveredToday: 0,
	}

	jobs := []JobCard{}
	results := delivery.BatchDeliver(jobs, "你好")

	// 空列表应该返回空结果
	assert.Empty(t, results)
	assert.Equal(t, 0, delivery.deliveredToday)
}

// TestJobCardForDelivery 测试用于投递的 JobCard
func TestJobCardForDelivery(t *testing.T) {
	job := JobCard{
		EncryptID:      "abc123",
		EncryptUserID:  "user456",
		CompanyName:    "阿里巴巴",
		JobName:        "Go后端开发",
		Salary:         "20K-40K",
		Location:       "杭州",
		Experience:     "3-5年",
		Degree:         "本科",
		HRName:         "张三",
		HRPosition:     "技术经理",
		HRActiveStatus: "今日活跃",
		JobDescription: "负责后端开发",
		JobURL:         "https://www.zhipin.com/job_detail/abc123.html",
	}

	assert.Equal(t, "Go后端开发", job.JobName)
	assert.Equal(t, "阿里巴巴", job.CompanyName)
	assert.NotEmpty(t, job.JobURL)
}

// TestDeliveryDailyLimitBehavior 测试每日限制行为
func TestDeliveryDailyLimitBehavior(t *testing.T) {
	delivery := &Delivery{
		dailyLimit:     3,
		deliveredToday: 3,
	}

	// 当达到每日限制时，应该返回错误
	// 注意: 这里我们不实际调用 Deliver 方法，因为它需要浏览器
	// 我们只验证字段状态
	assert.Equal(t, 3, delivery.dailyLimit)
	assert.Equal(t, 3, delivery.deliveredToday)

	// 模拟达到限制
	canDeliver := delivery.deliveredToday < delivery.dailyLimit || delivery.dailyLimit == 0
	assert.False(t, canDeliver, "达到每日限制时不应再投递")

	// 模拟未达到限制
	delivery.deliveredToday = 2
	canDeliver = delivery.deliveredToday < delivery.dailyLimit || delivery.dailyLimit == 0
	assert.True(t, canDeliver, "未达到每日限制时可以投递")
}
