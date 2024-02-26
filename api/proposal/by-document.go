package proposal

import (
	"fmt"
	"net/http"

	"github.com/csunibo/polleg/api"
	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
	"golang.org/x/exp/slog"
)

func ProposalByDocumentHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		deleteProposalByDocumentHandler(res, req)
	case http.MethodGet:
		getProposalByDocumentHandler(res, req)
	default:
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func getProposalByDocumentHandler(res http.ResponseWriter, req *http.Request) {
	// Check method GET is used
	if req.Method != http.MethodGet {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	docID := muxie.GetParam(res, "id")
	var questions []Proposal
	if err := db.Where(api.Question{Document: docID}).Find(&questions).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "db query failed")
		return
	}
	if len(questions) == 0 {
		util.WriteError(res, http.StatusInternalServerError, "Document not found")
		return
	}
	util.WriteJson(res, DocumentProposal{
		ID:        docID,
		Questions: questions,
	})
}

func deleteProposalByDocumentHandler(res http.ResponseWriter, req *http.Request) {
	if !auth.GetAdmin(req) {
		util.WriteError(res, http.StatusForbidden, "you are not admin")
		return
	}
	if req.Method != http.MethodDelete {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	docID := muxie.GetParam(res, "id")

	proposal := Proposal{Document: docID}

	if err := db.Where("document = ?", docID).Delete(&Proposal{}).Error; err != nil {
		fmt.Println()
		slog.Error("db query failed", "err", err)
		fmt.Println()
		util.WriteError(res, http.StatusInternalServerError, "db query failed")
		return
	}

	if err := util.WriteJson(res, proposal); err != nil {
		slog.Error("error while serializing the proposal", "err", err)
	}
}
