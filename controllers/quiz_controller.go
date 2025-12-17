package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StartQuizRequest struct {
	Subjects []string            `json:"subjects"`
	Topics   map[string][]string `json:"topics"`
	Count    map[string]int      `json:"count"`
}

type SubmitAnswer struct {
	QuestionID string `json:"question_id"`
	Selected   string `json:"selected"`
}

type SubmitQuizRequest struct {
	QuizID  string         `json:"quiz_id"`
	Answers []SubmitAnswer `json:"answers"`
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToUpper(s))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func StartQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req StartQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.Subjects) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no subjects selected"})
		return
	}

	rand.Seed(time.Now().UnixNano())

	finalIDs := []string{}
	meta := map[string]map[string]int{}

	for _, subj := range req.Subjects {

		subj = normalize(subj)
		requested := req.Count[subj]

		if requested <= 0 {
			continue
		}

		// ✅ MATCH DB SCHEMA
		filter := bson.M{
			"subject_id": subj,
		}

		if topics, ok := req.Topics[subj]; ok && len(topics) > 0 {
			normTopics := []string{}
			for _, t := range topics {
				normTopics = append(normTopics, normalize(t))
			}
			filter["topic_ids"] = bson.M{"$in": normTopics}
		}

		var questions []models.Question
		if err := mgm.Coll(&models.Question{}).SimpleFind(&questions, filter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		ids := []string{}
		for _, q := range questions {
			ids = append(ids, q.ID.Hex())
		}

		meta[subj] = map[string]int{
			"requested": requested,
			"available": len(ids),
		}

		if len(ids) == 0 {
			continue
		}

		rand.Shuffle(len(ids), func(i, j int) {
			ids[i], ids[j] = ids[j], ids[i]
		})

		take := requested
		if take > len(ids) {
			take = len(ids)
		}

		finalIDs = append(finalIDs, ids[:take]...)
	}

	if len(finalIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "insufficient questions",
			"meta":  meta,
		})
		return
	}

	session := &models.QuizSession{
		UserID:         userID,
		Subjects:       req.Subjects,
		QuestionIDs:    finalIDs,
		RequestedCount: req.Count,
		TotalQuestions: len(finalIDs),
		Score:          0,
		StartedAt:      time.Now().UTC(),
	}

	if err := mgm.Coll(session).Create(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create quiz"})
		return
	}

	pageIDs := finalIDs[:min(10, len(finalIDs))]
	objIDs := []primitive.ObjectID{}

	for _, id := range pageIDs {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	var qs []models.Question
	_ = mgm.Coll(&models.Question{}).SimpleFind(&qs, bson.M{
		"_id": bson.M{"$in": objIDs},
	})

	out := []gin.H{}
	for _, q := range qs {
		out = append(out, gin.H{
			"id":      q.ID.Hex(),
			"text":    q.Text,
			"options": q.Options,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"quiz_id":         session.ID.Hex(),
		"total_questions": session.TotalQuestions,
		"questions":       out,
		"meta":            meta,
	})
}

func SubmitQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	session := &models.QuizSession{}
	if err := mgm.Coll(session).FindByID(req.QuizID, session); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	if session.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ansMap := map[string]string{}
	for _, a := range req.Answers {
		ansMap[a.QuestionID] = strings.TrimSpace(strings.ToUpper(a.Selected))
	}

	objIDs := []primitive.ObjectID{}
	for _, id := range session.QuestionIDs {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	var qs []models.Question
	_ = mgm.Coll(&models.Question{}).SimpleFind(&qs, bson.M{
		"_id": bson.M{"$in": objIDs},
	})

	correct := 0
	details := []gin.H{}

	for _, q := range qs {
		userAns := ansMap[q.ID.Hex()]
		isCorrect := userAns == strings.ToUpper(strings.TrimSpace(q.Answer))
		if isCorrect {
			correct++
		}

		details = append(details, gin.H{
			"question_id": q.ID.Hex(),
			"text":        q.Text,
			"selected":    userAns,
			"correct":     q.Answer,
			"is_correct":  isCorrect,
		})
	}

	now := time.Now().UTC()
	session.Score = correct
	session.SubmittedAt = &now
	if err := mgm.Coll(session).Update(session); err != nil {
		fmt.Println("quiz update failed:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz_id":  session.ID.Hex(),
		"total":    session.TotalQuestions,
		"correct":  correct,
		"wrong":    session.TotalQuestions - correct,
		"details":  details,
		"finished": true,
	})
}
