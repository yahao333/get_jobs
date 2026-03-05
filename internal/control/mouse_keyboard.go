package control

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"

	"github.com/yahao333/get_jobs/internal/config"
)

// Mouse 鼠标控制
type Mouse struct {
	// 基础延迟（毫秒）
	baseDelay time.Duration
}

// NewMouse 创建鼠标控制器
func NewMouse() *Mouse {
	return &Mouse{
		baseDelay: 50,
	}
}

// Position 获取当前鼠标位置
func (m *Mouse) Position() (x, y int) {
	return robotgo.GetMousePos()
}

// Move 移动鼠标到指定位置
func (m *Mouse) Move(x, y int) error {
	robotgo.MoveMouse(x, y)
	m.randomDelay()
	return nil
}

// MoveTo 移动鼠标到指定位置（带平滑移动）
func (m *Mouse) MoveTo(x, y int) error {
	// 使用 smooth move 进行更自然的移动
	robotgo.MoveSmooth(x, y)
	m.randomDelay()
	return nil
}

// Click 鼠标点击
func (m *Mouse) Click(x, y int, button ...string) error {
	// 先移动到目标位置
	if err := m.MoveTo(x, y); err != nil {
		return err
	}

	// 默认左键
	btn := "left"
	if len(button) > 0 {
		btn = button[0]
	}

	robotgo.Click(btn)
	m.randomDelay()
	return nil
}

// DoubleClick 双击
func (m *Mouse) DoubleClick(x, y int, button ...string) error {
	if err := m.MoveTo(x, y); err != nil {
		return err
	}

	btn := "left"
	if len(button) > 0 {
		btn = button[0]
	}

	robotgo.Click(btn, true)
	m.randomDelay()
	return nil
}

// RightClick 右键点击
func (m *Mouse) RightClick(x, y int) error {
	return m.Click(x, y, "right")
}

// Scroll 滚动
func (m *Mouse) Scroll(direction string, amount int) error {
	// direction: "up", "down", "left", "right"
	x, y := 0, 0
	switch direction {
	case "up":
		y = -amount
	case "down":
		y = amount
	case "left":
		x = -amount
	case "right":
		x = amount
	}
	robotgo.Scroll(x, y)
	m.randomDelay()
	return nil
}

// ScrollTo 滚动到指定位置
func (m *Mouse) ScrollTo(x, y int, amount int) error {
	// 先移动到滚动位置
	robotgo.MoveMouse(x, y)
	time.Sleep(m.baseDelay)
	return m.Scroll("down", amount)
}

// randomDelay 添加随机延迟，模拟人类操作
func (m *Mouse) randomDelay() {
	delay := m.baseDelay + time.Duration(randomInt(0, 100))*time.Millisecond
	time.Sleep(delay)
}

// randomInt 生成范围内的随机整数
func randomInt(min, max int) int {
	if max <= min {
		return min
	}
	// 使用时间戳作为简单随机数
	return min + int(time.Now().UnixNano()%int64(max-min))
}

// Keyboard 键盘控制
type Keyboard struct {
	baseDelay time.Duration
}

// NewKeyboard 创建键盘控制器
func NewKeyboard() *Keyboard {
	return &Keyboard{
		baseDelay: 30,
	}
}

// Type 输入文本
func (k *Keyboard) Type(text string) error {
	robotgo.TypeStr(text)
	k.randomDelay()
	return nil
}

// TypeOnce 输入文本（单次）
func (k *Keyboard) TypeOnce(text string) error {
	robotgo.TypeStr(text)
	k.randomDelay()
	return nil
}

// Press 按键
func (k *Keyboard) Press(keys ...string) error {
	for _, key := range keys {
		robotgo.KeyTap(key)
		k.randomDelay()
	}
	return nil
}

// Hold 按住按键
func (k *Keyboard) Hold(keys ...string) error {
	if len(keys) == 0 {
		return fmt.Errorf("至少需要传入一个按键")
	}
	robotgo.KeyToggle(keys[0], "down")
	k.randomDelay()
	return nil
}

// Release 释放按键
func (k *Keyboard) Release(keys ...string) error {
	if len(keys) == 0 {
		return fmt.Errorf("至少需要传入一个按键")
	}
	robotgo.KeyToggle(keys[0], "up")
	k.randomDelay()
	return nil
}

// KeyDown 按下按键
func (k *Keyboard) KeyDown(key string) error {
	return k.Hold(key)
}

// KeyUp 释放按键
func (k *Keyboard) KeyUp(key string) error {
	return k.Release(key)
}

// randomDelay 添加随机延迟
func (k *Keyboard) randomDelay() {
	delay := k.baseDelay + time.Duration(randomInt(0, 50))*time.Millisecond
	time.Sleep(delay)
}

// Control 组合控制
type Control struct {
	Mouse    *Mouse
	Keyboard *Keyboard
}

// NewControl 创建控制器
func NewControl() *Control {
	return &Control{
		Mouse:    NewMouse(),
		Keyboard: NewKeyboard(),
	}
}

// MoveAndClick 移动并点击
func (c *Control) MoveAndClick(x, y int) error {
	config.Debug(fmt.Sprintf("点击坐标: (%d, %d)", x, y))
	return c.Mouse.Click(x, y)
}

// MoveAndDoubleClick 移动并双击
func (c *Control) MoveAndDoubleClick(x, y int) error {
	config.Debug(fmt.Sprintf("双击坐标: (%d, %d)", x, y))
	return c.Mouse.DoubleClick(x, y)
}

// MoveAndType 移动到输入框并输入文本
func (c *Control) MoveAndType(x, y int, text string) error {
	if err := c.Mouse.Click(x, y); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return c.Keyboard.Type(text)
}

// SelectAll 全选
func (c *Control) SelectAll() error {
	// Ctrl+A (Mac: Command+A)
	return c.Keyboard.Press("command", "a")
}

// Copy 复制
func (c *Control) Copy() error {
	return c.Keyboard.Press("command", "c")
}

// Paste 粘贴
func (c *Control) Paste() error {
	return c.Keyboard.Press("command", "v")
}

// Cut 剪切
func (c *Control) Cut() error {
	return c.Keyboard.Press("command", "x")
}

// SelectAllWin Windows 全选
func (c *Control) SelectAllWin() error {
	return c.Keyboard.Press("ctrl", "a")
}

// CopyWin Windows 复制
func (c *Control) CopyWin() error {
	return c.Keyboard.Press("ctrl", "c")
}

// PasteWin Windows 粘贴
func (c *Control) PasteWin() error {
	return c.Keyboard.Press("ctrl", "v")
}

// ScrollDown 向下滚动
func (c *Control) ScrollDown(amount int) error {
	return c.Mouse.Scroll("down", amount)
}

// ScrollUp 向上滚动
func (c *Control) ScrollUp(amount int) error {
	return c.Mouse.Scroll("up", amount)
}

// Wait 等待
func (c *Control) Wait(duration time.Duration) {
	time.Sleep(duration)
}

// RandomWait 随机等待
func (c *Control) RandomWait(min, max time.Duration) {
	wait := min + time.Duration(randomInt(int(min.Milliseconds()), int(max.Milliseconds())))*time.Millisecond
	time.Sleep(wait)
}
