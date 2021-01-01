package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlMap := make(map[string]string)
		urlMap["/"] = "/"
		urlMap["/v1/user/login"] = "/v1/user/login"
		urlMap["/v1/user/logout"] = "/v1/user/logout"

		curURL := c.Request.URL.Path

		//验证token
		if _, ok := urlMap[curURL]; ok == false && strings.Index(curURL, "/v1") != -1 {
			token := c.Request.Header.Get("X-Token")

			if token == "" {
				response.FailWithMessage("帐号已过期，请重新登录-1", c)
				c.Abort()
				return
			}

			user := model.GetInfo(token)

			if user.ID == 0 {
				response.FailWithMessage("帐号已过期，请重新登录", c)
				c.Abort()
				return
			}

		}
	}
}
