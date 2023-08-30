package middleware

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"blog-admin-api/service/token"
)

func Auth(c *core.Context) {
	if _, err := token.ParseRequest(c.Context); err != nil {
		c.FailWithErrCode(errcode.AdminInvalidToken, err.Error())
	}
}
