package app

import (
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5PrizeApi struct{}

var h5PrizeService = service.ServiceGroupApp.AppServiceGroup.H5PrizeService

// GetPrizeList 获取奖品列表
// @Tags H5-奖品
// @Summary 获取奖品列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=[]appResp.H5PrizeResp}
// @Router /h5/prize/list/{activityId} [get]
func (a *H5PrizeApi) GetPrizeList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5PrizeService.GetPrizeList(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetPrizeDetail 获取奖品详情
// @Tags H5-奖品
// @Summary 获取奖品详情
// @Security ApiKeyAuth
// @Produce application/json
// @Param prizeId path int true "奖品ID"
// @Success 200 {object} response.Response{data=appResp.H5PrizeResp}
// @Router /h5/prize/{prizeId} [get]
func (a *H5PrizeApi) GetPrizeDetail(c *gin.Context) {
	prizeIdStr := c.Param("prizeId")
	prizeId, err := strconv.ParseUint(prizeIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5PrizeService.GetPrizeDetail(uint(prizeId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetMyWinnings 获取我的中奖记录
// @Tags H5-奖品
// @Summary 获取我的中奖记录
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=[]appResp.WinningResp}
// @Router /h5/prize/my/{activityId} [get]
func (a *H5PrizeApi) GetMyWinnings(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5PrizeService.GetMyWinnings(userId, uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetWinningDetail 获取中奖详情
// @Tags H5-奖品
// @Summary 获取中奖详情
// @Security ApiKeyAuth
// @Produce application/json
// @Param winnerId path int true "中奖记录ID"
// @Success 200 {object} response.Response{data=appResp.WinningResp}
// @Router /h5/prize/winning/{winnerId} [get]
func (a *H5PrizeApi) GetWinningDetail(c *gin.Context) {
	winnerIdStr := c.Param("winnerId")
	winnerId, err := strconv.ParseUint(winnerIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5PrizeService.GetWinningDetail(uint(winnerId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRecentWinnings 获取最新中奖记录
// @Tags H5-奖品
// @Summary 获取最新中奖记录
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.WinningResp}
// @Router /h5/prize/recent/{activityId} [get]
func (a *H5PrizeApi) GetRecentWinnings(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)

	result, err := h5PrizeService.GetRecentWinnings(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetReceiveQrCode 获取领奖二维码
// @Tags H5-奖品
// @Summary 获取领奖二维码
// @Security ApiKeyAuth
// @Produce application/json
// @Param winnerId path int true "中奖记录ID"
// @Success 200 {object} response.Response{data=string}
// @Router /h5/prize/qrcode/{winnerId} [get]
func (a *H5PrizeApi) GetReceiveQrCode(c *gin.Context) {
	winnerIdStr := c.Param("winnerId")
	winnerId, err := strconv.ParseUint(winnerIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5PrizeService.GetReceiveQrCode(userId, uint(winnerId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"qrcode": result}, c)
}
