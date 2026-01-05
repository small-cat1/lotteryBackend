package ws

import (
	"time"
)

// EventTrigger WebSocket事件触发器
// 用于在业务服务中触发WebSocket事件推送
type EventTrigger struct {
	broadcaster *Broadcaster
}

// 全局触发器实例
var eventTrigger *EventTrigger

// GetEventTrigger 获取事件触发器
func GetEventTrigger() *EventTrigger {
	if eventTrigger == nil {
		eventTrigger = &EventTrigger{
			broadcaster: GetBroadcaster(),
		}
	}
	return eventTrigger
}

// ==================== 签到事件 ====================

// TriggerCheckInStats 触发签到统计更新
func (t *EventTrigger) TriggerCheckInStats(activityId uint, total, pending, approved, rejected int) {
	payload := CheckInStatsPayload{
		Total:    total,
		Pending:  pending,
		Approved: approved,
		Rejected: rejected,
	}
	t.broadcaster.BroadcastCheckInStats(activityId, payload)
}

// ==================== 弹幕事件 ====================

// TriggerDanmaku 触发弹幕事件
// 在发送弹幕成功后调用
func (t *EventTrigger) TriggerDanmaku(activityId uint, danmakuId uint, user UserBrief, content, color string, createdAt time.Time) {
	payload := DanmakuPayload{
		ID:        danmakuId,
		User:      user,
		Content:   content,
		Color:     color,
		IsTop:     0,
		CreatedAt: createdAt,
	}
	t.broadcaster.BroadcastDanmaku(activityId, payload)
}

// TriggerTopDanmaku 触发置顶弹幕事件
func (t *EventTrigger) TriggerTopDanmaku(activityId uint, danmakuId uint, user UserBrief, content, color string) {
	payload := TopDanmakuPayload{
		ID:      danmakuId,
		User:    user,
		Content: content,
		Color:   color,
	}
	t.broadcaster.BroadcastTopDanmaku(activityId, payload)
}

// ==================== 摇一摇事件 ====================

// TriggerRankingUpdate 触发排名更新
// 在用户提交分数后调用
func (t *EventTrigger) TriggerRankingUpdate(activityId, roundId uint, ranking []RankingItem) {
	payload := RankingUpdatePayload{
		RoundId: roundId,
		Ranking: ranking,
	}
	t.broadcaster.BroadcastRankingUpdate(activityId, payload)
}

// ==================== 用户通知 ====================

// NotifyUser 发送消息给指定用户
func (t *EventTrigger) NotifyUser(userId uint, msgType string, payload interface{}) {
	t.broadcaster.SendToUser(userId, msgType, payload)
}

// NotifyWin 通知用户中奖
func (t *EventTrigger) NotifyWin(userId uint, winner WinnerInfo) {
	t.broadcaster.NotifyUserWin(userId, winner)
}
