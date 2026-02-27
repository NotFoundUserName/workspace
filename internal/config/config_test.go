package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 测试从默认配置加载
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证默认配置
	if cfg.BotName != "ChatBot" {
		t.Errorf("Expected BotName to be 'ChatBot', got '%s'", cfg.BotName)
	}

	if cfg.DefaultLang != "zh" {
		t.Errorf("Expected DefaultLang to be 'zh', got '%s'", cfg.DefaultLang)
	}

	// 验证中文配置存在
	if _, ok := cfg.Languages["zh"]; !ok {
		t.Error("Expected Chinese language config to exist")
	}

	// 验证英文配置存在
	if _, ok := cfg.Languages["en"]; !ok {
		t.Error("Expected English language config to exist")
	}
}

func TestLoadFromEnv(t *testing.T) {
	// 设置环境变量
	os.Setenv("CHATBOT_BOT_NAME", "TestBot")
	os.Setenv("CHATBOT_DEFAULT_LANG", "en")
	defer os.Unsetenv("CHATBOT_BOT_NAME")
	defer os.Unsetenv("CHATBOT_DEFAULT_LANG")

	// 加载配置
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证环境变量覆盖
	if cfg.BotName != "TestBot" {
		t.Errorf("Expected BotName to be 'TestBot', got '%s'", cfg.BotName)
	}

	if cfg.DefaultLang != "en" {
		t.Errorf("Expected DefaultLang to be 'en', got '%s'", cfg.DefaultLang)
	}
}

func TestValidateConfig(t *testing.T) {
	// 创建一个不完整的配置
	cfg := &Config{
		BotName: "",
		DefaultLang: "",
		Languages: nil,
	}

	// 验证配置
	err := validateConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to validate config: %v", err)
	}

	// 验证配置被修正
	if cfg.BotName != "ChatBot" {
		t.Errorf("Expected BotName to be 'ChatBot', got '%s'", cfg.BotName)
	}

	if cfg.DefaultLang != "zh" {
		t.Errorf("Expected DefaultLang to be 'zh', got '%s'", cfg.DefaultLang)
	}

	// 验证中文配置被添加
	if _, ok := cfg.Languages["zh"]; !ok {
		t.Error("Expected Chinese language config to be added")
	}
}
