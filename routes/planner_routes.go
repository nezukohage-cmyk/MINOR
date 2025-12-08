package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func PlannerRoutes(r *gin.Engine) {

	pl := r.Group("/planner", middleware.AuthRequired())

	pl.POST("/generate", controllers.GeneratePlan)
	// 	pl.GET("/list", controllers.ListPlans)
	// 	pl.GET("/:id", controllers.GetPlan)
	// 	pl.DELETE("/:id", controllers.DeletePlan)
}
