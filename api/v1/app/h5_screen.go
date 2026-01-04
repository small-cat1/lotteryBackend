package app

import (
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5ScreenApi struct{}

// GetScreenCheckInList 签到墙数据
// @Tags H5-大屏
// @Summary 签到墙数据
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.CheckInRecordResp}
// @Router /h5/screen/checkIn/{activityId} [get]
func (a *H5ScreenApi) GetScreenCheckInList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)
	if limitReq.Limit <= 0 {
		limitReq.Limit = 100
	}

	result, err := h5CheckInService.GetRecentCheckIns(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetScreenDanmakuList 弹幕墙数据
// @Tags H5-大屏
// @Summary 弹幕墙数据
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.DanmakuResp}
// @Router /h5/screen/danmaku/{activityId} [get]
func (a *H5ScreenApi) GetScreenDanmakuList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)
	if limitReq.Limit <= 0 {
		limitReq.Limit = 50
	}

	result, err := h5DanmakuService.GetRecentDanmaku(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetScreenShakeRanking 摇一摇排行榜
// @Tags H5-大屏
// @Summary 摇一摇排行榜
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.ShakeRankingResp}
// @Router /h5/screen/ranking/{roundId} [get]
func (a *H5ScreenApi) GetScreenShakeRanking(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, err := strconv.ParseUint(roundIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)
	if limitReq.Limit <= 0 {
		limitReq.Limit = 20
	}

	result, err := h5ShakeService.GetShakeRanking(uint(roundId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetScreenWinnerList 中奖名单
// @Tags H5-大屏
// @Summary 中奖名单
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.WinningResp}
// @Router /h5/screen/winners/{activityId} [get]
func (a *H5ScreenApi) GetScreenWinnerList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)
	if limitReq.Limit <= 0 {
		limitReq.Limit = 20
	}

	result, err := h5PrizeService.GetRecentWinnings(uint(activityId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
