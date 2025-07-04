package service

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/nk-nigeria/whot-module/conf"
)

// BotManagementService xử lý logic thêm bot vào trận đấu
// Sử dụng rule mapping từ Excel

type BotManagementService struct {
	db *sql.DB
}

func NewBotManagementService(db *sql.DB) *BotManagementService {
	return &BotManagementService{
		db: db,
	}
}

// Tìm rule join phù hợp
func (s *BotManagementService) FindBotJoinRule(betAmount int64, userCount int) *conf.BotJoinRule {
	for _, rule := range conf.GetBotConfig().BotJoinRules {
		if betAmount >= rule.MinBet && betAmount < rule.MaxBet &&
			userCount >= rule.MinUsers && userCount <= rule.MaxUsers {
			return &rule
		}
	}
	return nil
}

// Tìm rule leave phù hợp
func (s *BotManagementService) FindBotLeaveRule(betAmount int64, lastResult int) *conf.BotLeaveRule {
	for _, rule := range conf.GetBotConfig().BotLeaveRules {
		if betAmount >= rule.MinBet && betAmount < rule.MaxBet &&
			(rule.LastResult == 0 || rule.LastResult == lastResult) {
			return &rule
		}
	}
	return nil
}

// Tìm rule tạo bàn phù hợp
func (s *BotManagementService) FindBotCreateTableRule(betAmount int64, activeTables int) *conf.BotCreateTableRule {
	for _, rule := range conf.GetBotConfig().BotCreateTableRules {
		if betAmount >= rule.MinBet && betAmount < rule.MaxBet &&
			activeTables >= rule.MinActiveTables && activeTables <= rule.MaxActiveTables {
			return &rule
		}
	}
	return nil
}

// Tìm rule phân nhóm bot phù hợp
func (s *BotManagementService) FindBotGroupRule(vip int, mcb int64) *conf.BotGroupRule {
	for _, rule := range conf.GetBotConfig().BotGroupRules {
		if vip >= rule.VIPMin && vip <= rule.VIPMax &&
			mcb >= rule.MCBMin && mcb < rule.MCBMax {
			return &rule
		}
	}
	return nil
}

// Logic bot join bàn (áp dụng xác suất, random time)
func (s *BotManagementService) ShouldBotJoin(ctx context.Context, betAmount int64, userCount int) bool {
	rule := s.FindBotJoinRule(betAmount, userCount)
	if rule == nil {
		return false
	}
	if rand.Intn(100) >= rule.JoinPercent {
		return false
	}
	if rule.RandomTimeMax > rule.RandomTimeMin {
		delay := rand.Intn(rule.RandomTimeMax-rule.RandomTimeMin+1) + rule.RandomTimeMin
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return true
}

// Logic bot leave bàn (áp dụng xác suất)
func (s *BotManagementService) ShouldBotLeave(ctx context.Context, betAmount int64, lastResult int) bool {
	rule := s.FindBotLeaveRule(betAmount, lastResult)
	if rule == nil {
		return false
	}
	return rand.Intn(100) < rule.LeavePercent
}

// Logic bot tạo bàn (áp dụng số bàn active, thời gian chờ)
func (s *BotManagementService) ShouldBotCreateTable(ctx context.Context, betAmount int64, activeTables int) (bool, int, int) {
	rule := s.FindBotCreateTableRule(betAmount, activeTables)
	if rule == nil {
		return false, 0, 0
	}
	waitTime := rule.WaitTimeMin
	if rule.WaitTimeMax > rule.WaitTimeMin {
		waitTime = rand.Intn(rule.WaitTimeMax-rule.WaitTimeMin+1) + rule.WaitTimeMin
	}
	retryWait := rule.RetryWaitMin
	if rule.RetryWaitMax > rule.RetryWaitMin {
		retryWait = rand.Intn(rule.RetryWaitMax-rule.RetryWaitMin+1) + rule.RetryWaitMin
	}
	return true, waitTime, retryWait
}

// Logic phân nhóm bot
func (s *BotManagementService) GetBotGroup(vip int, mcb int64) *conf.BotGroupRule {
	return s.FindBotGroupRule(vip, mcb)
}

// Helper: random int trong khoảng [min, max]
func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
