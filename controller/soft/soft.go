package soft

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	resp "github.com/jinlicode/jinli-panel/model/response"
	"github.com/jinlicode/jinli-panel/tools"
	"github.com/jinlicode/jinli-panel/utils"
)

type softStruct struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
}

var softlist = []softStruct{
	{
		Name: "nginx:v1",
		Desc: "nginx v1",
	}, {
		Name: "mysql:latest",
		Desc: "mysql latest",
	}, {
		Name: "php:v5.6",
		Desc: "php5.6",
	}, {
		Name: "php:v5.6-sec",
		Desc: "php5.6 安全版本",
	}, {
		Name: "php:v7.0",
		Desc: "php7.0",
	}, {
		Name: "php:v7.0-sec",
		Desc: "php7.0 安全版本",
	}, {
		Name: "php:v7.1",
		Desc: "php7.1",
	}, {
		Name: "php:v7.1-sec",
		Desc: "php7.1 安全版本",
	}, {
		Name: "php:v7.2",
		Desc: "php7.2",
	}, {
		Name: "php:v7.2-sec",
		Desc: "php7.2 安全版本",
	}, {
		Name: "php:v7.3",
		Desc: "php7.3",
	}, {
		Name: "php:v7.3-sec",
		Desc: "php7.3 安全版本",
	},
}

// GetSoftList 获取所有的镜像
func GetSoftList(c *gin.Context) {

	//获取所有的镜像
	imageMapList := tools.GetDockerImages()

	//获取安装中的插件
	dockerList, _ := model.GetTaskByTypeList("docker-shell", "1")
	dockerInstallMap := make(map[string]string)

	for _, v := range dockerList.([]request.Task) {
		dockerInstallMap[v.Name] = v.Name
	}

	//判断是否已经安装了
	for k, v := range softlist {
		if _, ok := imageMapList["hub.jinli.plus/jinlicode/"+v.Name]; ok {
			softlist[k].Status = 1
		}

		//判断安装中
		if _, ok := dockerInstallMap[v.Name]; ok {
			softlist[k].Status = 2
		}
	}

	response.OkWithData(resp.PageResult{
		List: softlist,
	}, c)
}

// GetPHPList 获取所有的php镜像
func GetPHPList(c *gin.Context) {
	var phplist []softStruct

	//获取所有的镜像
	imageMapList := tools.GetDockerImages()

	//获取所有的phpMap
	phpMap := make(map[string]softStruct)
	for _, v := range softlist {
		if strings.Index(v.Name, "php:v") != -1 {
			phpMap[v.Name] = softStruct{
				Name: v.Name,
				Desc: v.Desc,
			}
		}
	}

	//判断是否有镜像
	for _, v := range phpMap {
		if _, ok := imageMapList["hub.jinli.plus/jinlicode/"+v.Name]; ok {
			phplist = append(phplist, softStruct{
				Name: v.Name,
				Desc: v.Desc,
			})
		}
	}

	response.OkWithData(resp.PageResult{
		List: phplist,
	}, c)
}

// InstallSoft 安装镜像
func InstallSoft(c *gin.Context) {
	var R softStruct
	_ = c.ShouldBindJSON(&R)

	SoftVerify := utils.Rules{
		"Name": {utils.NotEmpty()},
	}

	SoftVerifyErr := utils.Verify(R, SoftVerify)
	if SoftVerifyErr != nil {
		response.FailWithMessage(SoftVerifyErr.Error(), c)
		return
	}

	//获取所有的phpMap
	softMap := make(map[string]softStruct)
	for _, v := range softlist {
		softMap[v.Name] = softStruct{
			Name: v.Name,
			Desc: v.Desc,
		}
	}

	//加入安装队列
	if _, ok := softMap[R.Name]; ok {
		model.AddTask(request.Task{
			Name:    R.Name,
			Desc:    softMap[R.Name].Desc,
			Type:    "docker-shell",
			Execstr: "docker pull hub.jinli.plus/jinlicode/" + R.Name,
		})
		response.OkWithData("success", c)
		return
	}
	response.FailWithMessage("软件不存在", c)
}
