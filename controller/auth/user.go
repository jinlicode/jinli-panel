package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Login(c *gin.Context) {
	var R request.LoginStruct
	_ = c.ShouldBindJSON(&R)

	username := R.Username
	password := R.Password

	id, err := model.DoLogin(username, password)

	if id == -1 {
		c.JSON(200, gin.H{
			"code":    60204,
			"message": "帐号不存在",
		})
	} else if id == -2 {
		c.JSON(200, gin.H{
			"code":    60204,
			"message": "密码错误",
		})
	} else if id == -3 {
		c.JSON(200, gin.H{
			"code":    60204,
			"message": "登录错误太多，请过15分钟在尝试",
		})
	} else {
		data := make(map[string]string)
		data["token"] = err
		c.JSON(200, Response{
			20000,
			data,
			"OK",
		})
	}
}

// Info 检测token是否过期
func Info(c *gin.Context) {
	token := c.Query("token")
	user := model.GetInfo(token)

	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code":    50008,
			"message": "帐号已过期，请重新登录",
		})
	} else {
		data := make(map[string]interface{})

		data["roles"] = []string{"admin"}
		data["introduction"] = user.Name
		data["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
		data["name"] = user.Name

		c.JSON(200, Response{
			20000,
			data,
			"OK",
		})
	}

}

func Logout(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	model.Logout(token)

	c.JSON(200, gin.H{
		"code": 20000,
		"data": "success",
	})
}
