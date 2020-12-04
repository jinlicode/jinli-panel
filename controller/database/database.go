package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	resp "github.com/jinlicode/jinli-panel/model/response"
)

// GetLists 获取网站列表
func GetLists(c *gin.Context) {

	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)

	err, list, total := model.GetSiteList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}
