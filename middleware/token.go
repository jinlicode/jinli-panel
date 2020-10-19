package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.URL)
	}
}
