#!/bin/bash

source /etc/os-release

## Ubuntu 20.04安装环境
function ubuntu_2004_install {
    echo "Ubuntu"
}

## centos7安装环境
function centos_7_install {
#增加必要安装包
sudo yum install -y curl wget unzip lrzsz lsof
#安装源
yum install epel-release -y
#关闭防火墙
systemctl stop firewalld.service && systemctl disable firewalld.service && setenforce 0 && sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config

#检测端口
#lsof -i:80
#lsof -i:443
#lsof -i:3306
#lsof -i:9527

#执行安装docker
#step 1: 安装必要的一些系统工具
sudo yum install -y yum-utils device-mapper-persistent-data lvm2 git epel-*
#Step 2: 添加软件源信息
sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#Step 3: 更新并安装 Docker-CE
sudo yum makecache fast && sudo yum -y install docker-ce-19.03.12 docker-compose-1.18.0
#Step 4: 开启Docker服务
sudo systemctl start docker
#step 5: 设置开机启动
sudo systemctl enable docker

#重载docker
sudo systemctl daemon-reload && sudo systemctl restart docker

#创建 Nginx 网段
docker network create nginx_net
#创建 Mysql 网段
docker network create mysql_net
}

## 公共执行函数
function public {
#创建项目目录
mkdir -p /var/jinli_panel/
#创建代码目录
mkdir -p /var/jinli_panel/code/
#设置代码目录10000
chown -R 10000:10000 /var/jinli_panel/code/

#创建各配置项目录
mkdir -p /var/jinli_panel/config/cert
mkdir -p /var/jinli_panel/config/mysql
mkdir -p /var/jinli_panel/config/nginx
mkdir -p /var/jinli_panel/config/php
mkdir -p /var/jinli_panel/config/rewrite

#创建备份目录
mkdir -p /var/jinli_panel/backup/
mkdir -p /var/jinli_panel/backup/database
mkdir -p /var/jinli_panel/backup/site

#创建自动备份目录
mkdir -p /var/jinli_panel/autobackup/
mkdir -p /var/jinli_panel/autobackup/database
mkdir -p /var/jinli_panel/autobackup/site

#创建日志目录
mkdir -p /var/jinli_panel/log/
mkdir -p /var/jinli_panel/log/nginx
mkdir -p /var/jinli_panel/log/openrasp
}

function last {

}

function check_port {
    echo "正在检测端口......"
    netstat -tlpn | grep "\b$1\b"
}

case $ID in
ubuntu)
    if [ $VERSION_ID == "20.04" ];then
        ubuntu_2004_install
        exit 1
    fi
        echo $VERSION_ID
        exit 1
    ;;
centos)
    yumdnf="yum"
    if test "$(echo "$VERSION_ID >= 22" | bc)" -ne 0; then
        yumdnf="dnf"
    fi
    sudo $yumdnf install -y redhat-lsb-core
    ;;
*)
    echo "系统发行版本不匹配，请更换官方支持系统。"
    ;;
esac
