package model

import (
	"time"

	"github.com/jinlicode/jinli-panel/model/request"
)

// GetTaskList 获取task列表内容
func GetTaskList(status string) (list interface{}, err error) {
	var task []request.Task
	err = db.Where("status = ?", status).Find(&task).Error
	return task, err
}

// AddTask
func AddTask(info request.Task) (err error) {
	task := request.Task{Name: info.Name, Status: "0", Start: time.Now().Unix(), End: 0, Execstr: info.Execstr, Addtime: time.Now().Format("2006-01-02 15:04:05")}
	err = db.Create(&task).Error
	return err
}

// SetTaskStatus
func SetTaskStatus(taskid int, status string) (err error) {
	task := request.Task{ID: taskid}
	err = db.Model(&task).Update("status", status).Error
	return err
}
