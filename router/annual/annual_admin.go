package annual

import (
	"github.com/gin-gonic/gin"
	v1 "lotteryBackend/api/v1"
	"lotteryBackend/middleware"
)

type AnnualAdminRouter struct{}

func (r *AnnualAdminRouter) InitAnnualAdminRouter(Router *gin.RouterGroup) {
	annualApi := v1.ApiGroupApp.AnnualApiGroup

	// 需要认证的路由
	adminRouter := Router.Group("annual").Use(middleware.OperationRecord())
	adminRouterWithoutRecord := Router.Group("annual")
	// ========== 用户管理 ==========
	{
		adminRouter.POST("user/list", annualApi.AnnualUserApi.GetUserList)        // 获取用户列表
		adminRouter.DELETE("user", annualApi.AnnualUserApi.DeleteUser)            // 删除用户
		adminRouter.DELETE("user/batch", annualApi.AnnualUserApi.DeleteUserByIds) // 批量删除用户
	}

	// ========== 活动管理 ==========
	{
		adminRouter.POST("activity", annualApi.AnnualActivityApi.CreateActivity)              // 创建活动
		adminRouter.POST("activity/list", annualApi.AnnualActivityApi.GetActivityList)        // 获取活动列表
		adminRouter.GET("activity/:id", annualApi.AnnualActivityApi.GetActivityById)          // 获取活动详情
		adminRouter.PUT("activity", annualApi.AnnualActivityApi.UpdateActivity)               // 更新活动
		adminRouter.PUT("activity/status", annualApi.AnnualActivityApi.UpdateActivityStatus)  // 更新活动状态
		adminRouter.DELETE("activity", annualApi.AnnualActivityApi.DeleteActivity)            // 删除活动
		adminRouter.DELETE("activity/batch", annualApi.AnnualActivityApi.DeleteActivityByIds) // 批量删除
	}

	// ========== 签到管理 ==========
	{
		adminRouter.POST("checkIn/list", annualApi.AnnualCheckInApi.GetCheckInList)              // 获取签到列表
		adminRouter.GET("checkIn/stats/:activityId", annualApi.AnnualCheckInApi.GetCheckInStats) // 签到统计
		adminRouterWithoutRecord.GET("checkIn/export", annualApi.AnnualCheckInApi.ExportCheckIn) // 导出签到
	}

	// ========== 弹幕管理 ==========
	{
		adminRouter.POST("danmaku/list", annualApi.AnnualDanmakuApi.GetDanmakuList)        // 获取弹幕列表
		adminRouter.PUT("danmaku/audit", annualApi.AnnualDanmakuApi.AuditDanmaku)          // 审核弹幕
		adminRouter.PUT("danmaku/top", annualApi.AnnualDanmakuApi.TopDanmaku)              // 置顶弹幕
		adminRouter.DELETE("danmaku", annualApi.AnnualDanmakuApi.DeleteDanmaku)            // 删除弹幕
		adminRouter.DELETE("danmaku/batch", annualApi.AnnualDanmakuApi.DeleteDanmakuByIds) // 批量删除
	}

	// ========== 摇一摇场次管理 ==========
	{
		adminRouter.POST("shakeRound", annualApi.AnnualShakeRoundApi.CreateShakeRound)              // 创建场次
		adminRouter.POST("shakeRound/list", annualApi.AnnualShakeRoundApi.GetShakeRoundList)        // 获取场次列表
		adminRouter.GET("shakeRound/:id", annualApi.AnnualShakeRoundApi.GetShakeRoundById)          // 获取场次详情
		adminRouter.PUT("shakeRound", annualApi.AnnualShakeRoundApi.UpdateShakeRound)               // 更新场次
		adminRouter.PUT("shakeRound/start/:id", annualApi.AnnualShakeRoundApi.StartShakeRound)      // 开始游戏
		adminRouter.PUT("shakeRound/stop/:id", annualApi.AnnualShakeRoundApi.StopShakeRound)        // 结束游戏
		adminRouter.DELETE("shakeRound", annualApi.AnnualShakeRoundApi.DeleteShakeRound)            // 删除场次
		adminRouter.GET("shakeRound/scores/:roundId", annualApi.AnnualShakeRoundApi.GetShakeScores) // 获取成绩排行
	}

	// ========== 奖品管理 ==========
	{
		adminRouter.POST("prize", annualApi.AnnualPrizeApi.CreatePrize)              // 创建奖品
		adminRouter.POST("prize/list", annualApi.AnnualPrizeApi.GetPrizeList)        // 获取奖品列表
		adminRouter.GET("prize/:id", annualApi.AnnualPrizeApi.GetPrizeById)          // 获取奖品详情
		adminRouter.PUT("prize", annualApi.AnnualPrizeApi.UpdatePrize)               // 更新奖品
		adminRouter.DELETE("prize", annualApi.AnnualPrizeApi.DeletePrize)            // 删除奖品
		adminRouter.DELETE("prize/batch", annualApi.AnnualPrizeApi.DeletePrizeByIds) // 批量删除
	}

	// ========== 中奖记录管理 ==========
	{
		adminRouter.POST("winner/list", annualApi.AnnualWinnerApi.GetWinnerList)              // 获取中奖列表
		adminRouter.PUT("winner/receive", annualApi.AnnualWinnerApi.ConfirmReceive)           // 确认领奖
		adminRouter.DELETE("winner", annualApi.AnnualWinnerApi.DeleteWinner)                  // 删除中奖记录
		adminRouterWithoutRecord.GET("winner/export", annualApi.AnnualWinnerApi.ExportWinner) // 导出中奖
		adminRouter.POST("winner/draw", annualApi.AnnualWinnerApi.RandomDraw)                 // 随机抽奖
	}
	// ========== 统计面板 ==========
	{
		adminRouterWithoutRecord.GET("dashboard/stats/:activityId", annualApi.AnnualDashboardApi.GetDashboardStats)
		adminRouterWithoutRecord.GET("dashboard/checkInTrend/:activityId", annualApi.AnnualDashboardApi.GetCheckInTrend)
		adminRouterWithoutRecord.GET("dashboard/prizeStats/:activityId", annualApi.AnnualDashboardApi.GetPrizeStats)
		adminRouterWithoutRecord.GET("dashboard/recentWinners/:activityId", annualApi.AnnualDashboardApi.GetRecentWinners)
		adminRouterWithoutRecord.GET("dashboard/hotWords/:activityId", annualApi.AnnualDashboardApi.GetHotWords)
		adminRouterWithoutRecord.GET("dashboard/recentDanmaku/:activityId", annualApi.AnnualDashboardApi.GetRecentDanmaku)
		adminRouterWithoutRecord.GET("dashboard/shakeRanking/:activityId", annualApi.AnnualDashboardApi.GetShakeRanking)
	}

	// ========== 系统配置 ==========
	{
		adminRouter.GET("config/list", annualApi.AnnualConfigApi.GetConfigList)      // 获取配置列表
		adminRouter.GET("config/all", annualApi.AnnualConfigApi.GetAllConfig)        // 获取所有配置
		adminRouter.GET("config/:key", annualApi.AnnualConfigApi.GetConfigByKey)     // 获取单个配置
		adminRouter.POST("config", annualApi.AnnualConfigApi.SetConfig)              // 设置配置
		adminRouter.POST("config/batch", annualApi.AnnualConfigApi.BatchSetConfig)   // 批量设置配置
		adminRouter.DELETE("config/:key", annualApi.AnnualConfigApi.DeleteConfig)    // 删除配置
		adminRouter.GET("config/wechat", annualApi.AnnualConfigApi.GetWechatConfig)  // 获取微信配置
		adminRouter.POST("config/wechat", annualApi.AnnualConfigApi.SetWechatConfig) // 设置微信配置
		adminRouter.GET("config/site", annualApi.AnnualConfigApi.GetSiteConfig)      // 获取站点配置
		adminRouter.POST("config/site", annualApi.AnnualConfigApi.SetSiteConfig)     // 设置站点配置
	}
}
