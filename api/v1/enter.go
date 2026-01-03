package v1

import (
	"lotteryBackend/api/v1/example"
	"lotteryBackend/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}
