package annual

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type AnnualShakeRoundApi struct{}

var annualShakeRoundService = service.ServiceGroupApp.AnnualServiceGroup.AnnualShakeRoundService

// CreateShakeRound 创建摇一摇场次
func (a *AnnualShakeRoundApi) CreateShakeRound(c *gin.Context) {
	var round annual.AnnualShakeRound
	if err := c.ShouldBindJSON(&round); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualShakeRoundService.CreateShakeRound(round); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetShakeRoundList 获取场次列表
func (a *AnnualShakeRoundApi) GetShakeRoundList(c *gin.Context) {
	var pageInfo annualReq.ShakeRoundSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualShakeRoundService.GetShakeRoundList(pageInfo)
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

// GetShakeRoundById 获取场次详情
func (a *AnnualShakeRoundApi) GetShakeRoundById(c *gin.Context) {
	id := c.Param("id")
	round, err := annualShakeRoundService.GetShakeRoundById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(round, c)
}

// UpdateShakeRound 更新场次
func (a *AnnualShakeRoundApi) UpdateShakeRound(c *gin.Context) {
	var round annual.AnnualShakeRound
	if err := c.ShouldBindJSON(&round); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualShakeRoundService.UpdateShakeRound(round); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// StartShakeRound 开始游戏
func (a *AnnualShakeRoundApi) StartShakeRound(c *gin.Context) {
	id := c.Param("id")
	if err := annualShakeRoundService.StartShakeRound(id); err != nil {
		global.GVA_LOG.Error("开始失败!", zap.Error(err))
		response.FailWithMessage("开始失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("游戏已开始", c)
}

// StopShakeRound 结束游戏
func (a *AnnualShakeRoundApi) StopShakeRound(c *gin.Context) {
	id := c.Param("id")
	winners, err := annualShakeRoundService.StopShakeRound(id)
	if err != nil {
		global.GVA_LOG.Error("结束失败!", zap.Error(err))
		response.FailWithMessage("结束失败: "+err.Error(), c)
		return
	}
	response.OkWithData(winners, c)
}

// DeleteShakeRound 删除场次
func (a *AnnualShakeRoundApi) DeleteShakeRound(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualShakeRoundService.DeleteShakeRound(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetShakeScores 获取成绩排行
func (a *AnnualShakeRoundApi) GetShakeScores(c *gin.Context) {
	roundId := c.Param("roundId")
	scores, err := annualShakeRoundService.GetShakeScores(roundId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(scores, c)
}
