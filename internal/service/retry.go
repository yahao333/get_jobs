// Package service 提供业务逻辑服务
// 错误处理与重试机制模块
package service

import (
	"fmt"
	"time"

	"github.com/yahao333/get_jobs/internal/config"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	MaxDelay      time.Duration // 最大延迟
	Multiplier    float64       // 延迟倍数
	RetryableFunc func(error) bool // 判断错误是否可重试
}

// DefaultRetryConfig 获取默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		RetryableFunc: func(err error) bool {
			// 默认：网络超时、临时不可用等错误可重试
			if err == nil {
				return false
			}
			errStr := err.Error()
			// 网络错误
			if contains(errStr, "timeout", "connection", "network", "temporary") {
				return true
			}
			return false
		},
	}
}

// contains 检查字符串是否包含任意一个子串
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// RetryWithConfig 带配置的重试
func RetryWithConfig(cfg *RetryConfig, fn func() error) error {
	var lastErr error
	delay := cfg.InitialDelay

	for i := 0; i <= cfg.MaxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// 检查是否可重试
		if cfg.RetryableFunc != nil && !cfg.RetryableFunc(err) {
			return err
		}

		// 最后一次不重试
		if i == cfg.MaxRetries {
			break
		}

		config.Warn(fmt.Sprintf("操作失败，%v 后重试 (%d/%d)", delay, i+1, cfg.MaxRetries))
		time.Sleep(delay)

		// 延迟倍增
		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return fmt.Errorf("重试 %d 次后仍然失败: %w", cfg.MaxRetries, lastErr)
}

// Retry 便捷的重试函数
func Retry(fn func() error, maxRetries int) error {
	return RetryWithConfig(&RetryConfig{
		MaxRetries:  maxRetries,
		InitialDelay: 1 * time.Second,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
	}, fn)
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	onError    func(error)
	onRetry   func(int, error)
	onSuccess func()
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		onError: func(err error) {
			config.Error("操作失败: ", err)
		},
		onRetry: func(count int, err error) {
			config.Warn(fmt.Sprintf("重试 %d: %v", count, err))
		},
		onSuccess: func() {
			config.Debug("操作成功")
		},
	}
}

// SetOnError 设置错误回调
func (h *ErrorHandler) SetOnError(fn func(error)) *ErrorHandler {
	h.onError = fn
	return h
}

// SetOnRetry 设置重试回调
func (h *ErrorHandler) SetOnRetry(fn func(int, error)) *ErrorHandler {
	h.onRetry = fn
	return h
}

// SetOnSuccess 设置成功回调
func (h *ErrorHandler) SetOnSuccess(fn func()) *ErrorHandler {
	h.onSuccess = fn
	return h
}

// Execute 执行并处理错误
func (h *ErrorHandler) Execute(fn func() error) error {
	err := RetryWithConfig(DefaultRetryConfig(), fn)
	if err != nil {
		if h.onError != nil {
			h.onError(err)
		}
		return err
	}
	if h.onSuccess != nil {
		h.onSuccess()
	}
	return nil
}

// ExecuteWithRetry 执行并处理错误（带自定义重试次数）
func (h *ErrorHandler) ExecuteWithRetry(fn func() error, maxRetries int) error {
	config := DefaultRetryConfig()
	config.MaxRetries = maxRetries

	err := RetryWithConfig(config, fn)
	if err != nil {
		if h.onError != nil {
			h.onError(err)
		}
		return err
	}
	if h.onSuccess != nil {
		h.onSuccess()
	}
	return nil
}
