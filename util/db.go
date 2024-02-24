package util

import (
	"fmt"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func ConnectDb(ConnStr string) error {
	gormLogger := slogGorm.New(
		slogGorm.WithTraceAll(), // trace all messages
	)
	var err error
	db, err = gorm.Open(postgres.Open(ConnStr), &gorm.Config{
		Logger:      gormLogger,
		PrepareStmt: true, // optimize raw queries
	})
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	return nil
}

func GetDb() *gorm.DB {
	return db
}
