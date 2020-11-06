package site

import (
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	"github.com/jinlicode/jinli-panel/utils"
)

func DelSite(c *gin.Context) {

	var R request.Site
	_ = c.ShouldBindJSON(&R)

	SiteVerify := utils.Rules{
		"ID": {utils.NotEmpty()},
	}

	SiteVerifyErr := utils.Verify(R, SiteVerify)
	if SiteVerifyErr != nil {
		response.FailWithMessage(SiteVerifyErr.Error(), c)
		return
	}

	model.DelSite(R)

	response.OkWithData("success", c)
}
