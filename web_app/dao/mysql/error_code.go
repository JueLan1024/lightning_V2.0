package mysql

import "errors"

var (
	ErrorUsernameExist     = errors.New("用户名已存在")
	ErrorUsernameNotFound  = errors.New("用户名不存在")
	ErrorCommunityIDExist  = errors.New("社区已存在")
	ErrorCommunityNotExist = errors.New("社区不存在")
	ErrorPostNotExist      = errors.New("帖子不存在")
	ErrorUserNotFound      = errors.New("用户不存在")
)
