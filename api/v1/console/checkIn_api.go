package console

import (
	"github.com/gin-gonic/gin"
	"lotteryBackend/model/common/response"
	consoleReq "lotteryBackend/model/console/request"
	"strconv"
)

type CheckInApi struct {
}

// ==================== 签到管理 ====================

// GetCheckInStats 获取签到统计
// @Tags Console-签到
// @Summary 获取签到统计（包含状态）
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Success 200 {object} response.Response
// @Router /console/checkin/stats [get]
func (a *CheckInApi) GetCheckInStats(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := checkInService.GetCheckInStats(uint(activityId), 0)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// OpenCheckIn 开启签到
// @Tags Console-签到
// @Summary 开启签到
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.OpenCheckInReq true "开启签到"
// @Success 200 {object} response.Response
// @Router /console/checkin/open [post]
func (a *CheckInApi) OpenCheckIn(c *gin.Context) {
	var req consoleReq.OpenCheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := checkInService.OpenCheckIn(req.ActivityId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("签到已开启", c)
}

// CloseCheckIn 关闭签到
// @Tags Console-签到
// @Summary 关闭签到
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.CloseCheckInReq true "关闭签到"
// @Success 200 {object} response.Response
// @Router /console/checkin/close [post]
func (a *CheckInApi) CloseCheckIn(c *gin.Context) {
	var req consoleReq.CloseCheckInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := checkInService.CloseCheckIn(req.ActivityId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("签到已关闭", c)
}
