package app

import (
	"github.com/gin-gonic/gin"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
)

type H5DanmakuApi struct{}

var h5DanmakuService = service.ServiceGroupApp.AppServiceGroup.H5DanmakuService

// SendDanmaku 发送弹幕
// @Tags H5-弹幕
// @Summary 发送弹幕
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SendDanmakuReq true "弹幕信息"
// @Success 200 {object} response.Response{data=appResp.DanmakuResp}
// @Router /h5/danmaku [post]
func (a *H5DanmakuApi) SendDanmaku(c *gin.Context) {
	var req request.SendDanmakuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userId := c.GetUint("h5UserId")

	result, err := h5DanmakuService.SendDanmaku(userId, req.ActivityId, req.Content, req.Color)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// GetDanmakuList 获取弹幕列表
// @Tags H5-弹幕
// @Summary 获取弹幕列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param activityId path int true "活动ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=appResp.H5PageResult}
// @Router /h5/danmaku/list/{activityId} [get]
func (a *H5DanmakuApi) GetDanmakuList(c *gin.Context) {
	var req request.DanmakuListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	result, err := h5DanmakuService.GetDanmakuList(req.ActivityId, req.Page, req.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
