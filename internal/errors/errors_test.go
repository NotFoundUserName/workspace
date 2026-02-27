package errors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequest(t *testing.T) {
	// 创建一个测试请求
	req, err := http.NewRequest("GET", "/api/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 创建一个测试响应记录器
	w := httptest.NewRecorder()

	// 调用BadRequest函数
	BadRequest(w, nil, "Test error")

	// 验证响应
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", w.Header().Get("Content-Type"))
	}
}

func TestInternalServerError(t *testing.T) {
	// 创建一个测试请求
	req, err := http.NewRequest("GET", "/api/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 创建一个测试响应记录器
	w := httptest.NewRecorder()

	// 调用InternalServerError函数
	InternalServerError(w, nil, "Test error")

	// 验证响应
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", w.Header().Get("Content-Type"))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	// 创建一个测试请求
	req, err := http.NewRequest("GET", "/api/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 创建一个测试响应记录器
	w := httptest.NewRecorder()

	// 调用MethodNotAllowed函数
	MethodNotAllowed(w)

	// 验证响应
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", w.Header().Get("Content-Type"))
	}
}
