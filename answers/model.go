package answers

import (
	"github.com/csunibo/stackunibo/documents"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	Document  documents.Document `json:"document"`
	Question  documents.Question `json:"question"`
	Parent    *Answer            `json:"parent"`
	User      string             `json:"user"`
	Content   string             `json:"content"`
	Upvotes   uint32             `json:"upvotes"`
	Downvotes uint32             `json:"downvotes"`
}
