package logic

import "errors"

var (
	ErrorUserExist            = errors.New("用户名已存在")
	ErrorUserNotExist         = errors.New("用户不存在")
	ErrorInvalidPassword      = errors.New("密码错误")
	ErrorRefreshTokenNotExist = errors.New("RefreshToken不存在")
	ErrorWorngTokenType       = errors.New("错误的token类型")
	ErrorInvalidRefeshToken   = errors.New("invalid refreshToken")
	ErrorCommunityExist       = errors.New("社区已存在")
	ErrorCommunityNotExist    = errors.New("社区不存在")
	ErrorPostNotExist         = errors.New("帖子不存在")
	ErrorInvalidPageToken     = errors.New("invalid pageToken")
	ErrorVoteRepeated         = errors.New("重复投票")
)
