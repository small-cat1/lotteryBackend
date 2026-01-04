package app

import (
	"errors"
	"lotteryBackend/constants"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/service/app/dto"
	"lotteryBackend/ws"
	"math/rand"
	"time"
)

type ConsoleService struct{}

// ==================== 大屏数据 ====================

// GetCheckInStats 获取签到统计
func (s *ConsoleService) GetCheckInStats(activityId uint) (*dto.CheckInStatsResp, error) {
	// 已签到人数
	var checkedCount int64
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ?", activityId).
		Count(&checkedCount)

	// 总人数（已报名且通过审核）
	var totalCount int64
	global.GVA_DB.Model(&annual.AnnualUser{}).
		Where("status = ?", 1).
		Count(&totalCount)

	var checkRate float64 = 0
	if totalCount > 0 {
		checkRate = float64(checkedCount) / float64(totalCount) * 100
	}

	return &dto.CheckInStatsResp{
		CheckedCount: int(checkedCount),
		TotalCount:   int(totalCount),
		CheckRate:    checkRate,
	}, nil
}

// GetRoundList 获取场次列表
func (s *ConsoleService) GetRoundList(activityId uint) ([]dto.RoundWithPrize, error) {
	var rounds []annual.AnnualShakeRound
	global.GVA_DB.Where("activity_id = ?", activityId).
		Order("sort ASC, id ASC").
		Find(&rounds)

	result := make([]dto.RoundWithPrize, 0, len(rounds))
	for _, r := range rounds {
		item := dto.RoundWithPrize{
			ID:          r.ID,
			RoundName:   r.RoundName,
			Duration:    r.Duration,
			WinnerCount: r.WinnerCount,
			Status:      *r.Status,
		}

		// 获取奖品信息
		if r.PrizeId > 0 {
			var prize annual.AnnualPrize
			if err := global.GVA_DB.First(&prize, r.PrizeId).Error; err == nil {
				item.Prize = &dto.PrizeInfo{
					ID:    prize.ID,
					Name:  prize.Name,
					Image: prize.Image,
					Level: *prize.Level,
				}
			}
		}

		// 统计实际中奖人数
		var winnerCount int64
		global.GVA_DB.Model(&annual.AnnualWinner{}).
			Where("round_id = ?", r.ID).
			Count(&winnerCount)
		item.ActualWinners = int(winnerCount)

		result = append(result, item)
	}

	return result, nil
}

// DanmakuListResp 弹幕列表响应
type DanmakuListResp struct {
	List      []DanmakuItem `json:"list"`
	UserCount int           `json:"userCount"`
}

type DanmakuItem struct {
	ID        uint       `json:"id"`
	Content   string     `json:"content"`
	Color     string     `json:"color"`
	IsTop     int        `json:"isTop"`
	CreatedAt time.Time  `json:"createdAt"`
	User      *UserBrief `json:"user"`
}

type UserBrief struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RealName   string `json:"realName"`
	Department string `json:"department"`
}

