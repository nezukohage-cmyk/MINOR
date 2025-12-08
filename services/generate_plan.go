package services

import (
	"encoding/json"
	"fmt"
	"strings"
	//	"time"
)

// GenerateStudyPlan builds a strict JSON prompt and asks Gemini, then parses the returned JSON.
// subject is optional, rows is the entire table the user enters.
func GenerateStudyPlan(userID string, title string, rows interface{}) (map[string]interface{}, string, error) {
	// Build the JSON prompt: instruct model to reply with STRICT JSON only
	preamble := `You are an academic planner assistant. Based ONLY on the user's study table, 
generate a complete study roadmap.

User does NOT know how many hours are required. YOU must:
- Estimate realistic hours_needed for each topic.
- Determine hours_per_day automatically.
- Break topics into daily tasks.
- Ensure schedule fits before deadlines.
- Include prerequisites where helpful.

Respond ONLY with strict JSON using this schema:

{
  "title": "<string>",
  "summary": "<short one-line summary>",
  "overall_hours_needed": <integer>,
  "per_subject": [
    {
      "subject": "<subject>",
      "estimated_hours_per_day": <integer>,
      "topics": [
        {
          "name": "<topic>",
          "prerequisites": ["..."],
          "estimated_hours_needed": <integer>,
          "start_date": "<ISO>",
          "end_date": "<ISO>",
          "daily_tasks": [
            {"date":"YYYY-MM-DD", "tasks":["..."]}
          ]
        }
      ]
    }
  ],
  "daily_timetable": [
    {"date":"YYYY-MM-DD", "tasks":["..."], "total_hours": <number>}
  ]
}`

	// Convert rows to pretty JSON to include in prompt
	rowsJSONBytes, _ := json.MarshalIndent(rows, "", "  ")
	prompt := preamble + "\nTitle: " + title + "\nUserTable:\n" + string(rowsJSONBytes) + "\n\nReturn the JSON now."

	// Call your existing Gemini text function (which returns text)
	respText, err := AskGemini("", prompt)
	if err != nil {
		return nil, "", fmt.Errorf("Gemini call failed: %v", err)
	}

	// Try to find the JSON substring from respText - in case model returns commentary + JSON
	respTextTrim := strings.TrimSpace(respText)

	// Attempt direct unmarshal
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(respTextTrim), &parsed); err == nil {
		return parsed, respTextTrim, nil
	}

	// If direct unmarshal failed, attempt to extract first {...} block
	start := strings.Index(respTextTrim, "{")
	end := strings.LastIndex(respTextTrim, "}")
	if start != -1 && end != -1 && end > start {
		jsonCandidate := respTextTrim[start : end+1]
		if err := json.Unmarshal([]byte(jsonCandidate), &parsed); err == nil {
			return parsed, respTextTrim, nil
		}
	}

	// If still failed, return raw response and error
	return nil, respTextTrim, fmt.Errorf("failed to parse Gemini JSON output")
}
