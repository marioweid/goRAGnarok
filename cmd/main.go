package main

import (
	"fmt"
	"goRAGnarok/internal"
	"net/http"
)

func main() {
	http.HandleFunc("/health", internal.HealthCheckHandler)
	http.HandleFunc("/v1/generate", internal.PostHandler)
	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
