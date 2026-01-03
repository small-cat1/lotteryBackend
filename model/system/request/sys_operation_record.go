package request

import (
	"lotteryBackend/model/common/request"
	"lotteryBackend/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
