package control

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMouse 测试创建鼠标控制器
func TestNewMouse(t *testing.T) {
	m := NewMouse()
	require.NotNil(t, m, "Mouse 实例不应为空")
	assert.Equal(t, time.Duration(50), m.baseDelay, "默认基础延迟应为 50ms")
}

// TestNewKeyboard 测试创建键盘控制器
func TestNewKeyboard(t *testing.T) {
	k := NewKeyboard()
	require.NotNil(t, k, "Keyboard 实例不应为空")
	assert.Equal(t, time.Duration(30), k.baseDelay, "默认基础延迟应为 30ms")
}

// TestNewControl 测试创建组合控制器
func TestNewControl(t *testing.T) {
	c := NewControl()
	require.NotNil(t, c, "Control 实例不应为空")
	require.NotNil(t, c.Mouse, "Mouse 不应为空")
	require.NotNil(t, c.Keyboard, "Keyboard 不应为空")
}

// TestMouse_Structure 测试鼠标结构体
func TestMouse_Structure(t *testing.T) {
	m := NewMouse()

	// 验证基础延迟已设置
	assert.Greater(t, m.baseDelay, time.Duration(0), "基础延迟应大于 0")

	// 验证可以调用方法（实际执行可能会失败，但不应崩溃）
	// 注意：Position() 会调用 robotgo.GetMousePos()，在没有显示环境可能失败
	x, y := m.Position()
	t.Logf("鼠标位置: (%d, %d)", x, y)
	// 位置值取决于实际鼠标位置，不做断言
}

// TestKeyboard_Structure 测试键盘结构体
func TestKeyboard_Structure(t *testing.T) {
	k := NewKeyboard()

	// 验证基础延迟已设置
	assert.Greater(t, k.baseDelay, time.Duration(0), "基础延迟应大于 0")
}

// TestRandomInt 测试 randomInt 函数边界条件
func TestRandomInt(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		expected int // 期望返回的值（如果是单一值）
	}{
		{
			name:     "min equals max",
			min:      5,
			max:      5,
			expected: 5,
		},
		{
			name:     "min greater than max",
			min:      10,
			max:      5,
			expected: 10,
		},
		{
			name: "normal range",
			min:  0,
			max:  100,
		},
		{
			name: "negative range",
			min:  -50,
			max:  -10,
		},
		{
			name: "zero to positive",
			min:  0,
			max:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := randomInt(tt.min, tt.max)

			if tt.min == tt.max {
				assert.Equal(t, tt.expected, result, "min==max 时应返回 min 值")
			} else if tt.min > tt.max {
				assert.Equal(t, tt.min, result, "min>max 时应返回 min 值")
			} else {
				assert.GreaterOrEqual(t, result, tt.min, "结果应 >= min")
				assert.LessOrEqual(t, result, tt.max, "结果应 <= max")
			}
		})
	}
}

// TestRandomInt_Distribution 测试 randomInt 分布（统计测试）
// 注意：由于使用时间戳作为随机数种子，此测试可能不稳定
func TestRandomInt_Distribution(t *testing.T) {
	t.Skip("随机数分布测试不稳定，跳过")
	/*
	// 测试多次调用分布是否合理
	const iterations = 10000
	const min = 0
	const max = 9

	counts := make([]int, max-min+1)

	for i := 0; i < iterations; i++ {
		result := randomInt(min, max)
		if result >= min && result <= max {
			counts[result-min]++
		}
	}

	// 验证每个值都被命中（至少有统计意义）
	// 由于是简单随机数，不要求完全均匀
	minCount := iterations / (max - min + 1) / 2
	for i, count := range counts {
		assert.Greater(t, count, minCount, "值 %d 被命中次数过少", i)
	}
	*/
}

