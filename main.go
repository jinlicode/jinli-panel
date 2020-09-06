package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.StaticFS("/static", http.Dir("/jinli/html/static"))
	router.StaticFile("/favicon.ico", "/jinli/html/favicon.ico")
	router.StaticFile("/", "/jinli/html/index.html")

	v1 := router.Group("/v1")
	{
		v1.GET("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "login",
			})
		})
	}

	router.Run("0.0.0.0:9527")

}
