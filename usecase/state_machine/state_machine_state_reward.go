package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/nk-nigeria/whot-module/usecase/service"
)

type StateReward struct {
	StateBase
	botIntegration *service.WhotBotIntegration
}

func NewStateReward(fn FireFn) *StateReward {
	return &StateReward{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StateReward) Enter(ctx context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[reward] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	// setup reward timeout
	state := procPkg.GetState()
	state.SetUpCountDown(rewardTimeout)

	// Initialize bot integration once
	s.botIntegration = service.NewWhotBotIntegration(procPkg.GetDb())

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GAME_STATE_REWARD,
			CountDown: int64(rewardTimeout.Seconds()),
		},
	)

	// process finish
	procPkg.GetProcessor().ProcessFinishGame(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)

	return nil
}

func (s *StateReward) Exit(ctx context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[reward] exit")
	// clear result
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()

	procPkg.GetProcessor().ProcessKickUserNotInterac(log.GetLogger(), procPkg.GetDispatcher(), state)
	state.ResetMatch()
	return nil
}

func (s *StateReward) Process(ctx context.Context, args ...interface{}) error {
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	if remain := state.GetRemainCountDown(); remain <= 0 {
		s.Trigger(ctx, triggerRewardTimeout)
	} else {
		// log.GetLogger().Info("[reward] not timeout %v", remain)
		if state.IsNeedNotifyCountDown() {
			log.GetLogger().Info("[reward] Process - sending countdown update: %d", remain)
			procPkg.GetProcessor().NotifyUpdateGameState(
				state,
				procPkg.GetLogger(),
				procPkg.GetDispatcher(),
				&pb.UpdateGameState{
					State:     pb.GameState_GAME_STATE_REWARD,
					CountDown: int64(remain),
				},
			)

			state.SetLastCountDown(remain)
			log.GetLogger().Info("[reward] Process - updated lastCountDown to: %d", remain)

			botPresences := state.GetBotPresences()

			if remain == 10 {
				// Step 1: Make initial random decision for bot leave (only once per bot)
				log.GetLogger().Info("[reward] Making initial bot leave decisions at countdown 10")
				for _, botPresence := range botPresences {
					botUserID := botPresence.GetUserId()
					_, exists := state.BotResults[botUserID]
					if !exists {
						continue
					}
					log.GetLogger().Info("[reward] Making random leave decision for bot %s", botUserID)
					s.botIntegration.SetMatchState(
						"",
						int64(state.Label.GetMarkUnit()),
						state.GetPresenceSize(),
						state.BotResults[botUserID], // lastResult
						0,                           // activeTables
					)
					// This will only random once per bot and create pending request
					if err := s.botIntegration.GetBotHelper().ProcessBotLeaveLogic(ctx, botUserID); err != nil {
						log.GetLogger().Error("[reward] Failed to decide bot leave for bot %s: %v", botUserID, err)
					}
				}
			} else if remain < 10 {
				// Step 2: Check and kick bots based on their leave time
				log.GetLogger().Info("[reward] Checking bot kick times at countdown %d", remain)
				for _, botPresence := range botPresences {
					botUserID := botPresence.GetUserId()

					// Check if this bot should be kicked based on time
					kicked, err := s.botIntegration.GetBotHelper().CheckAndKickExpiredBots(ctx, botUserID)
					if err != nil {
						log.GetLogger().Error("[reward] Failed to check/kick bot %s: %v", botUserID, err)
						continue
					}

					if kicked {
						log.GetLogger().Info("[reward] Bot %s was kicked due to expired leave time", botUserID)
						// Bot was already removed in CheckAndKickExpiredBots
						continue
					}
				}
			}
		}
	}

	return nil
}
