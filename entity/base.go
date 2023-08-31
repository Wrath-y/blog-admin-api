package entity

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Base struct {
	Id         int       `json:"id"`
	UpdateTime time.Time `json:"update_time"`
	CreateTime time.Time `json:"create_time"`
}
