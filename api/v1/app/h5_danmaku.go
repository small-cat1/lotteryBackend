package app

import (
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5DanmakuApi struct{}

var h5DanmakuService = service.ServiceGroupApp.AppServiceGroup.H5DanmakuService

// SendDanmaku 发送弹幕
// @Tags H5-弹幕
// @Summary 发送弹幕
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SendDanmakuReq true "弹幕信息"
// @Success 200 {object} response.Response{data=appResp.DanmakuResp}
// @Router /h5/danmaku [post]
func (a *H5DanmakuApi) SendDanmaku(c *gin.Context) {
	var req request.SendDanmakuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5DanmakuService.SendDanmaku(userId, req.ActivityId, req.Content, req.Color)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetDanmakuList 获取弹幕列表
// @Tags H5-弹幕
// @Summary 获取弹幕列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=appResp.H5PageResult}
// @Router /h5/danmaku/list/{activityId} [get]
func (a *H5DanmakuApi) GetDanmakuList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var pageReq request.H5PageReq
	c.ShouldBindQuery(&pageReq)

	result, err := h5DanmakuService.GetDanmakuList(uint(activityId), pageReq.Page, pageReq.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRecentDanmaku 获取最新弹幕
// @Tags H5-弹幕
// @Summary 获取最新弹幕
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.DanmakuResp}
// @Router /h5/danmaku/recent/{activityId} [get]
func (a *H5DanmakuApi) GetRecentDanmaku(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)

	result, err := h5DanmakuService.GetRecentDanmaku(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetTopDanmaku 获取置顶弹幕
// @Tags H5-弹幕
// @Summary 获取置顶弹幕
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=appResp.DanmakuResp}
// @Router /h5/danmaku/top/{activityId} [get]
func (a *H5DanmakuApi) GetTopDanmaku(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5DanmakuService.GetTopDanmaku(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetMyDanmaku 获取我的弹幕
// @Tags H5-弹幕
// @Summary 获取我的弹幕
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=[]appResp.DanmakuResp}
// @Router /h5/danmaku/my/{activityId} [get]
func (a *H5DanmakuApi) GetMyDanmaku(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5DanmakuService.GetMyDanmaku(userId, uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
