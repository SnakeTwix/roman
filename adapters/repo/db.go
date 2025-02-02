package repo

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
	"log"
	"roman/adapters/repo/model"
	"roman/port"
)

func InitDB(config port.ConfigService) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.SqliteDbFile()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	err = db.AutoMigrate(&model.Birthday{})
	if err != nil {
		// Quite likely don't want to continue anything if these fail
		log.Fatal("Migrations failed")
	}

	return db
}
