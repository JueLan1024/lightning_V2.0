package routes

import (
	"web_app/controller"
	"web_app/logger"
	"web_app/middlewares"
	"web_app/settings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(mode string, cfg *settings.RatelimitConfig) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin 设置成发布模式
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v2 := r.Group("/api/v2")
	{
		// 注册功能
		v2.POST("/signup", controller.SignUpHandler)
		// 登录功能
		v2.POST("/login", controller.LoginHandler)
		// 查看社区列表功能
		v2.GET("/community", controller.GetCommunityListHandler)
		// 查看社区详细信息功能
		v2.GET("/community/:id", controller.GetCommunityDetailHandler)
		// 查看帖子功能
		v2.GET("/post/:id", controller.GetPostDetailHandler)
		// 查看帖子列表功能
		v2.GET("/posts", controller.GetPostListHandler)
		// 使用jwt认证中间件
		v2.Use(middlewares.JWTMiddleware(), middlewares.RateLimitMiddleware(cfg.FillInterval, cfg.Cap))
		// 创建帖子功能
		v2.POST("/post", controller.CreatePostHandler)
		// 投票功能
		v2.POST("/vote", controller.VoteForPostHandler)

	}
	// 管理员接口
	admin := r.Group("/admin")
	{
		// // 将mysql中的社区ids更新到redis中
		// admin.PUT("/set/community/ids", controller.SetCommunityIDsInRedisHandler)
		// 创建新社区
		admin.POST("/add/community", controller.CreateCommunityHandler)
	}
	return r
}
