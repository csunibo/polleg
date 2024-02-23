package answers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/csunibo/stackunibo/auth"
	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
)

type AnswerObj struct {
	Question uint   `json:"question"`
	Parent   uint   `json:"parent"`
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
	user := auth.GetUser(req)
	// Declare a new Person struct.
	var ans AnswerObj

	err := json.NewDecoder(req.Body).Decode(&ans)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "decode error")
		return
	}

	fmt.Println(user, ans.Question)
	var quest Question
	util.GetDb().First(&quest, ans.Question)

	// var doc Document
	// util.GetDb().First(&quest, quest.)
	//fmt.Println(doc)

	/*
		util.GetDb().Create(&Answer{
			Document: ans.Document,
			Question: ans.Question,

			Parent:  nil,
			User:    user.Username,
			Content: ans.Content,
		})
	*/
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	answer := muxie.GetParam(res, "id")
	slog.Info("Fetching answers", "doc", answer)
}
