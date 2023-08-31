package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"blog-admin-api/service/articleseo"
	"strconv"
	"time"
)

type ArticleSeoRequest struct {
	ArticleID int `json:"article_id"`
	Details   []*ArticleSeoReqDetail
}

type ArticleSeoReqDetail struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func SetArticleSeo(c *core.Context) {
	var r *ArticleSeoRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}

	details := make([]*entity.ArticleSeo, 0, len(r.Details))
	for _, v := range r.Details {
		details = append(details, &entity.ArticleSeo{
			ArticleID: r.ArticleID,
			Name:      v.Name,
			Content:   v.Content,
			Base: &entity.Base{
				UpdateTime: time.Now().In(c.TimeLocation),
				CreateTime: time.Now().In(c.TimeLocation),
			},
		})

	}

	if err := new(entity.ArticleSeo).Set(details); err != nil {
		c.ErrorL("创建文章SEO失败", details, err.Error())
		c.FailWithErrCode(errcode.ArticleSeoSetFailed, nil)
		return
	}

	if err := articleseo.DelList(r.ArticleID); err != nil {
		c.ErrorL("删除文章SEO失败", r, err.Error())
	}

	c.Success(nil)
}

func GetArticleSeo(c *core.Context) {
	articleID, _ := strconv.Atoi(c.Param("article_id"))
	data, err := new(entity.ArticleSeo).FindByArticleID(articleID)
	if err != nil {
		c.ErrorL("获取seo列表失败", articleID, err.Error())
		c.FailWithErrCode(errcode.ArticleSeoGetFailed, nil)
		return
	}

	c.Success(map[string]interface{}{
		"list": data,
	})
}