// GetDanmakuList 获取弹幕列表
func (s *ConsoleService) GetDanmakuList(activityId uint, limit int) (*DanmakuListResp, error) {
	var danmakus []annual.AnnualDanmaku
	global.GVA_DB.Where("activity_id = ? AND status = ?", activityId, 1).
		Order("created_at DESC").
		Limit(limit).
		Find(&danmakus)

	// 获取弹幕用户数（去重）
	var userCount int64
	global.GVA_DB.Model(&annual.AnnualDanmaku{}).
		Where("activity_id = ? AND status = ?", activityId, 1).
		Distinct("user_id").
		Count(&userCount)

	list := make([]DanmakuItem, 0, len(danmakus))
	for _, d := range danmakus {
		item := DanmakuItem{
			ID:        d.ID,
			Content:   d.Content,
			Color:     d.Color,
			IsTop:     *d.IsTop,
			CreatedAt: d.CreatedAt,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, d.UserId).Error; err == nil {
			item.User = &UserBrief{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		list = append(list, item)
	}

	return &DanmakuListResp{
		List:      list,
		UserCount: int(userCount),
	}, nil
}

// RankingResp 排行榜响应
type RankingResp struct {
	Ranking     []RankingItem `json:"ranking"`
	PlayerCount int           `json:"playerCount"`
}

type RankingItem struct {
	Rank     int        `json:"rank"`
	UserId   uint       `json:"userId"`
	Score    int        `json:"score"`
	IsWinner bool       `json:"isWinner"`
	User     *UserBrief `json:"user"`
}

// GetRanking 获取排行榜
func (s *ConsoleService) GetRanking(roundId uint, limit int) (*RankingResp, error) {
	// 获取场次信息
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	// 获取成绩
	var scores []annual.AnnualShakeScore
	global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Limit(limit).
		Find(&scores)

	// 总参与人数
	var playerCount int64
	global.GVA_DB.Model(&annual.AnnualShakeScore{}).
		Where("round_id = ?", roundId).
		Count(&playerCount)

	ranking := make([]RankingItem, 0, len(scores))
	for i, s := range scores {
		item := RankingItem{
			Rank:     i + 1,
			UserId:   s.UserId,
			Score:    s.Score,
			IsWinner: i < round.WinnerCount,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, s.UserId).Error; err == nil {
			item.User = &UserBrief{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		ranking = append(ranking, item)
	}

	return &RankingResp{
		Ranking:     ranking,
		PlayerCount: int(playerCount),
	}, nil
}

// RandomDrawConfigResp 随机抽奖配置响应
type RandomDrawConfigResp struct {
	Enabled   bool   `json:"enabled"`
	PrizeId   uint   `json:"prizeId"`
	PrizeName string `json:"prizeName"`
}

// GetRandomDrawConfig 获取随机抽奖配置
func (s *ConsoleService) GetRandomDrawConfig() *RandomDrawConfigResp {
	result := &RandomDrawConfigResp{
		Enabled: false,
	}

	// 读取配置
	var enabledConfig annual.AnnualConfig
	if err := global.GVA_DB.Where("config_key = ?", "random_draw_enabled").First(&enabledConfig).Error; err == nil {
		result.Enabled = enabledConfig.ConfigValue == "1"
	}

	var prizeConfig annual.AnnualConfig
	if err := global.GVA_DB.Where("config_key = ?", "random_draw_prize_id").First(&prizeConfig).Error; err == nil {
		if prizeConfig.ConfigValue != "" && prizeConfig.ConfigValue != "0" {
			// 获取奖品名称
			var prize annual.AnnualPrize
			if err := global.GVA_DB.First(&prize, prizeConfig.ConfigValue).Error; err == nil {
				result.PrizeId = prize.ID
				result.PrizeName = prize.Name
			}
		}
	}

	return result
}

// ==================== 游戏控制 ====================

// StartGame 开始游戏（验证密码）
func (s *ConsoleService) StartGame(roundId uint, password string) error {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return errors.New("场次不存在")
	}

	// 验证密码
	if round.Password != password {
		return errors.New("密码错误")
	}

	if *round.Status != 0 {
		return errors.New("场次状态不正确")
	}

	// 更新状态
	now := time.Now()
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":     1,
		"start_time": now,
	})

	// 获取奖品信息
	var prize annual.AnnualPrize
	global.GVA_DB.First(&prize, round.PrizeId)

	// 广播游戏开始
	controller := ws.GetGameController()
	roundInfo := ws.RoundInfo{
		ID:          round.ID,
		RoundName:   round.RoundName,
		Duration:    round.Duration,
		WinnerCount: round.WinnerCount,
		Status:      1,
		Prize: ws.PrizeBrief{
			ID:    prize.ID,
			Name:  prize.Name,
			Image: prize.Image,
			Level: *prize.Level,
		},
	}
	controller.StartRound(round.ActivityId, roundInfo)

	return nil
}

// GameResultResp 游戏结果响应
type GameResultResp struct {
	Ranking []RankingItem `json:"ranking"`
	Winners []WinnerItem  `json:"winners"`
}

type WinnerItem struct {
	Rank   int            `json:"rank"`
	UserId uint           `json:"userId"`
	Score  int            `json:"score"`
	User   *UserBrief     `json:"user"`
	Prize  *dto.PrizeInfo `json:"prize"`
}

// StopGame 立即停止游戏
func (s *ConsoleService) StopGame(roundId uint) (*GameResultResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	if *round.Status != 1 {
		return nil, errors.New("游戏未在进行中")
	}

	return s.SettleGame(roundId)
}

// SettleGame 结算游戏
func (s *ConsoleService) SettleGame(roundId uint) (*GameResultResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	// 获取所有成绩并排序
	var scores []annual.AnnualShakeScore
	global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Find(&scores)

	// 获取奖品信息
	var prize annual.AnnualPrize
	global.GVA_DB.First(&prize, round.PrizeId)

	prizeInfo := &dto.PrizeInfo{
		ID:    prize.ID,
		Name:  prize.Name,
		Image: prize.Image,
		Level: *prize.Level,
	}

	ranking := make([]RankingItem, 0, len(scores))
	winners := make([]WinnerItem, 0)
	wsWinners := make([]ws.WinnerInfo, 0)
	wsRanking := make([]ws.RankingItem, 0)

	for i, score := range scores {
		// 获取用户信息
		var user annual.AnnualUser
		global.GVA_DB.First(&user, score.UserId)

		userBrief := &UserBrief{
			ID:         user.ID,
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			RealName:   user.RealName,
			Department: user.Department,
		}

		isWinner := i < round.WinnerCount

		rankItem := RankingItem{
			Rank:     i + 1,
			UserId:   score.UserId,
			Score:    score.Score,
			IsWinner: isWinner,
			User:     userBrief,
		}
		ranking = append(ranking, rankItem)

		// 更新成绩记录
		global.GVA_DB.Model(&score).Updates(map[string]interface{}{
			"rank":      i + 1,
			"is_winner": boolToInt(isWinner),
		})

		// 中奖处理
		if isWinner {
			winType := constants.WinTypeShake
			winner := annual.AnnualWinner{
				ActivityId: round.ActivityId,
				UserId:     score.UserId,
				PrizeId:    round.PrizeId,
				RoundId:    roundId,
				WinType:    &winType, // 摇一摇
			}
			global.GVA_DB.Create(&winner)

			winners = append(winners, WinnerItem{
				Rank:   i + 1,
				UserId: score.UserId,
				Score:  score.Score,
				User:   userBrief,
				Prize:  prizeInfo,
			})

			wsWinners = append(wsWinners, ws.WinnerInfo{
				ID:        winner.ID,
				User:      ws.UserBrief{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar, RealName: user.RealName, Department: user.Department},
				Prize:     ws.PrizeBrief{ID: prize.ID, Name: prize.Name, Image: prize.Image, Level: *prize.Level},
				WinType:   1,
				CreatedAt: winner.CreatedAt,
			})
		}

		wsRanking = append(wsRanking, ws.RankingItem{
			Rank:     i + 1,
			UserId:   score.UserId,
			User:     ws.UserBrief{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar, RealName: user.RealName, Department: user.Department},
			Score:    score.Score,
			IsWinner: isWinner,
		})
	}

	// 更新场次状态
	now := time.Now()
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":   2,
		"end_time": now,
	})

	// 更新奖品剩余数量
	if len(winners) > 0 {
		global.GVA_DB.Model(&prize).Update("remain_count", prize.RemainCount-len(winners))
	}

	// 广播游戏结束
	controller := ws.GetGameController()
	controller.StopRound(round.ActivityId, wsRanking, wsWinners)

	return &GameResultResp{
		Ranking: ranking,
		Winners: winners,
	}, nil
}

