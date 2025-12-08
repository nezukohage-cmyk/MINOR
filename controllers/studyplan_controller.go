package controllers

import (
	"fmt"
	"net/http"
	"time"

	"lexxi/models"
	"lexxi/services"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

// POST /planner/generate
func GeneratePlan(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}

	var body struct {
		Title string                   `json:"title"`
		Rows  []models.StudyPlannerRow `json:"rows"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call Gemini
	parsed, raw, err := services.GenerateStudyPlan(userID, body.Title, body.Rows)
	if err != nil {
		plan := &models.StudyPlan{
			UserID:      userID,
			Title:       body.Title,
			InputRows:   body.Rows,
			PlanJSON:    nil,
			RawOutput:   raw,
			GeneratedAt: time.Now(),
		}
		_ = mgm.Coll(plan).Create(plan)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"raw":   raw,
		})
		return
	}

	// Build model
	plan := &models.StudyPlan{
		UserID:      userID,
		Title:       body.Title,
		InputRows:   body.Rows,
		PlanJSON:    parsed,
		RawOutput:   raw,
		DailyTable:  []models.DailyTask{},
		GeneratedAt: time.Now(),
	}

	// ---------------------------------------------------
	// Extract daily_timetable safely
	// ---------------------------------------------------
	if dt, ok := parsed["daily_timetable"].([]interface{}); ok {
		for _, entry := range dt {

			row, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}

			date, _ := row["date"].(string)

			// Convert tasks safely
			rawTasks, ok := row["tasks"].([]interface{})
			if !ok {
				continue
			}

			tasks := []string{}
			for _, t := range rawTasks {
				if s, ok := t.(string); ok {
					tasks = append(tasks, s)
				}
			}

			// Push into model
			plan.DailyTable = append(plan.DailyTable, models.DailyTask{
				Date:  date,
				Tasks: tasks,
			})
		}
	} else {
		fmt.Println("WARNING: NO daily_timetable found in Gemini response")
	}

	// ---------------------------------------------------
	// Also store parsed daily table back into PlanJSON
	// Ensures DB contains the extracted version too
	// ---------------------------------------------------
	plan.PlanJSON["daily_timetable"] = plan.DailyTable

	// ---------------------------------------------------
	// Save StudyPlan to MongoDB
	// ---------------------------------------------------
	if err := mgm.Coll(plan).Create(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save study plan"})
		return
	}

	// ---------------------------------------------------
	// Save checklist tasks
	// ---------------------------------------------------
	if err := services.SaveDailyTasks(userID, plan.ID.Hex(), *plan); err != nil {
		fmt.Println("WARNING: Failed to save tasks:", err)
	}

	// ---------------------------------------------------
	// Return final output
	// ---------------------------------------------------
	c.JSON(http.StatusCreated, gin.H{
		"message": "plan generated",
		"plan_id": plan.ID.Hex(),
		"plan":    parsed,
	})
}
