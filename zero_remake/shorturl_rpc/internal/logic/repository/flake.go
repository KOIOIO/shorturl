package repository

import "github.com/sony/sonyflake"

// flake 是一个全局的雪花ID生成器实例
var Flake = sonyflake.NewSonyflake(sonyflake.Settings{})
