package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
	"gorm.io/gorm/clause"
)

type VoteValue int8

const (
	VoteUp   VoteValue = 1
	VoteNone VoteValue = 0
	VoteDown VoteValue = -1
)

type PutVoteRequest struct {
	Vote VoteValue `json:"vote"`
}

// Insert a vote on a answer
func PostVote(res http.ResponseWriter, req *http.Request) {
	// Check method POST is used
	if req.Method != http.MethodPost {
		util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}
	db := util.GetDb()
	user := auth.GetUser(req)

	rawAnsID := muxie.GetParam(res, "id")
	ansID, err := strconv.ParseUint(rawAnsID, 10, 0)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "invalid answer id")
		return
	}

	var v PutVoteRequest
	err = json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, fmt.Sprintf("decode error: %v", err))
		return
	}

	var ans Answer
	if err = db.First(&ans, ansID).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "the referenced answer does not exist")
		return
	}
	if ans.Parent != nil {
		util.WriteError(res, http.StatusBadRequest, "cannot vote a reply to an answer")
		return
	}

	vote := Vote{
		Answer: ans.ID,
		User:   user.Username,
		Vote:   int8(v.Vote),
	}
	if v.Vote == VoteUp || v.Vote == VoteDown {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "answer"}, {Name: "user"}},
			DoUpdates: clause.AssignmentColumns([]string{"vote"}),
		}).Create(&vote).Error
		if err != nil {
			util.WriteError(res, http.StatusInternalServerError, "could not update your vote")
			return
		}
	} else if v.Vote == VoteNone {
		if err := db.Delete(Vote{Answer: ans.ID, User: user.Username}).Error; err != nil {
			util.WriteError(res, http.StatusBadRequest, "could not delete the previous vote")
			return
		}
	} else {
		util.WriteError(res, http.StatusBadRequest, "the vote value must be either 1 or -1")
		return
	}

	if err = util.WriteJson(res, vote); err != nil {
		slog.Error("error while serializing the vote", "err", err)
	}
}
