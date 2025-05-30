package state_machine

import (
	"context"

	pb "github.com/nakama-nigeria/cgp-common/proto/whot"
	"github.com/qmuntal/stateless"
)

type UseCase interface {
	FireProcessEvent(ctx context.Context, args ...interface{}) error
	MustState() stateless.State
	GetPbState() pb.GameState
	IsPlayingState() bool
	IsReward() bool
	Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error
	TriggerIdle(ctx context.Context, args ...interface{}) error
}
