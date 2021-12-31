package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func getDbInstance() *gorm.DB {
	if dbInstance == nil {
		db, err := gorm.Open(sqlite.Open("twitter_tools.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err != nil {
			log.Fatal("failed to connect database")
		}
		dbInstance = db
	}

	return dbInstance
}

func InitDb() {
	db := getDbInstance()

	db.AutoMigrate(&Friend{})
}
