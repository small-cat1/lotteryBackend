package app

import (
	"github.com/gin-gonic/gin"
	api "lotteryBackend/api/v1"
	"lotteryBackend/middleware"
)

type H5Router struct{}

func (r *H5Router) InitH5Router(Router *gin.RouterGroup) {
	h5Api := api.ApiGroupApp.H5ApiGroup

	// ========== 公开接口（无需登录） ==========
	h5Public := Router.Group("h5")
	{
		// 授权相关
		h5Public.POST("auth/wechat", h5Api.WechatLogin)    // 微信授权登录
		h5Public.GET("auth/wxConfig", h5Api.GetWxJsConfig) // 获取微信JS-SDK配置

		// 配置相关
		h5Public.GET("config/wechat", h5Api.GetWechatConfig) // 获取微信AppID等配置

		// 活动相关（公开）
		h5Public.GET("activity/current", h5Api.GetCurrentActivity) // 获取当前活动
		h5Public.GET("activity/list", h5Api.GetActivityList)       // 获取活动列表
		h5Public.GET("activity/:id", h5Api.GetActivityDetail)      // 获取活动详情

		// 大屏数据（公开）
		h5Public.GET("screen/checkIn/:activityId", h5Api.GetScreenCheckInList) // 签到墙数据
		h5Public.GET("screen/danmaku/:activityId", h5Api.GetScreenDanmakuList) // 弹幕墙数据
		h5Public.GET("screen/ranking/:roundId", h5Api.GetScreenShakeRanking)   // 摇一摇排行
		h5Public.GET("screen/winners/:activityId", h5Api.GetScreenWinnerList)  // 中奖名单
	}

	// ========== 需要登录的接口 ==========
	h5Auth := Router.Group("h5")
	h5Auth.Use(middleware.H5Auth())
	{
		// 授权相关
		h5Auth.POST("auth/refresh", h5Api.RefreshToken) // 刷新Token
		h5Auth.POST("auth/logout", h5Api.Logout)        // 退出登录

		// 用户相关
		h5Auth.GET("user/info", h5Api.GetUserInfo)       // 获取用户信息
		h5Auth.POST("user/register", h5Api.UserRegister) // 用户报名
		h5Auth.PUT("user/info", h5Api.UpdateUserInfo)    // 更新用户信息
		h5Auth.GET("user/audit", h5Api.GetAuditStatus)   // 获取审核状态
	}

	// ========== 需要登录且已报名的接口 ==========
	h5Private := Router.Group("h5")
	h5Private.Use(middleware.H5Auth(), middleware.H5RegisteredRequired())
	{
		// 签到相关
		h5Private.POST("checkIn", h5Api.CheckIn)                            // 用户签到
		h5Private.GET("checkIn/status/:activityId", h5Api.GetCheckInStatus) // 获取签到状态
		h5Private.GET("checkIn/list/:activityId", h5Api.GetCheckInList)     // 获取签到列表
		//h5Private.GET("checkIn/stats/:activityId", h5Api.GetCheckInStats)    // 获取签到统计
		h5Private.GET("checkIn/recent/:activityId", h5Api.GetRecentCheckIns) // 获取最新签到

		// 弹幕相关
		h5Private.POST("danmaku", h5Api.SendDanmaku) // 发送弹幕
		//h5Private.GET("danmaku/list/:activityId", h5Api.GetDanmakuList)     // 获取弹幕列表
		h5Private.GET("danmaku/recent/:activityId", h5Api.GetRecentDanmaku) // 获取最新弹幕
		h5Private.GET("danmaku/top/:activityId", h5Api.GetTopDanmaku)       // 获取置顶弹幕
		h5Private.GET("danmaku/my/:activityId", h5Api.GetMyDanmaku)         // 获取我的弹幕

		// 摇一摇相关
		h5Private.GET("shake/round/current/:activityId", h5Api.GetCurrentRound) // 获取当前场次
		//h5Private.GET("shake/round/list/:activityId", h5Api.GetRoundList)       // 获取场次列表
		h5Private.GET("shake/round/:roundId", h5Api.GetRoundDetail)    // 获取场次详情
		h5Private.POST("shake/join", h5Api.JoinGame)                   // 加入游戏
		h5Private.POST("shake/score", h5Api.SubmitScore)               // 提交分数
		h5Private.GET("shake/ranking/:roundId", h5Api.GetShakeRanking) // 获取实时排名
		h5Private.GET("shake/my/:roundId", h5Api.GetMyScore)           // 获取我的成绩
		h5Private.GET("shake/result/:roundId", h5Api.GetRoundResult)   // 获取最终结果

		// 奖品相关
		h5Private.GET("prize/list/:activityId", h5Api.GetPrizeList)        // 获取奖品列表
		h5Private.GET("prize/my/:activityId", h5Api.GetMyWinnings)         // 获取我的中奖记录
		h5Private.GET("prize/recent/:activityId", h5Api.GetRecentWinnings) // 获取最新中奖
		h5Private.GET("prize/winning/:winnerId", h5Api.GetWinningDetail)   // 获取中奖详情
		h5Private.GET("prize/qrcode/:winnerId", h5Api.GetReceiveQrCode)    // 获取领奖二维码
		h5Private.GET("prize/:prizeId", h5Api.GetPrizeDetail)              // 获取奖品详情
	}
}
