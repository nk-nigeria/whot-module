package processor

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nakama-nigeria/whot-module/entity"
	"google.golang.org/protobuf/proto"
)

type UseCase interface {
	// Game Flow
	ProcessNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState)
	ProcessFinishGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState)

	// Player Actions
	PlayCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData)        // Đánh 1 lá bài hợp lệ
	ChooseWhotShape(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) //  chọn hình sau khi đánh Whot
	DrawCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData)        // Rút bài nếu không đánh được

	// State update
	NotifyUpdateGameState(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message)
	NotifyUpdateTable(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message)

	// Presence (player join/leave)
	ProcessPresencesJoin(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence)
	ProcessPresencesLeave(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence)
	ProcessPresencesLeavePending(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence)
	ProcessApplyPresencesLeave(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState)
}
