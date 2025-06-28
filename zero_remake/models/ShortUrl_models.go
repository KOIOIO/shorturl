package models

import "gorm.io/gorm"

// Shorturl 定义了短链接的数据结构。
// 包含两个字段：Shorturl 和 Url。
//
// Shorturl: 短链接字符串，使用 varchar(20) 类型，不允许为空，用于存储生成的短链接。
// Url: 原始链接字符串，使用 varchar(200) 类型，不允许为空，用于存储需要缩短的原始链接。
type Shorturl struct {
	gorm.Model
	Shorturl string `gorm:"type:varchar(20);not null" json:"shorturl"`
	Url      string `gorm:"type:varchar(200);not null" json:"url"`
	ID       uint64 `gorm:"primaryKey;autoIncrement:false" json:"id"` // 使用雪花ID作为主键
}
