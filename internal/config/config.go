package config

import (
	"encoding/json"
	"os"
	"strconv"
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

type LanguageConfig struct {
	Greeting        string         `json:"greeting"`
	DefaultResponse string         `json:"default_response"`
	ResponseRules   []ResponseRule `json:"response_rules"`
}

type Config struct {
	BotName    string                `json:"bot_name"`
	DefaultLang string                `json:"default_lang"`
	Languages  map[string]LanguageConfig `json:"languages"`
	AiApi      AiApiConfig           `json:"ai_api"`
}

func Load() (*Config, error) {
	// 尝试从配置文件加载
	var config *Config

	if _, err := os.Stat("config.json"); err == nil {
		config, err = loadFromFile("config.json")
		if err != nil {
			return nil, err
		}
	} else {
		config = getDefaultConfig()
	}

	// 从环境变量加载配置（覆盖文件配置）
	config = loadFromEnv(config)

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
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

func loadFromEnv(config *Config) *Config {
	// 加载基本配置
	if botName := os.Getenv("CHATBOT_NAME"); botName != "" {
		config.BotName = botName
	}
	if lang := os.Getenv("CHATBOT_DEFAULT_LANG"); lang != "" {
		config.DefaultLang = lang
	}
	
	// 加载默认语言的配置
	if greeting := os.Getenv("CHATBOT_GREETING"); greeting != "" {
		if langConfig, ok := config.Languages[config.DefaultLang]; ok {
			// 获取结构体，修改后放回map
			langConfig.Greeting = greeting
			config.Languages[config.DefaultLang] = langConfig
		} else {
			// 如果默认语言不存在，创建一个新的语言配置
			config.Languages[config.DefaultLang] = LanguageConfig{
				Greeting:        greeting,
				DefaultResponse: "我不太明白你的意思。你能换一种说法吗？",
				ResponseRules:   getDefaultResponseRules(),
			}
		}
	}
	if defaultResponse := os.Getenv("CHATBOT_DEFAULT_RESPONSE"); defaultResponse != "" {
		if langConfig, ok := config.Languages[config.DefaultLang]; ok {
			// 获取结构体，修改后放回map
			langConfig.DefaultResponse = defaultResponse
			config.Languages[config.DefaultLang] = langConfig
		} else {
			// 如果默认语言不存在，创建一个新的语言配置
			config.Languages[config.DefaultLang] = LanguageConfig{
				Greeting:        "你好！今天我能帮你什么忙？",
				DefaultResponse: defaultResponse,
				ResponseRules:   getDefaultResponseRules(),
			}
		}
	}

	// 加载AI API配置
	if enabled := os.Getenv("CHATBOT_AI_ENABLED"); enabled != "" {
		if val, err := strconv.ParseBool(enabled); err == nil {
			config.AiApi.Enabled = val
		}
	}
	if apiKey := os.Getenv("CHATBOT_AI_API_KEY"); apiKey != "" {
		config.AiApi.ApiKey = apiKey
	}
	if apiUrl := os.Getenv("CHATBOT_AI_API_URL"); apiUrl != "" {
		config.AiApi.ApiUrl = apiUrl
	}
	if model := os.Getenv("CHATBOT_AI_MODEL"); model != "" {
		config.AiApi.Model = model
	}
	if maxTokens := os.Getenv("CHATBOT_AI_MAX_TOKENS"); maxTokens != "" {
		if val, err := strconv.Atoi(maxTokens); err == nil {
			config.AiApi.MaxTokens = val
		}
	}

	return config
}

func validateConfig(config *Config) error {
	// 验证基本配置
	if config.BotName == "" {
		config.BotName = "ChatBot"
	}
	if config.DefaultLang == "" {
		config.DefaultLang = "zh"
	}

	// 验证语言配置
	if config.Languages == nil {
		config.Languages = make(map[string]LanguageConfig)
	}

	// 确保至少有中文配置
	if _, ok := config.Languages["zh"]; !ok {
		config.Languages["zh"] = LanguageConfig{
			Greeting:        "你好！今天我能帮你什么忙？",
			DefaultResponse: "我不太明白你的意思。你能换一种说法吗？",
			ResponseRules:   getDefaultResponseRules(),
		}
	}

	// 验证AI API配置
	if config.AiApi.Enabled {
		if config.AiApi.ApiKey == "" || config.AiApi.ApiKey == "your_api_key_here" {
			config.AiApi.Enabled = false
		}
		if config.AiApi.ApiUrl == "" {
			config.AiApi.ApiUrl = "https://api.openai.com/v1/chat/completions"
		}
		if config.AiApi.Model == "" {
			config.AiApi.Model = "gpt-3.5-turbo"
		}
		if config.AiApi.MaxTokens <= 0 {
			config.AiApi.MaxTokens = 150
		}
	}

	return nil
}

func getDefaultResponseRules() []ResponseRule {
	return []ResponseRule {
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
	}
}

func getDefaultConfig() *Config {
	return &Config{
		BotName:    "ChatBot",
		DefaultLang: "zh",
		Languages: map[string]LanguageConfig{
			"zh": {
				Greeting:        "你好！今天我能帮你什么忙？",
				DefaultResponse: "我不太明白你的意思。你能换一种说法吗？",
				ResponseRules: []ResponseRule{
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
					{
						Keywords: []string{"再见", "bye", "goodbye"},
						Response: "再见！祝你有愉快的一天！",
					},
					{
						Keywords: []string{"谢谢", "thank you"},
						Response: "不客气！很高兴能帮到你。",
					},
				},
			},
			"en": {
				Greeting:        "Hello! How can I help you today?",
				DefaultResponse: "I'm not sure I understand. Can you please rephrase?",
				ResponseRules: []ResponseRule{
					{
						Keywords: []string{"hello", "hi"},
						Response: "Hello! How can I help you today?",
					},
					{
						Keywords: []string{"how are you"},
						Response: "I'm doing well, thank you! How about you?",
					},
					{
						Keywords: []string{"name"},
						Response: "My name is ChatBot. Nice to meet you!",
					},
					{
						Keywords: []string{"help"},
						Response: "I can help you with basic questions. Try asking me about myself or just say hello!",
					},
					{
						Keywords: []string{"time"},
						Response: "I don't have access to the current time, but I'm here to chat with you!",
					},
					{
						Keywords: []string{"weather"},
						Response: "I don't have access to weather information, but I'm happy to chat about other topics!",
					},
					{
						Keywords: []string{"bye", "goodbye"},
						Response: "Goodbye! Have a nice day!",
					},
					{
						Keywords: []string{"thank you", "thanks"},
						Response: "You're welcome! I'm glad I could help.",
					},
				},
			},
		},
		AiApi: AiApiConfig{
			Enabled:    false,
			ApiKey:     "your_api_key_here",
			ApiUrl:     "https://api.openai.com/v1/chat/completions",
			Model:      "gpt-3.5-turbo",
			MaxTokens:  150,
		},
	}
}
