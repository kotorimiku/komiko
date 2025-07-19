package db

import (
	"log"
	"os"

	"komiko/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dbPath := "data/komiko.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		os.MkdirAll("data", 0755)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	db.Exec("PRAGMA journal_mode=WAL;")

	db.AutoMigrate(&model.User{}, &model.Book{}, &model.Library{}, &model.Genre{}, &model.Person{}, &model.Series{}, &model.Progress{})

	return db
}
