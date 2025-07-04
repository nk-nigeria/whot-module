package api

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/whot-module/conf"
	"google.golang.org/protobuf/proto"
)

const (
	rpcGetBotConfig    = "get_bot_config"
	rpcUpdateBotConfig = "update_bot_config"
)

// RpcGetBotConfig lấy cấu hình bot hiện tại
func RpcGetBotConfig(marshaler *proto.MarshalOptions, unmarshaler *proto.UnmarshalOptions) func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	return func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
		config := conf.GetBotConfig()

		response := map[string]interface{}{
			"success": true,
			"data":    config,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			logger.Error("Failed to marshal bot config response: %v", err)
			return "", err
		}

		return string(responseBytes), nil
	}
}

// RpcUpdateBotConfig cập nhật cấu hình bot
func RpcUpdateBotConfig(marshaler *proto.MarshalOptions, unmarshaler *proto.UnmarshalOptions) func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	return func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
		// Parse request payload
		var request map[string]interface{}
		if err := json.Unmarshal([]byte(payload), &request); err != nil {
			logger.Error("Failed to unmarshal bot config request: %v", err)
			return "", err
		}

		// Parse config data
		configData, ok := request["config"].(map[string]interface{})
		if !ok {
			logger.Error("Invalid config data in request")
			return "", runtime.NewError("Invalid config data", 3)
		}

		// Convert to BotConfig
		configBytes, err := json.Marshal(configData)
		if err != nil {
			logger.Error("Failed to marshal config data: %v", err)
			return "", err
		}

		var config conf.BotConfig
		if err := json.Unmarshal(configBytes, &config); err != nil {
			logger.Error("Failed to unmarshal config: %v", err)
			return "", err
		}

		// Validate config
		if err := validateBotConfig(&config); err != nil {
			logger.Error("Invalid bot config: %v", err)
			return "", runtime.NewError("Invalid bot config: "+err.Error(), 3)
		}

		// Save to database
		if err := conf.SaveBotConfigToDB(ctx, logger, db, &config); err != nil {
			logger.Error("Failed to save bot config to database: %v", err)
			return "", err
		}

		// Update in-memory config
		conf.InitBotConfig()

		response := map[string]interface{}{
			"success": true,
			"message": "Bot config updated successfully",
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			logger.Error("Failed to marshal response: %v", err)
			return "", err
		}

		return string(responseBytes), nil
	}
}

// validateBotConfig kiểm tra tính hợp lệ của cấu hình bot
func validateBotConfig(config *conf.BotConfig) error {
	// Validate BotJoinRules
	for i, rule := range config.BotJoinRules {
		if rule.MinBet < 0 || rule.MaxBet < rule.MinBet {
			return runtime.NewError("Invalid bet range in bot_join_rules at index "+string(rune(i)), 3)
		}
		if rule.MinUsers < 1 || rule.MaxUsers < rule.MinUsers {
			return runtime.NewError("Invalid user range in bot_join_rules at index "+string(rune(i)), 3)
		}
		if rule.RandomTimeMin < 0 || rule.RandomTimeMax < rule.RandomTimeMin {
			return runtime.NewError("Invalid random time in bot_join_rules at index "+string(rune(i)), 3)
		}
		if rule.JoinPercent < 0 || rule.JoinPercent > 100 {
			return runtime.NewError("Invalid join_percent in bot_join_rules at index "+string(rune(i)), 3)
		}
	}
	// Validate BotLeaveRules
	for i, rule := range config.BotLeaveRules {
		if rule.MinBet < 0 || rule.MaxBet < rule.MinBet {
			return runtime.NewError("Invalid bet range in bot_leave_rules at index "+string(rune(i)), 3)
		}
		if rule.LeavePercent < 0 || rule.LeavePercent > 100 {
			return runtime.NewError("Invalid leave_percent in bot_leave_rules at index "+string(rune(i)), 3)
		}
		if rule.LastResult != 0 && rule.LastResult != 1 && rule.LastResult != -1 {
			return runtime.NewError("Invalid last_result in bot_leave_rules at index "+string(rune(i)), 3)
		}
	}
	// Validate BotCreateTableRules
	for i, rule := range config.BotCreateTableRules {
		if rule.MinBet < 0 || rule.MaxBet < rule.MinBet {
			return runtime.NewError("Invalid bet range in bot_create_table_rules at index "+string(rune(i)), 3)
		}
		if rule.MinActiveTables < 0 || rule.MaxActiveTables < rule.MinActiveTables {
			return runtime.NewError("Invalid active_tables in bot_create_table_rules at index "+string(rune(i)), 3)
		}
		if rule.WaitTimeMin < 0 || rule.WaitTimeMax < rule.WaitTimeMin {
			return runtime.NewError("Invalid wait_time in bot_create_table_rules at index "+string(rune(i)), 3)
		}
		if rule.RetryWaitMin < 0 || rule.RetryWaitMax < rule.RetryWaitMin {
			return runtime.NewError("Invalid retry_wait in bot_create_table_rules at index "+string(rune(i)), 3)
		}
	}
	// Validate BotGroupRules
	for i, rule := range config.BotGroupRules {
		if rule.VIPMin < 0 || rule.VIPMax < rule.VIPMin {
			return runtime.NewError("Invalid VIP range in bot_group_rules at index "+string(rune(i)), 3)
		}
		if rule.MCBMin < 0 || rule.MCBMax < rule.MCBMin {
			return runtime.NewError("Invalid MCB range in bot_group_rules at index "+string(rune(i)), 3)
		}
	}
	return nil
}
