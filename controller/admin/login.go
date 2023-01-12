package admin

import (
	"blog-admin-api/core"
	"blog-admin-api/entity"
	"blog-admin-api/errcode"
	"blog-admin-api/server/auth"
	"blog-admin-api/server/token"
)

type UserRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *core.Context) {
	var (
		r      UserRequest
		logMap = make(map[string]interface{})
	)
	if err := c.ShouldBindJSON(&r); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidParam, err.Error())
		return
	}
	logMap["req"] = r

	// Get the administrator information by the login username.
	user, err := entity.GetUserByName(r.Account)
	if err != nil {
		c.ErrorL("获取用户失败", r, err.Error())
		c.FailWithErrCode(errcode.AdminUserGetFailed, nil)
		return
	}
	logMap["user"] = user

	// Compare the login password with the administrator password.
	if err := auth.Compare(user.Password, r.Password); err != nil {
		c.ErrorL("密码错误", logMap, err.Error())
		c.FailWithErrCode(errcode.AdminUserPasswdErr, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c.Context, token.Context{ID: user.Id, Account: user.Account}, "")
	if err != nil {
		c.ErrorL("token错误", logMap, err.Error())
		c.FailWithErrCode(errcode.AdminInvalidToken, nil)
		return
	}

	c.Success(entity.Token{Token: t})
}
