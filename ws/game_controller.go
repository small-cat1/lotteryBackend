package ws

import (
	"sync"
	"time"
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

// UpdateRanking 更新排名
func (gc *GameController) UpdateRanking(activityId uint, roundId uint, ranking []RankingItem) {
	gc.broadcaster.BroadcastRankingUpdate(activityId, RankingUpdatePayload{
		RoundId: roundId,
		Ranking: ranking,
	})
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
