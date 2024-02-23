package answers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/csunibo/stackunibo/documents"
	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
	"gorm.io/gorm"
)

type AnswerObj struct {
	Document  string `json:"document"`
	Question  uint   `json:"question"`
	Parent    uint   `json:"parent"`
	User      string `json:"user"`
	Content   string `json:"content"`
	Upvotes   uint32 `json:"upvotes"`
	Downvotes uint32 `json:"downvotes"`
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
	// Declare a new Person struct.
	var ans AnswerObj

	err := json.NewDecoder(req.Body).Decode(&ans)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "specify a redirect_uri url param")
		return
	}

	// Do something with the Person struct...

	util.Get().Create(&Answer{
		Document: documents.Document{ID: ans.Document},
		Question: documents.Question{Model: gorm.Model{ID: ans.Question}},

		Parent:    nil,
		User:      ans.User,
		Content:   ans.Content,
		Upvotes:   ans.Upvotes,
		Downvotes: ans.Downvotes,
	})
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	answer := muxie.GetParam(res, "id")
	slog.Info("Fetching answers", "doc", answer)
}
