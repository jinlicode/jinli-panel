package site

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	//检测nginx mysql
	if tools.CheckDockerStatus("nginx") == false {
		response.FailWithMessage("nginx服务未启动", c)
		return
	}

	//检测nginx mysql
	if tools.CheckDockerStatus("mysql") == false {
		response.FailWithMessage("mysql服务未启动", c)
		return
	}

	//入库
	siteid, _ := model.CreateSite(R)

	newDomain := tools.DotToUnderline(R.Domain)
	Domain := R.Domain
	//自动创建以网站名字命名的程序目录
	os.Mkdir(global.BASEPATH+"code/"+newDomain, 0755)
	os.Mkdir(global.BASEPATH+"config/php/"+newDomain, 0755)

	//写入404以及index文件到置顶目录
	tools.WriteFile(global.BASEPATH+"code/"+newDomain+"/index.html", Template.HTMLIndex(R.Domain))
	tools.WriteFile(global.BASEPATH+"code/"+newDomain+"/404.html", Template.HTML404())

	//创建网站的配置文件到对应的config配置文件中
	os.Mkdir(global.BASEPATH+"config/php/"+newDomain, 0755)
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

	// 创建对应的mysql数据库
	configInfo, _ := model.GetConfigInfo()

	//加入数据库表
	dataPwd := tools.RandomString(16)
	tools.CreateDatabase(tools.GetDockerIP("mysql"), configInfo.Mysqlpwd, newDomain, newDomain, dataPwd)

	// 数据库帐号密码入库
	model.AddDatabase(request.Database{
		Pid:      siteid,
		Name:     newDomain,
		Username: newDomain,
		Password: dataPwd,
	})

	shellString := "docker run -d --name " + newDomain + " --network  " + newDomain + "_net --user 10000:10000 --restart unless-stopped --env TZ=Asia/Shanghai -v " + global.BASEPATH + "code/" + newDomain + ":/var/www/" + newDomain + " -v " + global.BASEPATH + "config/php/" + newDomain + "/php.ini:/usr/local/etc/php/php.ini -v " + global.BASEPATH + "config/php/" + newDomain + "/php-fpm.conf:/usr/local/etc/php-fpm.conf -v " + global.BASEPATH + "config/php/" + newDomain + "/www.conf:/usr/local/etc/php-fpm.d/www.conf -v " + global.BASEPATH + "log/openrasp/" + newDomain + ":/opt/rasp/logs/alarm hub.jinli.plus/jinlicode/" + R.PhpVersion + " && docker exec nginx nginx -s reload"

	// 入task
	task := request.Task{Name: "", Execstr: shellString, Type: "site-shell", Siteid: siteid}
	model.AddTask(task)

	//创建测试网站 执行nginx重启
	// tools.ExecLinuxCommand(shellString)

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

	info, _ := model.GetSiteInfo(R.ID)

	siteInfo := info.(request.Site)

	//删site记录
	model.DelSite(R)

	//删domain记录
	model.DelSiteDomainBySiteID(siteInfo.ID)

	//删除对应的数据库记录
	model.DelDatabaseBySiteID(siteInfo.ID)

	newDomain := tools.DotToUnderline(siteInfo.Domain)

	// 同时脱钩mysql网络和nginx网络
	tools.ExecLinuxCommand("docker network disconnect " + newDomain + "_net mysql")
	tools.ExecLinuxCommand("docker network disconnect " + newDomain + "_net nginx")

	//删除对应的网络
	tools.ExecLinuxCommand("docker network rm " + newDomain + "_net")

	//删除对应的docker
	tools.ExecLinuxCommand("docker stop " + newDomain + " && docker rm " + newDomain)

	//删除对应的程序
	tools.ExecLinuxCommand("rm -rf " + global.BASEPATH + "code/" + newDomain)

	//删除对应的nginx配置文件 重写文件
	tools.ExecLinuxCommand("rm -f " + global.BASEPATH + "config/rewrite/" + newDomain + ".conf")
	tools.ExecLinuxCommand("rm -f " + global.BASEPATH + "config/nginx/" + newDomain + ".conf")

	//删除对应的php配置
	tools.ExecLinuxCommand("rm -rf " + global.BASEPATH + "config/php/" + newDomain)

	//删除对应的数据库
	configInfo, _ := model.GetConfigInfo()
	tools.MysqlQuery(tools.GetDockerIP("mysql"), "root", configInfo.Mysqlpwd, "mysql", "DROP DATABASE "+newDomain)

	//重启nginx
	tools.ExecLinuxCommand("docker exec nginx nginx -s reload")

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

