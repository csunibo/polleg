package documents

import "gorm.io/gorm"

type Document struct {
	ID string `json:"id"`
}

type Question struct {
	gorm.Model
	Document Document `json:"document"`
}