package api

import (
	"gorm.io/gorm"
	"time"
)

type Answer struct {
	gorm.Model
	Document string `json:"document"`
	Question uint   `json:"question" gorm:"foreignKey:Question;references:ID"`
	Parent   *uint  `json:"parent"`

	User      string   `json:"user"`
	Content   string   `json:"content"`
	Upvotes   uint32   `json:"upvotes"`
	Downvotes uint32   `json:"downvotes"`
	Replies   []Answer `json:"replies" gorm:"foreignKey:Parent;references:ID"`
}

type Question struct {
	gorm.Model
	Document string   `json:"document"`
	Start    uint32   `json:"start"`
	End      uint32   `json:"end"`
	Answers  []Answer `json:"answers" gorm:"foreignKey:Question;references:ID"`
}

type Vote struct {
	Answer uint   `json:"answer" gorm:"primaryKey"`
	User   string `json:"user" gorm:"primaryKey"`
	Vote   int8   `json:"vote"`

	// taken from from gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
