package model

import "github.com/jinlicode/jinli-panel/model/request"

func Init() {
	db.AutoMigrate(
		&request.Database{},
		&request.Task{},
		&request.Site{},
		&request.Domain{},
		&request.User{},
		&Config{},
	)
}
