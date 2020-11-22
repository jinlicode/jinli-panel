package site

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/Template"
	"github.com/jinlicode/jinli-panel/global"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	resp "github.com/jinlicode/jinli-panel/model/response"
	"github.com/jinlicode/jinli-panel/tools"
	"github.com/jinlicode/jinli-panel/utils"
)

// GetLists 获取网站列表
func GetLists(c *gin.Context) {

	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)

	err, list, total := model.GetSiteList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}

// CreateSite 新建网站
func CreateSite(c *gin.Context) {

	var R request.Site
	_ = c.ShouldBindJSON(&R)

	SiteVerify := make(map[string][]string)
	if R.IsSsl == 1 {
		SiteVerify = utils.Rules{
			"Domain":     {utils.NotEmpty()},
			"Email":      {utils.NotEmpty()},
			"PhpVersion": {utils.NotEmpty()},
		}
	} else {
		SiteVerify = utils.Rules{
			"Domain":     {utils.NotEmpty()},
			"PhpVersion": {utils.NotEmpty()},
		}
	}

	SiteVerifyErr := utils.Verify(R, SiteVerify)
	if SiteVerifyErr != nil {
		response.FailWithMessage(SiteVerifyErr.Error(), c)
		return
	}

	isEx := model.CheckSiteByDomain(R.Domain)
	if isEx == false {
		response.FailWithMessage("域名已存在", c)
		return
	}

	//入库
	siteid, _ := model.CreateSite(R)

	newDomain := tools.DotToUnderline(R.Domain)
	Domain := R.Domain
	//自动创建以网站名字命名的程序目录
	tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "code/" + newDomain)
	tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/php/" + newDomain)

	//写入404以及index文件到置顶目录
	tools.WriteFile(global.BASEPATH+"code/"+newDomain+"/index.html", Template.HTMLIndex())
	tools.WriteFile(global.BASEPATH+"code/"+newDomain+"/404.html", Template.HTML404())

	//创建网站的配置文件到对应的config配置文件中
	tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/php/" + newDomain)
	tools.WriteFile(global.BASEPATH+"config/php/"+newDomain+"/www.conf", Template.PhpWww())
	tools.WriteFile(global.BASEPATH+"config/php/"+newDomain+"/php.ini", Template.PhpIni())
	tools.WriteFile(global.BASEPATH+"config/php/"+newDomain+"/php-fpm.conf", Template.PhpFpm())
	tools.WriteFile(global.BASEPATH+"config/rewrite/"+newDomain+".conf", "")

	//创建对应nginx.conf到对应目录
	if R.IsSsl == 0 {
		TemplateNginxHTTPString := Template.TemplateNginxHttp(newDomain, Domain)
		tools.WriteFile(global.BASEPATH+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPString)

	} else {
		TemplateNginxHTTPSString := Template.TemplateNginxHttps(newDomain, Domain)
		tools.WriteFile(global.BASEPATH+"config/nginx/"+newDomain+".conf", TemplateNginxHTTPSString)
	}

	// 创建域名对应的网络
	tools.ExecLinuxCommand("docker network create " + newDomain + "_net")

	// 同时加入mysql网络和nginx网络
	tools.ExecLinuxCommand("docker network connect " + newDomain + "_net mysql")
	tools.ExecLinuxCommand("docker network connect " + newDomain + "_net nginx")

	//创建测试网站
	// tools.ExecLinuxCommand()

	//执行nginx重启
	// tools.ExecLinuxCommand()

	shellString := "docker run -d --name " + newDomain + " --network  " + newDomain + "_net --user 10000:10000 --restart unless-stopped --env TZ=Asia/Shanghai -v " + global.BASEPATH + "code/" + newDomain + ":/var/www/" + newDomain + " -v " + global.BASEPATH + "config/php/" + newDomain + "/php.ini:/usr/local/etc/php/php.ini -v " + global.BASEPATH + "config/php/" + newDomain + "/php-fpm.conf:/usr/local/etc/php-fpm.conf -v " + global.BASEPATH + "config/php/" + newDomain + "/www.conf:/usr/local/etc/php-fpm.d/www.conf -v " + global.BASEPATH + "log/openrasp/" + newDomain + ":/opt/rasp/logs/alarm hub.jinli.plus/jinlicode/php:v" + R.PhpVersion + " && docker exec nginx nginx -s reload"

	// 入task
	task := request.Task{Name: "", Execstr: shellString, Type: "shell", Siteid: siteid}
	model.AddTask(task)

	response.OkWithData("success", c)

}

// DelSite 删除网站
func DelSite(c *gin.Context) {

	var R request.Site
	_ = c.ShouldBindJSON(&R)

	SiteVerify := utils.Rules{
		"ID": {utils.NotEmpty()},
	}

	SiteVerifyErr := utils.Verify(R, SiteVerify)
	if SiteVerifyErr != nil {
		response.FailWithMessage(SiteVerifyErr.Error(), c)
		return
	}

	model.DelSite(R)

	response.OkWithData("success", c)
}

// GetSiteConf 获取网站配置项
func GetSiteConf(c *gin.Context) {
	id := c.Query("id")
	idString, _ := strconv.Atoi(id)

	info, _ := model.GetSiteInfo(idString)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	newDomain := tools.DotToUnderline(siteInfo.Domain)
	confText := tools.ReadFile(global.BASEPATH + "config/nginx/" + newDomain + ".conf")

	response.OkWithData(resp.TextResult{
		Text: confText,
	}, c)
}

// GetSiteRewrite 获取伪静态重写规则
func GetSiteRewrite(c *gin.Context) {
	id := c.Query("id")
	idString, _ := strconv.Atoi(id)

	info, _ := model.GetSiteInfo(idString)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	newDomain := tools.DotToUnderline(siteInfo.Domain)
	RewriteText := tools.ReadFile(global.BASEPATH + "config/rewrite/" + newDomain + ".conf")

	response.OkWithData(resp.TextResult{
		Text: RewriteText,
	}, c)
}

// GetSitePhp 获取PHP版本
func GetSitePhp(c *gin.Context) {
	id := c.Query("id")
	idString, _ := strconv.Atoi(id)

	info, _ := model.GetSiteInfo(idString)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	response.OkWithData(resp.TextResult{
		Text: siteInfo.PhpVersion,
	}, c)
}

// GetSiteDomain 获取所有的绑定域名
func GetSiteDomain(c *gin.Context) {
	id := c.Query("id")
	idString, _ := strconv.Atoi(id)

	info, _ := model.GetSiteInfo(idString)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	// 获取domian域名
	err, list := model.GetSiteDomainList(siteInfo.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List: list,
		}, c)
	}
}

// GetSiteBasepath 获取网站所有目录
func GetSiteBasepath(c *gin.Context) {
	id := c.Query("id")
	idString, _ := strconv.Atoi(id)

	info, _ := model.GetSiteInfo(idString)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	newDomain := tools.DotToUnderline(siteInfo.Domain)

	//获取当前所有的目录的切片
	DirListSlice := tools.GetPathFiles(global.BASEPATH+"code/"+newDomain, true)
	DirListSlice = append(DirListSlice, "/")

	response.OkWithData(resp.PageResult{
		List: DirListSlice,
	}, c)
}
