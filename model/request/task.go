package request

// Task TaskStruct
type Task struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Start   int64  `json:"start"`
	End     int64  `json:"end"`
	Addtime string `json:"addtime"`
	Execstr string `json:"execstr"`
}
