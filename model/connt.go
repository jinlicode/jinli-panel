package model

import (
	"github.com/jinlicode/jinli-panel/global"
	"github.com/jinlicode/jinli-panel/model/request"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDbConnt() {
	//open a db connection
	var err error
	db, err = gorm.Open(sqlite.Open(global.BASEPATH+"config.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&request.Database{},
		&request.Task{},
		&request.Site{},
		&request.Domain{},
		&request.User{},
		&Config{},
	)

}
