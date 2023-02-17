package middleware

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"blog-admin-api/service/token"
)

func Auth(c *core.Context) {
	if _, err := token.ParseRequest(c.Context); err != nil {
		c.Fatal("解析token失败", c.Context, err.Error())
		c.FailWithErrCode(errcode.AdminInvalidToken, nil)
	}
}