// UpdateSiteConf 更新配置项
func UpdateSiteConf(c *gin.Context) {
	var R request.Site
	_ = c.ShouldBindJSON(&R)

	info, _ := model.GetSiteInfo(R.ID)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	// 获取原始的conf数据
	newDomain := tools.DotToUnderline(siteInfo.Domain)
	hostConfFilePath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"
	confOldText := tools.ReadFile(hostConfFilePath)

	if R.HostConf != confOldText {
		// 先把旧的存入数据库 然后检测是否配置正确
		tools.WriteFile(hostConfFilePath, R.HostConf)
		checkNginx := tools.ExecLinuxCommandReturn("docker exec nginx nginx -t")

		checkNginxOk := strings.Contains(checkNginx, "successful")

		if checkNginxOk == false {
			// 重新还原数据到文件
			tools.WriteFile(hostConfFilePath, confOldText)
			response.FailWithMessage(checkNginx[:strings.Index(checkNginx, "\n")], c)
			return
		}

		// 运行nginx -s 命令
		tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
	}

	response.OkWithData("success", c)
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

// UpdateSiteRewrite 更新伪静态
func UpdateSiteRewrite(c *gin.Context) {
	var R request.Site
	_ = c.ShouldBindJSON(&R)

	info, _ := model.GetSiteInfo(R.ID)

	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	// 获取原始的conf数据
	newDomain := tools.DotToUnderline(siteInfo.Domain)
	rewriteFilePath := global.BASEPATH + "config/rewrite/" + newDomain + ".conf"
	rewriteOldText := tools.ReadFile(rewriteFilePath)

	if R.RewriteConf != rewriteOldText {
		// 先把旧的存入数据库 然后检测是否配置正确
		tools.WriteFile(rewriteFilePath, R.RewriteConf)
		checkNginx := tools.ExecLinuxCommandReturn("docker exec nginx nginx -t")
		checkNginxOk := strings.Contains(checkNginx, "successful")

		if checkNginxOk == false {
			// 重新还原数据到文件
			tools.WriteFile(rewriteFilePath, rewriteOldText)
			response.FailWithMessage(checkNginx[:strings.Index(checkNginx, "\n")], c)
			return
		}

		// 运行nginx -s 命令
		tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
	}

	response.OkWithData("success", c)
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

// UpdateSitePhp 更新php版本
func UpdateSitePhp(c *gin.Context) {

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

// UpdateSiteDomain 更新绑定域名
func UpdateSiteDomain(c *gin.Context) {

	type updateDomain struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}

	var R updateDomain
	_ = c.ShouldBindJSON(&R)

	info, _ := model.GetSiteInfo(R.ID)
	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	textDomain := R.Text
	domainString := strings.ReplaceAll(textDomain, "\t", "")

	// 获取原始的conf数据
	newDomain := tools.DotToUnderline(siteInfo.Domain)
	hostConfFilePath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"

	confText := tools.ReadFile(hostConfFilePath)
	confOldText := confText

	reg := regexp.MustCompile(`server_name\s*(.*);`)
	confSilce := reg.FindStringSubmatch(confText)

	domainSilce := strings.Split(domainString, "\n")

	// 获取数据库所有的域名
	_, domainList := model.GetSiteDomainAllList()

	// 获取所有的存在的域名map
	domainListMap := make(map[string]int)
	for k, v := range domainList.([]request.Domain) {
		domainListMap[v.Name] = k
	}

	var domain []request.Domain

	newDomainTemp := ""
	for _, v := range domainSilce {
		newDomainTemp = strings.TrimSpace(v)

		if _, ok := domainListMap[newDomainTemp]; ok {
			response.FailWithMessage(newDomainTemp+"域名已存在", c)
			return
		}

		if tools.CheckDomain(newDomainTemp) == false {
			response.FailWithMessage(newDomainTemp+"域名格式不正确", c)
			return
		}

		domain = append(domain, request.Domain{
			Pid:     siteInfo.ID,
			Name:    newDomainTemp,
			Addtime: time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	//通过之后更改域名conf文件
	confText = strings.Replace(confText, confSilce[0], strings.TrimRight(confSilce[0], ";")+" "+strings.Trim(fmt.Sprint(domainSilce), "[]")+";", -1)

	if confText != confOldText {
		// 检测是否配置正确
		tools.WriteFile(hostConfFilePath, confText)
		checkNginx := tools.ExecLinuxCommandReturn("docker exec nginx nginx -t")
		checkNginxOk := strings.Contains(checkNginx, "successful")

		if checkNginxOk == false {
			// 重新还原数据到文件
			tools.WriteFile(hostConfFilePath, confOldText)
			response.FailWithMessage(checkNginx[:strings.Index(checkNginx, "\n")], c)
			return
		}

		// 运行nginx -s 命令
		tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
	}

	// 入库
	model.CreateSiteDomain(domain)
	response.OkWithData("success", c)

}

// DelSiteDomain 删除域名
func DelSiteDomain(c *gin.Context) {

	type updateDomain struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	}

	var R updateDomain
	_ = c.ShouldBindJSON(&R)

	domainInfo, _ := model.GetSiteDomainInfo(R.ID)

	// 数据异常
	if domainInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	// 删除当前的域名
	model.DelSiteDomain(R.ID)

	info, _ := model.GetSiteInfo(domainInfo.Pid)
	siteInfo := info.(request.Site)

	// 获取原始的conf数据
	newDomain := tools.DotToUnderline(siteInfo.Domain)
	hostConfFilePath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"

	confText := tools.ReadFile(hostConfFilePath)
	confOldText := confText

	reg := regexp.MustCompile(`server_name\s*(.*);`)
	confSilce := reg.FindStringSubmatch(confText)

	// 获取数据库所有的域名
	_, domainList := model.GetSiteDomainList(siteInfo.ID)

	// 获取所有的存在的域名map
	var domainSilce []string
	for _, v := range domainList.([]request.Domain) {
		domainSilce = append(domainSilce, v.Name)
	}

	//通过之后更改域名conf文件
	confText = strings.Replace(confText, confSilce[1], strings.Trim(fmt.Sprint(domainSilce), "[]"), -1)

	if confText != confOldText {
		// 检测是否配置正确
		tools.WriteFile(hostConfFilePath, confText)
		checkNginx := tools.ExecLinuxCommandReturn("docker exec nginx nginx -t")
		checkNginxOk := strings.Contains(checkNginx, "successful")

		if checkNginxOk == false {
			// 重新还原数据到文件
			tools.WriteFile(hostConfFilePath, confOldText)
			response.FailWithMessage(checkNginx[:strings.Index(checkNginx, "\n")], c)
			return
		}

		// 运行nginx -s 命令
		tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
	}
	response.OkWithData("success", c)
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
	confText := tools.ReadFile(global.BASEPATH + "config/nginx/" + newDomain + ".conf")

	//获取当前所有的目录的切片
	DirListSlice := tools.GetPathFiles(global.BASEPATH+"code/"+newDomain, true)
	for k, v := range DirListSlice {
		DirListSlice[k] = "/" + v
	}

	DirListSlice = append(DirListSlice, "/")

	type Result struct {
		Basepath string      `json:"basepath"`
		List     interface{} `json:"list"`
	}

	//返回第一个匹配的字符串的首末位置
	reg := regexp.MustCompile(`\s*root\s*\$base(.*);`)
	Basepath := "/"

	BasepathSilce := reg.FindStringSubmatch(confText)

	if len(BasepathSilce) >= 2 && BasepathSilce[1] != "" {
		Basepath = BasepathSilce[1]
	}

	response.OkWithData(Result{
		Basepath: Basepath,
		List:     DirListSlice,
	}, c)
}

// UpdateSiteBasepath 更新根目录
func UpdateSiteBasepath(c *gin.Context) {

	type Data struct {
		ID       int    `json:"id"`
		Basepath string `json:"basepath"`
	}

	var R Data
	_ = c.ShouldBindJSON(&R)

	info, _ := model.GetSiteInfo(R.ID)
	siteInfo := info.(request.Site)

	// 数据异常
	if siteInfo.ID == 0 {
		response.FailWithMessage("获取数据失败", c)
	}

	// 获取原始的conf数据
	newDomain := tools.DotToUnderline(siteInfo.Domain)
	hostConfFilePath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"

	confText := tools.ReadFile(hostConfFilePath)
	confOldText := confText

	reg := regexp.MustCompile(`(\s*root\s*\$base)(.*);`)
	BasepathSilce := reg.FindStringSubmatch(confText)

	if R.Basepath == "/" {
		R.Basepath = ""
	}

	confText = strings.Replace(confText, BasepathSilce[0], BasepathSilce[1]+R.Basepath+";", -1)

	if confText != confOldText {
		// 检测是否配置正确
		tools.WriteFile(hostConfFilePath, confText)
		checkNginx := tools.ExecLinuxCommandReturn("docker exec nginx nginx -t")
		checkNginxOk := strings.Contains(checkNginx, "successful")

		if checkNginxOk == false {
			// 重新还原数据到文件
			tools.WriteFile(hostConfFilePath, confOldText)
			response.FailWithMessage(checkNginx[:strings.Index(checkNginx, "\n")], c)
			return
		}

		// 运行nginx -s 命令
		tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
	}

	response.OkWithData("success", c)
}

// UpdateSiteStatus 设置网站状态
func UpdateSiteStatus(c *gin.Context) {

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

	info, _ := model.GetSiteInfo(R.ID)

	siteInfo := info.(request.Site)

	if siteInfo.ID > 0 {

		if siteInfo.Status == 1 {

			// 获取原始的conf数据
			newDomain := tools.DotToUnderline(siteInfo.Domain)
			hostConfPath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"
			rewriteConfPath := global.BASEPATH + "config/rewrite/" + newDomain + ".conf"
			hostConf := tools.ReadFile(hostConfPath)
			rewriteConf := tools.ReadFile(rewriteConfPath)

			//把配置内容先写入数据库
			model.SetSiteInfoByID(siteInfo.ID, "host_conf", hostConf)
			model.SetSiteInfoByID(siteInfo.ID, "rewrite_conf", rewriteConf)
			model.SetSiteStatus(siteInfo.ID, "2")

			//删除配置项
			tools.ExecLinuxCommand("rm -f " + hostConfPath)
			tools.ExecLinuxCommand("rm -f " + rewriteConfPath)

			//重启nginx
			tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")

		} else if siteInfo.Status == 2 {

			// 获取原始的conf数据
			newDomain := tools.DotToUnderline(siteInfo.Domain)
			hostConfPath := global.BASEPATH + "config/nginx/" + newDomain + ".conf"
			rewriteConfPath := global.BASEPATH + "config/rewrite/" + newDomain + ".conf"

			//把配置内容先写入数据库
			model.SetSiteStatus(siteInfo.ID, "1")

			//把配置内容读入文件
			tools.WriteFile(hostConfPath, siteInfo.HostConf)
			tools.WriteFile(rewriteConfPath, siteInfo.RewriteConf)

			//重启nginx
			tools.ExecLinuxCommandReturn("docker exec nginx nginx -s reload")
		}
	}

	response.OkWithData("success", c)
}
