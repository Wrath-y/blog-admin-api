package entity

import (
	"blog-admin-api/pkg/db"
)

type ArticleSeo struct {
	*Base
	ArticleID int    `json:"article_id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
}

func (*ArticleSeo) TableName() string {
	return "article_seo"
}

func (a *ArticleSeo) Set(details []*ArticleSeo) error {
	if len(details) == 0 {
		return nil
	}
	tx := db.Orm.Begin()
	if err := tx.Exec("delete from article_seo where article_id = ?", details[0].ArticleID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range details {
		if err := tx.Create(v).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (*ArticleSeo) FindByArticleID(articleID int) ([]*ArticleSeo, error) {
	var articleSelList []*ArticleSeo
	return articleSelList, db.Orm.Raw("select * from article_seo where article_id = ?", articleID).Find(&articleSelList).Error
}
