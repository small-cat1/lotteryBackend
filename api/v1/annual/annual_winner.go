package annual

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type AnnualWinnerApi struct{}

var annualWinnerService = service.ServiceGroupApp.AnnualServiceGroup.AnnualWinnerService

// GetWinnerList 获取中奖列表
func (a *AnnualWinnerApi) GetWinnerList(c *gin.Context) {
	var pageInfo annualReq.WinnerSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualWinnerService.GetWinnerList(pageInfo)
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

// ConfirmReceive 确认领奖
func (a *AnnualWinnerApi) ConfirmReceive(c *gin.Context) {
	var req annualReq.WinnerReceive
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualWinnerService.ConfirmReceive(req.Id, req.VerifyCode); err != nil {
		global.GVA_LOG.Error("核销失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("核销成功", c)
}

// DeleteWinner 删除中奖记录
func (a *AnnualWinnerApi) DeleteWinner(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualWinnerService.DeleteWinner(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ExportWinner 导出中奖记录
func (a *AnnualWinnerApi) ExportWinner(c *gin.Context) {
	var pageInfo annualReq.WinnerSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	filePath, err := annualWinnerService.ExportWinner(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败: "+err.Error(), c)
		return
	}
	response.OkWithData(filePath, c)
}

// RandomDraw 随机抽奖
func (a *AnnualWinnerApi) RandomDraw(c *gin.Context) {
	var req annualReq.RandomDraw
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	winners, err := annualWinnerService.RandomDraw(req)
	if err != nil {
		global.GVA_LOG.Error("抽奖失败!", zap.Error(err))
		response.FailWithMessage("抽奖失败: "+err.Error(), c)
		return
	}
	response.OkWithData(winners, c)
}
