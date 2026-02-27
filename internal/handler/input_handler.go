package handler

import (
	"bufio"
	"os"
	"strings"
)

type InputHandler struct {
	scanner *bufio.Scanner
}

func NewInputHandler() *InputHandler {
	return &InputHandler{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (h *InputHandler) GetInput() (string, error) {
	print("你: ")
	if h.scanner.Scan() {
		input := h.scanner.Text()
		return strings.TrimSpace(input), nil
	}
	if err := h.scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}
