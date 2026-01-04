package app

import (
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5CheckInApi struct{}

var h5CheckInService = service.ServiceGroupApp.AppServiceGroup.H5CheckInService

// CheckIn 用户签到
// @Tags H5-签到
// @Summary 用户签到
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CheckInReq true "签到信息"
// @Success 200 {object} response.Response{data=appResp.CheckInRecordResp}
// @Router /h5/checkIn [post]
func (a *H5CheckInApi) CheckIn(c *gin.Context) {
	var req request.CheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")
	ip := c.ClientIP()

	result, err := h5CheckInService.CheckIn(userId, req.ActivityId, ip)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetCheckInStatus 获取签到状态
// @Tags H5-签到
// @Summary 获取签到状态
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=appResp.CheckInStatusResp}
// @Router /h5/checkIn/status/{activityId} [get]
func (a *H5CheckInApi) GetCheckInStatus(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5CheckInService.GetCheckInStatus(userId, uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetCheckInList 获取签到列表
// @Tags H5-签到
// @Summary 获取签到列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=appResp.H5PageResult}
// @Router /h5/checkIn/list/{activityId} [get]
func (a *H5CheckInApi) GetCheckInList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var pageReq request.H5PageReq
	c.ShouldBindQuery(&pageReq)

	result, err := h5CheckInService.GetCheckInList(uint(activityId), pageReq.Page, pageReq.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetCheckInStats 获取签到统计
// @Tags H5-签到
// @Summary 获取签到统计
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=appResp.CheckInStatsResp}
// @Router /h5/checkIn/stats/{activityId} [get]
func (a *H5CheckInApi) GetCheckInStats(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5CheckInService.GetCheckInStats(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRecentCheckIns 获取最新签到
// @Tags H5-签到
// @Summary 获取最新签到
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.CheckInRecordResp}
// @Router /h5/checkIn/recent/{activityId} [get]
func (a *H5CheckInApi) GetRecentCheckIns(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)

	result, err := h5CheckInService.GetRecentCheckIns(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
