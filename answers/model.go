package answers

import (
	"github.com/csunibo/stackunibo/documents"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	Document  documents.Document `json:"document" gorm:"foreignKey:Document;references:ID"`
	Question  documents.Question `json:"question" gorm:"foreignKey:Question;references:ID"`
	Parent    *Answer            `json:"parent" gorm:"foreignKey:Parent;references:ID"`
	User      string             `json:"user"`
	Content   string             `json:"content"`
	Upvotes   uint32             `json:"upvotes"`
	Downvotes uint32             `json:"downvotes"`
}
