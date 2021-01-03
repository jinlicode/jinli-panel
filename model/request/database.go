package request

// Database Database
type Database struct {
	ID       int    `json:"id" gorm:"primarykey"`
	Pid      int    `json:"pid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Addtime  string `json:"addtime"`
}
