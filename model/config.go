package model

// Config  Struct
type Config struct {
	ID       int    `json:"id"`
	Mysqlpwd string `json:"mysqlpwd"`
}

// SetConfigMsqlpwd 设置数据库密码
func SetConfigMsqlpwd(mysqlpwd string) (err error) {
	config := Config{ID: 1}
	err = db.Model(&config).Update("mysqlpwd", mysqlpwd).Error
	return err
}

// GetConfigInfo 获取配置内容
func GetConfigInfo() (configInfo Config, err error) {
	var config Config
	err = db.First(&config).Error
	return config, err
}
