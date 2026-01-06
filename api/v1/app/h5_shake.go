package app

import (
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5ShakeApi struct{}

var h5ShakeService = service.ServiceGroupApp.AppServiceGroup.H5ShakeService

// GetCurrentRound 获取当前场次
// @Tags H5-摇一摇
// @Summary 获取当前场次
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=appResp.ShakeRoundResp}
// @Router /h5/shake/round/current/{activityId} [get]
func (a *H5ShakeApi) GetCurrentRound(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5ShakeService.GetCurrentRound(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRoundList 获取场次列表
// @Tags H5-摇一摇
// @Summary 获取场次列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response{data=[]appResp.ShakeRoundResp}
// @Router /h5/shake/round/list/{activityId} [get]
func (a *H5ShakeApi) GetRoundList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, err := strconv.ParseUint(activityIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5ShakeService.GetRoundList(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRoundDetail 获取场次详情
// @Tags H5-摇一摇
// @Summary 获取场次详情
// @Security ApiKeyAuth
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Success 200 {object} response.Response{data=appResp.ShakeRoundResp}
// @Router /h5/shake/round/{roundId} [get]
func (a *H5ShakeApi) GetRoundDetail(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, err := strconv.ParseUint(roundIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5ShakeService.GetRoundDetail(uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetShakeRanking 获取实时排名
// @Tags H5-摇一摇
// @Summary 获取实时排名
// @Security ApiKeyAuth
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]appResp.ShakeRankingResp}
// @Router /h5/shake/ranking/{roundId} [get]
func (a *H5ShakeApi) GetShakeRanking(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, err := strconv.ParseUint(roundIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	var limitReq request.H5LimitReq
	c.ShouldBindQuery(&limitReq)

	result, err := h5ShakeService.GetShakeRanking(uint(roundId), limitReq.Limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetMyScore 获取我的成绩
// @Tags H5-摇一摇
// @Summary 获取我的成绩
// @Security ApiKeyAuth
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Success 200 {object} response.Response{data=appResp.MyScoreResp}
// @Router /h5/shake/my/{roundId} [get]
func (a *H5ShakeApi) GetMyScore(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, err := strconv.ParseUint(roundIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5ShakeService.GetMyScore(userId, uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetRoundResult 获取最终结果
// @Tags H5-摇一摇
// @Summary 获取最终结果
// @Security ApiKeyAuth
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Success 200 {object} response.Response{data=appResp.RoundResultResp}
// @Router /h5/shake/result/{roundId} [get]
func (a *H5ShakeApi) GetRoundResult(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, err := strconv.ParseUint(roundIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5ShakeService.GetRoundResult(userId, uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
