package middleware

import (
	"blog-admin-api/core"
	"time"
)

func TimeLocation(c *core.Context) {
	c.TimeLocation = time.FixedZone("CST", 8*3600)
}
