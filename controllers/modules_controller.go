package controllers

import (
	"bytes"
	"encoding/json"

	//"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory stores for demo
var (
	noteMutex sync.Mutex
	notes     = []map[string]interface{}{}

	todoMutex sync.Mutex
	todos     = map[string]map[string]interface{}{} // id -> {id,text,completed}

	classrooms = []map[string]interface{}{
		{"id": "c1", "title": "V Semester"},
		{"id": "c2", "title": "Java 2025"},
	}
)

// ---------- NOTES ----------
func GetNotes(c *gin.Context) {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func CreateNote(c *gin.Context) {
	var p struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	noteMutex.Lock()
	defer noteMutex.Unlock()
	entry := map[string]interface{}{"id": time.Now().UnixNano(), "title": p.Title, "url": p.Url}
	notes = append([]map[string]interface{}{entry}, notes...)
	c.JSON(http.StatusCreated, entry)
}

// ---------- TODOS ----------
func GetTodos(c *gin.Context) {
	todoMutex.Lock()
	defer todoMutex.Unlock()
	out := []map[string]interface{}{}
	for _, v := range todos {
		out = append(out, v)
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func CreateTodo(c *gin.Context) {
	var p struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	id := time.Now().Format("20060102150405.000000")
	item := map[string]interface{}{"id": id, "text": p.Text, "completed": false}
	todoMutex.Lock()
	todos[id] = item
	todoMutex.Unlock()
	c.JSON(http.StatusCreated, item)
}

func CompleteTodo(c *gin.Context) {
	id := c.Param("id")
	todoMutex.Lock()
	defer todoMutex.Unlock()
	if t, ok := todos[id]; ok {
		t["completed"] = true
		c.JSON(http.StatusOK, t)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	todoMutex.Lock()
	defer todoMutex.Unlock()
	if _, ok := todos[id]; ok {
		delete(todos, id)
		c.JSON(http.StatusOK, gin.H{"deleted": id})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

// ---------- CLASSROOMS ----------
func GetClassrooms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": classrooms})
}

// ---------- SUMMARIZER ----------
// This is a demo proxy that forwards the text to a configured external summarizer
// Set env GEMINI_URL to your model endpoint or use a local summarizer service.
// It includes a simple retry on 503 and returns a clear error if provider is overloaded.
func Summarize(c *gin.Context) {
	var p struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Basic validation length
	if len(p.Text) < 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text too short to summarize"})
		return
	}

	geminiURL := os.Getenv("GEMINI_URL") // set this in your env
	if geminiURL == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "summarizer not configured"})
		return
	}

	// Build payload - this depends on your model's expected format. Adjust accordingly.
	body := map[string]interface{}{
		"input": p.Text,
	}
	payload, _ := json.Marshal(body)

	// Try up to 2 retries if 503 (model overloaded)
	var resp *http.Response
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		resp, err = http.Post(geminiURL, "application/json", bytes.NewReader(payload))
		if err != nil {
			// network error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			// read body for details
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			// if last attempt, return clear error
			if attempt == 1 {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "model overloaded", "raw": string(b)})
				return
			}
			// wait and retry
			time.Sleep(800 * time.Millisecond)
			continue
		}
		break
	}

	if resp == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no response from summarizer"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "summarizer returned error", "raw": string(b)})
		return
	}

	b, _ := io.ReadAll(resp.Body)
	// assume summarizer returns {"summary":"..."}
	var got map[string]interface{}
	_ = json.Unmarshal(b, &got)
	summary := ""
	if s, ok := got["summary"].(string); ok {
		summary = s
	} else {
		// fallback: return raw body
		summary = string(b)
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}
