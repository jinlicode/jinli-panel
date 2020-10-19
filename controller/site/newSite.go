package site

import (
	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	"github.com/jinlicode/jinli-panel/utils"
)

func CreateSite(c *gin.Context) {

	var R request.SiteStruct
	_ = c.ShouldBindJSON(&R)

	SiteVerify := make(map[string][]string)
	if R.IsSsl == 1 {
		SiteVerify = utils.Rules{
			"Domain":     {utils.NotEmpty()},
			"Email":      {utils.NotEmpty()},
			"PhpVersion": {utils.NotEmpty()},
		}
	} else {
		SiteVerify = utils.Rules{
			"Domain":     {utils.NotEmpty()},
			"PhpVersion": {utils.NotEmpty()},
		}
	}

	SiteVerifyErr := utils.Verify(R, SiteVerify)
	if SiteVerifyErr != nil {
		response.FailWithMessage(SiteVerifyErr.Error(), c)
		return
	}

	isEx := model.CheckSiteByDomain(R.Domain)
	if isEx == false {
		response.FailWithMessage("域名已存在", c)
		return
	}

	//入库
	model.CreateSite(R)

	response.OkWithData("success", c)

}
