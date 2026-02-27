package bot

import (
	"bytes"
	"chatbot/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Bot struct {
	config    *config.Config
	httpClient *http.Client
}

// OpenAI API请求结构体
type openAiRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI API响应结构体
type openAiResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func New(cfg *config.Config) *Bot {
	return &Bot{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *Bot) ProcessInput(input string) string {
	// 转换为小写以进行匹配
	input = strings.ToLower(input)

	// 如果AI API启用，尝试使用AI生成响应
	if b.config.AiApi.Enabled && b.config.AiApi.ApiKey != "your_api_key_here" {
		aiResponse, err := b.generateAiResponse(input)
		if err == nil && aiResponse != "" {
			return aiResponse
		}
		// 如果AI API调用失败，回退到规则匹配
	}

	// 处理输入并生成响应
	response := b.generateResponse(input)

	return response
}

// 使用AI API生成响应
func (b *Bot) generateAiResponse(input string) (string, error) {
	// 构建请求体
	reqBody := openAiRequest{
		Model: b.config.AiApi.Model,
		Messages: []message{
			{
				Role:    "system",
				Content: "你是一个友好的聊天机器人，用中文回答用户的问题。",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		MaxTokens: b.config.AiApi.MaxTokens,
	}

	// 序列化请求体
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", b.config.AiApi.ApiUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.config.AiApi.ApiKey))

	// 发送请求
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// 解析响应
	var aiResp openAiResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", err
	}

	// 提取响应内容
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI API")
}

func (b *Bot) generateResponse(input string) string {
	// 使用配置文件中的响应规则
	for _, rule := range b.config.ResponseRules {
		for _, keyword := range rule.Keywords {
			if strings.Contains(input, keyword) {
				return rule.Response
			}
		}
	}

	// 如果没有匹配的规则，返回默认响应
	return b.config.DefaultResponse
}
