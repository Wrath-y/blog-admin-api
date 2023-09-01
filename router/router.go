package router

import (
	"blog-admin-api/core"
	"blog-admin-api/errcode"
	"blog-admin-api/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery)
	r.Use(middleware.SetV())
	r.Use(core.Handle(middleware.CORS))

	r.NoRoute(NoRoute)

	r.Any("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	g := r.Group("/", core.Handle(middleware.TimeLocation))
	loadAdmin(g)

	return r
}

func NoRoute(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, errcode.LibNoRoute)
}
