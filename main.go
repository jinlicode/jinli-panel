package main

import (
	"fmt"
	"os"

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
	fmt.Println("检测系统环境......")
	// 检测锦鲤面板 安装目录
	if tools.CheckFileExist(global.BASEPATH+"install.lock") == false {

		fmt.Println("开始初始化安装......")

		portError := ""
		// 检测80 443端口
		isPort80 := tools.PortInUse(80)
		if isPort80 {
			portError += "80 "
		}

		isPort443 := tools.PortInUse(443)
		if isPort443 {
			portError += "443 "
		}
		// 检测3306端口
		isPort3306 := tools.PortInUse(3306)
		if isPort3306 {
			portError += "3306 "
		}
		// 检测9527端口
		isPort9527 := tools.PortInUse(9527)
		if isPort9527 {
			portError += "9527 "
		}

		if portError != "" {
			fmt.Println("请先关闭" + portError + "端口，再继续运行锦鲤面板程序")
			os.Exit(1)
		}

		//执行安装docker
		tools.ExecDockerInstall()

		//创建项目目录
		os.Mkdir(global.BASEPATH, 0755)
		//自动创建db内容
		model.Init()
		//创建代码目录
		os.Mkdir(global.BASEPATH+"code/", 0755)
		//创建各配置项目录
		os.Mkdir(global.BASEPATH+"config/", 0755)
		os.Mkdir(global.BASEPATH+"config/cert/", 0755)
		os.Mkdir(global.BASEPATH+"config/mysql/", 0755)
		os.Mkdir(global.BASEPATH+"config/nginx/", 0755)
		os.Mkdir(global.BASEPATH+"config/php/", 0755)
		os.Mkdir(global.BASEPATH+"config/rewrite/", 0755)

		//创建备份目录
		os.Mkdir(global.BASEPATH+"backup/", 0755)
		os.Mkdir(global.BASEPATH+"backup/database/", 0755)
		os.Mkdir(global.BASEPATH+"backup/site/", 0755)

		//创建自动备份目录
		os.Mkdir(global.BASEPATH+"autobackup/", 0755)
		os.Mkdir(global.BASEPATH+"autobackup/database/", 0755)
		os.Mkdir(global.BASEPATH+"autobackup/site/", 0755)

		//创建日志目录
		os.Mkdir(global.BASEPATH+"log/", 0755)
		os.Mkdir(global.BASEPATH+"log/nginx/", 0755)
		os.Mkdir(global.BASEPATH+"log/openrasp/", 0755)

		//设置代码目录为 10000,10000
		tools.ExecLinuxCommand("chown -R 10000:10000 " + global.BASEPATH + "code/")

		//创建mysql my.cnf
		tools.WriteFile(global.BASEPATH+"config/mysql/my.cnf", Template.MysqlCnf())

		//创建 Nginx 网段
		tools.ExecLinuxCommand("docker network create nginx_net")

		//创建 Mysql 网段
		tools.ExecLinuxCommand("docker network create mysql_net")
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

		tools.WriteFile(global.BASEPATH+"install.lock", "installed")

	}

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
