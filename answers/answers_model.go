package answers

import (
	"net/http"

	auth "github.com/csunibo/stackunibo/auth"
	"github.com/csunibo/stackunibo/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
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

type Answer struct {
	gorm.Model     `json:"model"`
	DocumentId     uuid.UUID  `json:"documentId"`
	QuestionId     uuid.UUID  `json:"questionId"`
	AnswerId       uuid.UUID  `json:"answerId"`
	ParentAnswerId *uuid.UUID `json:"parentAnswerId"`
	userId         string     `json:"userId"`
	content        string     `json:"content"`
	upvotes        uint32     `json:"upvotes"`
	downvotes      uint32     `json:"downvotes"`
}
