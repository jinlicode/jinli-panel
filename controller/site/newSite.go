package site

import (
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/Template"
	"github.com/jinlicode/jinli-panel/global"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	"github.com/jinlicode/jinli-panel/tools"
	"github.com/jinlicode/jinli-panel/utils"
)

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
	model.CreateSite(R)

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

	NewSiteHTTPS := "否"
	//创建对应nginx.conf到对应目录
	if NewSiteHTTPS == "否" {
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
	tools.ExecLinuxCommand("docker run -d --name " + newDomain + " --network  " + newDomain + "_net --user 10000:10000 --restart unless-stopped --env TZ=Asia/Shanghai -v " + global.BASEPATH + "code/" + newDomain + ":/var/www/" + newDomain + " -v " + global.BASEPATH + "config/php/" + newDomain + "/php.ini:/usr/local/etc/php/php.ini -v " + global.BASEPATH + "config/php/" + newDomain + "/php-fpm.conf:/usr/local/etc/php-fpm.conf -v " + global.BASEPATH + "config/php/" + newDomain + "/www.conf:/usr/local/etc/php-fpm.d/www.conf -v " + global.BASEPATH + "log/openrasp/" + newDomain + ":/opt/rasp/logs/alarm hub.jinli.plus/jinlicode/php:v" + R.PhpVersion)

	//执行nginx重启
	tools.ExecLinuxCommand("docker exec nginx nginx -s reload")

	response.OkWithData("success", c)

}
