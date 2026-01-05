package ws

import (
	"fmt"
	"time"
)

// Broadcaster 广播服务
type Broadcaster struct {
	hub *Hub
}

// 全局广播服务实例
var broadcaster *Broadcaster

// GetBroadcaster 获取广播服务实例
func GetBroadcaster() *Broadcaster {
	if broadcaster == nil {
		broadcaster = &Broadcaster{
			hub: GetHub(),
		}
	}
	return broadcaster
}

// ==================== 签到广播 ====================

// BroadcastCheckInStats 广播签到统计 只广播到主持人端
func (b *Broadcaster) BroadcastCheckInStats(activityId uint, payload CheckInStatsPayload) {
	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeCheckInStats, payload)
}

// ==================== 弹幕广播 ====================

// BroadcastDanmaku 广播新弹幕
func (b *Broadcaster) BroadcastDanmaku(activityId uint, payload DanmakuPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeDanmaku, activityId)
	b.hub.BroadcastToRoom(roomID, TypeNewDanmaku, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeNewDanmaku, payload)
}

// BroadcastTopDanmaku 广播置顶弹幕
func (b *Broadcaster) BroadcastTopDanmaku(activityId uint, payload TopDanmakuPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeDanmaku, activityId)
	b.hub.BroadcastToRoom(roomID, TypeTopDanmaku, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeTopDanmaku, payload)
}

// ==================== 摇一摇广播 ====================

// BroadcastGameReady 广播游戏准备
func (b *Broadcaster) BroadcastGameReady(activityId uint, payload GameReadyPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeShake, activityId)
	b.hub.BroadcastToRoom(roomID, TypeGameReady, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeGameReady, payload)
}

// BroadcastRoundStart 广播场次开始
func (b *Broadcaster) BroadcastRoundStart(activityId uint, payload RoundStartPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeShake, activityId)
	b.hub.BroadcastToRoom(roomID, TypeRoundStart, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeRoundStart, payload)
}

// BroadcastRoundEnd 广播场次结束
func (b *Broadcaster) BroadcastRoundEnd(activityId uint, payload RoundEndPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeShake, activityId)
	b.hub.BroadcastToRoom(roomID, TypeRoundEnd, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeRoundEnd, payload)
}

// BroadcastRankingUpdate 广播排名更新
func (b *Broadcaster) BroadcastRankingUpdate(activityId uint, payload RankingUpdatePayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeShake, activityId)
	b.hub.BroadcastToRoom(roomID, TypeRankingUpdate, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeRankingUpdate, payload)
}

// BroadcastCountdown 广播倒计时
func (b *Broadcaster) BroadcastCountdown(activityId uint, payload CountdownPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeShake, activityId)
	b.hub.BroadcastToRoom(roomID, TypeCountdown, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeCountdown, payload)
}

// ==================== 抽奖广播 ====================

// BroadcastDrawStart 广播抽奖开始
func (b *Broadcaster) BroadcastDrawStart(activityId uint, payload DrawStartPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeDraw, activityId)
	b.hub.BroadcastToRoom(roomID, TypeDrawStart, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeDrawStart, payload)
}

// BroadcastDrawResult 广播抽奖结果
func (b *Broadcaster) BroadcastDrawResult(activityId uint, payload DrawResultPayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeDraw, activityId)
	b.hub.BroadcastToRoom(roomID, TypeDrawResult, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeDrawResult, payload)
}

// BroadcastRollingUpdate 广播滚动更新
func (b *Broadcaster) BroadcastRollingUpdate(activityId uint, payload RollingUpdatePayload) {
	roomID := fmt.Sprintf("%s:%d", RoomTypeDraw, activityId)
	b.hub.BroadcastToRoom(roomID, TypeRollingUpdate, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeRollingUpdate, payload)
}

// BroadcastDrawReset 广播抽奖重置
func (b *Broadcaster) BroadcastDrawReset(activityId uint) {
	payload := DrawResetPayload{Message: "抽奖已重置"}

	roomID := fmt.Sprintf("%s:%d", RoomTypeDraw, activityId)
	b.hub.BroadcastToRoom(roomID, TypeDrawReset, payload)

	screenRoomID := fmt.Sprintf("%s:%d", RoomTypeScreen, activityId)
	b.hub.BroadcastToRoom(screenRoomID, TypeDrawReset, payload)
}

// ==================== 用户消息 ====================

// SendToUser 发送消息给指定用户
func (b *Broadcaster) SendToUser(userID uint, msgType string, payload interface{}) {
	b.hub.SendToUser(userID, msgType, payload)
}

// NotifyUserWin 通知用户中奖
func (b *Broadcaster) NotifyUserWin(userID uint, winInfo WinnerInfo) {
	b.hub.SendToUser(userID, "win_notify", winInfo)
}

// ==================== 辅助方法 ====================

// GetOnlineCount 获取在线人数
func (b *Broadcaster) GetOnlineCount() int {
	return b.hub.GetOnlineCount()
}

// GetRoomCount 获取房间人数
func (b *Broadcaster) GetRoomCount(roomType string, activityId uint) int {
	roomID := fmt.Sprintf("%s:%d", roomType, activityId)
	return b.hub.GetRoomClientCount(roomID)
}

// IsUserOnline 检查用户是否在线
func (b *Broadcaster) IsUserOnline(userID uint) bool {
	return b.hub.IsUserOnline(userID)
}

// ==================== 便捷构造方法 ====================

// NewUserBrief 从数据创建用户简要信息
func NewUserBrief(id uint, nickname, avatar, realName, department string) UserBrief {
	return UserBrief{
		ID:         id,
		Nickname:   nickname,
		Avatar:     avatar,
		RealName:   realName,
		Department: department,
	}
}

// NewPrizeBrief 从数据创建奖品简要信息
func NewPrizeBrief(id uint, name, image string, level int) PrizeBrief {
	return PrizeBrief{
		ID:    id,
		Name:  name,
		Image: image,
		Level: level,
	}
}

// NewDanmakuPayload 创建弹幕消息
func NewDanmakuPayload(id uint, user UserBrief, content, color string, isTop int, createdAt time.Time) DanmakuPayload {
	return DanmakuPayload{
		ID:        id,
		User:      user,
		Content:   content,
		Color:     color,
		IsTop:     isTop,
		CreatedAt: createdAt,
	}
}

// NewRankingItem 创建排名项
func NewRankingItem(rank int, userId uint, user UserBrief, score int, isWinner bool) RankingItem {
	return RankingItem{
		Rank:     rank,
		UserId:   userId,
		User:     user,
		Score:    score,
		IsWinner: isWinner,
	}
}

// NewWinnerInfo 创建中奖信息
func NewWinnerInfo(id uint, user UserBrief, prize PrizeBrief, winType int, createdAt time.Time) WinnerInfo {
	return WinnerInfo{
		ID:        id,
		User:      user,
		Prize:     prize,
		WinType:   winType,
		CreatedAt: createdAt,
	}
}
