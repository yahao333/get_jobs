package storage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCookieJSONMarshal 测试 Cookie 结构体的 JSON 序列化
func TestCookieJSONMarshal(t *testing.T) {
	cookie := Cookie{
		Name:     "session_id",
		Value:    "abc123xyz",
		Domain:   ".zhipin.com",
		Path:     "/",
		Expires:  time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		HttpOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	data, err := json.Marshal(cookie)
	require.NoError(t, err, "JSON 序列化应该成功")

	// 验证序列化后的 JSON 包含必要的字段
	assert.Contains(t, string(data), "session_id")
	assert.Contains(t, string(data), "abc123xyz")
	assert.Contains(t, string(data), ".zhipin.com")
}

// TestCookieJSONUnmarshal 测试 Cookie 结构体的 JSON 反序列化
func TestCookieJSONUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "session_id",
		"value": "abc123xyz",
		"domain": ".zhipin.com",
		"path": "/",
		"expires": "2025-12-31T23:59:59Z",
		"http_only": true,
		"secure": true,
		"same_site": "Lax"
	}`

	var cookie Cookie
	err := json.Unmarshal([]byte(jsonData), &cookie)
	require.NoError(t, err, "JSON 反序列化应该成功")

	assert.Equal(t, "session_id", cookie.Name)
	assert.Equal(t, "abc123xyz", cookie.Value)
	assert.Equal(t, ".zhipin.com", cookie.Domain)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, "Lax", cookie.SameSite)
}

// TestCookieSliceJSON 测试 Cookie 切片的 JSON 序列化
func TestCookieSliceJSON(t *testing.T) {
	cookies := []Cookie{
		{
			Name:     "session_id",
			Value:    "abc123",
			Domain:   ".zhipin.com",
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: "Lax",
		},
		{
			Name:     "user_token",
			Value:    "xyz789",
			Domain:   ".zhipin.com",
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: false,
			Secure:   true,
			SameSite: "Strict",
		},
	}

	data, err := json.Marshal(cookies)
	require.NoError(t, err, "JSON 序列化应该成功")

	// 反序列化验证
	var parsedCookies []Cookie
	err = json.Unmarshal(data, &parsedCookies)
	require.NoError(t, err, "JSON 反序列化应该成功")

	assert.Equal(t, 2, len(parsedCookies))
	assert.Equal(t, "session_id", parsedCookies[0].Name)
	assert.Equal(t, "user_token", parsedCookies[1].Name)
}

// TestCookieEmptyFields 测试空字段的 Cookie
func TestCookieEmptyFields(t *testing.T) {
	cookie := Cookie{
		Name:     "",
		Value:    "",
		Domain:   "",
		Path:     "",
		Expires:  time.Time{},
		HttpOnly: false,
		Secure:   false,
		SameSite: "",
	}

	data, err := json.Marshal(cookie)
	require.NoError(t, err, "空字段的 JSON 序列化应该成功")

	var parsed Cookie
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "空字段的 JSON 反序列化应该成功")

	assert.Equal(t, "", parsed.Name)
	assert.Equal(t, "", parsed.Value)
	assert.Equal(t, "", parsed.Domain)
}

// TestCookieFields 测试 Cookie 各字段的默认值
func TestCookieFields(t *testing.T) {
	cookie := Cookie{}

	// 验证默认值为零值
	assert.Equal(t, "", cookie.Name)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, "", cookie.Domain)
	assert.Equal(t, "", cookie.Path)
	assert.Equal(t, time.Time{}, cookie.Expires)
	assert.False(t, cookie.HttpOnly)
	assert.False(t, cookie.Secure)
	assert.Equal(t, "", cookie.SameSite)
}
