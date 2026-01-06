package ws

// 初始化WebSocket模块
// 在应用启动时调用

import (
	"go.uber.org/zap"
	"lotteryBackend/global"
)

// Init 初始化WebSocket模块
func Init() {
	// 启动Hub
	hub := GetHub()
	global.GVA_LOG.Info("WebSocket Hub 已启动",
		zap.Int("online", hub.GetOnlineCount()))

	// 初始化广播器
	GetBroadcaster()

	// 初始化游戏控制器
	GetGameController()

	// 初始化事件触发器
	GetEventTrigger()

	GetShakeHandler() // 添加这行
	global.GVA_LOG.Info("WebSocket 模块初始化完成")
}

// Shutdown 关闭WebSocket模块
func Shutdown() {
	// 可以在这里添加清理逻辑
	global.GVA_LOG.Info("WebSocket 模块已关闭")
}
