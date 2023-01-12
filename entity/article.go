package entity

import (
	"blog-admin-api/pkg/db"
	"time"
)

type Article struct {
	Base
	Title  string `json:"title"`
	Image  string `json:"image"`
	Html   string `json:"html"`
	Con    string `json:"con"`
	Tags   string `json:"tags"`
	Hits   string `json:"hits"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func (*Article) TableName() string {
	return "article"
}

func (a *Article) Create() error {
	a.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return db.Orm.Create(a).Error
}

func (a *Article) Update(id int) error {
	a.Id = id
	a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	return db.Orm.Model(a).Updates(a).Error
}

func (*Article) Delete(id int) error {
	a := Article{}
	a.Id = id

	return db.Orm.Delete(a).Error
}

func (*Article) FindWithPage(page, limit int) ([]*Article, int64, error) {
	var count int64
	if err := db.Orm.Model(&Article{}).Count(&count).Error; err != nil {
		return nil, count, err
	}

	articles := make([]*Article, 0)
	err := db.Orm.Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&articles).Error

	return articles, count, err
}

func (*Article) GetById(id int) (*Article, error) {
	articles := new(Article)
	return articles, db.Orm.First(&articles, id).Error
}
