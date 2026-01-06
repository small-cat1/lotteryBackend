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
		h5Public.GET("config/wechat", h5Api.GetWechatConfig)
		h5Public.GET("activity/:id", h5Api.GetActivityDetail)
	}

	// ========== 需要登录 ==========
	h5Auth := Router.Group("h5")
	h5Auth.Use(middleware.H5Auth())
	{
		h5Auth.GET("user/info", h5Api.GetUserInfo)
		h5Auth.POST("user/checkIn", h5Api.CheckIn) // 签到
		// 弹幕（登录就能用）
		h5Auth.POST("danmaku", h5Api.SendDanmaku)
	}

	// ========== 需要签到且审核通过 ==========
	h5Private := Router.Group("h5")
	h5Private.Use(middleware.H5Auth(), middleware.H5RegisteredRequired())
	{
		// 摇一摇
		h5Private.GET("shake/round/current/:activityId", h5Api.GetCurrentRound)
		h5Private.GET("shake/result/:roundId", h5Api.GetRoundResult)

		// 奖品
		h5Private.GET("prize/my/:activityId", h5Api.GetMyWinnings)
		h5Private.GET("prize/winning/:winnerId", h5Api.GetWinningDetail)
	}
}
