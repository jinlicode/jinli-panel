package model

import (
	"time"

	"github.com/jinlicode/jinli-panel/model/request"
)

// GetDatabaseList
func GetDatabaseList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var database []request.Database
	err = db.Limit(limit).Offset(offset).Find(&database).Error
	return err, database, total
}

// AddDatabase
func AddDatabase(info request.Database) (err error) {
	database := request.Database{
		Name:     info.Name,
		Pid:      info.Pid,
		Username: info.Username,
		Password: info.Password,
		Addtime:  time.Now().Format("2006-01-02 15:04:05"),
	}

	err = db.Create(&database).Error
	return err
}

// DelDatabaseBySiteID
func DelDatabaseBySiteID(siteID int) (err error) {
	var database request.Database
	err = db.Where("pid = ?", siteID).Delete(&database).Error
	return err
}
