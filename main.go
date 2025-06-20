package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/nk-nigeria/whot-module/constant"
	"github.com/nk-nigeria/whot-module/message_queue"
	mockcodegame "github.com/nk-nigeria/whot-module/mock_code_game"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/whot-module/entity"
	"google.golang.org/protobuf/proto"

	"github.com/nk-nigeria/whot-module/api"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	rpcIdGameList    = "list_game"
	rpcIdFindMatch   = "find_match"
	rpcIdCreateMatch = "create_match"
)

// noinspection GoUnusedExportedFunction
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &proto.MarshalOptions{}
	unmarshaler := &proto.UnmarshalOptions{
		DiscardUnknown: false,
	}
	message_queue.InitNatsService(logger, constant.NastEndpoint, marshaler)
	mockcodegame.InitMapMockCodeListCard(logger)
	// cgbdb.RunMigrations(ctx, logger, db)
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}

	// initializer.RegisterMatchmakerMatched(func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	// logger.Info("MatchmakerMatched triggered with %d players", len(entries))

	// props := entries[0].GetProperties()

	// name := "unknown"
	// if v, ok := props["name"]; ok {
	// 	if str, ok := v.(string); ok {
	// 		name = str
	// 	} else {
	// 		logger.Warn("name không phải string")
	// 	}
	// }

	// password := ""
	// if v, ok := props["password"]; ok {
	// 	if str, ok := v.(string); ok {
	// 		password = str
	// 	} else {
	// 		logger.Warn("password không phải string")
	// 	}
	// }

	// var markUnit int32 = 0
	// if v, ok := props["bet"]; ok {
	// 	if f64, ok := v.(float64); ok {
	// 		markUnit = int32(f64)
	// 	} else {
	// 		logger.Warn("bet không phải float64")
	// 	}
	// } else {
	// 	logger.Warn("Không tìm thấy 'bet'")
	// }

	// matchInfo := &pb.Match{
	// 	Size:     2,
	// 	MaxSize:  4,
	// 	Name:     name,
	// 	Open:     len(password) > 0,
	// 	Password: password,
	// 	NumBot:   1,
	// 	MarkUnit: markUnit,
	// }
	// data, err := utilities.EncodeBase64Proto(matchInfo)
	// if err != nil {
	// 	logger.Error("failed to encode match data: %v", err)
	// 	return "", err
	// }

	// args := map[string]any{"data": data}
	// matchId, err := nk.MatchCreate(ctx, entity.ModuleName, args)
	// if err != nil {
	// 	logger.Error("failed to create match: %v", err)
	// 	return "", err
	// }
	// return matchId, nil
	// })

	if err := api.RegisterSessionEvents(db, nk, initializer); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
