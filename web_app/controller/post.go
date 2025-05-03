package controller

import (
	"errors"
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子功能
// @Summary 创建帖子
// @Description 创建一个新的帖子
// @Tags 帖子相关接口
// @Accept json
// @Produce json
// @Param Authorization	header string false "Bearer 用户令牌"
// @Param post body models.ParamPost true "帖子参数"
// @Security ApiKeyAuth
// @Success 200 {object} _Response "成功创建帖子"
// @Failure 400 {object} _Response "参数错误"
// @Failure 401 {object} _Response "用户未登录"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/post [post]
func CreatePostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	// 参数获取和参数检验
	p := new(models.ParamPost)
	authorID, err := GetCurrentUserID(c) // 获得当前用户ID
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid params",
			zap.Int64("authorID", authorID),
			zap.Error(err),
		)
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务处理
	if err := logic.CreatePost(ctx, p, authorID); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获得帖子信息功能
// @Summary 获取帖子详情
// @Description 根据帖子ID获取帖子详情
// @Tags 帖子相关接口
// @Produce json
// @Param id path string true "帖子ID"
// @Success 200 {object} _ResponsePostDetail "成功返回帖子详情"
// @Failure 400 {object} _Response "参数错误"
// @Failure 404 {object} _Response "帖子不存在"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	ctx := c.Request.Context()
	// 参数获取和参数检验
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理
	data, err := logic.GetPostDetail(ctx, id)
	if err == logic.ErrorPostNotExist {
		ResponseError(c, CodePostNotExists)
		return
	}
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获得帖子列表功能
// @Summary 获取帖子列表
// @Description 获取按时间或分数排序的帖子列表
// @Tags 帖子相关接口
// @Produce json
// @Param order query string false "排序方式(time或score)"
// @Param token query string false "pageToken"
// @Param community_id query string true "社区ID"
// @Success 200 {object} _ResponsePosts "成功返回pageToken和帖子列表"
// @Failure 400 {object} _Response "参数错误"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/posts [get]
func GetPostListHandler(c *gin.Context) {
	ctx := c.Request.Context()
	// 参数获取和参数检验
	p := new(models.ParamGetPostsInOrder)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("Get post list with invalid params")
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 业务处理
	data, err := logic.GetPostList(ctx, p)
	if err != nil {
		if errors.Is(err, logic.ErrorInvalidPageToken) {
			ResponseError(c, CodeInvalidToken)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
