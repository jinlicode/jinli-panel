package routers

import (
	"fmt"
	"net/http"

	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/controller/auth"
	"github.com/jinlicode/jinli-panel/controller/site"
	_ "github.com/jinlicode/jinli-panel/docs"
	"github.com/jinlicode/jinli-panel/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.TokenAuth())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("/static", http.Dir("./html/static"))
	router.StaticFile("/favicon.ico", "./html/favicon.ico")
	router.StaticFile("/", "./html/index.html")

	fmt.Println(router)

	v1 := router.Group("/v1")
	{
		v1.POST("/user/login", auth.Login)
		v1.GET("/user/info", auth.Info)
		v1.POST("/user/logout", auth.Logout)

		v1.GET("/site/list", site.GetLists)
		v1.POST("/site/create", site.CreateSite)
		v1.POST("/site/delete", site.DelSite)
		v1.GET("/site/edit", site.EditSite)

	}

	router.Run("0.0.0.0:9527")
	return router
}
