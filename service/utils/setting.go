package utils

import (
	"fmt"
	"github.com/go-ini/ini"
)

type Server struct {
	AppMode  string
	HttpPort string
}

// DbHost 数据库主机
type Database struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
}

type Redis struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

var ServerConfig = Server{}
var DatabaseConfig = Database{}
var RedisConfig = Redis{}

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
	ServerConfig.AppMode = file.Section("server").Key("AppMode").MustString("debug")
	ServerConfig.HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
}

func LoadData(file *ini.File) {
	DatabaseConfig.DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DatabaseConfig.DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DatabaseConfig.DbUser = file.Section("database").Key("DbUser").MustString("root")
	DatabaseConfig.DbPassWord = file.Section("database").Key("DbPassWord").MustString("wwy040609")
	DatabaseConfig.DbName = file.Section("database").Key("DbName").MustString("shorturl")
}
func LoadRedis(file *ini.File) {
	RedisConfig.RedisHost = file.Section("redis").Key("RedisHost").MustString("localhost")
	RedisConfig.RedisPort = file.Section("redis").Key("RedisPort").MustString("6379")
	RedisConfig.RedisPassword = file.Section("redis").Key("RedisPassword").MustString("")
}
