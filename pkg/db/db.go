package db

import (
	"blackbox-v2/pkg/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var environment string
var DB *gorm.DB
var dbFileName string

func connect() {
	var err error
	if environment == "dev" {
		DB, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		DB, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	config := config.Config
	dbFileName = config.Database.File
	environment = config.Environment
	connect()
}
