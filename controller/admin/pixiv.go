package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"blog-admin-api/server/spider"
	"strconv"
)

type UpdateImgRequest struct {
	Cookie string `json:"cookie" binding:"required"`
}

func GetPixivs(c *core.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	list, err := spider.List(c.DefaultQuery("next_marker", ""), page)
	if err != nil {
		c.ErrorL("获取图片列表失败", page, err.Error())
		c.FailWithErrCode(errcode.AdminNetworkBusy, nil)
		return
	}

	c.Success(list)
}

func AddPixiv(c *core.Context) {
	var r UpdateImgRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}

	if err := spider.Get(c, r.Cookie); err != nil {
		c.ErrorL("爬虫执行失败", nil, err.Error())
		c.FailWithErrCode(errcode.AdminNetworkBusy, nil)
		return
	}

	c.Success(nil)
}

func DelPixiv(c *core.Context) {
	res, err := spider.Delete(c.Query("name"))
	if err != nil {
		c.ErrorL("删除失败", nil, err.Error())
		c.FailWithErrCode(errcode.AdminNetworkBusy, nil)
		return
	}

	c.Success(res)
}

func GetPixivCount(c *core.Context) {
	res, err := spider.Count()
	if err != nil {
		c.ErrorL("获取数量失败", nil, err.Error())
		c.FailWithErrCode(errcode.AdminNetworkBusy, nil)
		return
	}

	c.Success(res)
}