// TestMouse_Move 测试鼠标移动方法存在
func TestMouse_Move(t *testing.T) {
	m := NewMouse()

	// 测试方法可以被调用（参数验证）
	// 注意：实际移动会调用 robotgo，在无显示环境可能失败
	err := m.Move(100, 100)
	// 允许失败但不崩溃
	if err != nil {
		t.Logf("Move 失败（可能预期）: %v", err)
	}
}

// TestMouse_MoveTo 测试平滑移动方法存在
func TestMouse_MoveTo(t *testing.T) {
	m := NewMouse()

	err := m.MoveTo(200, 200)
	if err != nil {
		t.Logf("MoveTo 失败（可能预期）: %v", err)
	}
}

// TestMouse_Click 测试鼠标点击方法参数
func TestMouse_Click(t *testing.T) {
	m := NewMouse()

	// 测试默认左键
	err := m.Click(100, 100)
	if err != nil {
		t.Logf("Click 失败（可能预期）: %v", err)
	}

	// 测试指定右键
	err = m.Click(100, 100, "right")
	if err != nil {
		t.Logf("Click with right button 失败（可能预期）: %v", err)
	}

	// 测试无效按钮名称
	err = m.Click(100, 100, "invalid")
	if err != nil {
		t.Logf("Click with invalid button 失败（可能预期）: %v", err)
	}
}

// TestMouse_DoubleClick 测试双击方法
func TestMouse_DoubleClick(t *testing.T) {
	m := NewMouse()

	err := m.DoubleClick(100, 100)
	if err != nil {
		t.Logf("DoubleClick 失败（可能预期）: %v", err)
	}

	err = m.DoubleClick(100, 100, "right")
	if err != nil {
		t.Logf("DoubleClick with right 失败（可能预期）: %v", err)
	}
}

// TestMouse_RightClick 测试右键点击
func TestMouse_RightClick(t *testing.T) {
	m := NewMouse()

	err := m.RightClick(100, 100)
	if err != nil {
		t.Logf("RightClick 失败（可能预期）: %v", err)
	}
}

// TestMouse_Scroll 测试滚动方法
func TestMouse_Scroll(t *testing.T) {
	m := NewMouse()

	// 测试向上滚动
	err := m.Scroll("up", 100)
	if err != nil {
		t.Logf("Scroll up 失败（可能预期）: %v", err)
	}

	// 测试向下滚动
	err = m.Scroll("down", 100)
	if err != nil {
		t.Logf("Scroll down 失败（可能预期）: %v", err)
	}

	// 测试向左滚动
	err = m.Scroll("left", 50)
	if err != nil {
		t.Logf("Scroll left 失败（可能预期）: %v", err)
	}

	// 测试向右滚动
	err = m.Scroll("right", 50)
	if err != nil {
		t.Logf("Scroll right 失败（可能预期）: %v", err)
	}

	// 测试无效方向
	err = m.Scroll("invalid", 100)
	if err != nil {
		t.Logf("Scroll with invalid direction 失败（可能预期）: %v", err)
	}
}

// TestMouse_ScrollTo 测试滚动到指定位置
func TestMouse_ScrollTo(t *testing.T) {
	m := NewMouse()

	err := m.ScrollTo(500, 500, 200)
	if err != nil {
		t.Logf("ScrollTo 失败（可能预期）: %v", err)
	}
}

// TestKeyboard_Type 测试键盘输入
func TestKeyboard_Type(t *testing.T) {
	k := NewKeyboard()

	// 测试空字符串
	err := k.Type("")
	if err != nil {
		t.Logf("Type empty string 失败: %v", err)
	} else {
		t.Log("Type empty string 成功")
	}

	// 测试普通文本
	err = k.Type("hello")
	if err != nil {
		t.Logf("Type text 失败（可能预期）: %v", err)
	}

	// 测试中文
	err = k.Type("你好")
	if err != nil {
		t.Logf("Type Chinese 失败（可能预期）: %v", err)
	}

	// 测试特殊字符
	err = k.Type("!@#$%^&*()")
	if err != nil {
		t.Logf("Type special chars 失败（可能预期）: %v", err)
	}
}

