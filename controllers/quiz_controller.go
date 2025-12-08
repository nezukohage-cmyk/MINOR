// package controllers

// import (
// 	//"context"
// 	"math/rand"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"lexxi/models"

// 	"github.com/gin-gonic/gin"
// 	"github.com/kamva/mgm/v3"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// // helper: shuffle string slice
// func shuffleStrings(a []string) {
// 	rand.Seed(time.Now().UnixNano())
// 	for i := len(a) - 1; i > 0; i-- {
// 		j := rand.Intn(i + 1)
// 		a[i], a[j] = a[j], a[i]
// 	}
// }

// // Request payload to start a quiz
// type StartQuizReq struct {
// 	Subjects []string            `json:"subjects"`
// 	Topics   map[string][]string `json:"topics"`
// 	Count    map[string]int      `json:"count"`
// }

// type QuizQuestionResp struct {
// 	ID      string   `json:"id"`
// 	Text    string   `json:"text"`
// 	Options []string `json:"options"`
// 	Subject string   `json:"subject"`
// 	Topics  []string `json:"topics"`
// }

// func StartQuiz(c *gin.Context) {
// 	userID := c.GetString("user_id")
// 	if userID == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var req StartQuizReq
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
// 		return
// 	}

// 	var finalIDs []string
// 	meta := map[string]map[string]int{}
// 	for _, subject := range req.Subjects {
// 		need := 0
// 		if v, ok := req.Count[subject]; ok {
// 			need = v
// 		} else {
// 			need = 10
// 		}

// 		filter := bson.M{"subject_id": subject}
// 		if ts, ok := req.Topics[subject]; ok && len(ts) > 0 {
// 			filter["topic_ids"] = bson.M{"$in": ts}
// 		}

// 		var qs []models.Question
// 		if err := mgm.Coll(&models.Question{}).SimpleFind(&qs, filter); err != nil {
// 			// skip subject on error (but record 0 available)
// 			meta[subject] = map[string]int{"requested": need, "available": 0}
// 			continue
// 		}

// 		// collect IDs
// 		ids := make([]string, 0, len(qs))
// 		for _, q := range qs {
// 			ids = append(ids, q.ID.Hex())
// 		}

// 		// shuffle and take needed
// 		shuffleStrings(ids)
// 		take := need
// 		if take > len(ids) {
// 			take = len(ids)
// 		}
// 		if take > 0 {
// 			finalIDs = append(finalIDs, ids[:take]...)
// 		}
// 		meta[subject] = map[string]int{"requested": need, "available": len(ids)}
// 	}

// 	if len(finalIDs) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "no questions available for selected subjects/topics", "meta": meta})
// 		return
// 	}

// 	// shuffle overall final set to mix subjects
// 	shuffleStrings(finalIDs)

// 	// create session
// 	session := &models.QuizSession{
// 		UserID:         userID,
// 		QuestionIDs:    finalIDs,
// 		RequestedCount: map[string]int{},
// 		Meta:           map[string]any{},
// 		StartedAt:      time.Now().UTC(),
// 	}
// 	_ = mgm.Coll(session).Create(session)

// 	// return first page (page size 10)
// 	page := 1
// 	limit := 10
// 	start := (page - 1) * limit
// 	end := start + limit
// 	if end > len(finalIDs) {
// 		end = len(finalIDs)
// 	}
// 	pageIDs := finalIDs[start:end]

// 	// fetch actual question docs for response
// 	var questions []models.Question
// 	if err := mgm.Coll(&models.Question{}).SimpleFind(&questions, bson.M{"_id": bson.M{"$in": pageIDs}}); err != nil {
// 		// If fetch fails, still return session id and meta
// 		c.JSON(http.StatusCreated, gin.H{
// 			"quiz_id":         session.ID.Hex(),
// 			"total_questions": len(finalIDs),
// 			"meta":            meta,
// 		})
// 		return
// 	}

// 	qMap := make(map[string]models.Question)
// 	for _, q := range questions {
// 		qMap[q.ID.Hex()] = q
// 	}
// 	out := make([]QuizQuestionResp, 0, len(pageIDs))
// 	for _, id := range pageIDs {
// 		if q, ok := qMap[id]; ok {
// 			out = append(out, QuizQuestionResp{
// 				ID:      q.ID.Hex(),
// 				Text:    q.Text,
// 				Options: q.Options,
// 				Subject: q.SubjectID,
// 				Topics:  q.TopicIDs,
// 			})
// 		}
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"quiz_id":         session.ID.Hex(),
// 		"meta":            meta,
// 		"total_questions": len(finalIDs),
// 		"page":            1,
// 		"limit":           limit,
// 		"questions":       out,
// 	})
// }

