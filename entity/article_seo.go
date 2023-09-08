package entity

import (
	"blog-admin-api/pkg/db"
)

type ArticleSeo struct {
	*Base
	ArticleID   int    `json:"article_id"`
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

func (*ArticleSeo) TableName() string {
	return "article_seo"
}

func (a *ArticleSeo) Set(data ArticleSeo) error {
	return db.Orm.Save(&data).Error
}

func (*ArticleSeo) GetByArticleID(articleID int) (*ArticleSeo, error) {
	var articleSelList *ArticleSeo
	return articleSelList, db.Orm.Raw("select * from article_seo where article_id = ?", articleID).First(&articleSelList).Error
}

func (*ArticleSeo) Delete(articleId int) error {
	return db.Orm.Exec("delete from article_seo where article_id = ?", articleId).Error
}
