package annual

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type AnnualCheckInApi struct{}

var annualCheckInService = service.ServiceGroupApp.AnnualServiceGroup.AnnualCheckInService

// GetCheckInList 获取签到列表
func (a *AnnualCheckInApi) GetCheckInList(c *gin.Context) {
	var pageInfo annualReq.CheckInSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualCheckInService.GetCheckInList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetCheckInStats 签到统计
func (a *AnnualCheckInApi) GetCheckInStats(c *gin.Context) {
	activityId := c.Param("activityId")
	stats, err := annualCheckInService.GetCheckInStats(activityId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(stats, c)
}

// UpdateCheckIn 更新签到信息
func (a *AnnualCheckInApi) UpdateCheckIn(c *gin.Context) {
	var req annualReq.CheckInUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualCheckInService.UpdateCheckIn(req); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateCheckInStatus 更新签到状态
func (a *AnnualCheckInApi) UpdateCheckInStatus(c *gin.Context) {
	var req annualReq.CheckInStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualCheckInService.UpdateCheckInStatus(req); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("审核成功", c)
}

// DeleteCheckIn 删除签到
func (a *AnnualCheckInApi) DeleteCheckIn(c *gin.Context) {
	var req annualReq.CheckInDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualCheckInService.DeleteCheckIn(req); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ExportCheckIn 导出签到
func (a *AnnualCheckInApi) ExportCheckIn(c *gin.Context) {
	var pageInfo annualReq.CheckInSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	filePath, err := annualCheckInService.ExportCheckIn(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败: "+err.Error(), c)
		return
	}
	response.OkWithData(filePath, c)
}
