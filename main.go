package main

import (
	"flag"
	"fmt"

	"github.com/jinlicode/jinli-panel/Template"
	"github.com/jinlicode/jinli-panel/global"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	"github.com/jinlicode/jinli-panel/routers"
	"github.com/jinlicode/jinli-panel/tools"
	"github.com/robfig/cron"
)

// @title 锦鲤管理面板
// @version 2.0
// @description 锦鲤管理面板主打安全、可靠、性能调优。容器化管理
// @termsOfService https://github.com/jinlicode/jinli-panel
// @license.name GNU v3
// @license.url https://github.com/jinlicode/jinli-panel/blob/master/LICENSE
func main() {
	init := flag.String("init", "", "--init=all")
	flag.Parse()
	// 检测锦鲤面板 安装目录
	if *init == "all" && tools.CheckFileExist(global.BASEPATH+"install.lock") == false {

		//初始化链接db
		model.InitDbConnt()

		//创建mysql my.cnf
		tools.WriteFile(global.BASEPATH+"config/mysql/my.cnf", Template.MysqlCnf())

		//自动生成mysql密码
		mysqlRandPassword := tools.RandomString(16)

		// 后台默认拉基础镜像
		model.AddTask(request.Task{Name: "docker-nginx", Execstr: " docker run -d --name nginx --network nginx_net --restart always --env TZ=Asia/Shanghai -p 80:80 -p 443:443 -v " + global.BASEPATH + "config/nginx:/etc/nginx/conf.d -v " + global.BASEPATH + "code:/var/www -v " + global.BASEPATH + "log/nginx:/var/log/nginx -v " + global.BASEPATH + "config/cert:/etc/letsencrypt -v " + global.BASEPATH + "config/rewrite:/etc/nginx/rewrite hub.jinli.plus/jinlicode/nginx:v1", Type: "docker-shell"})

		model.AddTask(request.Task{Name: "docker-mysql", Execstr: "docker run -d --name mysql --network mysql_net --restart always --env TZ=Asia/Shanghai --env MYSQL_ROOT_PASSWORD=" + mysqlRandPassword + "  -p 3306:3306 -v " + global.BASEPATH + "db:/var/lib/mysql -v " + global.BASEPATH + "imput_db:/docker-entrypoint-initdb.d -v " + global.BASEPATH + "config/mysql/my.cnf:/etc/mysql/my.cnf hub.jinli.plus/jinlicode/mysql", Type: "docker-shell"})

		model.AddTask(request.Task{Name: "phpmyadmin:5.0.2", Execstr: "docker pull hub.jinli.plus/jinlicode/phpmyadmin:5.0.2", Type: "docker-shell"})
		model.AddTask(request.Task{Name: "redis:5.0.9", Execstr: "docker pull hub.jinli.plus/jinlicode/redis:5.0.9", Type: "docker-shell"})
		model.AddTask(request.Task{Name: "memcached:1.6.6", Execstr: "docker pull hub.jinli.plus/jinlicode/memcached:1.6.6", Type: "docker-shell"})

		//加入数据库表
		model.SetConfigMsqlpwd(mysqlRandPassword)

		//随机用户名 随机密码
		userName := tools.RandomString(12)
		userPassword := tools.RandomString(12)
		model.AddUser(request.User{Name: userName, Password: userPassword})

		tools.WriteFile(global.BASEPATH+"install.lock", "installed")

		fmt.Println("|-用户名: " + userName)
		fmt.Println("|-新密码: " + userPassword)
		return

	}

	//初始化链接db
	model.InitDbConnt()

	//跑定时任务
	c := cron.New()
	c.AddFunc("*/2 * * * * *", func() {
		tools.RunTask()
	})
	c.Start()

	// 启动gin
	routers.InitRouter()
	select {}
}
