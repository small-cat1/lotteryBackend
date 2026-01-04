package v1

import (
	"lotteryBackend/api/v1/annual"
	"lotteryBackend/api/v1/app"
	"lotteryBackend/api/v1/console"
	"lotteryBackend/api/v1/example"
	"lotteryBackend/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	AnnualApiGroup  annual.ApiGroup
	H5ApiGroup      app.ApiGroup
	ConsoleApi      console.ApiGroup
}
