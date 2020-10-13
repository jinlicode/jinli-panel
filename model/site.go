package model

import "github.com/jinlicode/jinli-panel/model/request"

type Site struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Email      string `json:"email"`
	PhpVersion string `json:"php_version"`
	IsSsl      int64  `json:"is_ssl"`
	Status     int64  `json:"status"`
}

// GetSiteList

func GetSiteList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var site []Site
	err = db.Limit(limit).Offset(offset).Find(&site).Error
	return err, site, total
}
