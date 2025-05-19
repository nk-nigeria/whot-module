package cgbdb

import (
	"context"
	"database/sql"
	"strings"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/jackc/pgtype"
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const JackpotTableName = "jackpot"

// CREATE TABLE
//   public.jackpot (
//     id bigint NOT NULL,
// 	   game character varying(128) NOT NULL,
// 	   UNIQUE(game),
//     chips bigint NOT NULL DEFAULT 0,
//     create_time timestamp
//     with
//       time zone NOT NULL DEFAULT now(),
//       update_time timestamp
//     with
//       time zone NOT NULL DEFAULT now()
//   );

// ALTER TABLE
//
//	public.jackpot
//
// ADD
//
//	CONSTRAINT jackpot_pkey PRIMARY KEY (id)
func AddOrUpdateChipJackpot(ctx context.Context, logger runtime.Logger, db *sql.DB, game string, chips int64) error {
	var err error
	var num int64
	defer func() {
		if err == nil {
			AddJackpotHistoryChipsChange(ctx, logger, db, game, chips)
		}
	}()
	num, err = incChipJackpot(ctx, logger, db, game, chips)
	if num > 0 {
		return nil
	}
	query := "INSERT INTO " + JackpotTableName +
		" (id, game, chips, create_time, update_time) " +
		" SELECT $1, $2, $3, now(), now() " +
		" WHERE NOT EXISTS (SELECT id FROM " + JackpotTableName + " WHERE game=$4) "
	result, err := db.ExecContext(ctx, query,
		entity.SnowlakeNode.Generate().Int64(),
		game, chips, game)
	if err != nil {
		logger.WithField("game", game).WithField("err", err.Error()).Error("Error err when insert jackpot")
		err = status.Error(codes.Internal, "Error insert jackpot")
		return err
	}
	if rowsAffectedCount, _ := result.RowsAffected(); rowsAffectedCount == 0 {
		// return 0, status.Error(codes.Internal, "Error update reward refer user.")
		err = status.Error(codes.Internal, "Error when insert jackpot")
	}
	return err
}

func incChipJackpot(ctx context.Context, logger runtime.Logger, db *sql.DB, game string, chips int64) (int64, error) {
	query := "UPDATE " + JackpotTableName +
		" SET chips= " +
		" CASE WHEN chips+$2 > 0 then chips+$3 else 0 END, update_time=now()" +
		" WHERE game=$1"
	result, err := db.ExecContext(ctx, query, game, chips, chips)
	if err != nil {
		logger.WithField("game", game).WithField("err", err.Error()).Error("Error when update jackpot")
		return 0, status.Error(codes.Internal, "Error update jackpot")
	}
	rowsAffectedCount, _ := result.RowsAffected()
	if rowsAffectedCount != 1 {
		return 0, status.Error(codes.Internal, "Error update not effect")
	}
	return rowsAffectedCount, nil
}

func GetJackpot(ctx context.Context, logger runtime.Logger, db *sql.DB, game string) (*pb.Jackpot, error) {
	query := "SELECT id, game, chips, create_time FROM " + JackpotTableName +
		" WHERE game=$1"
	var dbId, dbChips int64
	var dbGame string
	var dbCreateTime pgtype.Timestamptz
	err := db.QueryRowContext(ctx, query, game).
		Scan(&dbId, &dbGame, &dbChips, &dbCreateTime)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return &pb.Jackpot{}, nil
		}
		logger.WithField("game", game).WithField("err", err.Error()).Error("Query jackpot error")
		return nil, status.Error(codes.Internal, "Query jackpot error")
	}
	jackpot := &pb.Jackpot{
		Id:             dbId,
		GameCode:       dbGame,
		Chips:          dbChips,
		CreateTimeUnix: dbCreateTime.Time.Unix(),
	}
	return jackpot, nil
}
