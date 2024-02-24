package answers

import (
	"encoding/json"
	"net/http"

	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
)

func NewDocument(res http.ResponseWriter, req *http.Request) {
	// Check method put is used
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

func GetAnswerOfQuestion(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")

	var ans []Answer
	db.Where("question = ?", docId).Find(&ans)

	util.WriteJson(res, ans)
}
