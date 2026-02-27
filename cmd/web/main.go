package main

import (
	"chatbot/internal/bot"
	"chatbot/internal/config"
	"chatbot/internal/errors"
	"chatbot/internal/logger"
	"encoding/json"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed web/*
var webFS embed.FS

type chatRequest struct {
	Message string `json:"message"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

type infoResponse struct {
	BotName      string   `json:"bot_name"`
	Greeting     string   `json:"greeting"`
	CurrentLang  string   `json:"current_lang"`
	SupportedLangs []string `json:"supported_langs"`
}

type languageRequest struct {
	Language string `json:"language"`
}

type languageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	CurrentLang string `json:"current_lang"`
}

type messageResponse struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type historyResponse struct {
	Messages []messageResponse `json:"messages"`
}

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("HTTP request received", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		})

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger.Info("HTTP request completed", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration,
		})
	})
}

// CORS中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Println("Starting web server...")
	
	// 加载配置
	log.Println("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// 创建机器人实例
	log.Println("Creating bot instance...")
	b := bot.New(cfg)
	log.Println("Bot instance created successfully")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			errors.MethodNotAllowed(w)
			return
		}

		// 获取支持的语言列表
		supportedLangs := make([]string, 0, len(cfg.Languages))
		for lang := range cfg.Languages {
			supportedLangs = append(supportedLangs, lang)
		}

		// 获取当前语言
		currentLang := b.GetLanguage()

		// 获取当前语言的问候语
		greeting := ""
		if langConfig, ok := cfg.Languages[currentLang]; ok {
			greeting = langConfig.Greeting
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(infoResponse{
			BotName:      cfg.BotName,
			Greeting:     greeting,
			CurrentLang:  currentLang,
			SupportedLangs: supportedLangs,
		})
	})

	mux.HandleFunc("/api/language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errors.MethodNotAllowed(w)
			return
		}
		defer r.Body.Close()

		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			errors.BadRequest(w, err, "Failed to read request body")
			return
		}

		var req languageRequest
		if err := json.Unmarshal(body, &req); err != nil {
			errors.BadRequest(w, err, "Invalid JSON format")
			return
		}

		lang := strings.TrimSpace(req.Language)
		if lang == "" {
			errors.BadRequest(w, nil, "Language is required")
			return
		}

		// 设置语言
		success := b.SetLanguage(lang)
		currentLang := b.GetLanguage()

		message := "Language changed successfully"
		if !success {
			message = "Language not supported, using default language"
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(languageResponse{
			Success:     success,
			Message:     message,
			CurrentLang: currentLang,
		})
	})
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errors.MethodNotAllowed(w)
			return
		}
		defer r.Body.Close()

		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			errors.BadRequest(w, err, "Failed to read request body")
			return
		}

		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			errors.BadRequest(w, err, "Invalid JSON format")
			return
		}

		msg := strings.TrimSpace(req.Message)
		if msg == "" {
			errors.BadRequest(w, nil, "Message is required")
			return
		}

		reply := b.ProcessInput(msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(chatResponse{Reply: reply})
	})

	mux.HandleFunc("/api/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			errors.MethodNotAllowed(w)
			return
		}

		// 获取对话历史
		history := b.GetHistory()

		// 转换为响应格式
		messages := make([]messageResponse, len(history))
		for i, msg := range history {
			messages[i] = messageResponse{
				Role:      msg.Role,
				Content:   msg.Content,
				Timestamp: msg.Timestamp,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(historyResponse{Messages: messages})
	})

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("Failed to init embedded web fs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	addr := os.Getenv("CHATBOT_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(corsMiddleware(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	displayURL := addr
	if strings.HasPrefix(addr, ":") {
		displayURL = "http://localhost" + addr + "/"
	} else if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		displayURL = strings.TrimRight(addr, "/") + "/"
	} else {
		displayURL = "http://" + strings.TrimRight(addr, "/") + "/"
	}
	log.Printf("Web chat UI: %s (listening on %s)", displayURL, addr)
	log.Println("Server starting...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Println("Server stopped")
}

