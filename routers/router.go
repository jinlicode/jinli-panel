package routers

import (
	"fmt"

	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/controller/auth"
	"github.com/jinlicode/jinli-panel/controller/database"
	"github.com/jinlicode/jinli-panel/controller/site"
	"github.com/jinlicode/jinli-panel/controller/soft"
	_ "github.com/jinlicode/jinli-panel/docs"
	"github.com/jinlicode/jinli-panel/middleware"
	"github.com/rakyll/statik/fs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/jinlicode/jinli-panel/statik"
)

func InitRouter() *gin.Engine {

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.TokenAuth())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	statikFS, err := fs.New()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		router.StaticFS("/", statikFS)
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/user/login", auth.Login)
		v1.GET("/user/info", auth.Info)
		v1.POST("/user/logout", auth.Logout)

		v1.GET("/site/list", site.GetLists)
		v1.POST("/site/create", site.CreateSite)
		v1.POST("/site/delete", site.DelSite)

		v1.GET("/site/get_conf", site.GetSiteConf)
		v1.POST("/site/update_conf", site.UpdateSiteConf)

		v1.GET("/site/get_rewrite", site.GetSiteRewrite)
		v1.POST("/site/update_rewrite", site.UpdateSiteRewrite)

		v1.GET("/site/get_php", site.GetSitePhp)
		v1.POST("/site/update_php", site.UpdateSitePhp)

		v1.GET("/site/get_domain", site.GetSiteDomain)
		v1.POST("/site/update_domain", site.UpdateSiteDomain)
		v1.POST("/site/del_domain", site.DelSiteDomain)

		v1.GET("/site/get_basepath", site.GetSiteBasepath)
		v1.POST("/site/update_basepath", site.UpdateSiteBasepath)

		v1.POST("/site/update_status", site.UpdateSiteStatus)

		v1.GET("/soft/list", soft.GetSoftList)
		v1.GET("/soft/phplist", soft.GetPHPList)
		v1.POST("/soft/install", soft.InstallSoft)

		v1.GET("/database/list", database.GetLists)

	}

	router.Run("0.0.0.0:9527")
	return router
}
