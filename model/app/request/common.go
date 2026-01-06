package request

// ========== 通用 ==========

type H5PageReq struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type H5LimitReq struct {
	Limit int `form:"limit" json:"limit"`
}

// 定义请求结构体（可以放到 request 包里）
type DanmakuListReq struct {
	ActivityId uint `form:"activityId" binding:"required"`
	Page       int  `form:"page"`
	PageSize   int  `form:"pageSize"`
}
