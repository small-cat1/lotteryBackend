package ws

import (
	"fmt"
	"lotteryBackend/global"
	"sync"
	"time"

	"go.uber.org/zap"
)

// GameController 游戏控制器
type GameController struct {
	broadcaster *Broadcaster
	// 正在进行的游戏 activityId -> *GameSession
	activeSessions map[uint]*GameSession
	sessionsMutex  sync.RWMutex
}

// GameSession 游戏会话
type GameSession struct {
	ActivityId  uint
	RoundId     uint
	Round       RoundInfo
	Status      string // waiting, playing, finished
	StartTime   time.Time
	Duration    int
	RemainTime  int
	StopChan    chan struct{}
	CountdownWg sync.WaitGroup
}

// 游戏状态常量
const (
	GameStatusWaiting  = "waiting"
	GameStatusPlaying  = "playing"
	GameStatusFinished = "finished"
)

// 全局游戏控制器
var (
	gameController *GameController
	gameOnce       sync.Once
)

// GetGameController 获取游戏控制器实例
func GetGameController() *GameController {
	gameOnce.Do(func() {
		gameController = &GameController{
			broadcaster:    GetBroadcaster(),
			activeSessions: make(map[uint]*GameSession),
		}
	})
	return gameController
}

// ==================== 摇一摇游戏控制 ====================

// PrepareRound 准备场次（发送准备消息）
func (gc *GameController) PrepareRound(activityId uint, round RoundInfo, countdown int) error {
	gc.sessionsMutex.Lock()
	defer gc.sessionsMutex.Unlock()

	// 检查是否已有进行中的游戏
	if session, ok := gc.activeSessions[activityId]; ok {
		if session.Status == GameStatusPlaying {
			return fmt.Errorf("当前已有进行中的游戏")
		}
	}

	// 创建会话
	session := &GameSession{
		ActivityId: activityId,
		RoundId:    round.ID,
		Round:      round,
		Status:     GameStatusWaiting,
		Duration:   round.Duration,
		StopChan:   make(chan struct{}),
	}
	gc.activeSessions[activityId] = session

	// 广播准备消息
	gc.broadcaster.BroadcastGameReady(activityId, GameReadyPayload{
		Round:     round,
		Countdown: countdown,
	})

	global.GVA_LOG.Info("游戏准备",
		zap.Uint("activityId", activityId),
		zap.Uint("roundId", round.ID),
		zap.Int("countdown", countdown))

	return nil
}

// StartRound 开始场次
func (gc *GameController) StartRound(activityId uint, round RoundInfo) error {
	gc.sessionsMutex.Lock()

	session, ok := gc.activeSessions[activityId]
	if !ok {
		// 如果没有准备阶段，直接创建会话
		session = &GameSession{
			ActivityId: activityId,
			RoundId:    round.ID,
			Round:      round,
			Duration:   round.Duration,
			StopChan:   make(chan struct{}),
		}
		gc.activeSessions[activityId] = session
	}

	session.Status = GameStatusPlaying
	session.StartTime = time.Now()
	session.RemainTime = round.Duration

	gc.sessionsMutex.Unlock()

	// 广播开始消息
	gc.broadcaster.BroadcastRoundStart(activityId, RoundStartPayload{
		Round: round,
	})

	// 启动倒计时
	go gc.runCountdown(activityId, session)

	global.GVA_LOG.Info("游戏开始",
		zap.Uint("activityId", activityId),
		zap.Uint("roundId", round.ID),
		zap.Int("duration", round.Duration))

	return nil
}

