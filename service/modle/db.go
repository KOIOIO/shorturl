package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"shorturl/utils"

	"time"
)

// Db 用于全局访问数据库连接
var Db *gorm.DB

// err 用于全局错误处理
var err error

// InitDb 初始化数据库连接
func InitDb() {
	// 构建数据库连接字符串
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DatabaseConfig.DbUser,
		utils.DatabaseConfig.DbPassWord,
		utils.DatabaseConfig.DbHost,
		utils.DatabaseConfig.DbPort,
		utils.DatabaseConfig.DbName,
	)
	// 打开数据库连接并配置gorm
	Db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})

	// 检查数据库连接是否成功
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数：", err)
		os.Exit(1)
	}

	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	// 注意:初次运行后可注销此行
	_ = Db.AutoMigrate(&Shorturl{})

	// 获取底层sql.DB对象
	sqlDB, _ := Db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
