package state_machine

import (
	"context"
	"strings"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	log "github.com/nakamaFramework/whot-module/pkg/log"
	"github.com/nakamaFramework/whot-module/pkg/packager"
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
	state.SetUpCountDown(preparingTimeout)
	// remove all user not interact 2 game continue
	listPrecense := state.GetPresenceNotInteract(2)
	if len(listPrecense) > 0 {
		listUserId := make([]string, len(listPrecense))
		for _, p := range listPrecense {
			listUserId = append(listUserId, p.GetUserId())
		}
		procPkg.GetLogger().Info("Kick %d user from math %s",
			len(listPrecense), strings.Join(listUserId, ","))
		state.AddLeavePresence(listPrecense...)
	}
	procPkg.GetProcessor().ProcessApplyPresencesLeave(ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state,
	)

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: int64(state.GetRemainCountDown()),
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
