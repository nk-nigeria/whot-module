package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
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
			State:     pb.GameState_GameStateMatching,
			CountDown: int64(matchingTimeout.Seconds()),
		},
	)

	// Check bot join only if we haven't reached maximum players (4)
	if state.GetPresenceSize() < entity.MaxPresences {
		// Get match ID for context
		matchID := procPkg.GetContext().Value("match_id").(string)

		// Create bot integration and update match state
		botIntegration := service.NewWhotBotIntegration(procPkg.GetDb())
		botIntegration.SetMatchState(
			matchID,
			int64(state.Label.GetMarkUnit()),
			state.GetPresenceSize(),
			0, // lastResult
			0, // activeTables
		)

		// Process bot logic
		botCtx := packager.GetContextWithProcessorPackager(procPkg)
		if err := botIntegration.ProcessBotLogic(botCtx); err != nil {
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

	// Check bot join only if we haven't reached maximum players (4)
	if state.GetPresenceSize() < entity.MaxPresences {
		// Get match ID for context
		matchID := procPkg.GetContext().Value("match_id").(string)

		// Create bot integration and update match state
		botIntegration := service.NewWhotBotIntegration(procPkg.GetDb())
		botIntegration.SetMatchState(
			matchID,
			int64(state.Label.GetMarkUnit()),
			state.GetPresenceSize(),
			0, // lastResult
			0, // activeTables
		)

		// Process bot logic
		botCtx := packager.GetContextWithProcessorPackager(procPkg)
		if err := botIntegration.ProcessBotLogic(botCtx); err != nil {
			procPkg.GetLogger().Error("[matching] Bot logic error: %v", err)
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
