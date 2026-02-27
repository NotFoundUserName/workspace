package main

import (
	"fmt"
	"log"

	"chatbot/internal/bot"
	"chatbot/internal/config"
)

func main() {
	fmt.Println("Testing bot initialization...")
	
	// 加载配置
	fmt.Println("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("Configuration loaded successfully")

	// 创建机器人实例
	fmt.Println("Creating bot instance...")
	chatbot := bot.New(cfg)
	fmt.Println("Bot instance created successfully")

	// 测试机器人处理输入
	fmt.Println("Testing bot process input...")
	response := chatbot.ProcessInput("你好")
	fmt.Printf("Bot response: %s\n", response)

	fmt.Println("Test completed successfully!")
}
