package initialize

import (
	_ "lotteryBackend/source/example"
	_ "lotteryBackend/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
