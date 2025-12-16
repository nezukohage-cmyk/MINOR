package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func ClusterRoutes(r *gin.Engine) {
	cluster := r.Group("/clusters")
	cluster.Use(middleware.AuthRequired())

	cluster.POST("/create", controllers.CreateCluster)
	cluster.GET("", controllers.ListClusters)
	cluster.GET("/:id", controllers.GetCluster)

	cluster.POST("/upload", controllers.UploadClusterNote)
	cluster.GET("/:id/notes", controllers.ListClusterNotes)
	cluster.DELETE("/note/:id", controllers.DeleteClusterNote)
}
