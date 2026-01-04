package dto

// CheckInStatsResp 签到统计响应
type CheckInStatsResp struct {
	CheckedCount int     `json:"checkedCount"`
	TotalCount   int     `json:"totalCount"`
	CheckRate    float64 `json:"checkRate"`
}

// RoundWithPrize 场次带奖品信息
type RoundWithPrize struct {
	ID            uint       `json:"id"`
	RoundName     string     `json:"roundName"`
	Duration      int        `json:"duration"`
	WinnerCount   int        `json:"winnerCount"`
	Status        int        `json:"status"`
	ActualWinners int        `json:"actualWinners"`
	Prize         *PrizeInfo `json:"prize"`
}

type PrizeInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Level int    `json:"level"`
}
