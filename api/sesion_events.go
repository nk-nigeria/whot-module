package api

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func RegisterSessionEvents(db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := initializer.RegisterEventSessionStart(eventSessionStartFunc(nk)); err != nil {
		return err
	}
	if err := initializer.RegisterEventSessionEnd(eventSessionEndFunc(nk)); err != nil {
		return err
	}
	return nil
}

func eventSessionStartFunc(nk runtime.NakamaModule) func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	return func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
		// restoreMatchSession(ctx, logger, nk)

	}
}

func eventSessionEndFunc(nk runtime.NakamaModule) func(context.Context, runtime.Logger, *api.Event) {
	return func(ctx context.Context, logger runtime.Logger, evt *api.Event) {

	}
}

// func restoreMatchSession(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule) {
// 	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
// 	if !ok {
// 		logger.Error("context did not contain user ID.")
// 		return
// 	}
// 	profile, err := entity.GetProfileUser(ctx, nk, userID)
// 	if err != nil {
// 		return
// 	}
// 	if profile.GetPlayingMatch() == "" {
// 		return
// 	}
// 	match, err := nk.MatchGet(ctx, profile.GetPlayingMatch())
// 	if err != nil {
// 		return
// 	}
// 	match.GetMatchId()
// }