// ==================== 随机抽奖 ====================

// RandomDraw 随机抽奖（从弹幕用户中抽取）
func (s *ConsoleService) RandomDraw(activityId uint, count int) ([]WinnerItem, error) {
	// 检查配置
	config := s.GetRandomDrawConfig()
	if !config.Enabled {
		return nil, errors.New("随机抽奖未开启")
	}

	if config.PrizeId == 0 {
		return nil, errors.New("未配置抽奖奖品")
	}

	// 获取弹幕用户（去重，排除已中奖）
	var danmakuUserIds []uint
	global.GVA_DB.Model(&annual.AnnualDanmaku{}).
		Where("activity_id = ? AND status = ?", activityId, 1).
		Distinct().
		Pluck("user_id", &danmakuUserIds)

	if len(danmakuUserIds) == 0 {
		return nil, errors.New("暂无弹幕用户")
	}

	// 排除已中奖用户
	var wonUserIds []uint
	global.GVA_DB.Model(&annual.AnnualWinner{}).
		Where("activity_id = ?", activityId).
		Pluck("user_id", &wonUserIds)

	wonMap := make(map[uint]bool)
	for _, id := range wonUserIds {
		wonMap[id] = true
	}

	candidates := make([]uint, 0)
	for _, id := range danmakuUserIds {
		if !wonMap[id] {
			candidates = append(candidates, id)
		}
	}

	if len(candidates) < count {
		return nil, errors.New("符合条件的用户不足")
	}

	// 随机抽取
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	selectedIds := candidates[:count]

	// 获取奖品
	var prize annual.AnnualPrize
	global.GVA_DB.First(&prize, config.PrizeId)

	prizeInfo := &dto.PrizeInfo{
		ID:    prize.ID,
		Name:  prize.Name,
		Image: prize.Image,
		Level: *prize.Level,
	}

	// 创建中奖记录
	winners := make([]WinnerItem, 0, count)
	wsWinners := make([]ws.WinnerInfo, 0, count)

	for _, userId := range selectedIds {
		var user annual.AnnualUser
		global.GVA_DB.First(&user, userId)
		winType := constants.WinTypeRandom
		winner := annual.AnnualWinner{
			ActivityId: activityId,
			UserId:     userId,
			PrizeId:    config.PrizeId,
			WinType:    &winType, // 随机抽奖
		}
		global.GVA_DB.Create(&winner)

		userBrief := &UserBrief{
			ID:         user.ID,
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			RealName:   user.RealName,
			Department: user.Department,
		}

		winners = append(winners, WinnerItem{
			UserId: userId,
			User:   userBrief,
			Prize:  prizeInfo,
		})

		wsWinners = append(wsWinners, ws.WinnerInfo{
			ID:        winner.ID,
			User:      ws.UserBrief{ID: user.ID, Nickname: user.Nickname, Avatar: user.Avatar, RealName: user.RealName, Department: user.Department},
			Prize:     ws.PrizeBrief{ID: prize.ID, Name: prize.Name, Image: prize.Image, Level: *prize.Level},
			WinType:   2,
			CreatedAt: winner.CreatedAt,
		})
	}

	// 更新奖品剩余数量
	global.GVA_DB.Model(&prize).Update("remain_count", prize.RemainCount-count)

	// 广播抽奖结果
	controller := ws.GetGameController()
	controller.AnnounceDrawResult(activityId, ws.PrizeBrief{ID: prize.ID, Name: prize.Name, Image: prize.Image, Level: *prize.Level}, wsWinners)

	return winners, nil
}

// GetWinners 获取中奖名单
func (s *ConsoleService) GetWinners(activityId uint, winType int) ([]WinnerItem, error) {
	query := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("activity_id = ?", activityId)
	if winType > 0 {
		query = query.Where("win_type = ?", winType)
	}

	var winners []annual.AnnualWinner
	query.Order("created_at DESC").Find(&winners)

	result := make([]WinnerItem, 0, len(winners))
	for _, w := range winners {
		var user annual.AnnualUser
		global.GVA_DB.First(&user, w.UserId)

		var prize annual.AnnualPrize
		global.GVA_DB.First(&prize, w.PrizeId)

		result = append(result, WinnerItem{
			UserId: w.UserId,
			User: &UserBrief{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			},
			Prize: &dto.PrizeInfo{
				ID:    prize.ID,
				Name:  prize.Name,
				Image: prize.Image,
				Level: *prize.Level,
			},
		})
	}

	return result, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
