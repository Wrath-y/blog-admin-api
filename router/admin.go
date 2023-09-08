package router

import (
	"blog-admin-api/controller/admin"
	"blog-admin-api/core"
	"blog-admin-api/middleware"
	"github.com/gin-gonic/gin"
)

func loadAdmin(r *gin.RouterGroup) {
	noAuthApi := r.Group("/", core.Handle(middleware.Logging))
	{
		noAuthApi.POST("/login", core.Handle(admin.Login))
	}

	authApi := r.Group("/", core.Handle(middleware.Auth), core.Handle(middleware.Logging))
	{
		articles := authApi.Group("/articles")
		{
			articles.POST("", core.Handle(admin.AddArticle))
			articles.DELETE("/:id", core.Handle(admin.DelArticle))
			articles.PUT("/:id", core.Handle(admin.UpdateArticle))
			articles.GET("", core.Handle(admin.GetArticles))
			articles.GET("/:id", core.Handle(admin.GetArticle))
		}
		seo := authApi.Group("/article_seo")
		{
			seo.POST("/", core.Handle(admin.SetArticleSeo))
			seo.GET("/:article_id", core.Handle(admin.GetArticleSeo))
		}
		uploads := authApi.Group("/uploads")
		{
			uploads.GET("", core.Handle(admin.GetUpload))
		}
		pixivs := authApi.Group("/pixivs")
		{
			pixivs.GET("", core.Handle(admin.GetPixivs))
			pixivs.GET("count", core.Handle(admin.GetPixivCount))
			pixivs.POST("", core.Handle(admin.AddPixiv))
			pixivs.DELETE("/:id", core.Handle(admin.DelPixiv))
		}
		comments := authApi.Group("/comments")
		{
			comments.GET("", core.Handle(admin.GetComments))
			comments.DELETE("/:id", core.Handle(admin.DelComment))
		}
		friend := authApi.Group("/friend_links")
		{
			friend.POST("", core.Handle(admin.AddFriend))
			friend.DELETE("/:id", core.Handle(admin.DelFriend))
			friend.PUT("/:id", core.Handle(admin.UpdateFriend))
			friend.GET("", core.Handle(admin.GetFriends))
			friend.GET("/:id", core.Handle(admin.GetFriend))
		}
	}
}
