package api

import (
	"context"
	"database/sql"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/whot-module/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*entity.MatchState)
	logger.Info("match join attempt, state=%v, meta=%v", s, metadata)

	// check password
	if !s.Label.Open {
		logger.Info("match protect with password, check password")
		joinPassword := metadata["password"]
		if joinPassword != s.Label.Password {
			return s, false, "wrong password"
		}
	}

	// Check if it's a user attempting to rejoin after a disconnect.
	if pRaw, _ := s.Presences.Get(presence.GetUserId()); pRaw != nil {
		if p, ok := pRaw.(entity.MyPrecense); ok {
			uid := presence.GetUserId()
			deviceID := metadata["device_id"]
			if p.Presence.GetSessionId() == presence.GetSessionId() || p.DeviceID == deviceID {
				logger.Info("user %s rejoin after disconnect ", uid)
				s.JoinsInProgress++
				return s, true, "Rejoining match Successfully"
			} else {
				logger.Info("user %s join from 2 different devices at the same time. --> reject join match", uid)
				return s, false, "you were in another table"
			}
		} else {
			logger.Warn("Presence type assertion failed")
			return s, false, "internal error"
		}
	}

	// join as new user

	// Check if match is full.
	if s.Presences.Size()+s.JoinsInProgress >= entity.MaxPresences {
		return s, false, "match full"
	}
	// check chip balance in wallet before allow join
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, presence.GetUserId())
	if err != nil {
		return s, false, status.Error(codes.Internal, "read chip balance failed").Error()
	}
	if wallet.Chips < int64(s.Label.Bet.GetMarkUnit()) {
		logger.Warn("[Reject] reject allow user %s join game, not enough chip join game, balance user chip %d , game bet %d",
			presence.GetUserId(), wallet.Chips, s.Label.Bet)
		return s, false, status.Error(codes.Internal, "chip balance not enough").Error()
	}

	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, "Join match Successfully"
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)
	logger.Info("match join, state=%v, presences=%v", s, presences)

	m.processor.ProcessPresencesJoin(ctx, logger, nk, db, dispatcher, s, presences)

	// If we're in matching state, reset the countdown
	if m.machine.IsMatchingState() {
		s.SetUpCountDown(10 * time.Second)
	}
	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)

	logger.Info("match leave, state=%v, presences=%v", s, presences)

	if m.machine.IsPreparingState() || m.machine.IsPlayingState() || m.machine.IsReward() {
		m.processor.ProcessPresencesLeavePending(ctx, logger, nk, dispatcher, s, presences)
		return s
	}
	m.processor.ProcessPresencesLeave(ctx, logger, nk, db, dispatcher, s, presences)
	if m.machine.IsMatchingState() {
		s.SetUpCountDown(10 * time.Second)
	}
	return s
}
