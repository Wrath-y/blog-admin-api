package middleware

import (
	"blog-admin-api/pkg/def"
	"blog-admin-api/pkg/util"
	"github.com/gin-gonic/gin"
)

func SetV() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID := c.GetHeader(def.XRequestID)
		if xRequestID == "" {
			xRequestID = util.UUID()
		}
		c.Set(def.XRequestID, xRequestID)
		c.Set("v1", c.Request.URL.Path)
		c.Next()
	}
}
