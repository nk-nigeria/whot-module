package conf

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

// --- Struct mới mapping đúng Excel ---
type BotJoinRule struct {
	MinBet        int64 `json:"min_bet"`
	MaxBet        int64 `json:"max_bet"`
	MinUsers      int   `json:"min_users"`
	MaxUsers      int   `json:"max_users"`
	RandomTimeMin int   `json:"random_time_min"`
	RandomTimeMax int   `json:"random_time_max"`
	JoinPercent   int   `json:"join_percent"`
}

type BotLeaveRule struct {
	MinBet       int64 `json:"min_bet"`
	MaxBet       int64 `json:"max_bet"`
	LastResult   int   `json:"last_result"`
	LeavePercent int   `json:"leave_percent"`
}

type BotCreateTableRule struct {
	MinBet          int64 `json:"min_bet"`
	MaxBet          int64 `json:"max_bet"`
	MinActiveTables int   `json:"min_active_tables"`
	MaxActiveTables int   `json:"max_active_tables"`
	WaitTimeMin     int   `json:"wait_time_min"`
	WaitTimeMax     int   `json:"wait_time_max"`
	RetryWaitMin    int   `json:"retry_wait_min"`
	RetryWaitMax    int   `json:"retry_wait_max"`
}

type BotGroupRule struct {
	VIPMin int   `json:"vip_min"`
	VIPMax int   `json:"vip_max"`
	MCBMin int64 `json:"mcb_min"`
	MCBMax int64 `json:"mcb_max"`
}

type BotConfig struct {
	BotJoinRules        []BotJoinRule        `json:"bot_join_rules"`
	BotLeaveRules       []BotLeaveRule       `json:"bot_leave_rules"`
	BotCreateTableRules []BotCreateTableRule `json:"bot_create_table_rules"`
	BotGroupRules       []BotGroupRule       `json:"bot_group_rules"`
}

var botConfig *BotConfig

func InitBotConfig() {
	// Khởi tạo config mẫu (có thể để rỗng hoặc load từ file mẫu)
	botConfig = &BotConfig{
		BotJoinRules:        []BotJoinRule{},
		BotLeaveRules:       []BotLeaveRule{},
		BotCreateTableRules: []BotCreateTableRule{},
		BotGroupRules:       []BotGroupRule{},
	}
}

func GetBotConfig() *BotConfig {
	if botConfig == nil {
		InitBotConfig()
	}
	return botConfig
}

func LoadBotConfigFromDB(ctx context.Context, logger runtime.Logger, db *sql.DB) error {
	query := `SELECT config_data FROM bot_config WHERE is_active = true ORDER BY created_at DESC LIMIT 1`
	var configData []byte
	err := db.QueryRowContext(ctx, query).Scan(&configData)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("No bot config found in database, using default config")
			return nil
		}
		return fmt.Errorf("failed to load bot config from database: %v", err)
	}
	var config BotConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to unmarshal bot config: %v", err)
	}
	botConfig = &config
	logger.Info("Bot config loaded from database successfully")
	return nil
}

func SaveBotConfigToDB(ctx context.Context, logger runtime.Logger, db *sql.DB, config *BotConfig) error {
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal bot config: %v", err)
	}
	query := `INSERT INTO bot_config (config_data, is_active, created_at) VALUES ($1, true, NOW())`
	_, err = db.ExecContext(ctx, query, configData)
	if err != nil {
		return fmt.Errorf("failed to save bot config to database: %v", err)
	}
	logger.Info("Bot config saved to database successfully")
	return nil
}
