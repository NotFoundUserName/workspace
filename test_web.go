package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting simple web server...")
	
	// 注册一个简单的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// 启动服务器
	addr := ":8080"
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
