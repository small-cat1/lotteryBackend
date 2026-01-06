package app

import (
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type H5ActivityApi struct{}

var h5ActivityService = service.ServiceGroupApp.AppServiceGroup.H5ActivityService

// GetActivityDetail 获取活动详情
// @Tags H5-活动
// @Summary 获取活动详情
// @Produce application/json
// @Param id path int true "活动ID"
// @Success 200 {object} response.Response{data=appResp.H5ActivityResp}
// @Router /h5/activity/{id} [get]
func (a *H5ActivityApi) GetActivityDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	result, err := h5ActivityService.GetActivityDetail(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}
