package app

import (
	"github.com/gin-gonic/gin"
	api "lotteryBackend/api/v1"
	"lotteryBackend/middleware"
)

type H5Router struct{}

func (r *H5Router) InitH5Router(Router *gin.RouterGroup) {
	h5Api := api.ApiGroupApp.H5ApiGroup

	// ========== 公开接口 ==========
	h5Public := Router.Group("h5")
	{
		h5Public.POST("auth/wechat", h5Api.WechatLogin)
		h5Public.GET("auth/wxConfig", h5Api.GetWxJsConfig)
		h5Public.GET("config/wechat", h5Api.GetWechatConfig)

		h5Public.GET("activity/current", h5Api.GetCurrentActivity)
		h5Public.GET("activity/list", h5Api.GetActivityList)
		h5Public.GET("activity/:id", h5Api.GetActivityDetail)

		h5Public.GET("screen/checkIn/:activityId", h5Api.GetScreenCheckInList)
		h5Public.GET("screen/danmaku/:activityId", h5Api.GetScreenDanmakuList)
		h5Public.GET("screen/ranking/:roundId", h5Api.GetScreenShakeRanking)
		h5Public.GET("screen/winners/:activityId", h5Api.GetScreenWinnerList)
	}

	// ========== 需要登录 ==========
	h5Auth := Router.Group("h5")
	h5Auth.Use(middleware.H5Auth())
	{
		h5Auth.POST("auth/refresh", h5Api.RefreshToken)
		h5Auth.POST("auth/logout", h5Api.Logout)

		h5Auth.GET("user/info", h5Api.GetUserInfo)
		h5Auth.POST("user/checkIn", h5Api.CheckIn) // 签到
		h5Auth.GET("user/audit", h5Api.GetAuditStatus)

		// 弹幕（登录就能用）
		h5Auth.POST("danmaku", h5Api.SendDanmaku)
		h5Auth.GET("danmaku/recent/:activityId", h5Api.GetRecentDanmaku)
		h5Auth.GET("danmaku/top/:activityId", h5Api.GetTopDanmaku)
		h5Auth.GET("danmaku/my/:activityId", h5Api.GetMyDanmaku)
	}

	// ========== 需要签到且审核通过 ==========
	h5Private := Router.Group("h5")
	h5Private.Use(middleware.H5Auth(), middleware.H5RegisteredRequired())
	{
		// 摇一摇
		h5Private.GET("shake/round/current/:activityId", h5Api.GetCurrentRound)
		h5Private.GET("shake/round/:roundId", h5Api.GetRoundDetail)
		h5Private.POST("shake/join", h5Api.JoinGame)
		h5Private.POST("shake/score", h5Api.SubmitScore)
		h5Private.GET("shake/ranking/:roundId", h5Api.GetShakeRanking)
		h5Private.GET("shake/my/:roundId", h5Api.GetMyScore)
		h5Private.GET("shake/result/:roundId", h5Api.GetRoundResult)

		// 奖品
		h5Private.GET("prize/list/:activityId", h5Api.GetPrizeList)
		h5Private.GET("prize/my/:activityId", h5Api.GetMyWinnings)
		h5Private.GET("prize/recent/:activityId", h5Api.GetRecentWinnings)
		h5Private.GET("prize/winning/:winnerId", h5Api.GetWinningDetail)
		h5Private.GET("prize/qrcode/:winnerId", h5Api.GetReceiveQrCode)
		h5Private.GET("prize/:prizeId", h5Api.GetPrizeDetail)
	}
}
