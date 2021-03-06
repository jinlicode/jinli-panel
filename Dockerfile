FROM alpine:3.12

# 安装基础环境
RUN apk add --no-cache  docker-cli=19.03.12-r0 docker-compose=1.25.4-r2

COPY html /jinli/html
COPY jinli-panel /jinli/jinli-panel
COPY run.sh /run.sh
RUN chmod +x /run.sh
EXPOSE 9527
WORKDIR /jinli
CMD ["/run.sh"]