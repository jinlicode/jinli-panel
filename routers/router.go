package routers

import (
	"net/http"

	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.StaticFS("/static", http.Dir("./html/static"))
	router.StaticFile("/favicon.ico", "./html/favicon.ico")
	router.StaticFile("/", "./html/index.html")

	v1 := router.Group("/v1")
	{
		v1.GET("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "login",
			})
		})
	}
	router.Run("0.0.0.0:9527")
	return router
}
