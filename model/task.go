package model

import (
	"time"

	"github.com/jinlicode/jinli-panel/model/request"
)

// GetTaskFirst 获取最新一条任务
func GetTaskFirst(status string) (tk request.Task, err error) {
	var task request.Task
	err = db.Where("status = ?", status).First(&task).Error
	return task, err
}

// GetTaskList 获取task列表内容
func GetTaskList(status string) (list interface{}, err error) {
	var task []request.Task
	err = db.Where("status = ?", status).Find(&task).Error
	return task, err
}

// GetTaskList 获取task列表内容
func GetTaskByTypeList(Type string, Status string) (list interface{}, err error) {
	var task []request.Task
	err = db.Where("type = ? and status = ?", Type, Status).Find(&task).Error
	return task, err
}

// AddTask
func AddTask(info request.Task) (err error) {
	task := request.Task{
		Name:    info.Name,
		Desc:    info.Desc,
		Siteid:  info.Siteid,
		Type:    info.Type,
		Status:  "0",
		Start:   time.Now().Unix(),
		End:     0,
		Execstr: info.Execstr,
		Addtime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = db.Create(&task).Error
	return err
}

// SetTaskStatus
func SetTaskStatus(taskid int, status string) (err error) {
	task := request.Task{ID: taskid}
	err = db.Model(&task).Update("status", status).Error
	return err
}
