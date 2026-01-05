package request

// ========== 通用 ==========

type H5PageReq struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type H5LimitReq struct {
	Limit int `form:"limit" json:"limit"`
}
