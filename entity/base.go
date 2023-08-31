package entity

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Base struct {
	Id         int       `json:"id"`
	UpdateTime time.Time `json:"UpdateTime"`
	CreateTime time.Time `json:"create_time"`
}
