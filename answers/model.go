package answers

import (
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	Document string
	Question uint
	Parent   *uint

	User      string   `json:"user"`
	Content   string   `json:"content"`
	Upvotes   uint32   `json:"upvotes"`
	Downvotes uint32   `json:"downvotes"`
	Replies   []Answer `json:"replies" gorm:"foreignKey:Parent;references:ID"`
}
