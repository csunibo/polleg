package util

import (
	"fmt"

	"github.com/csunibo/stackunibo/documents"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func ConnectDb(ConnStr string) error {
	var err error
	db, err = gorm.Open(postgres.Open(ConnStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	TestInit()
	return nil
}

func GetDb() *gorm.DB {
	return db
}

func TestInit() {
	db.Create(&documents.Document{
		ID: "stringasha",
	})

	db.Create(&documents.Question{
		Model:    gorm.Model{ID: 123456},
		Document: documents.Document{ID: "stringasha"},
	})
}
