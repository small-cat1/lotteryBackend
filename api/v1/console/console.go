package console

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"lotteryBackend/model/common/response"
	consoleReq "lotteryBackend/model/console/request"
)

type ConsoleApi struct{}

// ==================== 活动配置 ====================

// GetActivityDetail 获取活动详情
// @Tags Console-活动
// @Summary 获取活动详情
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Success 200 {object} response.Response
// @Router /console/activity/{activityId} [get]
func (a *ConsoleApi) GetActivityDetail(c *gin.Context) {
	activityIdStr := c.Param("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	result, err := activityService.GetActivityDetail(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetPrizeList 获取奖品列表
// @Tags Console-活动
// @Summary 获取奖品列表
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Success 200 {object} response.Response
// @Router /console/prizes [get]
func (a *ConsoleApi) GetPrizeList(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := activityService.GetPrizeList(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetRoundList 获取场次列表
// @Tags Console-场次
// @Summary 获取场次列表
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Success 200 {object} response.Response
// @Router /console/rounds [get]
func (a *ConsoleApi) GetRoundList(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := activityService.GetRoundList(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetRoundDetail 获取场次详情
// @Tags Console-场次
// @Summary 获取场次详情
// @Produce application/json
// @Param roundId path int true "场次ID"
// @Success 200 {object} response.Response
// @Router /console/rounds/{roundId} [get]
func (a *ConsoleApi) GetRoundDetail(c *gin.Context) {
	roundIdStr := c.Param("roundId")
	roundId, _ := strconv.ParseUint(roundIdStr, 10, 64)

	result, err := activityService.GetRoundDetail(uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// ==================== 弹幕管理 ====================

// GetDanmakuList 获取弹幕列表
// @Tags Console-弹幕
// @Summary 获取弹幕列表
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Param limit query int false "数量限制"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.Response
// @Router /console/danmaku/list [get]
func (a *ConsoleApi) GetDanmakuList(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	statusStr := c.DefaultQuery("status", "0")
	status, _ := strconv.Atoi(statusStr)

	result, err := danmakuService.GetDanmakuList(uint(activityId), limit, status)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// ==================== 游戏控制 ====================

// StartGame 开始游戏
// @Tags Console-游戏
// @Summary 开始游戏（验证密码）
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.StartGameReq true "开始游戏"
// @Success 200 {object} response.Response
// @Router /console/game/start [post]
func (a *ConsoleApi) StartGame(c *gin.Context) {
	var req consoleReq.StartGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := gameService.StartGame(req.RoundId, req.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("游戏已开始", c)
}

// StopGame 立即停止游戏
// @Tags Console-游戏
// @Summary 立即停止游戏
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.StopGameReq true "停止游戏"
// @Success 200 {object} response.Response
// @Router /console/game/stop [post]
func (a *ConsoleApi) StopGame(c *gin.Context) {
	var req consoleReq.StopGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := gameService.StopGame(req.RoundId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// CancelGame 取消游戏
// @Tags Console-游戏
// @Summary 取消游戏
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.CancelGameReq true "取消游戏"
// @Success 200 {object} response.Response
// @Router /console/game/cancel [post]
func (a *ConsoleApi) CancelGame(c *gin.Context) {
	var req consoleReq.CancelGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := gameService.CancelGame(req.RoundId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("游戏已取消", c)
}

// GetGameStatus 获取游戏状态
// @Tags Console-游戏
// @Summary 获取游戏状态
// @Produce application/json
// @Param roundId query int true "场次ID"
// @Success 200 {object} response.Response
// @Router /console/game/status [get]
func (a *ConsoleApi) GetGameStatus(c *gin.Context) {
	roundIdStr := c.Query("roundId")
	roundId, _ := strconv.ParseUint(roundIdStr, 10, 64)

	result, err := gameService.GetGameStatus(uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetWinners 获取中奖名单
// @Tags Console-游戏
// @Summary 获取中奖名单
// @Produce application/json
// @Param roundId query int true "场次ID"
// @Success 200 {object} response.Response
// @Router /console/game/winners [get]
func (a *ConsoleApi) GetWinners(c *gin.Context) {
	roundIdStr := c.Query("roundId")
	roundId, _ := strconv.ParseUint(roundIdStr, 10, 64)

	result, err := gameService.GetWinners(uint(roundId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

func (a *ConsoleApi) GetCurrent(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := gameService.GetCurrent(uint(activityId))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// ==================== 抽奖 ====================

// RandomDraw 随机抽奖
// @Tags Console-抽奖
// @Summary 随机抽奖（从已签到用户中抽取）
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.ConsoleRandomDrawReq true "随机抽奖"
// @Success 200 {object} response.Response
// @Router /console/draw/random [post]
func (a *ConsoleApi) RandomDraw(c *gin.Context) {
	var req consoleReq.ConsoleRandomDrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := drawService.RandomDraw(req.ActivityId, req.Count, req.PrizeId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// DanmakuDraw 弹幕抽奖
// @Tags Console-抽奖
// @Summary 弹幕抽奖
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.DanmakuDrawReq true "弹幕抽奖"
// @Success 200 {object} response.Response
// @Router /console/draw/danmaku [post]
func (a *ConsoleApi) DanmakuDraw(c *gin.Context) {
	var req consoleReq.DanmakuDrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := drawService.DanmakuDraw(req.ActivityId, req.Count, req.PrizeId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetAllWinners 获取所有中奖名单
// @Tags Console-抽奖
// @Summary 获取所有中奖名单
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Param winType query int false "中奖类型"
// @Success 200 {object} response.Response
// @Router /console/winners [get]
func (a *ConsoleApi) GetAllWinners(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	winTypeStr := c.DefaultQuery("winType", "0")
	winType, _ := strconv.Atoi(winTypeStr)

	result, err := drawService.GetAllWinners(uint(activityId), winType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}
