package request

import "lotteryBackend/model/common/request"

// CheckInSearch 签到搜索
type CheckInSearch struct {
	request.PageInfo
	ActivityId uint   `json:"activityId" form:"activityId"`
	UserId     uint   `json:"userId" form:"userId"`
	RealName   string `json:"realName" form:"realName"`
	Phone      string `json:"phone" form:"phone"`
	Department string `json:"department" form:"department"`
	Status     *int   `json:"status" form:"status"` // 使用指针，区分0和未传
}

// CheckInUpdate 更新签到信息
type CheckInUpdate struct {
	Id           uint   `json:"id" binding:"required"`
	RealName     string `json:"realName"`
	Phone        string `json:"phone"`
	Department   string `json:"department"`
	EmployeeNo   string `json:"employeeNo"`
	Status       int    `json:"status"`
	RejectReason string `json:"rejectReason"`
}

// CheckInStatusUpdate 更新签到状态（支持批量）
type CheckInStatusUpdate struct {
	Id           uint   `json:"id"`  // 单个审核
	Ids          []uint `json:"ids"` // 批量审核
	Status       int    `json:"status" binding:"required,oneof=1 2"`
	RejectReason string `json:"rejectReason"`
}

// CheckInDelete 删除签到
type CheckInDelete struct {
	Id  uint   `json:"id"`
	Ids []uint `json:"ids"`
}
