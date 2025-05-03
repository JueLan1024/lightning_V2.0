package logic

import (
	"context"
	"errors"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
	"web_app/settings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

// encryptPassword 对原密码加密
func encryptPassword(oPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(oPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 注册业务处理
func SignUp(ctx context.Context, p *models.ParamSignUp) (err error) {
	// 判断用户存不存在
	if err := mysql.UsernameExists(ctx, p.Username); err != nil {
		if errors.Is(err, mysql.ErrorUsernameExist) {
			return ErrorUserExist
		}
		return err
	}
	// 加密密码
	hashedPassword, err := encryptPassword(p.Password)
	if err != nil {
		return err
	}
	// 将用户信息存入Mysql
	user := &models.User{
		UserID:   snowflake.GenID(),
		Username: p.Username,
		Password: hashedPassword,
	}

	// 保存进数据库
	return mysql.CreateUser(user)
}

// 登录业务处理
func Login(ctx context.Context, p *models.ParamLogin) (accessToken, refreshToken string, err error) {
	// 创建用户信息结构体
	user := new(models.User)
	// 通过用户名从Mysql中获取用户信息
	if err = mysql.GetUserByUsername(ctx, p.Username, user); err != nil {
		if errors.Is(err, mysql.ErrorUsernameNotFound) {
			return "", "", ErrorUserNotExist
		}
		zap.L().Error("failed to get user by username", zap.String("username", p.Username), zap.Error(err))
		return "", "", err
	}
	// 校验密码是否一致
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return "", "", ErrorInvalidPassword
	}
	// 获取jwtAccessToken和RefreshToken
	accessToken, err = genToken(user.UserID, user.Password, AccessTokenType)
	if err != nil {
		zap.L().Error("failed to generate access token", zap.Error(err))
		return "", "", err
	}
	refreshToken, err = genToken(user.UserID, user.Password, RefreshTokenType)
	if err != nil {
		zap.L().Error("failed to generate refresh token", zap.Error(err))
		return "", "", err
	}
	// 将refreshTokn存入redis
	if err = redis.CreateRereshToken(ctx, user.UserID, refreshToken, settings.Conf.RefreshTokenDuration); err != nil {
		zap.L().Error("failed to store refresh token in redis", zap.Error(err))
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// genToken 根据tokenType生成accessToken或refreshToken
func genToken(userID int64, username string, tokenType string) (token string, err error) {
	if tokenType == AccessTokenType {
		token, err = jwt.GenAccessToken(userID, username, settings.Conf.AccessTokenDuration)
		return
	}
	if tokenType == RefreshTokenType {
		token, err = jwt.GenRefreshToken(userID, username, settings.Conf.RefreshTokenDuration)
		return
	}
	return "", ErrorWorngTokenType
}
