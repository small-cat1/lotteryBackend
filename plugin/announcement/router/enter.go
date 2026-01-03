package router

import "lotteryBackend/plugin/announcement/api"

var (
	Router  = new(router)
	apiInfo = api.Api.Info
)

type router struct{ Info info }
