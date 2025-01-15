package model

import "gorm.io/gorm"

type Shorturl struct {
	gorm.Model
	Shorturl string `gorm:"type:varchar(20);not null" json:"shorturl"`
	Url      string `gorm:"type:varchar(200);not null" json:"url"`
}
