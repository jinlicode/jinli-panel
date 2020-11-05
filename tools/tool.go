package tools

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// RandomString 返回随机字符串
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune
	defaultLetters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// ReadMysqlRootPassword mysql root密码
func ReadMysqlRootPassword(basepath string) string {
	return ""
}

// ReadMysqlHost  读取compose文件中的mysql host
func ReadMysqlHost(basepath string) string {
	return ""
}

// ReadSiteMysqlInfo  读取compose文件中的网站 mysql 信息
func ReadSiteMysqlInfo(basepath string, dockerName string, readType string) string {
	return ""
}

// GetPathFiles 获取木下的所有文件切片
func GetPathFiles(path string, isDir bool) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	fs := []string{}
	for _, f := range files {
		if isDir == true {
			if f.IsDir() == true {
				fs = append(fs, f.Name())
			}
		} else {
			fs = append(fs, f.Name())
		}
	}
	return fs
}

//PHPChooseVersion php版本选择
func PHPChooseVersion() string {
	return ""
}
