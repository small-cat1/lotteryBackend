package cache

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"lotteryBackend/global"
//	"time"
//
//	"github.com/redis/go-redis/v9"
//)
//
//// ==================== Redis Key 定义 ====================
//
//const (
//	// 活动配置 - Hash
//	KeyActivityInfo = "annual:activity:%d" // annual:activity:{activityId}
//
//	// 状态开关 - String
//	KeyCheckInSwitch = "annual:activity:%d:checkin" // 签到开关 0/1
//	KeyDanmakuSwitch = "annual:activity:%d:danmaku" // 弹幕开关 0/1
//
//	// 签到用户 - Set
//	KeyCheckInUsers = "annual:activity:%d:checkin:users" // 已签到用户ID集合
//	// 用户信息缓存 - Hash
//	KeyUserInfo = "annual:user:%d" // annual:user:{userId}
//	// 用户签到信息缓存 - Hash
//	KeyUserCheckIn = "annual:user:%d:checkin:%d" // annual:user:{userId}:checkin:{activityId}
//
//	// 游戏状态 - String
//	KeyCurrentRound = "annual:activity:%d:current_round" // 当前进行中的场次ID
//	KeyRoundStatus  = "annual:round:%d:status"           // 场次状态 0/1/2
//
//	// 摇一摇积分 - ZSet
//	KeyRoundScores = "annual:round:%d:scores" // member=userId, score=积分
//)
//
//// 过期时间
//const (
//	ActivityCacheExpire = 5 * time.Minute  // 活动缓存5分钟
//	SwitchExpire        = 24 * time.Hour   // 开关状态24小时
//	CheckInUsersExpire  = 24 * time.Hour   // 签到用户集合24小时
//	RoundScoresExpire   = 2 * time.Hour    // 场次积分2小时
//	UserInfoExpire      = 30 * time.Minute // 用户信息缓存30分钟
//)
//
//// UserInfoCache 用户信息缓存结构
//type UserInfoCache struct {
//	ID         uint   `json:"id"`
//	Nickname   string `json:"nickname"`
//	Avatar     string `json:"avatar"`
//	RealName   string `json:"realName"`
//	Department string `json:"department"`
//}
//
//// ==================== Key 生成方法 ====================
//func GetUserInfoKey(userId uint) string {
//	return fmt.Sprintf(KeyUserInfo, userId)
//}
//
//func GetUserCheckInKey(userId, activityId uint) string {
//	return fmt.Sprintf(KeyUserCheckIn, userId, activityId)
//}
//
//func GetActivityInfoKey(activityId uint) string {
//	return fmt.Sprintf(KeyActivityInfo, activityId)
//}
//
//func GetCheckInSwitchKey(activityId uint) string {
//	return fmt.Sprintf(KeyCheckInSwitch, activityId)
//}
//
//func GetDanmakuSwitchKey(activityId uint) string {
//	return fmt.Sprintf(KeyDanmakuSwitch, activityId)
//}
//
//func GetCheckInUsersKey(activityId uint) string {
//	return fmt.Sprintf(KeyCheckInUsers, activityId)
//}
//
//func GetCurrentRoundKey(activityId uint) string {
//	return fmt.Sprintf(KeyCurrentRound, activityId)
//}
//
//func GetRoundStatusKey(roundId uint) string {
//	return fmt.Sprintf(KeyRoundStatus, roundId)
//}
//
//func GetRoundScoresKey(roundId uint) string {
//	return fmt.Sprintf(KeyRoundScores, roundId)
//}
//
//// ==================== 用户缓存操作 ====================
//
//// SetUserInfoCache 缓存用户信息
//func SetUserInfoCache(userId uint, info *UserInfoCache) error {
//	ctx := context.Background()
//	key := GetUserInfoKey(userId)
//
//	data, err := json.Marshal(info)
//	if err != nil {
//		return err
//	}
//
//	return global.GVA_REDIS.Set(ctx, key, data, UserInfoExpire).Err()
//}
//
//// GetUserInfoCache 获取用户信息缓存
//func GetUserInfoCache(userId uint) (*UserInfoCache, error) {
//	ctx := context.Background()
//	key := GetUserInfoKey(userId)
//
//	data, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	var cache UserInfoCache
//	if err := json.Unmarshal([]byte(data), &cache); err != nil {
//		return nil, err
//	}
//	return &cache, nil
//}
//
//// BatchGetUserInfoCache 批量获取用户信息缓存
//func BatchGetUserInfoCache(userIds []uint) (map[uint]*UserInfoCache, []uint, error) {
//	ctx := context.Background()
//	result := make(map[uint]*UserInfoCache)
//	var missedIds []uint
//
//	if len(userIds) == 0 {
//		return result, missedIds, nil
//	}
//
//	// 使用 Pipeline 批量获取
//	pipe := global.GVA_REDIS.Pipeline()
//	cmds := make(map[uint]*redis.StringCmd)
//
//	for _, userId := range userIds {
//		key := GetUserInfoKey(userId)
//		cmds[userId] = pipe.Get(ctx, key)
//	}
//
//	pipe.Exec(ctx)
//
//	for userId, cmd := range cmds {
//		data, err := cmd.Result()
//		if err == redis.Nil {
//			missedIds = append(missedIds, userId)
//			continue
//		}
//		if err != nil {
//			missedIds = append(missedIds, userId)
//			continue
//		}
//
//		var cache UserInfoCache
//		if err := json.Unmarshal([]byte(data), &cache); err != nil {
//			missedIds = append(missedIds, userId)
//			continue
//		}
//		result[userId] = &cache
//	}
//
//	return result, missedIds, nil
//}
//
//// BatchSetUserInfoCache 批量设置用户信息缓存
//func BatchSetUserInfoCache(users map[uint]*UserInfoCache) error {
//	ctx := context.Background()
//	pipe := global.GVA_REDIS.Pipeline()
//
//	for userId, info := range users {
//		key := GetUserInfoKey(userId)
//		data, _ := json.Marshal(info)
//		pipe.Set(ctx, key, data, UserInfoExpire)
//	}
//
//	_, err := pipe.Exec(ctx)
//	return err
//}
//
//// ==================== 活动缓存操作 ====================
//
//type ActivityCache struct {
//	ID             uint   `json:"id"`
//	Title          string `json:"title"`
//	Logo           string `json:"logo"`
//	Cover          string `json:"cover"`
//	Description    string `json:"description"`
//	CheckInEnabled int    `json:"checkInEnabled"`
//	DanmakuEnabled int    `json:"danmakuEnabled"`
//	DanmakuAudit   int    `json:"danmakuAudit"`
//	WinnerExclude  int    `json:"winnerExclude"`
//	Status         int    `json:"status"`
//}
//
//// GetActivityCache 获取活动缓存
//func GetActivityCache(activityId uint) (*ActivityCache, error) {
//	ctx := context.Background()
//	key := GetActivityInfoKey(activityId)
//
//	data, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return nil, nil // 缓存不存在
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	var cache ActivityCache
//	if err := json.Unmarshal([]byte(data), &cache); err != nil {
//		return nil, err
//	}
//	return &cache, nil
//}
//
//// SetActivityCache 设置活动缓存
//func SetActivityCache(activity *ActivityCache) error {
//	ctx := context.Background()
//	key := GetActivityInfoKey(activity.ID)
//
//	data, err := json.Marshal(activity)
//	if err != nil {
//		return err
//	}
//
//	return global.GVA_REDIS.Set(ctx, key, data, ActivityCacheExpire).Err()
//}
//
//// DeleteActivityCache 删除活动缓存
//func DeleteActivityCache(activityId uint) error {
//	ctx := context.Background()
//	key := GetActivityInfoKey(activityId)
//	return global.GVA_REDIS.Del(ctx, key).Err()
//}
//
//// ==================== 开关状态操作 ====================
//
//// GetCheckInSwitch 获取签到开关状态
//func GetCheckInSwitch(activityId uint) (bool, error) {
//	ctx := context.Background()
//	key := GetCheckInSwitchKey(activityId)
//
//	val, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return false, nil // 默认关闭
//	}
//	if err != nil {
//		return false, err
//	}
//	return val == "1", nil
//}
//
//// SetCheckInSwitch 设置签到开关状态
//func SetCheckInSwitch(activityId uint, enabled bool) error {
//	ctx := context.Background()
//	key := GetCheckInSwitchKey(activityId)
//
//	val := "0"
//	if enabled {
//		val = "1"
//	}
//	return global.GVA_REDIS.Set(ctx, key, val, SwitchExpire).Err()
//}
//
//// GetDanmakuSwitch 获取弹幕开关状态
//func GetDanmakuSwitch(activityId uint) (bool, error) {
//	ctx := context.Background()
//	key := GetDanmakuSwitchKey(activityId)
//
//	val, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return false, nil
//	}
//	if err != nil {
//		return false, err
//	}
//	return val == "1", nil
//}
//
//// SetDanmakuSwitch 设置弹幕开关状态
//func SetDanmakuSwitch(activityId uint, enabled bool) error {
//	ctx := context.Background()
//	key := GetDanmakuSwitchKey(activityId)
//
//	val := "0"
//	if enabled {
//		val = "1"
//	}
//	return global.GVA_REDIS.Set(ctx, key, val, SwitchExpire).Err()
//}
//
//// ==================== 签到用户操作 ====================
//
//// AddCheckInUser 添加签到用户
//func AddCheckInUser(activityId uint, userId uint) (bool, error) {
//	ctx := context.Background()
//	key := GetCheckInUsersKey(activityId)
//
//	// SADD 返回添加成功的数量，如果用户已存在返回0
//	added, err := global.GVA_REDIS.SAdd(ctx, key, userId).Result()
//	if err != nil {
//		return false, err
//	}
//
//	// 设置过期时间
//	global.GVA_REDIS.Expire(ctx, key, CheckInUsersExpire)
//
//	return added > 0, nil
//}
//
//// IsUserCheckIn 检查用户是否已签到
//func IsUserCheckIn(activityId uint, userId uint) (bool, error) {
//	ctx := context.Background()
//	key := GetCheckInUsersKey(activityId)
//	return global.GVA_REDIS.SIsMember(ctx, key, userId).Result()
//}
//
//// GetCheckInCount 获取签到人数
//func GetCheckInCount(activityId uint) (int64, error) {
//	ctx := context.Background()
//	key := GetCheckInUsersKey(activityId)
//	return global.GVA_REDIS.SCard(ctx, key).Result()
//}
//
//// GetCheckInUserIds 获取所有签到用户ID
//func GetCheckInUserIds(activityId uint) ([]string, error) {
//	ctx := context.Background()
//	key := GetCheckInUsersKey(activityId)
//	return global.GVA_REDIS.SMembers(ctx, key).Result()
//}
//
//// ==================== 游戏状态操作 ====================
//
//// SetCurrentRound 设置当前进行中的场次
//func SetCurrentRound(activityId, roundId uint, duration int) error {
//	key := fmt.Sprintf("annual:activity:%d:current_round", activityId)
//	return global.GVA_REDIS.Set(context.Background(), key, roundId, time.Duration(duration)*time.Second).Err()
//}
//
//// GetCurrentRound 获取当前进行中的场次
//func GetCurrentRound(activityId uint) (uint, error) {
//	ctx := context.Background()
//	key := GetCurrentRoundKey(activityId)
//
//	val, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return 0, nil
//	}
//	if err != nil {
//		return 0, err
//	}
//
//	var roundId uint
//	fmt.Sscanf(val, "%d", &roundId)
//	return roundId, nil
//}
//
//// GetCurrentRoundRemaining 获取当前场次剩余时间（秒）
//func GetCurrentRoundRemaining(activityId uint) (int, error) {
//	key := fmt.Sprintf("annual:activity:%d:current_round", activityId)
//	ttl, err := global.GVA_REDIS.TTL(context.Background(), key).Result()
//	if err != nil {
//		return 0, err
//	}
//	if ttl < 0 {
//		return 0, nil
//	}
//	return int(ttl.Seconds()), nil
//}
//
//// ClearCurrentRound 清除当前场次
//func ClearCurrentRound(activityId uint) error {
//	ctx := context.Background()
//	key := GetCurrentRoundKey(activityId)
//	return global.GVA_REDIS.Del(ctx, key).Err()
//}
//
//// SetRoundStatus 设置场次状态
//func SetRoundStatus(roundId uint, status int) error {
//	ctx := context.Background()
//	key := GetRoundStatusKey(roundId)
//	return global.GVA_REDIS.Set(ctx, key, status, RoundScoresExpire).Err()
//}
//
//// GetRoundStatus 获取场次状态
//func GetRoundStatus(roundId uint) (int, error) {
//	ctx := context.Background()
//	key := GetRoundStatusKey(roundId)
//
//	val, err := global.GVA_REDIS.Get(ctx, key).Result()
//	if err == redis.Nil {
//		return -1, nil // 不存在
//	}
//	if err != nil {
//		return -1, err
//	}
//
//	var status int
//	fmt.Sscanf(val, "%d", &status)
//	return status, nil
//}
//
//// ==================== 摇一摇积分操作 ====================
//
//// IncrUserScore 增加用户积分
//func IncrUserScore(roundId uint, userId uint, score float64) (float64, error) {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//
//	newScore, err := global.GVA_REDIS.ZIncrBy(ctx, key, score, fmt.Sprintf("%d", userId)).Result()
//	if err != nil {
//		return 0, err
//	}
//
//	// 设置过期时间
//	global.GVA_REDIS.Expire(ctx, key, RoundScoresExpire)
//
//	return newScore, nil
//}
//
//// GetUserScore 获取用户积分
//func GetUserScore(roundId uint, userId uint) (float64, error) {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//
//	score, err := global.GVA_REDIS.ZScore(ctx, key, fmt.Sprintf("%d", userId)).Result()
//	if err == redis.Nil {
//		return 0, nil
//	}
//	return score, err
//}
//
//// GetUserRank 获取用户排名（从0开始）
//func GetUserRank(roundId uint, userId uint) (int64, error) {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//
//	rank, err := global.GVA_REDIS.ZRevRank(ctx, key, fmt.Sprintf("%d", userId)).Result()
//	if err == redis.Nil {
//		return -1, nil // 不在排行榜中
//	}
//	return rank, err
//}
//
//// RankingItem 排行榜项
//type RankingItem struct {
//	UserId uint    `json:"userId"`
//	Score  float64 `json:"score"`
//	Rank   int     `json:"rank"`
//}
//
//// GetRanking 获取排行榜
//func GetRanking(roundId uint, limit int64) ([]RankingItem, error) {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//
//	results, err := global.GVA_REDIS.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	ranking := make([]RankingItem, 0, len(results))
//	for i, z := range results {
//		var userId uint
//		fmt.Sscanf(z.Member.(string), "%d", &userId)
//		ranking = append(ranking, RankingItem{
//			UserId: userId,
//			Score:  z.Score,
//			Rank:   i + 1,
//		})
//	}
//
//	return ranking, nil
//}
//
//// GetPlayerCount 获取参与人数
//func GetPlayerCount(roundId uint) (int64, error) {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//	return global.GVA_REDIS.ZCard(ctx, key).Result()
//}
//
//// ClearRoundScores 清除场次积分
//func ClearRoundScores(roundId uint) error {
//	ctx := context.Background()
//	key := GetRoundScoresKey(roundId)
//	return global.GVA_REDIS.Del(ctx, key).Err()
//}
