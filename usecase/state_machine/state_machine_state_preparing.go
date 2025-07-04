package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/entity"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/nk-nigeria/whot-module/usecase/service"
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
	// thêm user đang chờ matching sang playing
	state.SetupMatchPresence()

	state.SetUpCountDown(preparingTimeout)
	initialCountdown := int64(preparingTimeout.Seconds())
	log.GetLogger().Info("[preparing] Enter - sending initial countdown: %d", initialCountdown)

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: initialCountdown,
		},
	)

	// Check bot join only if we haven't reached maximum players (4)
	if state.GetPresenceSize() < entity.MaxPresences {
		botService := service.NewBotManagementService(procPkg.GetDb())
		betAmount := int64(state.Label.GetMarkUnit())
		userCount := state.GetPresenceNotBotSize()
		if botService.ShouldBotJoin(procPkg.GetContext(), betAmount, userCount) {
			procPkg.GetLogger().Info("[preparing] Bot join triggered by rule")
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
		procPkg.GetLogger().Info("[preparing] Skip bot join - maximum players reached (%d)", entity.MaxPresences)
	}

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
		debugInfo := state.DebugCountDown()
		log.GetLogger().Info("[preparing] Process - countdown debug: %+v", debugInfo)

		if state.IsNeedNotifyCountDown() {
			log.GetLogger().Info("[preparing] Process - sending countdown update: %d", remain)
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
			log.GetLogger().Info("[preparing] Process - updated lastCountDown to: %d", remain)
		} else {
			log.GetLogger().Info("[preparing] Process - skipping countdown update (no change)")
		}

		// Check bot join only if we haven't reached maximum players (4)
		if state.GetPresenceSize() < entity.MaxPresences {
			// Sử dụng BotManagementService để thêm bot thông minh trong preparing phase
			botService := service.NewBotManagementService(procPkg.GetDb())
			betAmount := int64(state.Label.GetMarkUnit())
			userCount := state.GetPresenceNotBotSize()
			if botService.ShouldBotJoin(procPkg.GetContext(), betAmount, userCount) {
				procPkg.GetLogger().Info("[preparing] Bot join triggered by rule (Process)")
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
			procPkg.GetLogger().Info("[preparing] Skip bot join - maximum players reached (%d)", entity.MaxPresences)
		}
	} else {
		// check preparing condition
		log.GetLogger().Info("[preparing] Process - countdown finished, remain: %d", remain)
		// log.GetLogger().Info("[preparing] preparing timeout check presence count")
		if state.IsReadyToPlay() {
			// change to play
			log.GetLogger().Info("[preparing] Process - ready to play, triggering preparingDone")
			s.Trigger(ctx, triggerPreparingDone)
		} else {
			// change to wait
			log.GetLogger().Info("[preparing] Process - not ready to play, triggering preparingFailed")
			s.Trigger(ctx, triggerPreparingFailed)
		}
	}

	return nil
}
