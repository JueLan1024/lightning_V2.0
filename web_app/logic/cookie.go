package logic

import (
	"context"
	"errors"
	"web_app/dao/redis"
	"web_app/pkg/jwt"

	"go.uber.org/zap"
)

const (
	RefreshCookieName     = "refresh_token"
	RefreshCookiePath     = "/"
	RefreshCookieDomain   = "127.0.0.1"
	RefreshCookieHttpOnly = true
)

// ValidateRefreshToken 认证RefreshToken
func ValidateRefreshToken(ctx context.Context, refreshToken string) (accessToken string, mc *jwt.MyClaims, err error) {
	// 解析当前refreshToken
	mc, err = jwt.ParseToken(refreshToken)
	if err != nil {
		return "", nil, err
	}
	// 根据mc中的userID,从redis中查询refreshToken
	tokenInRedis, err := redis.GetRefreshToken(ctx, mc.UserID)
	// 在redis中没有找到
	if err != nil {
		if errors.Is(err, redis.ErrorRefreshTokenNotFound) {
			return "", nil, ErrorRefreshTokenNotExist
		}
		zap.L().Error("failed to get refreshToken in redis", zap.Int64("userID", mc.UserID), zap.Error(err))
		return "", nil, err
	}
	// 比较两个token是否一致
	if tokenInRedis != refreshToken {
		return "", nil, ErrorInvalidRefeshToken
	}
	// 一致就生成accessToken
	accessToken, err = genToken(mc.UserID, mc.Username, AccessTokenType)
	if err != nil {
		zap.L().Error("failed to generate accessToken", zap.Int64("userID", mc.UserID), zap.Error(err))
		return "", nil, err
	}
	return accessToken, mc, nil

}
