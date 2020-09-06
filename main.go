package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.StaticFS("/static", http.Dir("./html/static"))
	router.StaticFile("/favicon.ico", "./html/favicon.ico")
	router.StaticFile("/", "./html/index.html")

	v1 := router.Group("/v1")
    {
        v1.GET("/login", func(c *gin.Context){
			c.JSON(200, gin.H{
				"message": "login",
			})
		})
    }

	router.Run("0.0.0.0:8080")
}
