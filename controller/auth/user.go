package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
)

func Login(c *gin.Context) {
	var R request.LoginStruct
	_ = c.ShouldBindJSON(&R)

	username := R.Username
	password := R.Password

	id, err := model.DoLogin(username, password)

	if id == -1 {
		response.FailWithMessage("帐号不存在", c)
	} else if id == -2 {
		response.FailWithMessage("密码错误", c)
	} else if id == -3 {
		response.FailWithMessage("登录错误太多，请过15分钟在尝试", c)
	} else {
		data := make(map[string]string)
		data["token"] = err
		response.OkWithData(data, c)
	}
}

// Info 检测token是否过期
func Info(c *gin.Context) {
	token := c.Query("token")
	user := model.GetInfo(token)

	if user.ID == 0 {
		response.FailWithMessage("帐号已过期，请重新登录", c)
	} else {
		data := make(map[string]interface{})

		data["roles"] = []string{"admin"}
		data["introduction"] = user.Name
		data["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
		data["name"] = user.Name
		response.OkWithData(data, c)
	}

}

func Logout(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	model.Logout(token)
	response.OkWithData("success", c)

}
