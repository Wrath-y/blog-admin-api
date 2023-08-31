package entity

import (
	"blog-admin-api/pkg/db"
)

type Article struct {
	*Base
	Title  string `json:"title"`
	Image  string `json:"image"`
	Intro  string `json:"intro"`
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
	return db.Orm.Create(a).Error
}

func (a *Article) Update(id int) error {
	a.Id = id

	return db.Orm.Model(a).Save(a).Error
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
