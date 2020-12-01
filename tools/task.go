package tools

import (
	"github.com/jinlicode/jinli-panel/model"
)

var taskRunLogs = "/tmp/taskRunLogs.log"

// RunTask 跑定时任务
func RunTask() {

	//获取未开始的任务列表
	waitInfo, _ := model.GetTaskFirst("0")

	//设置任务状态为1运行中
	model.SetTaskStatus(waitInfo.ID, "1")

	// 判断是否是站点
	if waitInfo.Type == "site-shell" && waitInfo.Execstr != "" {
		ExecLinuxCommand(waitInfo.Execstr + " &> " + taskRunLogs)

		if waitInfo.Siteid > 0 {
			//设置网站为运行中
			model.SetSiteStatus(waitInfo.Siteid, "1")
		}

		// 安装软件
	} else if waitInfo.Type == "docker-shell" && waitInfo.Execstr != "" {

		ExecLinuxCommand(waitInfo.Execstr + " &> " + taskRunLogs)
	}

	//设置任务状态为2已运行
	model.SetTaskStatus(waitInfo.ID, "2")
}
