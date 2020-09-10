package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open(sqlite.Open("db/config.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
