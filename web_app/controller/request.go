package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"

var ErrorNeedLogin = errors.New("未登录")

// GetCurrentUserID 获得当前用户ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uidStr, exists := c.Get(CtxUserIDKey)
	if !exists {
		err = ErrorNeedLogin
		return
	}
	userID, ok := uidStr.(int64)
	if !ok {
		err = ErrorNeedLogin
		return
	}
	return
}
