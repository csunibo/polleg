package documents

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Document Document `json:"document" gorm:"foreignKey:Document;references:ID"`
	Coord    uint32   `json:"coord"`
}

type DocReq struct {
	Document string   `json:"document"`
	Coords   []uint32 `json:"coords"`
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

	var data DocReq
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		_ = util.WriteError(res, http.StatusBadRequest, "couldn't decode body")
		return
	}

	// save questions
	var questions []Question
	for _, coord := range data.Coords {
		question := Question{
			Document: Document{ID: data.Document},
			Coord:    coord,
		}
		questions = append(questions, question)
		if err := db.Save(question).Error; err != nil {
			_ = util.WriteError(res, http.StatusBadRequest, "couldn't create doc")
			return
		}
	}

	// save document
	doc := Document{
		ID:        data.Document,
		Questions: questions,
	}

	db.Save(doc)
}

func Get(res http.ResponseWriter, req *http.Request) {
	var docs Document
	db := util.GetDb()
	docId := muxie.GetParam(res, "id")

	db.First(&docs, docId)
	fmt.Print(docs)
	util.WriteJson(res, docs)
}
