package entity

import (
	"blog-admin-api/pkg/db"
)

type Administrator struct {
	*Base
	Account  string `json:"account"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func (*Administrator) TableName() string {
	return "administrator"
}

func GetUserByName(account string) (*Administrator, error) {
	a := new(Administrator)

	return a, db.Orm.Where("account = ?", account).First(&a).Error
}
