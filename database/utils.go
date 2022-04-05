package database

import (
	"github.com/tupini07/twitter-tools/print_utils"
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
			print_utils.Fatal("failed to connect database")
		}
		dbInstance = db
	}

	return dbInstance
}

func InitDb() {
	db := getDbInstance()

	db.AutoMigrate(&Friend{})
}
