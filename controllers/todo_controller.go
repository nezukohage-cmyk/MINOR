package controllers

import (
	"lexxi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /todo/create
func CreateTask(c *gin.Context) {
	userID := c.GetString("user_id")

	var body struct {
		Date     string  `json:"date"`
		Task     string  `json:"task"`
		Deadline *string `json:"deadline"` // optional ISO8601 or RFC3339
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if body.Date == "" || body.Task == "" {
		c.JSON(400, gin.H{"error": "date and task are required"})
		return
	}

	newTask := &models.TodoTask{
		UserID: userID,
		Date:   body.Date,
		Task:   body.Task,
		Done:   false,
	}

	// ---- DEADLINE PARSING BLOCK ----
	if body.Deadline != nil && *body.Deadline != "" {

		var parsed time.Time
		var err error

		deadlineStr := *body.Deadline

		// Try strict RFC3339 first
		parsed, err = time.Parse(time.RFC3339, deadlineStr)

		if err != nil {
			// Try ISO8601 WITHOUT timezone (Flutter default)
			parsed, err = time.Parse("2006-01-02T15:04:05.000", deadlineStr)

			if err != nil {
				// Try ISO8601 without milliseconds
				parsed, err = time.Parse("2006-01-02T15:04:05", deadlineStr)

				if err != nil {
					c.JSON(400, gin.H{
						"error": "invalid deadline format; expected RFC3339 or ISO8601",
					})
					return
				}
			}
		}

		newTask.Deadline = &parsed
	}

	// ---- SAVE TASK ----
	if err := mgm.Coll(newTask).Create(newTask); err != nil {
		c.JSON(500, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(201, gin.H{"task": newTask})
}

// GET /todo/date/:date
func GetTasksByDate(c *gin.Context) {
	userID := c.GetString("user_id")
	date := c.Param("date")

	var tasks []models.TodoTask
	err := mgm.Coll(&models.TodoTask{}).SimpleFind(
		&tasks,
		bson.M{"user_id": userID, "date": date},
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch tasks"})
		return
	}

	c.JSON(200, gin.H{"tasks": tasks})
}

// PATCH /todo/:id/toggle
func ToggleTask(c *gin.Context) {
	taskID := c.Param("id")

	var task models.TodoTask
	if err := mgm.Coll(&task).FindByID(taskID, &task); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	task.Done = !task.Done

	if err := mgm.Coll(&task).Update(&task); err != nil {
		c.JSON(500, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(200, gin.H{"done": task.Done})
}
