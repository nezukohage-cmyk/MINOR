package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ListGeminiModels() (string, error) {
	//var apiKey string
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is missing")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models?key=" + apiKey
	// GET the /models list
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// Return raw JSON so you can inspect in terminal
	return string(body), nil
}
