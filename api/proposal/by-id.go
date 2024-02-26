package proposal

import (
	"net/http"
	"strconv"

	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
	"golang.org/x/exp/slog"
)

func ProposalByIdHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		deleteProposalByIdHandler(res, req)
	case http.MethodGet:
		getProposalByIdHandler(res, req)
	default:
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
	}
}

func getProposalByIdHandler(res http.ResponseWriter, req *http.Request) {
	// Check method GET is used
	if req.Method != http.MethodGet {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	proposalID := muxie.GetParam(res, "id")
	propID, err := strconv.ParseUint(proposalID, 10, 0)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "invalid answer id")
		return
	}
	var props Proposal
	if err := db.Where(Proposal{ID: propID}).Take(&props).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "Not found")
		return
	}
	util.WriteJson(res, props)
}

func deleteProposalByIdHandler(res http.ResponseWriter, req *http.Request) {
	if !auth.GetAdmin(req) {
		util.WriteError(res, http.StatusForbidden, "you are not admin")
		return
	}
	db := util.GetDb()
	proposalID := muxie.GetParam(res, "id")
	propID, err := strconv.ParseUint(proposalID, 10, 0)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "invalid answer id")
		return
	}

	proposal := Proposal{ID: propID}

	if err := db.Delete(&Proposal{}, propID).Error; err != nil {
		util.WriteError(res, http.StatusInternalServerError, "db query failed")
		return
	}

	if err := util.WriteJson(res, proposal); err != nil {
		slog.Error("error while serializing the proposal", "err", err)
	}
}
