package answers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/csunibo/stackunibo/auth"
	"github.com/csunibo/stackunibo/util"
)

type AnswerObj struct {
	Question uint   `json:"question"`
	Parent   *uint  `json:"parent"`
	Content  string `json:"content"`
}

func AnswerHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		handleGet(res, req)
	case http.MethodPut:
		handlePut(res, req)
	default:
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func handlePut(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()
	user := auth.GetUser(req)

	// Declare a new Person struct.
	var ans AnswerObj

	err := json.NewDecoder(req.Body).Decode(&ans)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "decode error")
		return
	}

	fmt.Println(user)
	var quest Question
	if err := db.First(&quest, ans.Question).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "no Question associated with request (or other Error)")
		return
	}

	if ans.Parent != nil {
		var Parent Answer
		if err = db.First(&Parent, ans.Parent).Error; err != nil {
			util.WriteError(res, http.StatusBadRequest, "parent is given but none found")
			return
		}
		if Parent.Question != quest.ID {
			util.WriteError(res, http.StatusBadRequest, "mismatch between parent question and this question")
			return
		}
	}

	util.GetDb().Create(&Answer{
		Question:  ans.Question,
		Parent:    ans.Parent,
		User:      user.Username,
		Content:   ans.Content,
		Upvotes:   0,
		Downvotes: 0,
	})

}

func handleGet(res http.ResponseWriter, req *http.Request) {
	// var quest []Question
	// db := util.GetDb()
	// docId := muxie.GetParam(res, "id")
	// db.Where("id = ?", docId).Find(&quest)
	// util.WriteJson(res, docs)
}
