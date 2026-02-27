package main

import (
	"fmt"
	"log"

	"chatbot/internal/config"
)

func main() {
	fmt.Println("Testing config loading...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("Config loaded successfully!")
	fmt.Printf("BotName: %s\n", cfg.BotName)
	fmt.Printf("DefaultLang: %s\n", cfg.DefaultLang)
	fmt.Printf("Supported languages: %v\n", cfg.Languages)
}
