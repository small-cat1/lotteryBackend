package app

import (
	"github.com/gin-gonic/gin"
	api "lotteryBackend/api/v1"
)

type ConsoleRouter struct{}

func (r *ConsoleRouter) InitConsoleRouter(Router *gin.RouterGroup) {
	consoleApi := api.ApiGroupApp.H5ApiGroup.ConsoleApi

	// ========== 大屏数据接口（公开） ==========
	screenGroup := Router.Group("h5/screen")
	{
		// 签到
		screenGroup.GET("/checkIn/stats/:activityId", consoleApi.GetCheckInStats)

		// 场次
		screenGroup.GET("/rounds/:activityId", consoleApi.GetRoundList)

		// 弹幕
		screenGroup.GET("/danmaku/:activityId", consoleApi.GetDanmakuList)

		// 排行榜
		screenGroup.GET("/ranking/:roundId", consoleApi.GetRanking)

		// 中奖名单
		screenGroup.GET("/winners/:activityId", consoleApi.GetWinners)

		// 配置
		screenGroup.GET("/config/randomDraw", consoleApi.GetRandomDrawConfig)
	}

	// ========== 控制台接口（公开，场次密码验证） ==========
	consoleGroup := Router.Group("console")
	{
		// 游戏控制
		consoleGroup.POST("/game/start", consoleApi.StartGame)   // 开始游戏（验证密码）
		consoleGroup.POST("/game/stop", consoleApi.StopGame)     // 立即停止
		consoleGroup.POST("/game/settle", consoleApi.SettleGame) // 结算

		// 随机抽奖
		consoleGroup.POST("/randomDraw", consoleApi.RandomDraw)
	}
}
