package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"blog-admin-api/service/article"
	"strconv"
	"time"
)

type ArticleRequest struct {
	Title  string `json:"title" binding:"required"`
	Image  string `json:"image"`
	Intro  string `json:"intro"`
	Html   string `json:"html" binding:"required"`
	Con    string `json:"con" binding:"required"`
	Tags   string `json:"tags" binding:"required"`
	Status int    `json:"status"`
	Source int    `json:"source"`
}

func AddArticle(c *core.Context) {
	var (
		r      ArticleRequest
		logMap = make(map[string]interface{})
	)
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, err.Error())
		return
	}
	logMap["r"] = r

	res := &entity.Article{
		Title:  r.Title,
		Image:  r.Image,
		Intro:  r.Intro,
		Html:   r.Html,
		Con:    r.Con,
		Tags:   r.Tags,
		Status: r.Status,
		Base: &entity.Base{
			UpdateTime: time.Now().In(c.TimeLocation),
			CreateTime: time.Now().In(c.TimeLocation),
		},
	}
	logMap["res"] = res
	if err := res.Create(); err != nil {
		c.ErrorL("创建文章失败", logMap, err.Error())
		c.FailWithErrCode(errcode.ArticleCreateFailed, nil)
		return
	}

	if err := article.DelList(); err != nil {
		c.ErrorL("删除缓存失败", logMap, err.Error())
	}

	c.Success(res)
}

func DelArticle(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	logMap := make(map[string]interface{})
	logMap["id"] = id
	if err := new(entity.Article).Delete(id); err != nil {
		c.ErrorL("删除失败", logMap, err.Error())
		c.FailWithErrCode(errcode.ArticleDelFailed, nil)
		return
	}

	if err := new(entity.ArticleSeo).Delete(id); err != nil {
		c.ErrorL("删除seo失败", logMap, err.Error())
	}

	if err := article.DelById(id); err != nil {
		c.ErrorL("删除缓存失败", logMap, err.Error())
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

	articleInfo, err := new(entity.Article).GetById(id)
	if err != nil {
		c.ErrorL("获取文章失败", id, err.Error())
		c.FailWithErrCode(errcode.ArticleGetFailed, nil)
		return
	}
	logMap["articleInfo"] = articleInfo

	data := &entity.Article{
		Title:  r.Title,
		Image:  r.Image,
		Intro:  r.Intro,
		Html:   r.Html,
		Con:    r.Con,
		Status: r.Status,
		Tags:   r.Tags,
		Base: &entity.Base{
			CreateTime: articleInfo.CreateTime,
			UpdateTime: time.Now().In(c.TimeLocation),
		},
	}
	logMap["data"] = data
	if err := data.Update(id); err != nil {
		c.ErrorL("更新文章失败", data, err.Error())
		c.FailWithErrCode(errcode.ArticleUpdateFailed, nil)
		return
	}

	if err := article.DelList(); err != nil {
		c.ErrorL("删除缓存失败", logMap, err.Error())
	}
	if err := article.DelById(id); err != nil {
		c.ErrorL("删除缓存失败", logMap, err.Error())
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
