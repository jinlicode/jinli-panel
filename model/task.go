package model

import (
	"time"

	"github.com/jinlicode/jinli-panel/model/request"
)

// AddTask
func AddTask(info request.Task) (err error) {
	task := request.Task{Name: info.Name, Status: "0", Start: time.Now().Unix(), End: 0, Execstr: info.Execstr, Addtime: time.Now().Format("2006-01-02 15:04:05")}
	err = db.Create(&task).Error
	return err
}