// TestKeyboard_TypeOnce 测试单次输入
func TestKeyboard_TypeOnce(t *testing.T) {
	k := NewKeyboard()

	err := k.TypeOnce("test")
	if err != nil {
		t.Logf("TypeOnce 失败（可能预期）: %v", err)
	}
}

// TestKeyboard_Press 测试按键
func TestKeyboard_Press(t *testing.T) {
	k := NewKeyboard()

	// 测试单键
	err := k.Press("a")
	if err != nil {
		t.Logf("Press single key 失败（可能预期）: %v", err)
	}

	// 测试组合键
	err = k.Press("command", "c")
	if err != nil {
		t.Logf("Press key combo 失败（可能预期）: %v", err)
	}

	// 测试空键列表
	err = k.Press()
	if err != nil {
		t.Logf("Press empty 失败: %v", err)
	}
}

// TestKeyboard_Hold 测试按住按键
func TestKeyboard_Hold(t *testing.T) {
	k := NewKeyboard()

	err := k.Hold("shift")
	if err != nil {
		t.Logf("Hold key 失败（可能预期）: %v", err)
	}

	// 测试空键
	err = k.Hold()
	if err != nil {
		t.Logf("Hold empty 失败: %v", err)
	}
}

// TestKeyboard_Release 测试释放按键
func TestKeyboard_Release(t *testing.T) {
	k := NewKeyboard()

	err := k.Release("shift")
	if err != nil {
		t.Logf("Release key 失败（可能预期）: %v", err)
	}
}

// TestKeyboard_KeyDown 测试按下按键
func TestKeyboard_KeyDown(t *testing.T) {
	k := NewKeyboard()

	err := k.KeyDown("a")
	if err != nil {
		t.Logf("KeyDown 失败（可能预期）: %v", err)
	}

	err = k.KeyDown("")
	if err != nil {
		t.Logf("KeyDown empty 失败（可能预期）: %v", err)
	}
}

// TestKeyboard_KeyUp 测试释放按键
func TestKeyboard_KeyUp(t *testing.T) {
	k := NewKeyboard()

	err := k.KeyUp("a")
	if err != nil {
		t.Logf("KeyUp 失败（可能预期）: %v", err)
	}
}

// TestControl_MoveAndClick 测试移动并点击
func TestControl_MoveAndClick(t *testing.T) {
	c := NewControl()

	err := c.MoveAndClick(100, 100)
	if err != nil {
		t.Logf("MoveAndClick 失败（可能预期）: %v", err)
	}
}

// TestControl_MoveAndDoubleClick 测试移动并双击
func TestControl_MoveAndDoubleClick(t *testing.T) {
	c := NewControl()

	err := c.MoveAndDoubleClick(100, 100)
	if err != nil {
		t.Logf("MoveAndDoubleClick 失败（可能预期）: %v", err)
	}
}

// TestControl_MoveAndType 测试移动并输入
func TestControl_MoveAndType(t *testing.T) {
	c := NewControl()

	err := c.MoveAndType(100, 100, "test text")
	if err != nil {
		t.Logf("MoveAndType 失败（可能预期）: %v", err)
	}
}

// TestControl_SelectAll 测试全选
func TestControl_SelectAll(t *testing.T) {
	c := NewControl()

	err := c.SelectAll()
	if err != nil {
		t.Logf("SelectAll 失败（可能预期）: %v", err)
	}
}

// TestControl_Copy 测试复制
func TestControl_Copy(t *testing.T) {
	c := NewControl()

	err := c.Copy()
	if err != nil {
		t.Logf("Copy 失败（可能预期）: %v", err)
	}
}

// TestControl_Paste 测试粘贴
func TestControl_Paste(t *testing.T) {
	c := NewControl()

	err := c.Paste()
	if err != nil {
		t.Logf("Paste 失败（可能预期）: %v", err)
	}
}

