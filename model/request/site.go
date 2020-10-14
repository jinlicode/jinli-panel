package request

// site SiteStruct
type SiteStruct struct {
	ID         int    `json:"id"`
	Domain     string `json:"domain"`
	Email      string `json:"email"`
	PhpVersion string `json:"php_version"`
	IsSsl      int64  `json:"is_ssl"`
	Status     int64  `json:"status"`
}
