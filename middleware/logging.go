package middleware

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"blog-admin-api/pkg/def"
	"blog-admin-api/pkg/logging"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

const bodyLimitKB = 5000

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logging(c *core.Context) {
	start := time.Now()

	raw, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	w := &BodyLogWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
	c.Writer = w

	logger := logging.New()
	logger.SetRequestID(c.GetString(def.XRequestID))
	logger.Setv1(c.GetString("v1"))
	c.Logger = logger

	rawKB := len(raw) / 1024 // => to KB
	if rawKB > bodyLimitKB {
		c.Logger.Info("接口请求与响应", string(raw[:1024]), nil)
		c.FailWithErrCode(errcode.AdminBodyTooLarge.WithDetail(fmt.Sprintf("消息限制%dKB, 本消息%dKB", bodyLimitKB, rawKB)), nil)
		return
	}

	c.Next()

	logger.Setv2(c.GetString("v2"))
	logger.Setv3(c.GetString("v3"))

	request := map[string]any{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"query":  logging.SpreadMaps(c.Request.URL.Query()),
		"header": logging.SpreadMaps(c.Request.Header),
		"body":   string(raw),
	}
	c.Logger.Info("接口请求与响应", request, w.body.Bytes(), logging.AttrOption{StartTime: &start})
}
