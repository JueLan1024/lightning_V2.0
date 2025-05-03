package controller

import (
	"errors"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// VoteForPostHandler 给帖子投票功能
// @Summary 给帖子投票
// @Description 用户可以对帖子进行投票（点赞或点踩）
// @Tags 投票相关接口
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param vote body models.ParamVoteForPost true "投票参数"
// @Security	ApiKeyAuth
// @Success 200 {object} _Response "投票成功"
// @Failure 400 {object} _Response "参数错误"
// @Failure 401 {object} _Response "用户未登录"
// @Failure 403 {object} _Response "重复投票"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/vote [post]
func VoteForPostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	//参数获取和参数检验
	p := new(models.ParamVoteForPost)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("VoteForPost with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	userID, err := GetCurrentUserID(c) // 获得当前用户ID
	if err != nil {
		zap.L().Error("GetCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 业务处理
	if err := logic.VoteForPost(ctx, userID, p); err != nil {
		if errors.Is(err, logic.ErrorVoteRepeated) {
			ResponseError(c, CodeVoteRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}
