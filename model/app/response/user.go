package response

import "time"

// H5UserResp 用户信息响应
type H5UserResp struct {
	ID       uint   `json:"ID"`
	OpenId   string `json:"openId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	// 签到信息（如果传了活动ID）
	CheckIn *H5CheckInInfo `json:"checkIn"`
}

// H5CheckInInfo 签到信息
type H5CheckInInfo struct {
	IsCheckedIn  bool      `json:"isCheckedIn"`  // 是否已签到
	RealName     string    `json:"realName"`     // 真实姓名
	Phone        string    `json:"phone"`        // 手机号
	Department   string    `json:"department"`   // 部门
	EmployeeNo   string    `json:"employeeNo"`   // 工号
	Status       int       `json:"status"`       // 审核状态 0待审核 1通过 2拒绝
	RejectReason string    `json:"rejectReason"` // 拒绝原因
	CheckInTime  time.Time `json:"checkInTime"`  // 签到时间
}

// H5UserBrief 用户简要信息（用于列表展示）

// CheckInRecordResp 签到记录响应
type CheckInRecordResp struct {
	ID          uint        `json:"id"`
	User        H5UserBrief `json:"user"`
	RealName    string      `json:"realName"`
	Department  string      `json:"department"`
	CheckInTime time.Time   `json:"checkInTime"`
}

// AuditStatusResp 审核状态响应
type AuditStatusResp struct {
	Status       int    `json:"status"`       // 0待审核 1已通过 2已拒绝
	RejectReason string `json:"rejectReason"` // 拒绝原因
}
