package response

// H5PageResult 分页结果
type H5PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// ========== 通用 ==========

// H5UserBrief 用户简要信息
type H5UserBrief struct {
	ID         uint   `json:"ID"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RealName   string `json:"realName"`
	Department string `json:"department"`
}
