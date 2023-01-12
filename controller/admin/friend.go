package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"strconv"
	"time"
)

type FriendRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Url   string `json:"url" binding:"required"`
}

func AddFriend(c *core.Context) {
	var r FriendRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, nil)
		return
	}

	res := &entity.FriendLink{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
		Base: &entity.Base{
			UpdatedAt: time.Now().In(c.TimeLocation),
			CreatedAt: time.Now().In(c.TimeLocation),
		},
	}
	if err := res.Create(); err != nil {
		c.ErrorL("创建友链失败", res, err.Error())
		c.FailWithErrCode(errcode.FriendLinkCreateFailed, nil)
		return
	}

	c.Success(res)
}

func DelFriend(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := new(entity.FriendLink).Delete(id); err != nil {
		c.ErrorL("删除友链失败", id, err.Error())
		c.FailWithErrCode(errcode.FriendLinkDelFailed, nil)
		return
	}

	c.Success(nil)
}

func UpdateFriend(c *core.Context) {
	var (
		r      FriendRequest
		logMap = make(map[string]interface{})
	)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidSign, nil)
		return
	}
	logMap["id"] = id
	logMap["req"] = r

	comment, err := new(entity.FriendLink).GetById(id)
	if err != nil {
		c.ErrorL("获取友链失败", logMap, err.Error())
		c.FailWithErrCode(errcode.FriendLinkGetFailed, nil)
		return
	}

	res := &entity.FriendLink{
		Name:  r.Name,
		Email: r.Email,
		Url:   r.Url,
		Base: &entity.Base{
			UpdatedAt: time.Now().In(c.TimeLocation),
			CreatedAt: comment.CreatedAt,
		},
	}
	if err := res.Update(id); err != nil {
		c.ErrorL("更新友链失败", logMap, err.Error())
		c.FailWithErrCode(errcode.FriendLinkUpdateFailed, nil)
		return
	}

	c.Success(res)
}

func GetFriends(c *core.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}

	data, count, err := new(entity.FriendLink).FindWithPage(page, 15)
	if err != nil {
		c.ErrorL("获取友链列表失败", page, err.Error())
		c.FailWithErrCode(errcode.FriendLinkGetFailed, nil)
		return
	}

	c.Success(map[string]interface{}{
		"list":  data,
		"count": count,
	})
}

func GetFriend(c *core.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := new(entity.FriendLink).GetById(id)
	if err != nil {
		c.ErrorL("获取友链失败", id, err.Error())
		c.FailWithErrCode(errcode.FriendLinkGetFailed, nil)
		return
	}

	c.Success(res)
}