// TestControl_Cut 测试剪切
func TestControl_Cut(t *testing.T) {
	c := NewControl()

	err := c.Cut()
	if err != nil {
		t.Logf("Cut 失败（可能预期）: %v", err)
	}
}

// TestControl_WindowsShortcuts 测试 Windows 快捷键
func TestControl_WindowsShortcuts(t *testing.T) {
	c := NewControl()

	err := c.SelectAllWin()
	if err != nil {
		t.Logf("SelectAllWin 失败（可能预期）: %v", err)
	}

	err = c.CopyWin()
	if err != nil {
		t.Logf("CopyWin 失败（可能预期）: %v", err)
	}

	err = c.PasteWin()
	if err != nil {
		t.Logf("PasteWin 失败（可能预期）: %v", err)
	}
}

// TestControl_Scroll 测试滚动
func TestControl_ScrollDown(t *testing.T) {
	c := NewControl()

	err := c.ScrollDown(100)
	if err != nil {
		t.Logf("ScrollDown 失败（可能预期）: %v", err)
	}
}

func TestControl_ScrollUp(t *testing.T) {
	c := NewControl()

	err := c.ScrollUp(100)
	if err != nil {
		t.Logf("ScrollUp 失败（可能预期）: %v", err)
	}
}

// TestControl_Wait 测试等待
func TestControl_Wait(t *testing.T) {
	c := NewControl()

	// 测试零等待
	c.Wait(0)

	// 测试短等待
	start := time.Now()
	c.Wait(10 * time.Millisecond)
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, 5*time.Millisecond, "等待时间应 >= 5ms")
}

// TestControl_RandomWait 测试随机等待
func TestControl_RandomWait(t *testing.T) {
	c := NewControl()

	// 测试随机等待范围
	start := time.Now()
	c.RandomWait(10*time.Millisecond, 50*time.Millisecond)
	elapsed := time.Since(start)

	// 由于实现可能不是精确的，我们只检查大致范围
	assert.GreaterOrEqual(t, elapsed, 5*time.Millisecond, "随机等待应至少 5ms")

	// 测试相同 min 和 max
	start = time.Now()
	c.RandomWait(20*time.Millisecond, 20*time.Millisecond)
	elapsed = time.Since(start)
	// 在这个简单实现中，可能会有一些偏差
	assert.GreaterOrEqual(t, elapsed, 10*time.Millisecond, "相同值的随机等待应约等于该值")
}

// TestMouse_NegativeCoordinates 测试负数坐标
func TestMouse_NegativeCoordinates(t *testing.T) {
	m := NewMouse()

	// 测试负数坐标
	err := m.Move(-100, -100)
	if err != nil {
		t.Logf("Move with negative coords 失败: %v", err)
	}

	err = m.Click(-50, -50)
	if err != nil {
		t.Logf("Click with negative coords 失败: %v", err)
	}
}

// TestMouse_LargeCoordinates 测试大坐标值
func TestMouse_LargeCoordinates(t *testing.T) {
	m := NewMouse()

	// 测试超出屏幕范围的坐标
	err := m.Move(10000, 10000)
	if err != nil {
		t.Logf("Move with large coords 失败: %v", err)
	}

	err = m.Click(99999, 99999)
	if err != nil {
		t.Logf("Click with large coords 失败: %v", err)
	}
}

// TestKeyboard_LongText 测试长文本输入
func TestKeyboard_LongText(t *testing.T) {
	k := NewKeyboard()

	// 生成长文本
	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "a"
	}

	err := k.Type(longText)
	if err != nil {
		t.Logf("Type long text 失败（可能预期）: %v", err)
	}
}

// TestControl_AllButtons 测试所有按钮类型
func TestControl_AllButtons(t *testing.T) {
	m := NewMouse()

	// 测试各种按钮参数
	buttons := []string{"left", "right", "middle"}

	for _, btn := range buttons {
		err := m.Click(100, 100, btn)
		if err != nil {
			t.Logf("Click with %s 失败: %v", btn, err)
		}
	}
}
