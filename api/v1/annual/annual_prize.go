package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnnualPrizeApi struct{}

var annualPrizeService = service.ServiceGroupApp.AnnualServiceGroup.AnnualPrizeService

// CreatePrize 创建奖品
func (a *AnnualPrizeApi) CreatePrize(c *gin.Context) {
	var prize annual.AnnualPrize
	if err := c.ShouldBindJSON(&prize); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 初始化剩余数量
	prize.RemainCount = prize.TotalCount
	if err := annualPrizeService.CreatePrize(prize); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetPrizeList 获取奖品列表
func (a *AnnualPrizeApi) GetPrizeList(c *gin.Context) {
	var pageInfo annualReq.PrizeSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualPrizeService.GetPrizeList(pageInfo)
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

// GetPrizeById 获取奖品详情
func (a *AnnualPrizeApi) GetPrizeById(c *gin.Context) {
	id := c.Param("id")
	prize, err := annualPrizeService.GetPrizeById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(prize, c)
}

// UpdatePrize 更新奖品
func (a *AnnualPrizeApi) UpdatePrize(c *gin.Context) {
	var prize annual.AnnualPrize
	if err := c.ShouldBindJSON(&prize); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualPrizeService.UpdatePrize(prize); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeletePrize 删除奖品
func (a *AnnualPrizeApi) DeletePrize(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualPrizeService.DeletePrize(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePrizeByIds 批量删除奖品
func (a *AnnualPrizeApi) DeletePrizeByIds(c *gin.Context) {
	var req annualReq.GetByIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualPrizeService.DeletePrizeByIds(req.Ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
