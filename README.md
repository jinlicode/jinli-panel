# 锦鲤管理面板

web管理端

## 开发者

    golang > 1.14
    node > 12
    yarn > 1.22.5

### 启动开发环境

```bash
bash dev-go.sh
bash dev-vue.sh
```

注意需要两个窗口，修改代码后会自动编译，不需要再次编译，方便快速开发

### 静态文件编译

需要执行 `~/go/bin/statik -src=./html/`

### 接口文档地址

    http://xxx.com/swagger/index.html

注释说明地址：

    https://swaggo.github.io/swaggo.io/declarative_comments_format/general_api_info.html