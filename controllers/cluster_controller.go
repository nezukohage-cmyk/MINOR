package controllers

import (
	"fmt"
	"lexxi/models"
	"lexxi/services"
	"net/http"

	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateCluster(c *gin.Context) {
	var body struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cluster := &models.Cluster{
		Name:        body.Name,
		Description: body.Description,
		Tags:        body.Tags,
	}

	if err := mgm.Coll(cluster).Create(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cluster"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cluster created",
		"data":    cluster,
	})
}

func ListClusters(c *gin.Context) {
	search := c.Query("q")

	filter := bson.M{}

	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": search, "$options": "i"}},
				{"tags": bson.M{"$regex": search, "$options": "i"}},
				{"description": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	var clusters []models.Cluster
	err := mgm.Coll(&models.Cluster{}).SimpleFind(&clusters, filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load clusters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": clusters})
}

func UploadClusterNote(c *gin.Context) {
	clusterID := c.PostForm("cluster_id")
	title := c.PostForm("title")
	uploaderID := c.GetString("user_id")

	fmt.Println("UPLOAD RECEIVED FOR CLUSTER:", clusterID)

	if clusterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cluster_id required"})
		return
	}

	file, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF file required"})
		return
	}

	url, err := services.UploadFile(file, "clusters/"+clusterID)
	if err != nil {
		fmt.Println("UPLOAD FAILED:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary upload failed"})
		return
	}

	note := &models.ClusterNote{
		ClusterID:  clusterID,
		UploaderID: uploaderID,
		Title:      title,
		FileURL:    url,
	}

	if err := mgm.Coll(note).Create(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful", "data": note})
}

// //////////////////////////////////////////////////////////////////////////////////
// 4) LIST ALL NOTES IN A CLUSTER
// //////////////////////////////////////////////////////////////////////////////////
func ListClusterNotes(c *gin.Context) {
	clusterID := c.Param("id")

	var notes []models.ClusterNote
	err := mgm.Coll(&models.ClusterNote{}).SimpleFind(&notes, bson.M{"cluster_id": clusterID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}

// //////////////////////////////////////////////////////////////////////////////////
// 5) GET SINGLE CLUSTER DETAILS
// //////////////////////////////////////////////////////////////////////////////////
func GetCluster(c *gin.Context) {
	id := c.Param("id")

	cluster := &models.Cluster{}
	err := mgm.Coll(cluster).FindByID(id, cluster)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cluster})
}

// //////////////////////////////////////////////////////////////////////////////////
// 6) DELETE A NOTE (OPTIONAL)
// //////////////////////////////////////////////////////////////////////////////////
func DeleteClusterNote(c *gin.Context) {
	noteID := c.Param("id")
	userID := c.GetString("user_id")

	note := &models.ClusterNote{}
	err := mgm.Coll(note).FindByID(noteID, note)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Only uploader can delete
	if note.UploaderID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	if err := mgm.Coll(note).Delete(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
func DeleteCluster(c *gin.Context) {
	id := c.Param("id")

	cluster := &models.Cluster{}
	err := mgm.Coll(cluster).FindByID(id, cluster)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	if err := mgm.Coll(cluster).Delete(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	// Also delete notes belonging to it (optional)
	mgm.Coll(&models.ClusterNote{}).DeleteMany(mgm.Ctx(), bson.M{"cluster_id": id})

	c.JSON(http.StatusOK, gin.H{"message": "Cluster deleted"})
}
