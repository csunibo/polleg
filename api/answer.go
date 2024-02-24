package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
	"golang.org/x/exp/slog"
)

type PutAnswerRequest struct {
	Question uint   `json:"question"`
	Parent   *uint  `json:"parent"`
	Content  string `json:"content"`
}

var (
	VOTES_QUERY = fmt.Sprintf(`
  select   votes.answer,
           count(case votes.vote when %d then 1 else null end) as upvotes,
           count(case votes.vote when %d then 1 else null end) as downvotes
  from     votes
  group by answer
`, VoteUp, VoteDown)
	ANSWERS_QUERY = fmt.Sprintf(`
  select   *
  from     answers
  full join     (%s) on answer = answers.id
  where    answers.parent is NULL and answers.question = ? 
`, VOTES_QUERY)
	REPLIES_QUERY = `
  select   *
  from     answers
  where    answers.parent IN ?
`
)

// Insert a new answer under a question
func PutAnswerHandler(res http.ResponseWriter, req *http.Request) {
	// Check method PUT is used
	if req.Method != http.MethodPut {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	user := auth.GetUser(req)

	var ans PutAnswerRequest
	err := json.NewDecoder(req.Body).Decode(&ans)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, fmt.Sprintf("decode error: %v", err))
		return
	}

	var quest Question
	if err := db.First(&quest, ans.Question).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "the referenced question does not exist")
		return
	}

	if ans.Parent != nil {
		var Parent Answer
		if err = db.First(&Parent, ans.Parent).Error; err != nil {
			util.WriteError(res, http.StatusBadRequest, "the referenced parent does not exist")
			return
		}
		if Parent.Question != quest.ID {
			util.WriteError(res, http.StatusBadRequest, "mismatch between parent question and this question")
			return
		}
	}

	// TODO: upvotes and downvotes should really be just the result of a
	// COUNT() aggregator on the votes table
	answer := Answer{
		Question:  ans.Question,
		Parent:    ans.Parent,
		User:      user.Username,
		Content:   ans.Content,
		Upvotes:   0,
		Downvotes: 0,
	}
	err = db.Create(&answer).Error
	if err != nil {
		slog.Error("error while creating the answer", "answer", answer, "err", err)
		util.WriteError(res, http.StatusBadRequest, "could not insert the answer")
		return
	}

	if err = util.WriteJson(res, answer); err != nil {
		slog.Error("error while serializing the answer", "err", err)
	}
}

// Given a question ID, return the question and all its answers
func GetQuestionHandler(res http.ResponseWriter, req *http.Request) {
	// Check method GET is used
	if req.Method != http.MethodGet {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	rawQID := muxie.GetParam(res, "id")
	qID, err := strconv.ParseUint(rawQID, 10, 0)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "invalid question id")
		return
	}

	var question Question
	if err := db.First(&question, uint(qID)).Error; err != nil {
		slog.Error("question not found", "err", err)
		util.WriteError(res, http.StatusInternalServerError, "question not found")
		return
	}
	var answers []Answer
	fmt.Println(ANSWERS_QUERY)
	if err := db.Raw(ANSWERS_QUERY, question.ID).Scan(&answers).Error; err != nil {
		slog.Error("could not fetch answers", "err", err)
		util.WriteError(res, http.StatusInternalServerError, "could not fetch answers")
		return
	}
	answersIDs := []uint{}
	answersIndex := map[uint]int{}
	for i, answer := range answers {
		answersIDs = append(answersIDs, answer.ID)
		answersIndex[answer.ID] = i
	}
	var replies []Answer
	if err := db.Raw(REPLIES_QUERY, answersIDs).Scan(&replies).Error; err != nil {
		slog.Error("could not fetch replies", "err", err)
		util.WriteError(res, http.StatusInternalServerError, "could not fetch replies")
		return
	}
	for _, reply := range replies {
		index := answersIndex[*reply.Parent]
		answers[index].Replies = append(answers[index].Replies, reply)
	}

	question.Answers = answers
	if err := util.WriteJson(res, question); err != nil {
		slog.Error("error while serializing the question", "err", err)
	}
}
