package bot

import (
	"chatbot/internal/config"
	"testing"
)

func TestProcessInput(t *testing.T) {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 创建机器人实例
	chatbot := New(cfg)

	// 测试用例
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello", "你好！今天我能帮你什么忙？"},
		{"Hi", "你好！今天我能帮你什么忙？"},
		{"How are you?", "我很好，谢谢！你呢？"},
		{"What's your name?", "我的名字是ChatBot。很高兴认识你！"},
		{"Help", "我可以帮助你回答一些基本问题。试着问我关于我自己的问题，或者只是打个招呼！"},
		{"What time is it?", "我无法获取当前时间，但我可以陪你聊天！"},
		{"What's the weather like?", "我无法获取天气信息，但我很乐意和你聊其他话题！"},
		{"Random question", "我不太明白你的意思。你能换一种说法吗？"},
	}

	// 运行测试
	for _, tc := range testCases {
		result := chatbot.ProcessInput(tc.input)
		if result != tc.expected {
			t.Errorf("Input: %q, Expected: %q, Got: %q", tc.input, tc.expected, result)
		}
	}
}
