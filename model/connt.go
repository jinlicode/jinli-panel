package model

import (
	"github.com/jinlicode/jinli-panel/global"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open(sqlite.Open(global.BASEPATH+"config.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
