package main

import (
	"fmt"
	"os"

	"github.com/jinlicode/jinli-panel/Template"
	"github.com/jinlicode/jinli-panel/global"
	"github.com/jinlicode/jinli-panel/routers"
	"github.com/jinlicode/jinli-panel/tools"
)

// @title 锦鲤管理面板
// @version 2.0
// @description 锦鲤管理面板主打安全、可靠、性能调优。容器化管理
// @termsOfService https://github.com/jinlicode/jinli-panel
// @license.name GNU v3
// @license.url https://github.com/jinlicode/jinli-panel/blob/master/LICENSE
func main() {

	// 检测锦鲤面板 安装目录
	if tools.CheckFileExist(global.BASEPATH+"install.lock") == false {

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

		// 后台默认拉基础镜像
		tools.ExecLinuxCommand("nohup docker pull hub.jinli.plus/jinlicode/mysql:latest > /dev/null 2>&1 & ")
		tools.ExecLinuxCommand("nohup docker pull hub.jinli.plus/jinlicode/nginx:v1 > /dev/null 2>&1 & ")
		tools.ExecLinuxCommand("nohup docker pull hub.jinli.plus/jinlicode/memcached:1.6.6 > /dev/null 2>&1 & ")
		tools.ExecLinuxCommand("nohup docker pull hub.jinli.plus/jinlicode/redis:5.0.9 > /dev/null 2>&1 & ")
		tools.ExecLinuxCommand("nohup docker pull hub.jinli.plus/jinlicode/phpmyadmin:5.0.2 > /dev/null 2>&1 & ")

		//创建项目目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH)
		//创建代码目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "code/")
		//创建各配置项目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/cert/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/mysql/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/nginx/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/php/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "config/rewrite/")

		//创建备份目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "backup/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "backup/database/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "backup/site/")

		//创建自动备份目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "autobackup/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "autobackup/database/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "autobackup/site/")

		//创建日志目录
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "log/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "log/nginx/")
		tools.ExecLinuxCommand("mkdir " + global.BASEPATH + "log/openrasp/")

		//设置代码目录为 10000,10000
		tools.ExecLinuxCommand("chown -R 10000:10000 " + global.BASEPATH + "code/")

		//创建mysql my.cnf
		tools.WriteFile(global.BASEPATH+"config/mysql/my.cnf", Template.MysqlCnf())

		//创建 Nginx 网段
		tools.ExecLinuxCommand("docker network create nginx_net")
		//创建nginx
		tools.ExecLinuxCommand("nohup docker run -d --name nginx --network nginx_net --restart always --env TZ=Asia/Shanghai -p 80:80 -p 443:443 -v " + global.BASEPATH + "config/nginx:/etc/nginx/conf.d -v " + global.BASEPATH + "code:/var/www -v " + global.BASEPATH + "log/nginx:/var/log/nginx -v " + global.BASEPATH + "config/cert:/etc/letsencrypt -v " + global.BASEPATH + "config/rewrite:/etc/nginx/rewrite hub.jinli.plus/jinlicode/nginx:v1 > /dev/null 2>&1 & ")

		//创建 Mysql 网段
		tools.ExecLinuxCommand("docker network create mysql_net")
		//自动生成mysql密码
		mysqlRandPassword := tools.RandomString(16)
		tools.ExecLinuxCommand("nohup docker run -d --name mysql --network mysql_net --restart always --env TZ=Asia/Shanghai --env MYSQL_ROOT_PASSWORD=" + mysqlRandPassword + "  -p 3306:3306 -v " + global.BASEPATH + "db:/var/lib/mysql -v " + global.BASEPATH + "imput_db:/docker-entrypoint-initdb.d -v " + global.BASEPATH + "config/mysql/my.cnf:/etc/mysql/my.cnf hub.jinli.plus/jinlicode/mysql > /dev/null 2>&1 & ")

		tools.WriteFile(global.BASEPATH+"install.lock", "installed")

	}

	// 启动gin
	routers.InitRouter()
}
