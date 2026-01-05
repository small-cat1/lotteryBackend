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
// GetUserInfo 获取用户信息
func (a *H5UserApi) GetUserInfo(c *gin.Context) {
	var req request.GetUserInfoReq
	_ = c.ShouldBindQuery(&req)

	userId := c.GetUint("h5UserId")

	result, err := h5UserService.GetUserInfo(userId, req.ActivityId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// CheckIn 用户签到
func (a *H5UserApi) CheckIn(c *gin.Context) {
	var req request.CheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userId := c.GetUint("h5UserId")
	ip := c.ClientIP()

	result, err := h5UserService.CheckIn(userId, req, ip)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetAuditStatus 获取审核状态
func (a *H5UserApi) GetAuditStatus(c *gin.Context) {
	userId := c.GetUint("h5UserId")
	activityId := c.GetUint("activityId") // 从query获取

	result, err := h5UserService.GetAuditStatus(userId, activityId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
