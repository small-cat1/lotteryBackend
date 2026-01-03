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
		adminRouter.POST("user/list", annualApi.AnnualUserApi.GetUserList)              // 获取用户列表
		adminRouter.GET("user/:id", annualApi.AnnualUserApi.GetUserById)                // 获取用户详情
		adminRouter.PUT("user", annualApi.AnnualUserApi.UpdateUser)                     // 更新用户
		adminRouter.PUT("user/status", annualApi.AnnualUserApi.UpdateUserStatus)        // 更新用户状态
		adminRouter.DELETE("user", annualApi.AnnualUserApi.DeleteUser)                  // 删除用户
		adminRouter.DELETE("user/batch", annualApi.AnnualUserApi.DeleteUserByIds)       // 批量删除用户
		adminRouterWithoutRecord.GET("user/export", annualApi.AnnualUserApi.ExportUser) // 导出用户
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
}
