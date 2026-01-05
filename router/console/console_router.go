package console

import (
	"github.com/gin-gonic/gin"
	api "lotteryBackend/api/v1"
)

type ConsoleRouter struct{}

func (r *ConsoleRouter) InitConsoleRouter(Router *gin.RouterGroup) {
	consoleApi := api.ApiGroupApp.ConsoleApi.ConsoleApi
	checkInApi := api.ApiGroupApp.ConsoleApi.CheckInApi

	// ========== 控制台接口 ==========
	consoleGroup := Router.Group("console")
	{
		// ---------- 活动配置 ----------
		consoleGroup.GET("/activity/:activityId", consoleApi.GetActivityDetail) // 获取活动详情
		consoleGroup.GET("/prizes", consoleApi.GetPrizeList)                    // 获取奖品列表

		// ---------- 签到管理 ----------
		consoleGroup.GET("/checkin/stats", checkInApi.GetCheckInStats) // 获取签到统计（包含状态）
		consoleGroup.POST("/checkin/open", checkInApi.OpenCheckIn)     // 开启签到
		consoleGroup.POST("/checkin/close", checkInApi.CloseCheckIn)   // 关闭签到

		// ---------- 弹幕管理 ----------
		consoleGroup.GET("/danmaku/list", consoleApi.GetDanmakuList) // 获取弹幕列表
		consoleGroup.POST("/danmaku/audit", consoleApi.AuditDanmaku) // 审核弹幕

		// ---------- 场次管理 ----------
		consoleGroup.GET("/rounds", consoleApi.GetRoundList)            // 获取场次列表
		consoleGroup.GET("/rounds/:roundId", consoleApi.GetRoundDetail) // 获取场次详情

		// ---------- 游戏控制 ----------
		consoleGroup.POST("/game/start", consoleApi.StartGame)     // 开始游戏（验证密码）
		consoleGroup.POST("/game/stop", consoleApi.StopGame)       // 停止游戏
		consoleGroup.POST("/game/cancel", consoleApi.CancelGame)   // 取消游戏
		consoleGroup.GET("/game/status", consoleApi.GetGameStatus) // 获取游戏状态
		consoleGroup.GET("/game/ranking", consoleApi.GetRanking)   // 获取排行榜
		consoleGroup.GET("/game/winners", consoleApi.GetWinners)   // 获取中奖名单
		consoleGroup.GET("/game/current", consoleApi.GetCurrent)   // 获取游戏状态

		// ---------- 抽奖 ----------
		consoleGroup.POST("/draw/random", consoleApi.RandomDraw)   // 随机抽奖
		consoleGroup.POST("/draw/danmaku", consoleApi.DanmakuDraw) // 弹幕抽奖
	}
}
