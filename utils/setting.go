package utils

import (
	"fmt"
	"github.com/go-ini/ini"
)

var AppMode string
var HttpPort string

// DbHost 数据库主机
var DbHost string
var DbPort string
var DbUser string
var DbPassWord string
var DbName string

var RedisHost string
var RedisPort string
var RedisPassword string

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
	LoadRedis(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
}

func LoadData(file *ini.File) {
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("wwy040609")
	DbName = file.Section("database").Key("DbName").MustString("shorturl")
}
func LoadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").MustString("localhost")
	RedisPort = file.Section("redis").Key("RedisPort").MustString("6379")
	RedisPassword = file.Section("redis").Key("RedisPassword").MustString("")
}
