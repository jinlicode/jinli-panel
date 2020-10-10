package routers

import (
	"net/http"

	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/controller/auth"
	_ "github.com/jinlicode/jinli-panel/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("/static", http.Dir("./html/static"))
	router.StaticFile("/favicon.ico", "./html/favicon.ico")
	router.StaticFile("/", "./html/index.html")

	v1 := router.Group("/v1/user")
	{
		v1.POST("/login", auth.Login)
		v1.GET("/info", auth.Info)
	}

	router.Run("0.0.0.0:9527")
	return router
}
