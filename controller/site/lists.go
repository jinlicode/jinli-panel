package site

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinlicode/jinli-panel/global/response"
	"github.com/jinlicode/jinli-panel/model"
	"github.com/jinlicode/jinli-panel/model/request"
	resp "github.com/jinlicode/jinli-panel/model/response"
)

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

	// data := make(map[string]string)

	// c.JSON(200, Response{
	// 	20000,
	// 	data,
	// 	"OK",
	// })

	// {
	// 	"code": 20000,
	// 	"data": {
	// 		"total": 100,
	// 		"items": [
	// 			{
	// 				"id": 1,
	// 				"php_version": "5.6",
	// 				"url": "siytsdp.ne",
	// 				"email": "b.ehmoqvpuk@oimzbi.gov.cn",
	// 				"status": 0,
	// 				"is_ssl": 1
	// 			}
	// 		]
	// 	}
	// }

}
