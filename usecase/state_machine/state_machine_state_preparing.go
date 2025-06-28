package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
)

type StatePreparing struct {
	StateBase
}

func NewStatePreparing(fn FireFn) *StatePreparing {
	return &StatePreparing{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StatePreparing) Enter(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[preparing] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.GetLogger().Info("state %v", state.Presences)

	procPkg.GetProcessor().ProcessApplyPresencesLeave(ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state,
	)
	state.SetUpCountDown(preparingTimeout)
	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: int64(preparingTimeout.Seconds()),
		},
	)

	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[preparing] exit")
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
	// log.GetLogger().Info("[preparing] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	if remain := state.GetRemainCountDown(); remain > 0 {
		if state.IsNeedNotifyCountDown() {
			procPkg.GetProcessor().NotifyUpdateGameState(
				state,
				procPkg.GetLogger(),
				procPkg.GetDispatcher(),
				&pb.UpdateGameState{
					State:     pb.GameState_GameStatePreparing,
					CountDown: int64(remain),
				},
			)

			state.SetLastCountDown(remain)
		}
	} else {
		// check preparing condition
		// log.GetLogger().Info("[preparing] preparing timeout check presence count")
		if state.GetPresenceNotBotSize() == 0 {
			s.Trigger(ctx, triggerIdle)
			log.GetLogger().Info("[preparing] no one presence, trigger idle")
			return nil
		}
		if state.IsReadyToPlay() {
			// change to play
			s.Trigger(ctx, triggerPreparingDone)
		} else {
			// change to wait
			s.Trigger(ctx, triggerPreparingFailed)
		}
	}

	return nil
}
