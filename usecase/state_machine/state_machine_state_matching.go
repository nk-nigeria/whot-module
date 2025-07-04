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
		botService := service.NewBotManagementService(procPkg.GetDb())
		betAmount := int64(state.Label.GetMarkUnit())
		userCount := state.GetPresenceSize()
		if botService.ShouldBotJoin(procPkg.GetContext(), betAmount, userCount) {
			procPkg.GetLogger().Info("[matching] Bot join triggered by rule")
			procPkg.GetProcessor().AddBotToMatch(
				procPkg.GetContext(),
				procPkg.GetLogger(),
				procPkg.GetNK(),
				procPkg.GetDb(),
				procPkg.GetDispatcher(),
				state,
				1,
			)
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
	// log.GetLogger().Info("[matching] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	remain := state.GetRemainCountDown()

	// Check bot join only if we haven't reached maximum players (4)
	if state.GetPresenceSize() < entity.MaxPresences {
		botService := service.NewBotManagementService(procPkg.GetDb())
		betAmount := int64(state.Label.GetMarkUnit())
		userCount := state.GetPresenceSize()
		if botService.ShouldBotJoin(procPkg.GetContext(), betAmount, userCount) {
			procPkg.GetLogger().Info("[matching] Bot join triggered by rule (Process)")
			procPkg.GetProcessor().AddBotToMatch(
				procPkg.GetContext(),
				procPkg.GetLogger(),
				procPkg.GetNK(),
				procPkg.GetDb(),
				procPkg.GetDispatcher(),
				state,
				1,
			)
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
