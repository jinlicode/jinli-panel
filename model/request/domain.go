package request

// Domain SiteStruct
type Domain struct {
	ID      int    `json:"id" gorm:"primarykey"`
	Pid     int    `json:"pid"`
	Name    string `json:"name"`
	Addtime string `json:"addtime"`
}
