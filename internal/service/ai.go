package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/loks666/get_jobs/internal/config"
)

// AIProvider AI 服务提供商接口
type AIProvider interface {
	AnalyzeImage(imageData []byte, prompt string) (string, error)
	GenerateText(prompt string) (string, error)
}

// QwenVL 阿里 Qwen VL 服务
type QwenVL struct {
	apiKey  string
	model   string
	endpoint string
	client  *http.Client
}

// NewQwenVL 创建 Qwen VL 服务
func NewQwenVL(apiKey string, model string) *QwenVL {
	return &QwenVL{
		apiKey:  apiKey,
		model:   model,
		endpoint: "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// AnalyzeImage 分析图片
func (q *QwenVL) AnalyzeImage(imageData []byte, prompt string) (string, error) {
	// 将图片转换为 base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求
	reqBody := map]interface{}{
		"model": q.model,
		"input": map]interface{}{
			"messages": []map]interface{}{
				{
					"role": "user",
					"content": []map]interface{}{
						{
							"image": "data:image/png;base64," + imageBase64,
						},
						{
							"text": prompt,
						},
					},
				},
			},
		},
		"parameters": map]interface{}{
			"result_format": "message",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", q.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+q.apiKey)

	// 发送请求
	resp, err := q.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误: %d - %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	var result map]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取回复内容
	output, ok := result["output"].(map]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: 缺少 output 字段")
	}

	choices, ok := output["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("响应格式错误: 缺少 choices")
	}

	firstChoice, ok := choices[0].(map]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: choices 格式错误")
	}

	message, ok := firstChoice["message"].(map]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: 缺少 message")
	}

	content, ok := message["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("响应格式错误: 缺少 content")
	}

	// 获取文本内容
	firstContent, ok := content[0].(map]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误: content 格式错误")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("响应格式错误: 缺少 text 字段")
	}

	return text, nil
}

// GenerateText 生成文本
func (q *QwenVL) GenerateText(prompt string) (string, error) {
	// 构建请求
	reqBody := map]interface{}{
		"model": q.model,
		"input": map]interface{}{
			"messages": []map]interface{}{
				{
					"role": "user",
					"content": []map]interface{}{
						{
							"text": prompt,
						},
					},
				},
			},
		},
		"parameters": map]interface{}{
			"result_format": "message",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %w", err)
	}

	// 使用文本生成 API 端点
	textEndpoint := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	req, err := http.NewRequest("POST", textEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+q.apiKey)

	resp, err := q.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误: %d - %s", resp.StatusCode, string(body))
	}

	var result map]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	output, ok := result["output"].(map]interface{})
	if !ok {
		return "", fmt.Errorf("响应格式错误")
	}

	text, ok := output["text"].(string)
	if !ok {
		return "", fmt.Errorf("响应格式错误: 缺少 text")
	}

	return text, nil
}

// AIVisualAnalyzer AI 视觉分析器
type AIVisualAnalyzer struct {
	provider AIProvider
}

// NewAIVisualAnalyzer 创建 AI 视觉分析器
func NewAIVisualAnalyzer(provider AIProvider) *AIVisualAnalyzer {
	return &AIVisualAnalyzer{
		provider: provider,
	}
}

// ElementPosition 元素位置
type ElementPosition struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Text   string `json:"text"`
	Type   string `json:"type"`
}

// AnalyzePageElements 分析页面元素
func (a *AIVisualAnalyzer) AnalyzePageElements(imageData []byte, description string) ([]ElementPosition, error) {
	prompt := fmt.Sprintf(`分析这张网页截图，找出所有可交互元素的位置。
需要找出以下类型的元素：
1. 按钮（clickable buttons）
2. 链接（links）
3. 输入框（input fields）
4. 文本框（text areas）

对于每个元素，请提供：
- 元素类型（button/link/input）
- 元素的中心坐标（x, y）
- 元素的宽高（width, height）
- 元素的文本内容或描述

请用 JSON 数组格式返回，格式如下：
[{"type": "button", "x": 100, "y": 200, "width": 80, "height": 30, "text": "提交"}]

如果找不到任何元素，返回空数组：[]

%s`, description)

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 结果
	var positions []ElementPosition
	if err := json.Unmarshal([]byte(result), &positions); err != nil {
		// 尝试提取 JSON 部分
		return a.parseJSONFromText(result)
	}

	return positions, nil
}

// FindButton 查找按钮位置
func (a *AIVisualAnalyzer) FindButton(imageData []byte, buttonText string) (*ElementPosition, error) {
	prompt := fmt.Sprintf(`分析这张网页截图，找出包含 "%s" 文本的按钮位置。
返回 JSON 格式：
{"type": "button", "x": 100, "y": 200, "width": 80, "height": 30, "text": "按钮文本"}

如果没有找到，返回 {"error": "not found"}`, buttonText)

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	var pos ElementPosition
	if err := json.Unmarshal([]byte(result), &pos); err != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err)
	}

	if pos.X == 0 && pos.Y == 0 {
		return nil, fmt.Errorf("未找到按钮: %s", buttonText)
	}

	return &pos, nil
}

// FindInputBox 查找输入框位置
func (a *AIVisualAnalyzer) FindInputBox(imageData []byte, boxDescription string) (*ElementPosition, error) {
	prompt := fmt.Sprintf(`分析这张网页截图，找出输入框的位置。
%s
返回 JSON 格式：
{"type": "input", "x": 100, "y": 200, "width": 300, "height": 40}`, boxDescription)

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	var pos ElementPosition
	if err := json.Unmarshal([]byte(result), &pos); err != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err)
	}

	return &pos, nil
}

