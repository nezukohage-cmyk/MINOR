package controllers

import (
	"net/http"
	"strings"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /communities/search?q=sdm
func SearchCommunities(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query 'q' is required"})
		return
	}

	regex := bson.M{"$regex": q, "$options": "i"}

	filter := bson.M{
		"$or": []bson.M{
			{"name": regex},
			{"region": regex},
		},
	}

	var communities []models.Community
	err := mgm.Coll(&models.Community{}).SimpleFind(&communities, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search communities"})
		return
	}

	if len(communities) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "college doesn't exist",
			"results":    []models.Community{},
			"totalFound": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":    communities,
		"totalFound": len(communities),
	})
}

// POST /communities/join/:name
func JoinCommunity(c *gin.Context) {
	userID := c.GetString("user_id")

	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "community name required"})
		return
	}

	// find community by name (case-insensitive)
	filter := bson.M{"name": bson.M{"$regex": "^" + name + "$", "$options": "i"}}

	var community models.Community
	err := mgm.Coll(&models.Community{}).First(filter, &community)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	// load user
	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	// add community to user's joined list
	communityID := community.ID.Hex()
	if !stringInSlice(communityID, user.CommunityIDs) {
		user.CommunityIDs = append(user.CommunityIDs, communityID)
		if err := mgm.Coll(user).Update(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user communities"})
			return
		}
	}

	// add user to community members list
	if !stringInSlice(userID, community.Members) {
		community.Members = append(community.Members, userID)
		if err := mgm.Coll(&community).Update(&community); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update community members"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "joined community",
		"community": community,
	})
}

func MyCommunities(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	if len(user.CommunityIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"communities": []models.Community{},
			"total":       0,
		})
		return
	}

	// convert string IDs → ObjectIDs
	objIDs := make([]primitive.ObjectID, 0, len(user.CommunityIDs))
	for _, id := range user.CommunityIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	var communities []models.Community
	if err := mgm.Coll(&models.Community{}).SimpleFind(&communities, filter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch communities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"communities": communities,
		"total":       len(communities),
	})
}

// POST /communities/create
func CreateCommunity(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Name   string `json:"name"`
		Region string `json:"region"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	// check name uniqueness
	count, _ := mgm.Coll(&models.Community{}).CountDocuments(
		mgm.Ctx(),
		bson.M{"name": bson.M{"$regex": "^" + req.Name + "$", "$options": "i"}},
	)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "community already exists"})
		return
	}

	community := &models.Community{
		Name:      req.Name,
		Region:    req.Region,
		Members:   []string{},
		CreatedBy: userID,
	}

	if err := mgm.Coll(community).Create(community); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create community"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "community created",
		"community": community,
	})
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
func LeaveCommunity(c *gin.Context) {
	userID := c.GetString("user_id")
	name := strings.TrimSpace(c.Param("name"))

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "community name required"})
		return
	}

	// find community
	filter := bson.M{"name": bson.M{"$regex": "^" + name + "$", "$options": "i"}}
	var community models.Community

	err := mgm.Coll(&models.Community{}).First(filter, &community)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	// load user
	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	// remove community from user list
	user.CommunityIDs = removeString(user.CommunityIDs, community.ID.Hex())
	if err := mgm.Coll(user).Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user membership"})
		return
	}

	// remove user from community members
	community.Members = removeString(community.Members, userID)
	if err := mgm.Coll(&community).Update(&community); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update community members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "left community",
		"community": community,
	})
}
func removeString(list []string, value string) []string {
	newList := []string{}
	for _, item := range list {
		if item != value {
			newList = append(newList, item)
		}
	}
	return newList
}
