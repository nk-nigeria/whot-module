package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/whot-module/entity"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/nk-nigeria/whot-module/usecase/service"
)

type StateMatching struct {
	StateBase
}

func NewStateMatching(fn FireFn) *StateMatching {
	return &StateMatching{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StateMatching) Enter(ctx context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[matching] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.SetUpCountDown(matchingTimeout)
	procPkg.GetLogger().Info("apply leave presence")
	procPkg.GetProcessor().ProcessApplyPresencesLeave(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)
	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GAME_STATE_MATCHING,
			CountDown: int64(matchingTimeout.Seconds()),
		},
	)

	// Initialize bot integration if not already done
	botIntegration := GetGlobalBotIntegration()
	if botIntegration == nil {
		botIntegration = service.NewWhotBotIntegration(procPkg.GetDb())
		SetGlobalBotIntegration(botIntegration)
	}

	// Step 1: Make initial bot join decision (only once)
	if state.GetPresenceSize() < entity.MaxPresences {
		// Update match state
		botIntegration.SetMatchState(
			state.Label.GetMatchId(),
			int64(state.Label.GetMarkUnit()),
			state.GetPresenceSize(),
			0, // lastResult
			0, // activeTables
		)

		// Make initial bot join decision (random once)
		botCtx := packager.GetContextWithProcessorPackager(procPkg)
		log.GetLogger().Info("[matching] Making initial bot join decision")
		if err := botIntegration.ProcessJoinBotLogic(botCtx); err != nil {
			procPkg.GetLogger().Error("[matching] Bot logic error: %v", err)
		}
	} else {
		procPkg.GetLogger().Info("[matching] Skip bot join - maximum players reached (%d)", entity.MaxPresences)
	}

	return nil
}

func (s *StateMatching) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[matching] exit")
	return nil
}

func (s *StateMatching) Process(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[matching] processing - checking bot join")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	remain := state.GetRemainCountDown()

	// Step 2: Check and join bots based on time (continuously)
	if state.GetPresenceSize() < entity.MaxPresences {
		// Check if bot should join based on time
		botCtx := packager.GetContextWithProcessorPackager(procPkg)
		joined, err := GetGlobalBotIntegration().CheckAndJoinExpiredBots(botCtx)
		if err != nil {
			procPkg.GetLogger().Error("[matching] Bot join error: %v", err)
		} else if joined {
			procPkg.GetLogger().Info("[matching] Bot joined based on time")
		}
	} else {
		procPkg.GetLogger().Info("[matching] Skip bot join - maximum players reached (%d)", entity.MaxPresences)
	}

	if remain > 0 {
		return nil
	}

	if state.GetPresenceNotBotSize() == 0 {
		s.Trigger(ctx, triggerIdle)
		log.GetLogger().Info("[matching] no one presence, trigger idle")
		return nil
	}

	presenceCount := state.GetPresenceSize()
	if state.IsReadyToPlay() {
		s.Trigger(ctx, triggerPresenceReady)
	} else if presenceCount <= 0 {
		s.Trigger(ctx, triggerIdle)
	}

	return nil
}
