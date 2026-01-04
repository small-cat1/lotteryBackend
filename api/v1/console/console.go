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

// ==================== 签到管理 ====================

// GetCheckInStats 获取签到统计
// @Tags Console-签到
// @Summary 获取签到统计（包含状态）
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Success 200 {object} response.Response
// @Router /console/checkin/stats [get]
func (a *ConsoleApi) GetCheckInStats(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)

	result, err := checkInService.GetCheckInStats(uint(activityId))
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
func (a *ConsoleApi) OpenCheckIn(c *gin.Context) {
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
func (a *ConsoleApi) CloseCheckIn(c *gin.Context) {
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

// GetCheckInList 获取签到列表
// @Tags Console-签到
// @Summary 获取签到列表
// @Produce application/json
// @Param activityId query int true "活动ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response
// @Router /console/checkin/list [get]
func (a *ConsoleApi) GetCheckInList(c *gin.Context) {
	activityIdStr := c.Query("activityId")
	activityId, _ := strconv.ParseUint(activityIdStr, 10, 64)
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	keyword := c.Query("keyword")

	result, err := checkInService.GetCheckInList(uint(activityId), page, pageSize, keyword)
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

// AuditDanmaku 审核弹幕
// @Tags Console-弹幕
// @Summary 审核弹幕
// @Accept application/json
// @Produce application/json
// @Param data body consoleReq.AuditDanmakuReq true "审核弹幕"
// @Success 200 {object} response.Response
// @Router /console/danmaku/audit [post]
func (a *ConsoleApi) AuditDanmaku(c *gin.Context) {
	var req consoleReq.AuditDanmakuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := danmakuService.AuditDanmaku(req.DanmakuId, req.Status)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("审核成功", c)
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

// GetRanking 获取排行榜
// @Tags Console-游戏
// @Summary 获取排行榜
// @Produce application/json
// @Param roundId query int true "场次ID"
// @Param limit query int false "数量限制"
// @Success 200 {object} response.Response
// @Router /console/game/ranking [get]
func (a *ConsoleApi) GetRanking(c *gin.Context) {
	roundIdStr := c.Query("roundId")
	roundId, _ := strconv.ParseUint(roundIdStr, 10, 64)
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	result, err := gameService.GetRanking(uint(roundId), limit)
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
