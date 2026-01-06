package app

import (
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
