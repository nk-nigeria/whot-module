package state_machine

import (
	"context"
	"time"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	log "github.com/nakamaFramework/whot-module/pkg/log"
	"github.com/nakamaFramework/whot-module/pkg/packager"
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
	state.SetUpCountDown(1 * time.Second)
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
			State: pb.GameState_GameStateMatching,
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
