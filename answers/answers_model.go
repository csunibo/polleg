package answers

import (
	"net/http"

	auth "github.com/csunibo/stackunibo/auth"
	"github.com/csunibo/stackunibo/util"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

func AnswerHandler(res http.ResponseWriter, req *http.Request) {
	parsedToken, err := auth.RequireJWTCookie(res, req)
	if err != nil {
		return
	}

	user := parsedToken.Claims.(jwt.MapClaims)["user"]

	err = util.WriteJson(res, user)
	if err != nil {
		_ = util.WriteError(res, http.StatusInternalServerError, "")
		slog.Error("could not encode json:", "error", err)
	}
}

/*
type Answer struct {
	document_id: uuid
	question_id: uuid
	answer_id: uuid
	parent_answer_id: option<uuid> // si fa con un puntatore in go
	user_id: string // username di github
	content: string
	upvotes: uint32
	dowvotes: uint32
  }*/
