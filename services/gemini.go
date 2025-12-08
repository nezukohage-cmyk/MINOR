package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const GEMINI_MODEL = "models/gemini-flash-latest"

func AskGemini(subject string, question string) (string, error) {

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is missing")
	}

	// Build a prompt. Subject is optional.
	var prompt string
	if subject == "" {
		prompt = fmt.Sprintf("You are an academic tutor. Answer clearly and concisely.\nQuestion: %s", question)
	} else {
		prompt = fmt.Sprintf("You are an academic tutor specializing in %s. Explain concepts as if teaching a student.\nQuestion: %s", subject, question)
	}

	// Request body format for v1beta
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	// Correct v1beta endpoint with your selected model
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s",
		GEMINI_MODEL,
		apiKey,
	)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	fmt.Println("GEMINI RAW:", string(raw))

	// Parse response
	var output struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(raw, &output); err != nil {
		return "", fmt.Errorf("invalid Gemini response")
	}

	if len(output.Candidates) == 0 {
		return "", fmt.Errorf("Gemini returned no candidates")
	}

	return output.Candidates[0].Content.Parts[0].Text, nil
}
