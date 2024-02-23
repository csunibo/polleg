package answers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
