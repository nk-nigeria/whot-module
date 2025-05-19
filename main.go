package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/nakamaFramework/whot-module/constant"
	"github.com/nakamaFramework/whot-module/message_queue"
	mockcodegame "github.com/nakamaFramework/whot-module/mock_code_game"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nakamaFramework/whot-module/entity"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/nakamaFramework/whot-module/api"
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

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	message_queue.InitNatsService(logger, constant.NastEndpoint, marshaler)
	mockcodegame.InitMapMockCodeListCard()
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}

	if err := api.RegisterSessionEvents(db, nk, initializer); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
