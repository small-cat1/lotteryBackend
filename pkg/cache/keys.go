package cache

import (
	"fmt"
	"time"
)

// ==================== Redis Key 定义 ====================

const (
	// 活动配置 - Hash
	KeyActivityInfo = "annual:activity:%d" // annual:activity:{activityId}

	// 状态开关 - String
	KeyCheckInSwitch = "annual:activity:%d:checkin" // 签到开关 0/1
	KeyDanmakuSwitch = "annual:activity:%d:danmaku" // 弹幕开关 0/1

	// 签到用户 - Set
	KeyCheckInUsers = "annual:activity:%d:checkin:users" // 已签到用户ID集合

	// 用户信息缓存 - Hash
	KeyUserInfo = "annual:user:%d" // annual:user:{userId}

	// 用户签到信息缓存 - Hash
	KeyUserCheckIn = "annual:user:%d:checkin:%d" // annual:user:{userId}:checkin:{activityId}

	// 游戏状态 - String
	KeyCurrentRound = "annual:activity:%d:current_round" // 当前进行中的场次ID
	KeyRoundStatus  = "annual:round:%d:status"           // 场次状态 0/1/2

	// 摇一摇积分 - ZSet
	KeyRoundScores = "annual:round:%d:scores" // member=userId, score=积分
)

// 过期时间
const (
	ActivityCacheExpire = 5 * time.Minute  // 活动缓存5分钟
	SwitchExpire        = 24 * time.Hour   // 开关状态24小时
	CheckInUsersExpire  = 24 * time.Hour   // 签到用户集合24小时
	RoundScoresExpire   = 2 * time.Hour    // 场次积分2小时
	UserInfoExpire      = 30 * time.Minute // 用户信息缓存30分钟
)

// ==================== Key 生成方法 ====================

func GetUserInfoKey(userId uint) string {
	return fmt.Sprintf(KeyUserInfo, userId)
}

func GetUserCheckInKey(userId, activityId uint) string {
	return fmt.Sprintf(KeyUserCheckIn, userId, activityId)
}

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

func GetRoundStatusKey(roundId uint) string {
	return fmt.Sprintf(KeyRoundStatus, roundId)
}

func GetRoundScoresKey(roundId uint) string {
	return fmt.Sprintf(KeyRoundScores, roundId)
}
