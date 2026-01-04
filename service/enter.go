package service

import (
	"lotteryBackend/service/annual"
	"lotteryBackend/service/app"
	"lotteryBackend/service/console"
	"lotteryBackend/service/example"
	"lotteryBackend/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	AnnualServiceGroup  annual.ServiceGroup
	AppServiceGroup     app.H5ServiceGroup
	ConsoleServiceGroup console.ServiceGroup
}
