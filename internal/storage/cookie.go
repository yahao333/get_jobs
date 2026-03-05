package storage

import (
	"time"
)

// Cookie Cookie 存储结构
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	HttpOnly bool      `json:"http_only"`
	Secure   bool      `json:"secure"`
	SameSite string    `json:"same_site"`
}
