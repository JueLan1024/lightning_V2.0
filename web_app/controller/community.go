package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetCommunityListHandler 获得社区列表功能
// @Summary 获取社区列表
// @Description 获取所有社区的列表
// @Tags 社区相关接口
// @Produce json
// @Success 200 {object} _ResponseCommunityList "成功返回社区列表"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/community [get]
func GetCommunityListHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// 业务处理
	data, err := logic.GetCommunityList(ctx)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// CreateCommunityHandler 创建新社区功能
// @Summary 创建社区
// @Description 创建一个新的社区
// @Tags 社区相关接口
// @Accept json
// @Produce json
// @Param community body models.ParamCommunity true "社区参数"
// @Success 200 {object} _Response "成功创建社区"
// @Failure 400 {object} _Response "参数错误"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /admin/add/community [post]
func CreateCommunityHandler(c *gin.Context) {
	ctx := c.Request.Context()
	// 参数获取和参数检验
	p := new(models.ParamCommunity)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Create community with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务处理
	if err := logic.CreateCommunity(ctx, p); err != nil {
		zap.L().Error("failed create community", zap.Error(err))
		if err == logic.ErrorCommunityExist {
			ResponseError(c, CodeCommunityExists)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)

}

// GetCommunityDetailHandler 获得社区详细信息功能
// @Summary 获取社区详情
// @Description 根据社区ID获取社区的详细信息
// @Tags 社区相关接口
// @Produce json
// @Param id path string true "社区ID"
// @Success 200 {object} _ResponseCommunityDetail "成功返回社区详情"
// @Failure 400 {object} _Response "参数错误"
// @Failure 404 {object} _Response "社区不存在"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/community/{id} [get]
func GetCommunityDetailHandler(c *gin.Context) {
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
	data, err := logic.GetCommunityDetail(ctx, id)
	if err == logic.ErrorCommunityNotExist {
		ResponseError(c, CodeCommunityNotExists)
		return
	}
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// // SetCommunityIDsInRedisHandler 将mysql中的社区IDs更新到redis中
// func SetCommunityIDsInRedisHandler(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	err := logic.SetCommunityIDs(ctx)
// 	if err != nil {
// 		zap.L().Error("failed to Set community ids in redis", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuccess(c, nil)
// }
