package console

import "lotteryBackend/service"

type ApiGroup struct {
	ConsoleApi
}

var (
	activityService = service.ServiceGroupApp.ConsoleServiceGroup.ActivityService
	checkInService  = service.ServiceGroupApp.ConsoleServiceGroup.CheckInService
	danmakuService  = service.ServiceGroupApp.ConsoleServiceGroup.DanmakuService
	gameService     = service.ServiceGroupApp.ConsoleServiceGroup.GameService
	drawService     = service.ServiceGroupApp.ConsoleServiceGroup.DrawService
)
