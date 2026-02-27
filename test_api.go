package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Testing API endpoints...")

	// 测试/api/info端点
	fmt.Println("Testing /api/info endpoint...")
	resp, err := http.Get("http://localhost:8080/api/info")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", body)
}