// StopRound 停止场次
func (gc *GameController) StopRound(activityId uint, ranking []RankingItem, winners []WinnerInfo) error {
	gc.sessionsMutex.Lock()
	session, ok := gc.activeSessions[activityId]
	if !ok {
		gc.sessionsMutex.Unlock()
		return fmt.Errorf("没有进行中的游戏")
	}

	session.Status = GameStatusFinished
	close(session.StopChan)

	gc.sessionsMutex.Unlock()

	// 广播结束消息
	gc.broadcaster.BroadcastRoundEnd(activityId, RoundEndPayload{
		RoundId: session.RoundId,
		Ranking: ranking,
		Winners: winners,
	})

	// 通知中奖用户
	for _, winner := range winners {
		gc.broadcaster.NotifyUserWin(winner.User.ID, winner)
	}

	global.GVA_LOG.Info("游戏结束",
		zap.Uint("activityId", activityId),
		zap.Uint("roundId", session.RoundId),
		zap.Int("winnersCount", len(winners)))

	// 清理会话
	gc.sessionsMutex.Lock()
	delete(gc.activeSessions, activityId)
	gc.sessionsMutex.Unlock()

	return nil
}

// UpdateRanking 更新排名
func (gc *GameController) UpdateRanking(activityId uint, roundId uint, ranking []RankingItem) {
	gc.broadcaster.BroadcastRankingUpdate(activityId, RankingUpdatePayload{
		RoundId: roundId,
		Ranking: ranking,
	})
}

// runCountdown 运行倒计时
func (gc *GameController) runCountdown(activityId uint, session *GameSession) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-session.StopChan:
			return

		case <-ticker.C:
			gc.sessionsMutex.Lock()
			session.RemainTime--
			remainTime := session.RemainTime
			gc.sessionsMutex.Unlock()

			// 广播倒计时
			gc.broadcaster.BroadcastCountdown(activityId, CountdownPayload{
				RoundId:    session.RoundId,
				RemainTime: remainTime,
			})

			// 时间到，触发结束回调
			if remainTime <= 0 {
				// 这里应该触发结束逻辑，但实际结束应该由外部服务调用StopRound
				// 以便获取最终排名和生成中奖记录
				global.GVA_LOG.Info("倒计时结束，等待结算",
					zap.Uint("activityId", activityId),
					zap.Uint("roundId", session.RoundId))
				return
			}
		}
	}
}

// GetActiveSession 获取活跃会话
func (gc *GameController) GetActiveSession(activityId uint) *GameSession {
	gc.sessionsMutex.RLock()
	defer gc.sessionsMutex.RUnlock()
	return gc.activeSessions[activityId]
}

// IsRoundActive 检查场次是否活跃
func (gc *GameController) IsRoundActive(activityId uint) bool {
	gc.sessionsMutex.RLock()
	defer gc.sessionsMutex.RUnlock()
	session, ok := gc.activeSessions[activityId]
	return ok && session.Status == GameStatusPlaying
}

// ==================== 抽奖控制 ====================

// StartDraw 开始抽奖
func (gc *GameController) StartDraw(activityId uint, prize PrizeBrief, drawCount int, candidates []UserBrief) {
	gc.broadcaster.BroadcastDrawStart(activityId, DrawStartPayload{
		Prize:      prize,
		DrawCount:  drawCount,
		Candidates: candidates,
	})

	global.GVA_LOG.Info("抽奖开始",
		zap.Uint("activityId", activityId),
		zap.String("prizeName", prize.Name),
		zap.Int("drawCount", drawCount))
}

// AnnounceDrawResult 公布抽奖结果
func (gc *GameController) AnnounceDrawResult(activityId uint, prize PrizeBrief, winners []WinnerInfo) {
	gc.broadcaster.BroadcastDrawResult(activityId, DrawResultPayload{
		Prize:   prize,
		Winners: winners,
	})

	// 通知中奖用户
	for _, winner := range winners {
		gc.broadcaster.NotifyUserWin(winner.User.ID, winner)
	}

	global.GVA_LOG.Info("抽奖结果公布",
		zap.Uint("activityId", activityId),
		zap.Int("winnersCount", len(winners)))
}

// UpdateRolling 更新滚动动画
func (gc *GameController) UpdateRolling(activityId uint, users []UserBrief) {
	gc.broadcaster.BroadcastRollingUpdate(activityId, RollingUpdatePayload{
		Users: users,
	})
}

// ResetDraw 重置抽奖
func (gc *GameController) ResetDraw(activityId uint) {
	gc.broadcaster.BroadcastDrawReset(activityId)
}
