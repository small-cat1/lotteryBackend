package app

import (
	"github.com/gin-gonic/gin"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"
)

type ConsoleApi struct{}

var consoleService = service.ServiceGroupApp.AppServiceGroup.ConsoleService

// ==================== 大屏数据接口 ====================

// GetCheckInStats 获取签到统计
// @Tags 控制台
// @Summary 获取签到统计
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response
// @Router /h5/screen/checkIn/stats/{activityId} [get]
func (a *ConsoleApi) GetCheckInStats(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := consoleService.GetCheckInStats(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetRoundList 获取场次列表
// @Tags 控制台
// @Summary 获取场次列表
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response
// @Router /h5/screen/rounds/{activityId} [get]
func (a *ConsoleApi) GetRoundList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := consoleService.GetRoundList(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetDanmakuList 获取弹幕列表（含用户统计）
// @Tags 控制台
// @Summary 获取弹幕列表
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response
// @Router /h5/screen/danmaku/{activityId} [get]
func (a *ConsoleApi) GetDanmakuList(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	result, err := consoleService.GetDanmakuList(uint(activityId), limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetRanking 获取排行榜
// @Tags 控制台
// @Summary 获取排行榜
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response
// @Router /h5/screen/ranking/{roundId} [get]
func (a *ConsoleApi) GetRanking(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, _ := strconv.ParseUint(roundIdStr, 10, 64)
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	result, err := consoleService.GetRanking(uint(roundId), limit)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetRandomDrawConfig 获取随机抽奖配置
// @Tags 控制台
// @Summary 获取随机抽奖配置
// @Produce application/json
// @Success 200 {object} response.Response
// @Router /h5/screen/config/randomDraw [get]
func (a *ConsoleApi) GetRandomDrawConfig(c *gin.Context) {
	result := consoleService.GetRandomDrawConfig()
	response.OkWithData(result, c)
}

// ==================== 游戏控制接口 ====================

// StartGame 开始游戏（验证密码）
// @Tags 控制台
// @Summary 开始游戏
// @Accept application/json
// @Produce application/json
// @Param data body request.StartGameReq true "开始游戏"
// @Success 200 {object} response.Response
// @Router /console/game/start [post]
func (a *ConsoleApi) StartGame(c *gin.Context) {
	var req request.StartGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := consoleService.StartGame(req.RoundId, req.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("游戏已开始", c)
}

// StopGame 立即停止游戏
// @Tags 控制台
// @Summary 立即停止游戏
// @Accept application/json
// @Produce application/json
// @Param data body request.StopGameReq true "停止游戏"
// @Success 200 {object} response.Response
// @Router /console/game/stop [post]
func (a *ConsoleApi) StopGame(c *gin.Context) {
	var req request.StopGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := consoleService.StopGame(req.RoundId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// SettleGame 结算游戏
// @Tags 控制台
// @Summary 结算游戏
// @Accept application/json
// @Produce application/json
// @Param data body request.SettleGameReq true "结算游戏"
// @Success 200 {object} response.Response
// @Router /console/game/settle [post]
func (a *ConsoleApi) SettleGame(c *gin.Context) {
	var req request.SettleGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := consoleService.SettleGame(req.RoundId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// ==================== 随机抽奖接口 ====================

// RandomDraw 随机抽奖（从弹幕用户中抽取）
// @Tags 控制台
// @Summary 随机抽奖
// @Accept application/json
// @Produce application/json
// @Param data body request.RandomDrawReq true "随机抽奖"
// @Success 200 {object} response.Response
// @Router /console/randomDraw [post]
func (a *ConsoleApi) RandomDraw(c *gin.Context) {
	var req request.RandomDrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := consoleService.RandomDraw(req.ActivityId, req.Count)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetWinners 获取中奖名单
// @Tags 控制台
// @Summary 获取中奖名单
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param winType query int false "中奖类型 1摇一摇 2随机"
// @Success 200 {object} response.Response
// @Router /h5/screen/winners/{activityId} [get]
func (a *ConsoleApi) GetWinners(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	winTypeStr := c.Query("winType")
	winType, _ := strconv.Atoi(winTypeStr)

	result, err := consoleService.GetWinners(uint(activityId), winType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}