// // GET /quiz/:quiz_id/questions?page=1
// func QuizPage(c *gin.Context) {
// 	userID := c.GetString("user_id")
// 	if userID == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	quizID := c.Param("quiz_id")
// 	if quizID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz id required"})
// 		return
// 	}

// 	// load session
// 	session := &models.QuizSession{}
// 	if err := mgm.Coll(session).FindByID(quizID, session); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
// 		return
// 	}

// 	// ensure ownership
// 	if session.UserID != userID {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "not your quiz"})
// 		return
// 	}

// 	// pagination
// 	page := 1
// 	limit := 10
// 	if p := c.Query("page"); p != "" {
// 		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
// 			page = parsed
// 		}
// 	}
// 	start := (page - 1) * limit
// 	if start >= len(session.QuestionIDs) {
// 		c.JSON(http.StatusOK, gin.H{"questions": []QuizQuestionResp{}, "page": page, "limit": limit})
// 		return
// 	}
// 	end := start + limit
// 	if end > len(session.QuestionIDs) {
// 		end = len(session.QuestionIDs)
// 	}
// 	pageIDs := session.QuestionIDs[start:end]

// 	// fetch questions
// 	var questions []models.Question
// 	if err := mgm.Coll(&models.Question{}).SimpleFind(&questions, bson.M{"_id": bson.M{"$in": pageIDs}}); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load questions"})
// 		return
// 	}

// 	// order preserved
// 	qMap := make(map[string]models.Question)
// 	for _, q := range questions {
// 		qMap[q.ID.Hex()] = q
// 	}
// 	out := []QuizQuestionResp{}
// 	for _, id := range pageIDs {
// 		if q, ok := qMap[id]; ok {
// 			out = append(out, QuizQuestionResp{
// 				ID:      q.ID.Hex(),
// 				Text:    q.Text,
// 				Options: q.Options,
// 				Subject: q.SubjectID,
// 				Topics:  q.TopicIDs,
// 			})
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"quiz_id":   session.ID.Hex(),
// 		"page":      page,
// 		"limit":     limit,
// 		"questions": out,
// 		"total":     len(session.QuestionIDs),
// 	})
// }

// type SubmitAnswer struct {
// 	QuestionID string `json:"question_id"`
// 	Selected   string `json:"selected"`
// }

// type SubmitReq struct {
// 	QuizID       string         `json:"quiz_id"`
// 	Answers      []SubmitAnswer `json:"answers"`
// 	TimeTakenSec int            `json:"time_taken_sec"`
// }

// // POST /quiz/submit
// func SubmitQuiz(c *gin.Context) {
// 	userID := c.GetString("user_id")
// 	if userID == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	var req SubmitReq
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
// 		return
// 	}

// 	session := &models.QuizSession{}
// 	if err := mgm.Coll(session).FindByID(req.QuizID, session); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
// 		return
// 	}
// 	if session.UserID != userID {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "not your quiz"})
// 		return
// 	}

// 	// Build map of answers sent
// 	ansMap := make(map[string]string)
// 	for _, a := range req.Answers {
// 		ansMap[a.QuestionID] = strings.TrimSpace(strings.ToLower(a.Selected))
// 	}

// 	// fetch all questions in session
// 	var qdocs []models.Question
// 	if err := mgm.Coll(&models.Question{}).SimpleFind(&qdocs, bson.M{"_id": bson.M{"$in": session.QuestionIDs}}); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load questions"})
// 		return
// 	}
// 	// map by id
// 	qMap := make(map[string]models.Question)
// 	for _, q := range qdocs {
// 		qMap[q.ID.Hex()] = q
// 	}

// 	total := len(session.QuestionIDs)
// 	correct := 0
// 	byTopic := map[string]int{}
// 	for _, qid := range session.QuestionIDs {
// 		q, ok := qMap[qid]
// 		if !ok {
// 			continue
// 		}
// 		userAns, answered := ansMap[qid]
// 		correctAns := strings.TrimSpace(strings.ToLower(q.Answer))
// 		if answered && userAns == correctAns {
// 			correct++
// 		} else {

