package app

import (
	"github.com/gin-gonic/gin"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type H5UserApi struct{}

var h5UserService = service.ServiceGroupApp.AppServiceGroup.H5UserService

// GetUserInfo 获取用户信息
// @Tags H5-用户
// @Summary 获取用户信息
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=appResp.H5UserResp}
// @Router /h5/user/info [get]
func (a *H5UserApi) GetUserInfo(c *gin.Context) {
	userId := c.GetUint("h5UserId")

	result, err := h5UserService.GetUserInfo(userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// UserRegister 用户报名
// @Tags H5-用户
// @Summary 用户报名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UserRegisterReq true "报名信息"
// @Success 200 {object} response.Response{data=appResp.H5UserResp}
// @Router /h5/user/register [post]
func (a *H5UserApi) UserRegister(c *gin.Context) {
	var req request.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5UserService.UserRegister(userId, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// UpdateUserInfo 更新用户信息
// @Tags H5-用户
// @Summary 更新用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateUserInfoReq true "用户信息"
// @Success 200 {object} response.Response{data=appResp.H5UserResp}
// @Router /h5/user/info [put]
func (a *H5UserApi) UpdateUserInfo(c *gin.Context) {
	var req request.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5UserService.UpdateUserInfo(userId, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetAuditStatus 获取审核状态
// @Tags H5-用户
// @Summary 获取审核状态
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=appResp.AuditStatusResp}
// @Router /h5/user/audit [get]
func (a *H5UserApi) GetAuditStatus(c *gin.Context) {
	userId := c.GetUint("h5UserId")

	result, err := h5UserService.GetAuditStatus(userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
