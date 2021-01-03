package request

// site SiteStruct
type Site struct {
	ID          int    `json:"id" gorm:"primarykey"`
	Domain      string `json:"domain"`
	Email       string `json:"email"`
	PhpVersion  string `json:"php_version"`
	RewriteConf string `json:"rewrite_conf"`
	HostConf    string `json:"host_conf"`
	IsSsl       int64  `json:"is_ssl"`
	Status      int64  `json:"status"`
	Addtime     string `json:"addtime"`
}
