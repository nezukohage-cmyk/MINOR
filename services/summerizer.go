package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"lexxi/models"
	//"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

const SUMMARIZER_MODEL = "models/gemini-2.5-flash"

// ---------- AI SUMMARIZER ----------
func Summarize(text string) (string, error) {

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is missing")
	}

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "Summarize the following text clearly:\n\n" + text},
				},
			},
		},
	}

	data, _ := json.Marshal(reqBody)

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s",
		SUMMARIZER_MODEL,
		apiKey,
	)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
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
		return "", fmt.Errorf("Gemini returned no candidates – model may be overloaded")
	}

	return output.Candidates[0].Content.Parts[0].Text, nil
}

// ---------- SAVE SUMMARY ----------
func SaveSummary(userID, inputText, summary string) error {
	record := &models.Summary{
		UserID:    userID,
		InputText: inputText, // FIXED field name
		Summary:   summary,
		CreatedAt: time.Now(),
	}

	return mgm.Coll(&models.Summary{}).Create(record)
}

// ---------- FETCH SUMMARIES ----------
func FetchSummaries(userID string) ([]models.Summary, error) {
	var items []models.Summary

	err := mgm.Coll(&models.Summary{}).SimpleFind(
		&items,
		bson.M{"user_id": userID},
	)

	return items, err
}
