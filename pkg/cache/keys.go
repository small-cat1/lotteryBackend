package cache

import (
	"fmt"
	"time"
)

// ==================== Redis Key 定义 ====================

const (
	// 活动相关
	KeyActivityInfo  = "annual:activity:%d:info"          // 活动信息缓存
	KeyCheckInSwitch = "annual:activity:%d:checkin"       // 签到开关 0/1
	KeyDanmakuSwitch = "annual:activity:%d:danmaku"       // 弹幕开关 0/1
	KeyCheckInUsers  = "annual:activity:%d:checkin:users" // 已签到用户ID集合 (Set)
	KeyCurrentRound  = "annual:activity:%d:current_round" // 当前进行中的场次ID

	// 用户相关
	KeyUserInfo    = "annual:user:%d"                // 用户信息缓存
	KeyUserCheckIn = "annual:user:%d:checkin:%d"     // 用户签到信息 user:{userId}:checkin:{activityId}

	// 场次相关
	KeyRoundStatus = "annual:round:%d:status" // 场次状态 0/1/2
	KeyRoundInfo   = "annual:round:%d:info"   // 场次信息缓存
	KeyRoundScores = "annual:round:%d:scores" // 场次积分排行榜 (ZSet)
)

// ==================== 过期时间 ====================

const (
	ActivityCacheExpire = 5 * time.Minute  // 活动缓存 5分钟
	SwitchExpire        = 24 * time.Hour   // 开关状态 24小时
	CheckInUsersExpire  = 24 * time.Hour   // 签到用户集合 24小时
	UserInfoExpire      = 30 * time.Minute // 用户信息缓存 30分钟
	RoundInfoExpire     = 2 * time.Hour    // 场次信息缓存 2小时
	RoundScoresExpire   = 2 * time.Hour    // 场次积分 2小时
)

// ==================== Key 生成方法 ====================

// --- 活动相关 ---

func GetActivityInfoKey(activityId uint) string {
	return fmt.Sprintf(KeyActivityInfo, activityId)
}

func GetCheckInSwitchKey(activityId uint) string {
	return fmt.Sprintf(KeyCheckInSwitch, activityId)
}

func GetDanmakuSwitchKey(activityId uint) string {
	return fmt.Sprintf(KeyDanmakuSwitch, activityId)
}

func GetCheckInUsersKey(activityId uint) string {
	return fmt.Sprintf(KeyCheckInUsers, activityId)
}

func GetCurrentRoundKey(activityId uint) string {
	return fmt.Sprintf(KeyCurrentRound, activityId)
}

// --- 用户相关 ---

func GetUserInfoKey(userId uint) string {
	return fmt.Sprintf(KeyUserInfo, userId)
}

func GetUserCheckInKey(userId, activityId uint) string {
	return fmt.Sprintf(KeyUserCheckIn, userId, activityId)
}

// --- 场次相关 ---

func GetRoundStatusKey(roundId uint) string {
	return fmt.Sprintf(KeyRoundStatus, roundId)
}

func GetRoundInfoKey(roundId uint) string {
	return fmt.Sprintf(KeyRoundInfo, roundId)
}

func GetRoundScoresKey(roundId uint) string {
	return fmt.Sprintf(KeyRoundScores, roundId)
}
