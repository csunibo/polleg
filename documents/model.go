package documents

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/csunibo/stackunibo/answers"
	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
	"gorm.io/gorm"
)

type Document struct {
	ID        string     `json:"id"`
	Questions []Question `json:"questions" gorm:"foreignKey:Document;references:ID"`
}

type Question struct {
	gorm.Model
	Document string           `json:"document"`
	Start    uint32           `json:"start"`
	End      uint32           `json:"end"`
	Answers  []answers.Answer `json:"answers" gorm:"foreignKey:Question;references:ID"`
}

type Coord struct {
	gorm.Model
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}

type DocReq struct {
	Document string  `json:"document"`
	Coords   []Coord `json:"coords"`
}

func Handler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		handlePut(res, req)
	default:
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func handlePut(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()

	// decode data
	var data DocReq
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		_ = util.WriteError(res, http.StatusBadRequest, "couldn't decode body")
		return
	}

	// save document
	doc := Document{
		ID: data.Document,
	}
	if err := db.Save(doc).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "couldn't create doc")
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

func Get(res http.ResponseWriter, req *http.Request) {
	var docs Document
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")

	db.First(&docs, docId)
	fmt.Print(docs)
	util.WriteJson(res, docs)
}
