FROM centos:7

# 安装基础环境
RUN yum install -y yum-utils \
    device-mapper-persistent-data \
    lvm2 \
    epel-release && \
    yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo &&\
    yum makecache fast && \
    yum -y install docker-ce-19.03.12 docker-compose-1.18.0