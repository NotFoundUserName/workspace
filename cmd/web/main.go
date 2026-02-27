package main

import (
	"chatbot/internal/bot"
	"chatbot/internal/config"
	"embed"
	"encoding/json"
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
	BotName  string `json:"bot_name"`
	Greeting string `json:"greeting"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b := bot.New(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(infoResponse{
			BotName:  cfg.BotName,
			Greeting: cfg.Greeting,
		})
	})
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		msg := strings.TrimSpace(req.Message)
		if msg == "" {
			http.Error(w, "message is required", http.StatusBadRequest)
			return
		}

		reply := b.ProcessInput(msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(chatResponse{Reply: reply})
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
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
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
	log.Fatal(srv.ListenAndServe())
}

