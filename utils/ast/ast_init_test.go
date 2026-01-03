package ast

import (
	"lotteryBackend/global"
	"path/filepath"
)

func init() {
	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("../../../")
	global.GVA_CONFIG.AutoCode.Server = "server"
}
