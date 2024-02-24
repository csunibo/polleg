package api

import (
	"encoding/json"
	"net/http"

	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
)

type Document struct {
	ID        string     `json:"id"`
	Questions []Question `json:"questions"`
}

type Coord struct {
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}

type PutDocumentRequest struct {
	ID     string  `json:"id"`
	Coords []Coord `json:"coords"`
}

// Insert a new document with all the questions
func PutDocumentHandler(res http.ResponseWriter, req *http.Request) {
	// Check method PUT is used
	if req.Method != http.MethodPut {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()

	// decode data
	var data PutDocumentRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		_ = util.WriteError(res, http.StatusBadRequest, "couldn't decode body")
		return
	}

	// save questions
	var questions []Question
	for _, coord := range data.Coords {
		q := Question{
			Document: data.ID,
			Start:    coord.Start,
			End:      coord.End,
		}
		questions = append(questions, q)
	}

	if err := db.Save(questions).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "couldn't create questions")
		return
	}

	util.WriteJson(res, Document{
		ID:        data.ID,
		Questions: questions,
	})
}

// Given a document's ID, return all the questions
func GetDocumentHandler(res http.ResponseWriter, req *http.Request) {
	// Check method GET is used
	if req.Method != http.MethodGet {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	var questions []Question
	db := util.GetDb()
	docID := muxie.GetParam(res, "id")
	db.Where(Question{Document: docID}).Find(&questions)
	util.WriteJson(res, Document{
		ID:        docID,
		Questions: questions,
	})
}
