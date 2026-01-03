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

type AnnualUserApi struct{}

var annualUserService = service.ServiceGroupApp.AnnualServiceGroup.AnnualUserService

// GetUserList 获取用户列表
func (a *AnnualUserApi) GetUserList(c *gin.Context) {
	var pageInfo annualReq.UserSearch
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := annualUserService.GetUserList(pageInfo)
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

// GetUserById 获取用户详情
func (a *AnnualUserApi) GetUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := annualUserService.GetUserById(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(user, c)
}

// UpdateUser 更新用户
func (a *AnnualUserApi) UpdateUser(c *gin.Context) {
	var user annual.AnnualUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualUserService.UpdateUser(user); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateUserStatus 更新用户状态
func (a *AnnualUserApi) UpdateUserStatus(c *gin.Context) {
	var req annualReq.UpdateStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualUserService.UpdateUserStatus(req.ID, req.Status); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteUser 删除用户
func (a *AnnualUserApi) DeleteUser(c *gin.Context) {
	var req annualReq.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualUserService.DeleteUser(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteUserByIds 批量删除用户
func (a *AnnualUserApi) DeleteUserByIds(c *gin.Context) {
	var req annualReq.GetByIds
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := annualUserService.DeleteUserByIds(req.Ids); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ExportUser 导出用户
func (a *AnnualUserApi) ExportUser(c *gin.Context) {
	var pageInfo annualReq.UserSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	filePath, err := annualUserService.ExportUser(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败: "+err.Error(), c)
		return
	}
	response.OkWithData(filePath, c)
}
