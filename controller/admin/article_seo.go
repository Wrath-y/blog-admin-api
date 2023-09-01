package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"blog-admin-api/service/article"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ArticleSeoRequest struct {
	entity.Base
	ArticleID   int    `json:"article_id"`
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

func SetArticleSeo(c *core.Context) {
	var r *ArticleSeoRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}

	c.Info("设置文章SEO", r, "")

	data := entity.ArticleSeo{
		Base: &entity.Base{
			Id:         r.Id,
			UpdateTime: time.Now().In(c.TimeLocation),
			CreateTime: r.CreateTime,
		},
		ArticleID:   r.ArticleID,
		Title:       r.Title,
		Keywords:    r.Keywords,
		Description: r.Description,
	}

	if err := new(entity.ArticleSeo).Set(data); err != nil {
		c.ErrorL("创建文章SEO失败", data, err.Error())
		c.FailWithErrCode(errcode.ArticleSeoSetFailed, nil)
		return
	}

	if err := article.DelById(r.ArticleID); err != nil {
		c.ErrorL("删除seo缓存失败", r, err.Error())
	}

	c.Success(nil)
}

func GetArticleSeo(c *core.Context) {
	articleID, _ := strconv.Atoi(c.Param("article_id"))
	data, err := new(entity.ArticleSeo).GetByArticleID(articleID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.ErrorL("获取seo列表失败", articleID, err.Error())
		c.FailWithErrCode(errcode.ArticleSeoGetFailed, nil)
		return
	}

	c.Success(data)
}
