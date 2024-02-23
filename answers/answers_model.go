package answers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model     `json:"model"`
	DocumentID     uuid.UUID  `json:"document_ID"`
	QuestionID     uuid.UUID  `json:"question_ID"`
	AnswerID       uuid.UUID  `json:"answer_ID"`
	ParentAnswerID *uuid.UUID `json:"parentAnswer_ID"`
	userID         string     `json:"user_ID"`
	content        string     `json:"content"`
	upvotes        uint32     `json:"upvotes"`
	downvotes      uint32     `json:"downvotes"`
}
