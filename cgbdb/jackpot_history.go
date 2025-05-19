package cgbdb

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nakamaFramework/whot-module/entity"
)

const JackpotHisotryTableName = "jackpot_history"

func AddJackpotHistoryChipsChange(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB, game string,
	chips int64,
) error {
	dbId := entity.SnowlakeNode.Generate().Int64()
	query := "INSERT INTO " + JackpotHisotryTableName +
		" (id, game, chips, metadata, create_time, update_time) " +
		" VALUES ($1, $2, $3, $4, now(), now())"
	_, err := db.ExecContext(ctx, query, dbId, game, chips, "")
	if err != nil && err != context.DeadlineExceeded {
		logger.WithField("err", err).Error("add jackpot history failed.")
		return err
	}
	return nil
}

func AddJackpotHistoryUserWin(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB, game string,
	user string,
	chips int64,
) error {
	dbId := entity.SnowlakeNode.Generate().Int64()
	query := "INSERT INTO " + JackpotHisotryTableName +
		" (id, game, chips, metadata, create_time, update_time) " +
		` VALUES ($1, $2, $3, jsonb_build_object('uid_win', '` + user + `'), now(), now())`
	_, err := db.ExecContext(ctx, query, dbId, game, chips)
	if err != nil && err != context.DeadlineExceeded {
		logger.WithField("err", err).Error("add jackpot history failed.")
		return err
	}
	return nil
}
