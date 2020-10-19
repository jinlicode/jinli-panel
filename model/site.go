package model

import (
	"github.com/jinlicode/jinli-panel/model/request"
)

type Site struct {
	ID         int    `json:"id"`
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

// CheckSiteByDomain
func CheckSiteByDomain(domain string) bool {
	var site Site

	db.Where("domain = ?", domain).First(&site).Scan(&site)

	if site.ID == 0 {
		return true
	}

	return false
}

// CreateSite
func CreateSite(info request.SiteStruct) (err error) {
	site := Site{Domain: info.Domain, Email: info.Email, PhpVersion: info.PhpVersion, IsSsl: info.IsSsl, Status: 0}

	err = db.Create(&site).Error
	return err
}

// DelSite
func DelSite(info request.SiteStruct) (err error) {
	site := Site{ID: info.ID}
	err = db.Where("id = ?", site.ID).Delete(&site).Error
	return err
}
