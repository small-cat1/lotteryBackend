package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/model/common/response"
	"lotteryBackend/service"
	"lotteryBackend/ws"

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
	// 如果是审核通过（status == 1），先查询弹幕信息用于推送
	var danmakuList []annual.AnnualDanmaku
	if req.Status == 1 {
		global.GVA_DB.Preload("User").Where("id IN ?", req.Ids).Find(&danmakuList)
	}
	if err := annualDanmakuService.AuditDanmaku(req.Ids, req.Status); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败: "+err.Error(), c)
		return
	}
	// 审核通过后，通过 WebSocket 推送弹幕
	if req.Status == 1 && len(danmakuList) > 0 {
		broadcaster := ws.GetBroadcaster()
		// 收集所有 activityId + userId 组合
		var userIds []uint
		var activityId uint
		for _, d := range danmakuList {
			userIds = append(userIds, d.UserId)
			activityId = d.ActivityId // 同一批弹幕应该是同一个活动
		}

		// 批量查询签到信息
		var checkIns []annual.AnnualCheckIn
		global.GVA_DB.Where("activity_id = ? AND user_id IN ?", activityId, userIds).Find(&checkIns)

		// 构建 userId -> checkIn 的映射
		checkInMap := make(map[uint]annual.AnnualCheckIn)
		for _, ci := range checkIns {
			checkInMap[ci.UserId] = ci
		}

		for _, danmaku := range danmakuList {
			checkIn := checkInMap[danmaku.UserId]

			payload := ws.NewDanmakuPayload(
				danmaku.ID,
				ws.NewUserBrief(
					danmaku.User.ID,
					danmaku.User.Nickname,
					danmaku.User.Avatar,
					checkIn.RealName,
					checkIn.Department,
				),
				danmaku.Content,
				danmaku.Color,
				*danmaku.IsTop,
				danmaku.CreatedAt,
			)
			broadcaster.BroadcastDanmaku(danmaku.ActivityId, payload)
		}
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
