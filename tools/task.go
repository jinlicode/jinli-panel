package tools

import (
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
)

// RunTask 跑定时任务
func RunTask() {

	//获取未开始的任务列表
	waitList, _ := model.GetTaskList("0")

	for _, taskv := range waitList.([]request.Task) {
		// 判断是否是
		if taskv.Type == "shell" && taskv.Execstr != "" {
			ExecLinuxCommand(taskv.Execstr)
		}

		if taskv.Siteid > 0 {
			//设置网站为运行中
			model.SetSiteStatus(taskv.Siteid, "1")
		}

		//设置任务状态为1已运行
		model.SetTaskStatus(taskv.ID, "1")

	}
}
