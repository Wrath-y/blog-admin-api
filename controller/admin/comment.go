package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"strconv"
)

func GetComments(c *core.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	data, count, err := new(entity.Comment).FindWithPage(page, 15)
	if err != nil {
		c.ErrorL("获取评论列表失败", page, err.Error())
		c.FailWithErrCode(errcode.CommentGetFailed, nil)
		return
	}

	c.Success(map[string]interface{}{
		"list":  data,
		"count": count,
	})
}

func DelComment(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := new(entity.Comment).Delete(id); err != nil {
		c.ErrorL("删除评论失败", id, err.Error())
		c.FailWithErrCode(errcode.CommentDelFailed, nil)
	}

	c.Success(nil)
}