// ConfirmAction 确认操作结果
func (a *AIVisualAnalyzer) ConfirmAction(beforeImage, afterImage []byte, action string) (bool, string, error) {
	prompt := fmt.Sprintf(`比较操作前后的两张网页截图，判断操作是否成功。

操作描述：%s

请分析：
1. 页面是否发生了预期的变化
2. 是否出现了错误提示
3. 操作是否完成

返回 JSON 格式：
{"success": true, "reason": "操作成功，页面已更新"}`, action)

	result, err := a.provider.AnalyzeImage(afterImage, prompt)
	if err != nil {
		return false, "", err
	}

	// 解析结果
	var response struct {
		Success bool   `json:"success"`
		Reason  string `json:"reason"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return true, result, nil // 无法解析时假设成功
	}

	return response.Success, response.Reason, nil
}

// FindJobCards 查找岗位卡片
func (a *AIVisualAnalyzer) FindJobCards(imageData []byte) ([]ElementPosition, error) {
	prompt := `分析这张网页截图，找出所有招聘岗位卡片的位置。

每个岗位卡片通常包含：
- 职位名称
- 公司名称
- 薪资范围
- 工作地点

返回 JSON 数组格式：
[{"type": "job_card", "x": 100, "y": 200, "width": 400, "height": 120, "text": "高级Go开发工程师 - 字节跳动 - 25K-50K"}]`

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	var positions []ElementPosition
	if err := json.Unmarshal([]byte(result), &positions); err != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err)
	}

	return positions, nil
}

// FindChatButton 查找聊天/沟通按钮
func (a *AIVisualAnalyzer) FindChatButton(imageData []byte) (*ElementPosition, error) {
	prompt := `分析这张网页截图，找出"立即沟通"或"boss直聊"按钮的位置。
返回 JSON 格式：
{"type": "button", "x": 100, "y": 200, "width": 100, "height": 40, "text": "立即沟通"}

如果没有找到，返回 {"error": "not found"}`

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	var pos ElementPosition
	if err := json.Unmarshal([]byte(result), &pos); err != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err)
	}

	if pos.X == 0 && pos.Y == 0 {
		return nil, fmt.Errorf("未找到聊天按钮")
	}

	return &pos, nil
}

// FindSendButton 查找发送按钮
func (a *AIVisualAnalyzer) FindSendButton(imageData []byte) (*ElementPosition, error) {
	prompt := `分析这张网页截图，找出消息发送按钮的位置。
可能是"发送"按钮、飞机图标按钮等。
返回 JSON 格式：
{"type": "button", "x": 100, "y": 200, "width": 60, "height": 30, "text": "发送"}

如果没有找到，返回 {"error": "not found"}`

	result, err := a.provider.AnalyzeImage(imageData, prompt)
	if err != nil {
		return nil, err
	}

	var pos ElementPosition
	if err := json.Unmarshal([]byte(result), &pos); err != nil {
		return nil, fmt.Errorf("解析结果失败: %w", err)
	}

	if pos.X == 0 && pos.Y == 0 {
		return nil, fmt.Errorf("未找到发送按钮")
	}

	return &pos, nil
}

// parseJSONFromText 从文本中提取 JSON
func (a *AIVisualAnalyzer) parseJSONFromText(text string) ([]ElementPosition, error) {
	// 尝试找到 JSON 数组的开始和结束
	start := -1
	end := -1

	for i := 0; i < len(text); i++ {
		if text[i] == '[' {
			start = i
		}
		if text[i] == ']' && start != -1 {
			end = i + 1
			break
		}
	}

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("无法从文本中提取 JSON: %s", text)
	}

	jsonStr := text[start:end]
	var positions []ElementPosition
	if err := json.Unmarshal([]byte(jsonStr), &positions); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return positions, nil
}

// GreetingGenerator 打招呼语生成器
type GreetingGenerator struct {
	provider AIProvider
	template string
}

// NewGreetingGenerator 创建打招呼语生成器
func NewGreetingGenerator(provider AIProvider, template string) *GreetingGenerator {
	if template == "" {
		template = "请基于以下信息生成简洁友好的中文打招呼语，不超过60字：\n个人介绍：%s\n关键词：%s\n职位名称：%s\n职位描述：%s"
	}
	return &GreetingGenerator{
		provider: provider,
		template: template,
	}
}

// Generate 生成打招呼语
func (g *GreetingGenerator) Generate(introduce, keyword, jobName, jobDesc string) (string, error) {
	prompt := fmt.Sprintf(g.template, introduce, keyword, jobName, jobDesc)
	return g.provider.GenerateText(prompt)
}

// InitAIService 初始化 AI 服务
func InitAIService() (AIProvider, error) {
	// 从配置获取 API Key
	apiKey := config.GetString("ai.qwen.api_key")
	if apiKey == "" {
		// 尝试从环境变量获取
		apiKey = "your-api-key"
		config.Warn("AI API Key 未配置，将使用占位符")
	}

	model := config.GetString("ai.qwen.model")
	if model == "" {
		model = "qwen-vl-plus"
	}

	config.Info("初始化 AI 服务: ", model)
	return NewQwenVL(apiKey, model), nil
}
