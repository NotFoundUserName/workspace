package bot

import (
	"bytes"
	"chatbot/internal/config"
	"chatbot/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type Bot struct {
	config    *config.Config
	httpClient *http.Client
	history   []Message
	historyLimit int
	currentLang string
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
	// 创建自定义的HTTP客户端，配置连接池和超时
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Bot{
		config:    cfg,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
		history:   make([]Message, 0, 20), // 初始化历史记录，容量为20
		historyLimit: 20, // 历史记录限制为20条
		currentLang: cfg.DefaultLang, // 设置默认语言
	}
}

// 设置当前语言
func (b *Bot) SetLanguage(lang string) bool {
	// 检查语言是否存在
	if _, ok := b.config.Languages[lang]; ok {
		b.currentLang = lang
		return true
	}
	return false
}

// 获取当前语言
func (b *Bot) GetLanguage() string {
	return b.currentLang
}

// 获取当前语言的配置
func (b *Bot) getCurrentLanguageConfig() config.LanguageConfig {
	if langConfig, ok := b.config.Languages[b.currentLang]; ok {
		return langConfig
	}
	// 如果当前语言不存在，返回默认语言配置
	if langConfig, ok := b.config.Languages[b.config.DefaultLang]; ok {
		return langConfig
	}
	// 如果默认语言也不存在，返回一个空的配置
	return config.LanguageConfig{}
}

func (b *Bot) ProcessInput(input string) string {
	// 转换为小写以进行匹配
	input = strings.ToLower(input)

	// 记录用户输入
	logger.Info("User input received", map[string]interface{}{
		"input": input,
	})

	// 将用户输入添加到对话历史
	b.addToHistory("user", input)

	// 如果AI API启用，尝试使用AI生成响应
	if b.config.AiApi.Enabled && b.config.AiApi.ApiKey != "your_api_key_here" {
		logger.Info("Using AI API for response generation")
		aiResponse, err := b.generateAiResponse(input)
		if err == nil && aiResponse != "" {
			// 将AI响应添加到对话历史
			b.addToHistory("assistant", aiResponse)
			logger.Info("AI API response generated successfully")
			return aiResponse
		}
		// 如果AI API调用失败，回退到规则匹配
		logger.Warn("AI API call failed, falling back to rule-based response", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 处理输入并生成响应
	response := b.generateResponse(input)

	// 将响应添加到对话历史
	b.addToHistory("assistant", response)

	// 记录生成的响应
	logger.Info("Response generated", map[string]interface{}{
		"response": response,
	})

	return response
}

// 添加消息到对话历史
func (b *Bot) addToHistory(role, content string) {
	// 创建新消息
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	// 添加到历史记录
	b.history = append(b.history, message)

	// 如果超过历史记录限制，删除最早的消息
	if len(b.history) > b.historyLimit {
		b.history = b.history[1:]
	}
}

// 获取对话历史
func (b *Bot) GetHistory() []Message {
	return b.history
}

// 使用AI API生成响应
func (b *Bot) generateAiResponse(input string) (string, error) {
	// 根据当前语言设置系统提示
	systemPrompt := "你是一个友好的聊天机器人，用中文回答用户的问题。"
	if b.currentLang == "en" {
		systemPrompt = "You are a friendly chatbot that answers user questions in English."
	}

	// 构建消息列表，包括系统消息和对话历史
	messages := []message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// 添加最近的对话历史（最多5轮对话）
	historyLimit := 10 // 最多10条消息（5轮对话）
	startIdx := 0
	if len(b.history) > historyLimit {
		startIdx = len(b.history) - historyLimit
	}

	for i := startIdx; i < len(b.history); i++ {
		msg := b.history[i]
		messages = append(messages, message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 构建请求体
	reqBody := openAiRequest{
		Model:      b.config.AiApi.Model,
		Messages:   messages,
		MaxTokens:  b.config.AiApi.MaxTokens,
	}

	// 序列化请求体
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", b.config.AiApi.ApiUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.config.AiApi.ApiKey))
	// 添加User-Agent头
	req.Header.Set("User-Agent", "ChatBot/1.0")

	// 发送请求
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应
		errorBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status code: %d, response: %s", resp.StatusCode, string(errorBody))
	}

	// 解析响应
	var aiResp openAiResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 提取响应内容
	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI API")
}

func (b *Bot) generateResponse(input string) string {
	// 获取当前语言的配置
	langConfig := b.getCurrentLanguageConfig()

	// 使用当前语言的响应规则
	for _, rule := range langConfig.ResponseRules {
		for _, keyword := range rule.Keywords {
			if strings.Contains(input, keyword) {
				return rule.Response
			}
		}
	}

	// 如果没有匹配的规则，返回当前语言的默认响应
	return langConfig.DefaultResponse
}
