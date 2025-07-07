package state_machine

import (
	"context"

	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/nk-nigeria/whot-module/usecase/service"
)

type StateReward struct {
	StateBase
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

	// process finish
	procPkg.GetProcessor().ProcessFinishGame(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)

	// Check if bots should leave after game ends
	botIntegration := service.NewWhotBotIntegration(procPkg.GetDb())
	betAmount := int64(state.Label.GetMarkUnit())

	// Get last game result (1 for win, -1 for lose, 0 for draw)
	lastResult := 0 // Default to draw
	if state.GetBalanceResult() != nil {
		// Determine result based on balance - this is a simplified logic
		// You can implement more sophisticated logic based on your game rules
		balanceResult := state.GetBalanceResult()
		if balanceResult.GetUpdates()[0].GetAmountChipAdd() > 0 {
			lastResult = 1 // Win
		} else {
			lastResult = -1 // Lose
		}
		// If both are 0, it's a draw (lastResult = 0)
	}

	// Update match state for bot decision
	botIntegration.SetMatchState(
		"", // matchID not needed for leave logic
		betAmount,
		state.GetPresenceSize(),
		lastResult,
		0, // activeTables
	)

	if botIntegration.GetBotHelper().ShouldBotLeave(procPkg.GetContext()) {
		log.GetLogger().Info("[reward] Bot leave triggered by rule after game end, result=%d", lastResult)
		// Remove one bot from the match
		botPresences := state.GetBotPresences()
		if len(botPresences) > 0 {
			// Remove the first bot
			botToRemove := botPresences[0]
			state.RemovePresence(botToRemove)
			log.GetLogger().Info("[reward] Removed bot %s from match", botToRemove.GetUserId())

			// Free the bot back to the pool
			botIntegration.GetBotHelper().FreeBot(botToRemove.GetUserId())
		}
	} else {
		log.GetLogger().Info("[reward] No bot leave triggered, result=%d", lastResult)
	}

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
	}

	return nil
}
