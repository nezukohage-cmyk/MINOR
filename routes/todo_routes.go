package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(r *gin.Engine) {
	todo := r.Group("/todo", middleware.AuthRequired())

	todo.GET("/date/:date", controllers.GetTasksByDate)
	todo.PATCH("/:id/toggle", controllers.ToggleTask)

	todo.GET("/:date", controllers.GetTasksByDate)
	todo.POST("/:id/toggle", controllers.ToggleTask)
}
