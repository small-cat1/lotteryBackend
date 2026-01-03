package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnnualDashboardApi struct{}

var annualDashboardService = service.ServiceGroupApp.AnnualServiceGroup.AnnualDashboardService

// GetDashboardStats 获取统计概览数据
func (a *AnnualDashboardApi) GetDashboardStats(c *gin.Context) {
	activityId := c.Param("activityId")
	stats, err := annualDashboardService.GetDashboardStats(activityId)
	if err != nil {
		global.GVA_LOG.Error("获取统计数据失败!", zap.Error(err))
		response.FailWithMessage("获取统计数据失败: "+err.Error(), c)
		return
	}
	response.OkWithData(stats, c)
}

// GetCheckInTrend 获取签到趋势
func (a *AnnualDashboardApi) GetCheckInTrend(c *gin.Context) {
	activityId := c.Param("activityId")
	list, err := annualDashboardService.GetCheckInTrend(activityId)
	if err != nil {
		global.GVA_LOG.Error("获取签到趋势失败!", zap.Error(err))
		response.FailWithMessage("获取签到趋势失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetPrizeStats 获取奖品统计
func (a *AnnualDashboardApi) GetPrizeStats(c *gin.Context) {
	activityId := c.Param("activityId")
	list, err := annualDashboardService.GetPrizeStats(activityId)
	if err != nil {
		global.GVA_LOG.Error("获取奖品统计失败!", zap.Error(err))
		response.FailWithMessage("获取奖品统计失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetRecentWinners 获取最新中奖记录
func (a *AnnualDashboardApi) GetRecentWinners(c *gin.Context) {
	activityId := c.Param("activityId")
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	list, err := annualDashboardService.GetRecentWinners(activityId, limit)
	if err != nil {
		global.GVA_LOG.Error("获取最新中奖记录失败!", zap.Error(err))
		response.FailWithMessage("获取最新中奖记录失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetHotWords 获取弹幕热词
func (a *AnnualDashboardApi) GetHotWords(c *gin.Context) {
	activityId := c.Param("activityId")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	list, err := annualDashboardService.GetHotWords(activityId, limit)
	if err != nil {
		global.GVA_LOG.Error("获取弹幕热词失败!", zap.Error(err))
		response.FailWithMessage("获取弹幕热词失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetRecentDanmaku 获取最新弹幕
func (a *AnnualDashboardApi) GetRecentDanmaku(c *gin.Context) {
	activityId := c.Param("activityId")
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	list, err := annualDashboardService.GetRecentDanmaku(activityId, limit)
	if err != nil {
		global.GVA_LOG.Error("获取最新弹幕失败!", zap.Error(err))
		response.FailWithMessage("获取最新弹幕失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetShakeRanking 获取摇一摇排行榜
func (a *AnnualDashboardApi) GetShakeRanking(c *gin.Context) {
	activityId := c.Param("activityId")
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	list, err := annualDashboardService.GetShakeRanking(activityId, limit)
	if err != nil {
		global.GVA_LOG.Error("获取摇一摇排行榜失败!", zap.Error(err))
		response.FailWithMessage("获取摇一摇排行榜失败: "+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}