// 			for _, t := range q.TopicIDs {
// 				byTopic[t] = byTopic[t] + 1
// 			}
// 		}
// 	}

// 	accuracy := 0.0
// 	if total > 0 {
// 		accuracy = (float64(correct) / float64(total)) * 100.0
// 	}

// 	attempt := &models.QuizAttempt{
// 		UserID:       userID,
// 		QuizSession:  session.ID.Hex(),
// 		Total:        total,
// 		Correct:      correct,
// 		Accuracy:     accuracy,
// 		TimeTakenSec: req.TimeTakenSec,
// 		ByTopic:      map[string]any{},
// 		CreatedAt:    time.Now().UTC(),
// 	}

// 	byTopicAny := map[string]any{}
// 	for k, v := range byTopic {
// 		byTopicAny[k] = v
// 	}
// 	attempt.ByTopic = byTopicAny
// 	_ = mgm.Coll(attempt).Create(attempt)

// 	now := time.Now().UTC()
// 	session.SubmittedAt = &now
// 	session.TimeTakenSec = req.TimeTakenSec
// 	_ = mgm.Coll(session).Update(session)

// 	weak := []string{}
// 	for k := range byTopicAny {
// 		weak = append(weak, k)
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"score":       correct,
// 		"total":       total,
// 		"accuracy":    accuracy,
// 		"weak_topics": weak,
// 	})
// }

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
	Subjects []string            `json:"subjects"` // e.g. ["DSA","OS"]
	Topics   map[string][]string `json:"topics"`   // optional: { "DSA": ["Linked List"], "OS": ["Paging"] }
	Count    map[string]int      `json:"count"`    // requested counts per subject, e.g. {"DSA":10, "OS":5}
	// optional: page size can be added later
}

// Answer payload for submitting a quiz
type SubmitAnswer struct {
	QuestionID string `json:"question_id"`
	Selected   string `json:"selected"`
}

type SubmitQuizRequest struct {
	QuizID  string         `json:"quiz_id"`
	Answers []SubmitAnswer `json:"answers"`
}

// GET helper: returns first N of a slice
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// normalize tag strings for consistent matching
func normalizeTag(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// POST /quiz/start
func StartQuiz(c *gin.Context) {
	var req StartQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if len(req.Subjects) == 0 || len(req.Count) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjects and count required"})
		return
	}

	// normalize inputs
	for i := range req.Subjects {
		req.Subjects[i] = normalizeTag(req.Subjects[i])
	}
	normalizedTopics := make(map[string][]string)
	for subj, tlist := range req.Topics {
		ns := normalizeTag(subj)
		for _, t := range tlist {
			normalizedTopics[ns] = append(normalizedTopics[ns], normalizeTag(t))
		}
	}

	// for each subject, gather candidate question IDs
	bySubjectIDs := make(map[string][]string)
	meta := make(map[string]map[string]int) // meta[subject] = {"available":n, "requested":x}

	for _, subj := range req.Subjects {
		requested := req.Count[subj]
		if requested == 0 {
			// maybe user passed lowercase key; try normalized access
			requested = req.Count[strings.ToLower(subj)]
			if requested == 0 {
				// fallback: if not specified, skip
				requested = 0
			}
		}

		// build filter: subject match and (if given) topics $in
		filter := bson.M{"subject_id": subj}

		if topicsForSubj, ok := normalizedTopics[subj]; ok && len(topicsForSubj) > 0 {
			// query topic ids case-insensitively by matching uppercased stored values as well --
			// we assume topic IDs were stored normalized. If not, consider indexing or storing normalized versions.
			filter["topic_ids"] = bson.M{"$in": topicsForSubj}
		}

		var qdocs []models.Question
		err := mgm.Coll(&models.Question{}).SimpleFind(&qdocs, filter)
		if err != nil {
			// if anything goes wrong, treat as zero available
			bySubjectIDs[subj] = []string{}
			meta[subj] = map[string]int{"available": 0, "requested": requested}
			continue
		}

		// collect IDs (strings)
		ids := make([]string, 0, len(qdocs))
		for _, q := range qdocs {
			ids = append(ids, q.ID.Hex())
		}

		meta[subj] = map[string]int{"available": len(ids), "requested": requested}
		bySubjectIDs[subj] = ids
	}

	// choose questions per subject
	finalIDs := []string{}
	rand.Seed(time.Now().UnixNano())
	for subj, ids := range bySubjectIDs {
		reqCount := req.Count[subj]
		// safety: if request omitted or zero, skip
		if reqCount <= 0 {
			continue
		}

		// shuffle and take up to requested
		if len(ids) <= reqCount {
			finalIDs = append(finalIDs, ids...)
		} else {
			// shuffle a copy
			tmp := make([]string, len(ids))
			copy(tmp, ids)
			rand.Shuffle(len(tmp), func(i, j int) { tmp[i], tmp[j] = tmp[j], tmp[i] })
			finalIDs = append(finalIDs, tmp[:reqCount]...)
		}
	}

	if len(finalIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "no questions available for selected subjects/topics",
			"meta":  meta,
		})
		return
	}

	// create session and save meta + requested counts
	session := &models.QuizSession{
		UserID:         userID,
		QuestionIDs:    finalIDs,
		RequestedCount: req.Count,
		Meta:           map[string]interface{}{},
		StartedAt:      time.Now().UTC(),
	}
	// store meta as-is (map[string]map[string]int) if session.Meta is interface-based
	for k, v := range meta {
		session.Meta[k] = v
	}

	if err := mgm.Coll(session).Create(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create quiz session"})
		return
	}

	// prepare first page of questions (page size = 10)
	pageSize := 10
	pageIDs := finalIDs[:min(pageSize, len(finalIDs))]

	// convert pageIDs to ObjectIDs for query
	objIDs := make([]primitive.ObjectID, 0, len(pageIDs))
	for _, s := range pageIDs {
		if oid, err := primitive.ObjectIDFromHex(s); err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	var questions []models.Question
	if err := mgm.Coll(&models.Question{}).SimpleFind(&questions, bson.M{"_id": bson.M{"$in": objIDs}}); err != nil {
		// can't fetch docs — still return session id and meta
		c.JSON(http.StatusCreated, gin.H{
			"quiz_id":         session.ID.Hex(),
			"total_questions": len(finalIDs),
			"meta":            meta,
		})
		return
	}

	// return session id and first page of questions
	c.JSON(http.StatusCreated, gin.H{
		"quiz_id":         session.ID.Hex(),
		"total_questions": len(finalIDs),
		"meta":            meta,
		"questions":       questions,
	})
}

