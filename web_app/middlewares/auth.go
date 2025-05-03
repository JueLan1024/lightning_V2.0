package middlewares

import (
	"context"
	"errors"
	"strings"
	"web_app/controller"
	"web_app/logic"
	"web_app/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTMiddleware Token认证中间件
func JWTMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context() //获取上下文

		AuthorHead := c.Request.Header.Get("Authorization")
		var mc *jwt.MyClaims
		// 如果 Authorization 为空，检查 refreshToken
		if AuthorHead == "" {
			if handleRefreshToken(ctx, c) { // 如果 refreshToken 验证成功
				c.Next() //执行后续handler
				return   //结束Token认证
			}
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// Authorization不为空,解析accessToken
		parts := strings.SplitN(AuthorHead, " ", 2) //将Authorization按空格分割
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// 解析accessToken
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			// 如果accessToken不合规,检查 refreshToken
			if handleRefreshToken(ctx, c) { // 如果 refreshToken 验证成功
				c.Next() //执行后续handler
				return   //结束Token认证
			}
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将用户id存入上下文
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}

func handleRefreshToken(ctx context.Context, c *gin.Context) bool {
	// 从当前客户端Cookie获取refreshToken
	refreshToken, err := c.Cookie(logic.RefreshCookieName)
	if err != nil {
		zap.L().Error("failed to get refresh_token in Cookie", zap.Error(err))
		return false
	}
	// 验证 refreshToken 并生成新的 accessToken
	accessToken, mc, err := logic.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, logic.ErrorRefreshTokenNotExist) || errors.Is(err, logic.ErrorInvalidRefeshToken) {
			zap.L().Warn("RefreshToken expired")
		}
		return false
	}
	// 返回成功响应，和新accessToken
	controller.ResponseSuccess(c, accessToken)
	// 将用户id存入上下文
	c.Set(controller.CtxUserIDKey, mc.UserID)
	return true
}
