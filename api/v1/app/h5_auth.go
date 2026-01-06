package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type H5AuthApi struct{}

var h5AuthService = service.ServiceGroupApp.AppServiceGroup.H5AuthService

// WechatLogin 微信授权登录
// @Tags H5-授权
// @Summary 微信授权登录
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appResp.WechatLoginResp}
// @Router /h5/auth/wechat [post]
func (a *H5AuthApi) WechatLogin(c *gin.Context) {
	var req request.WechatLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	result, err := h5AuthService.WechatLogin(req.Code, req.ActivityId)
	if err != nil {
		global.GVA_LOG.Error("微信登录失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetWxJsConfig 获取微信JS-SDK配置
// @Tags H5-授权
// @Summary 获取微信JS-SDK配置
// @accept application/json
// @Produce application/json
// @Param url query string true "当前页面URL"
// @Success 200 {object} response.Response{data=appResp.WxJsConfigResp}
// @Router /h5/auth/wxConfig [get]
func (a *H5AuthApi) GetWxJsConfig(c *gin.Context) {
	var req request.WxJsConfigReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5AuthService.GetWxJsConfig(req.Url)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetWechatConfig 获取微信配置
// @Tags H5-授权
// @Summary 获取微信配置(AppID)
// @Produce application/json
// @Success 200 {object} response.Response{data=appResp.WechatConfigResp}
// @Router /h5/config/wechat [get]
func (a *H5AuthApi) GetWechatConfig(c *gin.Context) {
	result, err := h5AuthService.GetWechatConfig()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// RefreshToken 刷新Token
// @Tags H5-授权
// @Summary 刷新Token
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=string}
// @Router /h5/auth/refresh [post]
func (a *H5AuthApi) RefreshToken(c *gin.Context) {
	userId := c.GetUint("h5UserId")

	token, err := h5AuthService.RefreshToken(userId)
	if err != nil {
		response.FailWithMessage("刷新失败", c)
		return
	}

	response.OkWithData(gin.H{"token": token}, c)
}
