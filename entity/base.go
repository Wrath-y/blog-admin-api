package entity

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Base struct {
	Id        int       `json:"id"`
	UpdatedAt time.Time `json:"update_time"`
	CreatedAt time.Time `json:"create_time"`
}
