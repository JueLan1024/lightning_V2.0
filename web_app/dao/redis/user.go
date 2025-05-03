package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// CreateRereshToken 将新创建的RereshToken存入redis
func CreateRereshToken(ctx context.Context, userID int64, token string, time_duration time.Duration) (err error) {
	return rdb.Set(ctx, GetKeyUserRefreshToken(userID), token, time_duration).Err()
}

// GetRefreshToken 通过用户id获取refreshToken
func GetRefreshToken(ctx context.Context, userID int64) (token string, err error) {
	token, err = rdb.Get(ctx, GetKeyUserRefreshToken(userID)).Result()
	if err == redis.Nil {
		return "", ErrorRefreshTokenNotFound
	}
	return token, err
}
