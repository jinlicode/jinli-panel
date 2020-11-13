package model

import (
	"time"

	"github.com/jinlicode/jinli-panel/model/request"
)

// GetSiteList
func GetSiteList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var site []request.Site
	err = db.Limit(limit).Offset(offset).Find(&site).Error
	return err, site, total
}

// CheckSiteByDomain
func CheckSiteByDomain(name string) bool {
	var domain request.Domain

	db.Where("name = ?", name).First(&domain).Scan(&domain)

	if domain.ID == 0 {
		return true
	}

	return false
}

// CreateSite
func CreateSite(info request.Site) (id int, err error) {
	site := request.Site{Domain: info.Domain, Email: info.Email, PhpVersion: info.PhpVersion, IsSsl: info.IsSsl, Status: 0, Addtime: time.Now().Format("2006-01-02 15:04:05")}

	err = db.Create(&site).Error
	domain := request.Domain{Name: info.Domain, Pid: 0, Addtime: time.Now().Format("2006-01-02 15:04:05")}
	domain.ID = site.ID
	err = db.Create(&domain).Error
	return site.ID, err
}

// DelSite
func DelSite(info request.Site) (err error) {
	site := request.Site{ID: info.ID}
	err = db.Where("id = ?", site.ID).Delete(&site).Error
	return err
}

// SetSiteStatus
func SetSiteStatus(siteid int, status string) (err error) {
	site := request.Site{ID: siteid}
	err = db.Model(&site).Update("status", status).Error
	return err
}
