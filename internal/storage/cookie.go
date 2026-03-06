// Package storage 提供数据持久化存储功能
// 包括 Cookie、黑名单、岗位数据、投递记录等
package storage

import (
	"time"
)

// Cookie 表示 HTTP Cookie 的存储结构
// 用于持久化保存浏览器登录状态，实现免登录功能
type Cookie struct {
	Name     string    `json:"name"`      // Cookie 名称
	Value    string    `json:"value"`     // Cookie 值
	Domain   string    `json:"domain"`    // 所属域名（如 .zhipin.com）
	Path     string    `json:"path"`      // 生效路径（如 /）
	Expires  time.Time `json:"expires"`   // 过期时间
	HttpOnly bool      `json:"http_only"` // 是否仅 HTTP 可访问
	Secure   bool      `json:"secure"`    // 是否仅 HTTPS 传输
	SameSite string    `json:"same_site"` // SameSite 属性（Strict/Lax/None）
}
