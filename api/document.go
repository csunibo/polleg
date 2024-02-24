package api

import (
	"encoding/json"
	"net/http"

	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
)

// Insert a new document with all the questions
func PutDocument(res http.ResponseWriter, req *http.Request) {
	// Check method PUT is used
	if req.Method != http.MethodPut {
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
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

	util.WriteJson(res, util.Res{Res: "OK"})
}

// Get a question by an ID
func GetQuestionsById(res http.ResponseWriter, req *http.Request) {
	var docs Question
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")
	db.Where("id = ?", docId).Find(&docs)
	util.WriteJson(res, docs)
}

// Given a document's ID, return all the question
func GetQuestionsByDoc(res http.ResponseWriter, req *http.Request) {
	var docs []Question
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")
	db.Where("document = ?", docId).Find(&docs)
	util.WriteJson(res, docs)
}
