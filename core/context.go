package core

import (
	"blog-admin-api/errcode"
	"blog-admin-api/pkg/def"
	"blog-admin-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Context struct {
	Env string
	*gin.Context
	*logging.Logger
	TimeLocation *time.Location
}

func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func (c *Context) Fail(code int, msg string, detail, data interface{}) {
	if c.Env == def.EnvProduction {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":   code,
			"msg":    msg,
			"detail": detail,
			"data":   data,
		})
	}
}

func (c *Context) FailWithErrCode(err *errcode.ErrCode, data interface{}) {
	c.Fail(err.Code, err.Msg, err.Detail, data)
}
