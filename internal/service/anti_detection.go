// Package service 提供业务逻辑服务
// 反检测模块：模拟人类行为，避免被网站检测为机器人
package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/yahao333/get_jobs/internal/config"
)

// 全局随机数生成器，避免每次创建 NewSource 的开销
var (
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu sync.Mutex
)

// AntiDetection 反检测器
// 通过添加随机延迟、模拟人类行为等方式，降低被检测为机器人的风险
type AntiDetection struct {
	minDelay time.Duration // 最小延迟
	maxDelay time.Duration // 最大延迟
	enabled  bool          // 是否启用
}

// NewAntiDetection 创建反检测器
func NewAntiDetection() *AntiDetection {
	return &AntiDetection{
		minDelay: 100 * time.Millisecond, // 最小 100ms
		maxDelay: 500 * time.Millisecond, // 最大 500ms
		enabled:  true,
	}
}

// SetDelayRange 设置延迟范围
func (a *AntiDetection) SetDelayRange(min, max time.Duration) {
	a.minDelay = min
	a.maxDelay = max
}

// RandomDelay 随机延迟
// 在 minDelay 和 maxDelay 之间随机等待一段时间
func (a *AntiDetection) RandomDelay() {
	if !a.enabled {
		return
	}
	delay := a.randomDuration(a.minDelay, a.maxDelay)
	time.Sleep(delay)
}

// HumanClick 模拟人类点击
// 包含移动鼠标和点击操作
func (a *AntiDetection) HumanClick(fn func() error) error {
	// 先随机延迟
	a.RandomDelay()

	// 执行操作
	err := fn()
	if err != nil {
		return err
	}

	// 操作后随机延迟
	a.RandomDelay()

	return nil
}

// HumanType 模拟人类输入
// 逐字符输入，带有随机延迟
func (a *AntiDetection) HumanType(text string, typeFn func(string) error) error {
	for _, char := range text {
		// 每个字符输入后随机延迟
		a.RandomDelay()

		err := typeFn(string(char))
		if err != nil {
			return err
		}
	}
	return nil
}

// HumanScroll 模拟人类滚动
// 滚动带有随机停顿
func (a *AntiDetection) HumanScroll(scrollFn func() error) error {
	// 滚动前延迟
	a.RandomDelay()

	err := scrollFn()
	if err != nil {
		return err
	}

	// 滚动后延迟
	a.RandomDelay()

	return nil
}

// RandomDuration 生成随机时长
func (a *AntiDetection) randomDuration(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}
	delta := max.Nanoseconds() - min.Nanoseconds()
	
	mu.Lock()
	defer mu.Unlock()
	return min + time.Duration(r.Int63n(delta))
}

// BetweenActions 操作间隔
// 在两个操作之间添加随机延迟，模拟人类思考时间
func (a *AntiDetection) BetweenActions() {
	// 人类思考时间通常在 1-3 秒之间
	delay := a.randomDuration(1*time.Second, 3*time.Second)
	config.Debug(fmt.Sprintf("等待 %.1f 秒...", delay.Seconds()))
	time.Sleep(delay)
}

// BeforeJobApplication 投递前延迟
// 在投递简历前添加较长的延迟，模拟人类决策过程
func (a *AntiDetection) BeforeJobApplication() {
	// 人类在投递前会浏览岗位信息，这个过程通常在 5-15 秒
	delay := a.randomDuration(5*time.Second, 15*time.Second)
	config.Debug(fmt.Sprintf("浏览岗位信息中，等待 %.1f 秒...", delay.Seconds()))
	time.Sleep(delay)
}

// AfterJobApplication 投递后延迟
// 投递完成后添加延迟，模拟人类行为
func (a *AntiDetection) AfterJobApplication() {
	// 投递后人类会查看结果或继续浏览
	delay := a.randomDuration(2*time.Second, 5*time.Second)
	time.Sleep(delay)
}

// DailyLimitReached 达到每日上限
// 模拟人类在达到目标后的行为
func (a *AntiDetection) DailyLimitReached() {
	config.Info("今日投递已达上限，停止投递")
	// 人类不会在一天内投递太多简历
	// 这里可以记录状态，等待明天
}

// WeekLimitReached 达到每周上限
func (a *AntiDetection) WeekLimitReached() {
	config.Warn("本周投递已达上限，建议休息")
}

// HumanBehaviorPatterns 人类行为模式
// 定义一些典型的人类行为模式
var HumanBehaviorPatterns = map[string]time.Duration{
	"快速浏览": 500 * time.Millisecond,
	"仔细阅读": 3 * time.Second,
	"填写表单": 5 * time.Second,
	"发送消息": 2 * time.Second,
	"思考回复": 10 * time.Second,
	"切换页面": 1 * time.Second,
}

// ApplyPattern 应用行为模式
// 根据指定的行为模式添加延迟
func (a *AntiDetection) ApplyPattern(patternName string) {
	if duration, ok := HumanBehaviorPatterns[patternName]; ok {
		// 添加一定的随机性
		variance := float64(duration) * 0.3
		
		mu.Lock()
		val := r.Float64()
		mu.Unlock()
		
		actualDuration := duration + time.Duration(val*variance*2-variance)
		config.Debug(fmt.Sprintf("执行行为模式 [%s]，等待 %.1f 秒", patternName, actualDuration.Seconds()))
		time.Sleep(actualDuration)
	}
}
