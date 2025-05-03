package controller

import (
	"errors"
	"web_app/logic"
	"web_app/models"
	"web_app/settings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param user body models.ParamSignUp true "注册参数"
// @Success 200 {object} _Response "注册成功"
// @Failure 400 {object} _Response "参数错误"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/signup [post]
func SignUpHandler(c *gin.Context) {
	// 获取上下文
	ctx := c.Request.Context()

	// 1.获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator类型的
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return

	}
	// 2.业务处理
	if err := logic.SignUp(ctx, p); err != nil {
		if errors.Is(err, logic.ErrorUserExist) {
			ResponseError(c, CodeUsernameExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求的函数
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户相关接口
// @Accept json
// @Produce json
// @Param user body models.ParamLogin true "登录参数"
// @Success 200 {object} _Response "登录成功，返回 accessToken"
// @Failure 400 {object} _Response "参数错误"
// @Failure 401 {object} _Response "用户名或密码错误"
// @Failure 500 {object} _Response "服务器繁忙"
// @Router /api/v2/login [post]
func LoginHandler(c *gin.Context) {
	ctx := c.Request.Context() //获取上下文

	// 1.参数获取和参数校验
	p := new(models.ParamLogin)
	// 校验格式是否正确
	if err := c.ShouldBindJSON(&p); err != nil {
		// 记录错误日志
		zap.L().Error("LoginHandler invalid param", zap.Error(err))
		// 判断错误类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2.业务处理,在logic层校验用户名是否存在，密码是否正确
	accessToken, refreshToken, err := logic.Login(ctx, p)
	if err != nil {
		// 用户名不存在或密码错误
		if errors.Is(err, logic.ErrorUserNotExist) || errors.Is(err, logic.ErrorInvalidPassword) {
			ResponseError(c, CodeUsernameOrPasswordWrong)
			return
		}
		// 其他错误
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	// 设置 HttpOnly Cookie
	secure := settings.Conf.Mode == "release" // 判断是不是发布者模式
	c.SetCookie(
		logic.RefreshCookieName,
		refreshToken,
		int(settings.Conf.RefreshTokenDuration.Seconds()),
		logic.RefreshCookiePath,
		logic.RefreshCookieDomain,
		secure,
		logic.RefreshCookieHttpOnly,
	)
	// 返回 accessToken
	ResponseSuccess(c, accessToken)
}
