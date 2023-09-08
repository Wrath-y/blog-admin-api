package entity

import (
	"blog-admin-api/pkg/db"
)

type Comment struct {
	*Base
	Name      string `json:"name" gorm:"name"`
	Email     string `json:"email" gorm:"email"`
	Url       string `json:"url" gorm:"url"`
	Type      int    `json:"type" gorm:"type"`
	Content   string `json:"content" gorm:"content"`
	ArticleId int    `json:"article_id" gorm:"article_id"`
	Pid       int    `json:"pid" gorm:"pid"`
	Ppid      int    `json:"ppid" gorm:"ppid"`
}

func (*Comment) TableName() string {
	return "comment"
}

func (*Comment) FindWithPage(page, limit int) ([]*Comment, int64, error) {
	if limit == 0 {
		limit = 6
	}

	var count int64
	if err := db.Orm.Model(&Comment{}).Count(&count).Error; err != nil {
		return nil, count, err
	}

	comments := make([]*Comment, 0)
	err := db.Orm.Offset((page - 1) * limit).Limit(limit).Find(&comments).Error

	return comments, count, err
}

func (*Comment) Delete(id int) error {
	return db.Orm.Exec("delete from comment where id = ?", id).Error
}

func (*Comment) GetById(id int) (*Comment, error) {
	c := new(Comment)
	return c, db.Orm.First(&c, id).Error
}
