package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto"
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
			State:     pb.GameState_GAME_STATE_PREPARING,
			CountDown: initialCountdown,
		},
	)

	// Use global bot integration from matching state instead of creating new one
	botIntegration := GetGlobalBotIntegration()

	if botIntegration == nil {
		// Fallback: create new instance if global one is not available
		botIntegration = service.NewWhotBotIntegration(procPkg.GetDb())
		log.GetLogger().Info("[preparing] Created new bot integration (global not available)")
	} else {
		log.GetLogger().Info("[preparing] Using global bot integration from matching state")
	}

	// Check for pending join requests from matching state
	if state.GetPresenceSize() < entity.MaxPresences {

		// Check if there are pending join requests from matching state
		remainingTime, hasRequest := botIntegration.GetBotHelper().GetBotJoinRemainingTime(ctx)
		if hasRequest {
			log.GetLogger().Info("[preparing] Enter - Found pending join request from matching with remaining time: %v", remainingTime)
		} else {
			log.GetLogger().Info("[preparing] Enter - No pending join request from matching")
			SetGlobalBotIntegration(nil)
		}

		// Don't execute pending requests here, let Process() handle it
	} else {
		procPkg.GetLogger().Info("[preparing] Skip bot join - maximum players reached (%d)", entity.MaxPresences)
		SetGlobalBotIntegration(nil)
	}

	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[preparing] exit")
	// Clear global bot integration
	if GetGlobalBotIntegration() != nil {
		SetGlobalBotIntegration(nil)
	}
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
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
					State:     pb.GameState_GAME_STATE_PREPARING,
					CountDown: int64(remain),
				},
			)

			state.SetLastCountDown(remain)
			log.GetLogger().Info("[preparing] Process - updated lastCountDown to: %d", remain)
		} else {
			log.GetLogger().Info("[preparing] Process - skipping countdown update (no change)")
		}
		if GetGlobalBotIntegration() == nil {
			log.GetLogger().Info("[preparing] Skip bot join - global bot integration not available")
			return nil // skip bot join
		}

		// Check bot join only if we haven't reached maximum players (4)
		if state.GetPresenceSize() < entity.MaxPresences {

			// Check if bot should join based on time
			botCtx := packager.GetContextWithProcessorPackager(procPkg)
			joined, err := GetGlobalBotIntegration().CheckAndJoinExpiredBots(botCtx)
			if err != nil {
				procPkg.GetLogger().Error("[preparing] Bot join error: %v", err)
			} else if joined {
				procPkg.GetLogger().Info("[preparing] Bot joined based on time")
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