// POST /quiz/submit
func SubmitQuiz(c *gin.Context) {
	var req SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if req.QuizID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id required"})
		return
	}

	// load session by id
	sessObjID, err := primitive.ObjectIDFromHex(req.QuizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	session := &models.QuizSession{}
	if err := mgm.Coll(session).FindByID(sessObjID, session); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	// ensure it's the same user (optional)
	if session.UserID != userID {
		// you can allow admins etc; but for now restrict
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// map answers for quick lookup
	ansMap := make(map[string]string)
	for _, a := range req.Answers {
		ansMap[a.QuestionID] = a.Selected
	}

	// convert session.QuestionIDs -> []primitive.ObjectID
	allObjIDs := make([]primitive.ObjectID, 0, len(session.QuestionIDs))
	for _, sid := range session.QuestionIDs {
		if oid, err := primitive.ObjectIDFromHex(sid); err == nil {
			allObjIDs = append(allObjIDs, oid)
		}
	}

	// fetch all question docs
	var qdocs []models.Question
	if err := mgm.Coll(&models.Question{}).SimpleFind(&qdocs, bson.M{"_id": bson.M{"$in": allObjIDs}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load questions"})
		return
	}

	// score
	total := len(qdocs)
	correct := 0
	details := make([]map[string]interface{}, 0, total)

	for _, q := range qdocs {
		userAns := ansMap[q.ID.Hex()]
		isCorrect := strings.TrimSpace(strings.ToUpper(userAns)) == strings.TrimSpace(strings.ToUpper(q.Answer))
		if isCorrect {
			correct++
		}
		details = append(details, map[string]interface{}{
			"question_id": q.ID.Hex(),
			"text":        q.Text,
			"selected":    userAns,
			"correct":     q.Answer,
			"is_correct":  isCorrect,
		})
	}

	// update session
	session.SubmittedAt = time.Now().UTC()
	session.Score = correct
	if err := mgm.Coll(session).Update(session); err != nil {
		// non-fatal; return result anyway
		fmt.Println("warning: failed to update quiz session with score:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"correct":  correct,
		"wrong":    total - correct,
		"score":    correct,
		"details":  details,
		"quiz_id":  req.QuizID,
		"finished": true,
	})
}
