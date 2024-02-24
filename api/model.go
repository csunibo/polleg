package api

import (
	"gorm.io/gorm"
	"time"
)

type Answer struct {
	// taken from from gorm.Model, so we can json strigify properly
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Question uint  `json:"question" gorm:"foreignKey:Question;references:ID"`
	Parent   *uint `json:"-"`

	User      string   `json:"user"`
	Content   string   `json:"content"`
	Upvotes   uint32   `json:"upvotes"`
	Downvotes uint32   `json:"downvotes"`
	Replies   []Answer `json:"replies" gorm:"foreignKey:Parent;references:ID"`
}

type Question struct {
	// taken from from gorm.Model, so we can json strigify properly
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

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
