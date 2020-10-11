package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinlicode/jinli-panel/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int
	Name       string
	Password   string
	Token      string
	ExpireTime int64
	FailNum    int64
	FailTime   int64
}

func GetAtuh() {
	// Atuh := db.First(&user)
	// return Atuh
}

// DoLogin 登录操作
func DoLogin(name string, password string) (int, string) {
	var user User

	// 从新生成密码
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	encodePW := string(hash)
	fmt.Println(encodePW)

	// 检测帐号密码是否正确
	db.First(&user).Scan(&user)

	nowTime := time.Now().Unix()

	//判断是否过期错误次数超期
	if nowTime-user.FailTime < 900 && user.FailNum >= 5 {
		return -3, "" //超期登录
	}

	if user.Name != name {
		// 如果不正确直接增加一次错误登录 更新一次错误时间
		db.Model(User{}).Where("id = ?", user.ID).Updates(User{FailTime: nowTime, FailNum: user.FailNum + 1})
		return -1, "" //帐号不存在
	}

	// 正确密码验证
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 如果不正确直接增加一次错误登录 更新一次错误时间
		db.Model(User{}).Where("id = ?", user.ID).Updates(User{FailTime: nowTime, FailNum: user.FailNum + 1})
		return -2, "" //密码错误
	}

	// 生成token
	tokenMd5 := utils.MD5V([]byte(strconv.Itoa(user.ID) + encodePW + strconv.FormatInt(nowTime, 10)))

	// 如果正确设置token值 和过期时间
	db.Model(User{}).Select("expire_time", "fail_time", "fail_num", "token").Where("id = ?", user.ID).Updates(User{ExpireTime: nowTime + 86400, FailTime: 0, FailNum: 0, Token: tokenMd5})

	return user.ID, tokenMd5
}

// GetInfo 获取用户信息
func GetInfo(token string) User {
	// 检测帐号密码是否正确
	var user User

	db.Where("token = ?", token).First(&user).Scan(&user)

	return user
}

// Logout 退出登录 清空token
func Logout(token string) bool {
	// 检测帐号密码是否正确
	db.Model(User{}).Select("token", "fail_time", "fail_num").Where("token = ?", token).Updates(User{Token: "", FailTime: 0, FailNum: 0})

	return true
}

// CheckToken 检测token是否过期
func CheckToken() bool {
	return true
}
