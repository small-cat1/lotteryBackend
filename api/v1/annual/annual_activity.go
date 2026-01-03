package annual

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"lotteryBackend/utils"
)

type AnnualActivityApi struct{}

var annualActivityService = service.ServiceGroupApp.AnnualServiceGroup.AnnualActivityService

// CreateActivity 创建活动
func (a *AnnualActivityApi) CreateActivity(c *gin.Context) {
	var activity annual.AnnualActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	activity.CreatedBy = utils.GetUserID(c)
	if err := annualActivityService.CreateActivity(activity); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetActivityList 获取活动列表
func (a *AnnualActivityApi) GetActivityList(c *gin.Context) {
	var pageInfo annualReq.ActivitySearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualActivityService.GetActivityList(pageInfo)
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

// GetActivityById 获取活动详情
func (a *AnnualActivityApi) GetActivityById(c *gin.Context) {
	id := c.Param("id")
	activity, err := annualActivityService.GetActivityById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(activity, c)
}

// UpdateActivity 更新活动
func (a *AnnualActivityApi) UpdateActivity(c *gin.Context) {
	var activity annual.AnnualActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualActivityService.UpdateActivity(activity); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateActivityStatus 更新活动状态
func (a *AnnualActivityApi) UpdateActivityStatus(c *gin.Context) {
	var req annualReq.UpdateStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualActivityService.UpdateActivityStatus(req.ID, req.Status); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteActivity 删除活动
func (a *AnnualActivityApi) DeleteActivity(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualActivityService.DeleteActivity(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteActivityByIds 批量删除活动
func (a *AnnualActivityApi) DeleteActivityByIds(c *gin.Context) {
	var req annualReq.GetByIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualActivityService.DeleteActivityByIds(req.Ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
