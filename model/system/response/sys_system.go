package response

import "lotteryBackend/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
