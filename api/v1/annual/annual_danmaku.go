package annual

import (
	"lotteryBackend/global"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnnualDanmakuApi struct{}

var annualDanmakuService = service.ServiceGroupApp.AnnualServiceGroup.AnnualDanmakuService

// GetDanmakuList 获取弹幕列表
func (a *AnnualDanmakuApi) GetDanmakuList(c *gin.Context) {
	var pageInfo annualReq.DanmakuSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualDanmakuService.GetDanmakuList(pageInfo)
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

// AuditDanmaku 审核弹幕
func (a *AnnualDanmakuApi) AuditDanmaku(c *gin.Context) {
	var req annualReq.DanmakuAudit
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualDanmakuService.AuditDanmaku(req.Ids, req.Status); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("审核成功", c)
}

// TopDanmaku 置顶弹幕
func (a *AnnualDanmakuApi) TopDanmaku(c *gin.Context) {
	var req annualReq.DanmakuTop
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualDanmakuService.TopDanmaku(req.ID, req.IsTop); err != nil {
		global.GVA_LOG.Error("操作失败!", zap.Error(err))
		response.FailWithMessage("操作失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// DeleteDanmaku 删除弹幕
func (a *AnnualDanmakuApi) DeleteDanmaku(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualDanmakuService.DeleteDanmaku(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDanmakuByIds 批量删除弹幕
func (a *AnnualDanmakuApi) DeleteDanmakuByIds(c *gin.Context) {
	var req annualReq.GetByIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualDanmakuService.DeleteDanmakuByIds(req.Ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
