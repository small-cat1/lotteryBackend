package router

import (
	"lotteryBackend/router/annual"
	"lotteryBackend/router/app"
	"lotteryBackend/router/console"
	"lotteryBackend/router/example"
	"lotteryBackend/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Annual  annual.RouterGroup
	App     app.H5Router
	Console console.ConsoleRouter
}
