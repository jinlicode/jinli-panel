package model

import "github.com/jinlicode/jinli-panel/model/request"

func InitDb() {
	db.AutoMigrate(
		&request.Database{},
		&request.Task{},
		&request.Site{},
		&request.Domain{},
		&request.User{},
		&Config{},
	)
}
