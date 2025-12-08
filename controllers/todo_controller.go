package controllers

import (
	"lexxi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

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
