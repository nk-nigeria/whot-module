package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
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
			CountDown: int64(state.GetRemainCountDown()),
		},
	)
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
	} else {
		// log.GetLogger().Info("state idle presences size %v", presenceCount)
	}

	return nil
}
