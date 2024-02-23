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

type Question struct {
	gorm.Model
	Document string   `json:"document"`
	Start    uint32   `json:"start"`
	End      uint32   `json:"end"`
	Answers  []Answer `json:"answers" gorm:"foreignKey:Question;references:ID"`
}

type Coord struct {
	gorm.Model
	Start uint32 `json:"start"`
	End   uint32 `json:"end"`
}

type DocReq struct {
	Document string  `json:"document"`
	Coords   []Coord `json:"coords"`
}

type Vote struct {
	gorm.Model
	Answer uint   `json:"answer"`
	User   string `json:"user"`
	Vote   int8   `json:"vote"`
}
