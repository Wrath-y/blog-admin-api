package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"strconv"
)

type ArticleRequest struct {
	Title  string `json:"title" binding:"required"`
	Image  string `json:"image"`
	Html   string `json:"html" binding:"required"`
	Con    string `json:"con" binding:"required"`
	Tags   string `json:"tags" binding:"required"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func AddArticle(c *core.Context) {
	var r ArticleRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, err.Error())
		return
	}

	res := &entity.Article{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
	}
	if err := res.Create(); err != nil {
		c.ErrorL("创建文章失败", res, nil)
		c.FailWithErrCode(errcode.ArticleCreateFailed, nil)
		return
	}

	c.Success(res)
}

func DelArticle(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := new(entity.Article).Delete(id); err != nil {
		c.FailWithErrCode(errcode.ArticleDelFailed, nil)
		return
	}

	c.Success(nil)
}

func UpdateArticle(c *core.Context) {
	var (
		r      ArticleRequest
		logMap = make(map[string]interface{})
	)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}
	logMap["id"] = id
	data := &entity.Article{
		Title: r.Title,
		Image: r.Image,
		Html:  r.Html,
		Con:   r.Con,
		Tags:  r.Tags,
	}
	logMap["data"] = data
	if err := data.Update(id); err != nil {
		c.ErrorL("更新文章失败", data, err.Error())
		c.FailWithErrCode(errcode.ArticleUpdateFailed, nil)
		return
	}

	c.Success(data)
}

func GetArticles(c *core.Context) {
	logMap := make(map[string]interface{})

	pageStr := c.DefaultQuery("page", "0")
	logMap["pageStr"] = pageStr
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.ErrorL("解析参数失败", logMap, err.Error())
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}
	logMap["page"] = page

	articles, count, err := new(entity.Article).FindWithPage(page, 6)
	if err != nil {
		c.ErrorL("获取文章列表失败", logMap, err.Error())
		return
	}

	c.Success(map[string]interface{}{
		"list":  articles,
		"count": count,
	})
}

func GetArticle(c *core.Context) {
	logMap := make(map[string]interface{})

	idStr := c.Param("id")
	logMap["idStr"] = idStr
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.ErrorL("解析参数失败", logMap, err.Error())
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}
	logMap["id"] = id

	res, err := new(entity.Article).GetById(id)
	if err != nil {
		c.ErrorL("获取文章失败", logMap, err.Error())
		c.FailWithErrCode(errcode.ArticleGetFailed, nil)
		return
	}

	c.Success(res)
}
