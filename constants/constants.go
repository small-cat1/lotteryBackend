package constants

// 用户状态
const (
	UserStatusPending  = 0 // 待审核
	UserStatusApproved = 1 // 已通过
	UserStatusRejected = 2 // 已拒绝
)

// 用户报名状态
const (
	UserNotRegistered = 0 // 游客
	UserRegistered    = 1 // 已报名
)

// 活动状态
const (
	ActivityStatusNotStarted = 0 // 未开始
	ActivityStatusOngoing    = 1 // 进行中
	ActivityStatusEnded      = 2 // 已结束
)

// 弹幕状态
const (
	DanmakuStatusPending  = 0 // 待审核
	DanmakuStatusApproved = 1 // 已通过
	DanmakuStatusRejected = 2 // 已拒绝
)

// 摇一摇场次状态
const (
	ShakeRoundStatusNotStarted = 0 // 未开始
	ShakeRoundStatusOngoing    = 1 // 进行中
	ShakeRoundStatusEnded      = 2 // 已结束
)

// 奖品等级
const (
	PrizeLevelSpecial = 1 // 特等奖
	PrizeLevelFirst   = 2 // 一等奖
	PrizeLevelSecond  = 3 // 二等奖
	PrizeLevelThird   = 4 // 三等奖
	PrizeLevelPartic  = 5 // 参与奖
)

// 中奖方式
const (
	WinTypeShake   = 1 // 摇一摇
	WinTypeRandom  = 2 // 随机抽奖
	WinTypeDanmaku = 3 // 弹幕抽奖
)

// 领奖状态
const (
	WinnerStatusNotReceived = 0 // 未领取
	WinnerStatusReceived    = 1 // 已领取
)

//```
//
//---
//
//## 📂 文件结构建议
//```
//server/model/annual/
//├── annual_user.go           // 用户模型
//├── annual_activity.go       // 活动模型
//├── annual_check_in.go       // 签到模型
//├── annual_danmaku.go        // 弹幕模型
//├── annual_shake_round.go    // 摇一摇场次模型
//├── annual_shake_score.go    // 摇一摇成绩模型
//├── annual_prize.go          // 奖品模型
//├── annual_winner.go         // 中奖记录模型
//├── annual_config.go         // 配置模型
//└── constant.go              // 常量定义
//
//server/model/annual/request/
//└── annual.go                // 请求结构体
