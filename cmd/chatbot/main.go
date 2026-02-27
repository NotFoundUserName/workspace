package main

import (
	"fmt"
	"log"
	"os"

	"chatbot/internal/bot"
	"chatbot/internal/config"
	"chatbot/internal/handler"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建机器人实例
	chatbot := bot.New(cfg)

	// 创建输入处理器
	inputHandler := handler.NewInputHandler()

	// 启动聊天循环
	fmt.Println("聊天机器人已启动。输入'exit'或'退出'退出程序。")
	fmt.Println("----------------------------------------")

	for {
		// 获取用户输入
		input, err := inputHandler.GetInput()
		if err != nil {
			log.Printf("Error getting input: %v", err)
			continue
		}

		// 检查是否退出
		if input == "exit" || input == "退出" {
			fmt.Println("再见！")
			os.Exit(0)
		}

		// 处理输入并获取响应
		response := chatbot.ProcessInput(input)

		// 输出响应
		fmt.Printf("Bot: %s\n", response)
	}
}
