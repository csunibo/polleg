package answers

import (
	"encoding/json"
	"net/http"

	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
)

func DocHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		handleDocPut(res, req)
	default:
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func handleDocPut(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()

	// decode data
	var data DocReq
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		_ = util.WriteError(res, http.StatusBadRequest, "couldn't decode body")
		return
	}

	// save questions
	var questions []Question
	for _, coord := range data.Coords {
		q := Question{
			Document: data.Document,
			Start:    coord.Start,
			End:      coord.End,
		}
		questions = append(questions, q)
	}

	if err := db.Save(questions).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "couldn't create questions")
		return
	}
}

func GetQuestionsById(res http.ResponseWriter, req *http.Request) {
	var docs Question
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")
	db.Where("id = ?", docId).Find(&docs)
	util.WriteJson(res, docs)
}

func GetQuestionsByDoc(res http.ResponseWriter, req *http.Request) {
	var docs []Question
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")
	db.Where("document = ?", docId).Find(&docs)
	util.WriteJson(res, docs)
}
