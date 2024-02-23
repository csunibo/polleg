package answers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/csunibo/stackunibo/auth"
	"github.com/csunibo/stackunibo/util"
	"github.com/kataras/muxie"
)

type AnswerObj struct {
	Question uint   `json:"question"`
	Parent   *uint  `json:"parent"`
	Content  string `json:"content"`
}

func PutAnswerHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	// Check method put is used
	if req.Method != http.MethodPut {
		_ = util.WriteError(res, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	db := util.GetDb()
	user := auth.GetUser(req)

	// Declare a new Person struct.
	var ans AnswerObj

	err := json.NewDecoder(req.Body).Decode(&ans)
	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "decode error")
		return
	}

	// var quest Question
	// if err := db.First(&quest, ans.Question).Error; err != nil {
	// 	util.WriteError(res, http.StatusBadRequest, "no Question associated with request (or other Error)")
	// 	return
	// }

	// if ans.Parent != nil {
	// 	var Parent Answer
	// 	if err = db.First(&Parent, ans.Parent).Error; err != nil {
	// 		util.WriteError(res, http.StatusBadRequest, "parent is given but none found")
	// 		return
	// 	}
	// 	if Parent.Question != quest.ID {
	// 		util.WriteError(res, http.StatusBadRequest, "mismatch between parent question and this question")
	// 		return
	// 	}
	// }

	err = db.Create(&Answer{
		Question:  ans.Question,
		Parent:    ans.Parent,
		User:      user.Username,
		Content:   ans.Content,
		Upvotes:   0,
		Downvotes: 0,
	}).Error

	if err != nil {
		util.WriteError(res, http.StatusBadRequest, "create error")
		return
	}

	r := Res{
		Res: "OK",
	}
	util.WriteJson(res, r)
}

func GetAnswerById(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()
	id := muxie.GetParam(res, "id")

	var ans Answer
	if err := db.First(&ans, id).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "Answer not found")
		return
	}

	util.WriteJson(res, ans)
}

type VoteObj struct {
	Vote int8 `json:"vote"`
}

func PostVote(res http.ResponseWriter, req *http.Request) {
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
	db.Where("votes.user = ? AND votes.answer = ?", user.Username, ans.ID).Find(&voteRecord)
	fmt.Println(voteRecord)
	if len(voteRecord) == 0 {
		db.Create(&Vote{
			User:   user.Username,
			Answer: ans.ID,
			Vote:   vote.Vote,
		})

		if vote.Vote > 0 {
			ans.Upvotes += 1
		} else {
			ans.Downvotes += 1
		}
		db.Save(&ans)
		return
	}

	if voteRecord[0].Vote == vote.Vote {
		if voteRecord[0].Vote > 0 {
			ans.Upvotes -= 1
		} else {
			ans.Downvotes -= 1
		}
		db.Save(&ans)
		db.Delete(&voteRecord[0])
		return
	}

	voteRecord[0].Vote = vote.Vote
	if voteRecord[0].Vote > 0 {
		ans.Upvotes += 1
		ans.Downvotes -= 1
	} else {
		ans.Upvotes -= 1
		ans.Downvotes += 1
	}
	db.Save(&ans)
	db.Save(&voteRecord[0])
}

func GetAnswersByQuestion(res http.ResponseWriter, req *http.Request) {
	db := util.GetDb()
	qid := muxie.GetParam(res, "question")

	var ans []Answer
	if err := db.Where("question = ?", qid).Find(&ans).Error; err != nil {
		util.WriteError(res, http.StatusBadRequest, "Answer not found")
		return
	}

	util.WriteJson(res, ans)
}
