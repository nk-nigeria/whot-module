package cgbdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/whot-module/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const FeeGameTableName = "feegame"

// CREATE TABLE
//   public.feegame (
//     id bigint NOT NULL,
// 	   user_id character varying(128) NOT NULL,
// 	   game character varying(128) NOT NULL,
//     fee bigint NOT NULL DEFAULT 0,
//     create_time timestamp
//     with
//       time zone NOT NULL DEFAULT now(),
//       update_time timestamp
//     with
//       time zone NOT NULL DEFAULT now()
//   );

// ALTER TABLE
//   public.feegame
// ADD
//   CONSTRAINT feegame_pkey PRIMARY KEY (id)

func AddNewFeeGame(ctx context.Context, logger runtime.Logger, db *sql.DB, feeGame entity.FeeGame) (int64, error) {
	feeGame.Id = entity.SnowlakeNode.Generate().Int64()
	if feeGame.Game == "" {
		feeGame.Game = entity.ModuleName
	}
	query := "INSERT INTO " + FeeGameTableName +
		" (id, user_id, game, fee,create_time, update_time) " +
		" VALUES ( $1, $2, $3, $4, now(), now())"
	result, err := db.ExecContext(ctx, query,
		feeGame.Id, feeGame.UserID, feeGame.Game, feeGame.Fee)
	if err != nil {
		logger.Error("Add new fee game user %s, game %s, fee %d err %s",
			feeGame.UserID, feeGame.Game, feeGame.Fee, err.Error())
		return 0, status.Error(codes.Internal, "Error add fee game.")
	}
	if rowsAffectedCount, _ := result.RowsAffected(); rowsAffectedCount != 1 {
		logger.Error("Did not insert new free game, user %s, game %s, fee %d ",
			feeGame.UserID, feeGame.Game, feeGame.Fee)
		return 0, status.Error(codes.Internal, "Error add fee game.")
	}
	return feeGame.Id, nil
}

func AddNewMultiFeeGame(ctx context.Context, logger runtime.Logger, db *sql.DB, listFeeGame []entity.FeeGame) error {
	if len(listFeeGame) == 0 {
		return nil
	}
	query := "INSERT INTO " + FeeGameTableName +
		" (id, user_id, game, fee,create_time, update_time)  VALUES "

	args := make([]interface{}, 0)
	lenList := len(listFeeGame)
	for idx, feeGame := range listFeeGame {
		feeGame.Id = entity.SnowlakeNode.Generate().Int64()
		if feeGame.Game == "" {
			feeGame.Game = entity.ModuleName
		}
		args = append(args, feeGame.Id, feeGame.UserID, feeGame.Game, feeGame.Fee)
		paramIdx := idx*4 + 1
		query += fmt.Sprintf("( $%d, $%d, $%d, $%d, now(), now())", paramIdx, paramIdx+1, paramIdx+2, paramIdx+3)
		if idx < lenList-1 {
			query += ","
		}
	}
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("Add new list fee game  err %s",
			err.Error())
		return status.Error(codes.Internal, "Error add list fee game.")
	}
	if rowsAffectedCount, _ := result.RowsAffected(); rowsAffectedCount == 0 {
		logger.Error("Did not insert new list free game")
		return status.Error(codes.Internal, "Error add list fee game.")
	}
	return nil
}
