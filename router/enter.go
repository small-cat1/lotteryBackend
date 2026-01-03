package router

import (
	"lotteryBackend/router/annual"
	"lotteryBackend/router/example"
	"lotteryBackend/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Annual  annual.RouterGroup
}
