package config

import (
	"encoding/json"
	"os"
)

type ResponseRule struct {
	Keywords []string `json:"keywords"`
	Response string   `json:"response"`
}

type AiApiConfig struct {
	Enabled    bool   `json:"enabled"`
	ApiKey     string `json:"api_key"`
	ApiUrl     string `json:"api_url"`
	Model      string `json:"model"`
	MaxTokens  int    `json:"max_tokens"`
}

type Config struct {
	BotName         string         `json:"bot_name"`
	Greeting        string         `json:"greeting"`
	DefaultResponse string         `json:"default_response"`
	AiApi           AiApiConfig    `json:"ai_api"`
	ResponseRules   []ResponseRule `json:"response_rules"`
}

func Load() (*Config, error) {
	// 尝试从配置文件加载
	if _, err := os.Stat("config.json"); err == nil {
		return loadFromFile("config.json")
	}

	// 如果配置文件不存在，返回默认配置
	return getDefaultConfig(), nil
}

func loadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getDefaultConfig() *Config {
	return &Config{
		BotName:         "ChatBot",
		Greeting:        "你好！今天我能帮你什么忙？",
		DefaultResponse: "我不太明白你的意思。你能换一种说法吗？",
		AiApi: AiApiConfig{
			Enabled:    false,
			ApiKey:     "your_api_key_here",
			ApiUrl:     "https://api.openai.com/v1/chat/completions",
			Model:      "gpt-3.5-turbo",
			MaxTokens:  150,
		},
		ResponseRules: []ResponseRule {
			{
				Keywords: []string{"hello", "hi", "你好", "嗨"},
				Response: "你好！今天我能帮你什么忙？",
			},
			{
				Keywords: []string{"how are you", "你好吗", "怎么样"},
				Response: "我很好，谢谢！你呢？",
			},
			{
				Keywords: []string{"name", "名字"},
				Response: "我的名字是ChatBot。很高兴认识你！",
			},
			{
				Keywords: []string{"help", "帮助"},
				Response: "我可以帮助你回答一些基本问题。试着问我关于我自己的问题，或者只是打个招呼！",
			},
			{
				Keywords: []string{"time", "时间"},
				Response: "我无法获取当前时间，但我可以陪你聊天！",
			},
			{
				Keywords: []string{"weather", "天气"},
				Response: "我无法获取天气信息，但我很乐意和你聊其他话题！",
			},
		},
	}
}
