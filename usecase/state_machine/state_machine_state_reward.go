package state_machine

import (
	"context"

	log "github.com/nakamaFramework/whot-module/pkg/log"
	"github.com/nakamaFramework/whot-module/pkg/packager"
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

	return nil
}

func (s *StateReward) Exit(ctx context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[reward] exit")
	// clear result
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.ResetBalanceResult()
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
