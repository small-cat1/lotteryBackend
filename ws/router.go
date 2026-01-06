package ws

import (
	"github.com/gin-gonic/gin"
)

// InitWebSocketRouter 初始化WebSocket路由
func InitWebSocketRouter(Router *gin.RouterGroup) {
	handler := NewHandler()

	wsGroup := Router.Group("ws")
	{
		// 通用WebSocket连接
		wsGroup.GET("", handler.HandleConnection)

		// H5前端连接（需要Token）
		wsGroup.GET("/h5", handler.HandleH5Connection)

		// 大屏连接（无需Token）
		wsGroup.GET("/screen", handler.HandleScreenConnection)

		// 按类型分的连接（可选）
		wsGroup.GET("/checkin", handler.HandleScreenConnection) // 签到大屏
		wsGroup.GET("/danmaku", handler.HandleScreenConnection) // 弹幕大屏
		wsGroup.GET("/shake", handler.HandleScreenConnection)   // 摇一摇大屏
	}
}
