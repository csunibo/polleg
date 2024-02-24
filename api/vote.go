package api

import (
	"encoding/json"
	"net/http"

	"github.com/csunibo/polleg/auth"
	"github.com/csunibo/polleg/util"
	"github.com/kataras/muxie"
)

type VoteObj struct {
	Vote int8 `json:"vote"`
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
	id := muxie.GetParam(res, "id")

	// Declare a new Person struct.
	var vote VoteObj
	err := json.NewDecoder(req.Body).Decode(&vote)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "decode error")
		return
	}

	if vote.Vote != 1 && vote.Vote != -1 {
		util.WriteError(res, http.StatusBadRequest, "body 'vote' must be either 1 or -1")
		return
	}

	var ans Answer
	if err = db.First(&ans, id).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "no question associated with request id")
		return
	}

	var voteRecord []Vote
	err = db.Where("votes.user = ? AND votes.answer = ?", user.Username, ans.ID).Find(&voteRecord).Error

	if err != nil {
		util.WriteError(res, http.StatusInternalServerError, "couldn't find user's vote")
		return
	}

	if len(voteRecord) == 0 {
		// user vote doesn't exist
		err = db.Create(&Vote{
			User:   user.Username,
			Answer: ans.ID,
			Vote:   vote.Vote,
		}).Error

		if err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't create user's vote")
			return
		}

		if vote.Vote > 0 {
			ans.Upvotes += 1
		} else {
			ans.Downvotes += 1
		}

		if err = db.Save(&ans).Error; err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't save user's vote")
			return
		}
	} else if voteRecord[0].Vote == vote.Vote {
		// user already voted and vote the same option
		if voteRecord[0].Vote > 0 {
			ans.Upvotes -= 1
		} else {
			ans.Downvotes -= 1
		}

		if err = db.Save(&ans).Error; err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't save user's same vote")
			return
		}

		if err = db.Delete(&voteRecord[0]).Error; err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't delete user's vote")
			return
		}
	} else {
		// user already voted and vote the other option
		voteRecord[0].Vote = vote.Vote
		if voteRecord[0].Vote > 0 {
			ans.Upvotes += 1
			ans.Downvotes -= 1
		} else {
			ans.Upvotes -= 1
			ans.Downvotes += 1
		}

		if err = db.Save(&ans).Error; err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't save user's other vote")
			return
		}

		if err = db.Save(&voteRecord[0]).Error; err != nil {
			util.WriteError(res, http.StatusInternalServerError, "couldn't save user's old vote")
			return
		}
	}

	if err = util.WriteJson(res, util.Res{Res: "OK"}); err != nil {
		util.WriteError(res, http.StatusInternalServerError, "couldn't write response")
	}
}
