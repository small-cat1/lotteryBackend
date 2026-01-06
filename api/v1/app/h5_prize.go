package app

import (
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5PrizeApi struct{}

var h5PrizeService = service.ServiceGroupApp.AppServiceGroup.H5PrizeService

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
