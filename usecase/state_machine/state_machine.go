package state_machine

import (
	"context"
	"time"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/qmuntal/stateless"
)

const (
	stateInit      = "Init"      // Khởi tạo
	stateIdle      = "Idle"      // Chờ người chơi đầu tiên vào bàn
	stateMatching  = "Matching"  // Chờ đủ người chơi vào bàn
	statePreparing = "Preparing" // Đếm ngược và chuẩn bị bắt đầu ván chơi
	statePlay      = "Play"      // Chơi bài, chờ người chơi đánh bài
	stateReward    = "Reward"    // Tính điểm và trả thưởng cho người chơi
	stateFinish    = "Finish"    // Kết thúc ván chơi, trả về trạng thái ban đầu
)

const (
	triggerIdle            = "GameIdle"
	triggerMatching        = "GameMatching"
	triggerPresenceReady   = "GamePresenceReady"
	triggerPreparingDone   = "GamePreparingDone"
	triggerPreparingFailed = "GamePreparingFailed"
	triggerPlayTimeout     = "GamePlayTimeout"
	triggerRewardTimeout   = "GameRewardTimeout"
	triggerNoOne           = "GameNoOne"

	triggerProcess = "GameProcess"
)

const (
	idleTimeout      = time.Second * 30 // Thời gian chờ không có người chơi vào bàn
	matchingTimeout  = time.Second * 10 // Thời gian chờ đủ người chơi vào bàn
	preparingTimeout = time.Second * 10 // Thời gian chờ đếm ngược và chuẩn bị bắt đầu ván chơi
	//playTimeout      = time.Second * 10
	rewardTimeout = time.Second * 20 // Thời gian chờ tính điểm và trả thưởng cho người chơi
)

type Machine struct {
	state *stateless.StateMachine
}

func (m *Machine) configure() {
	fireCtx := m.state.FireCtx

	// init state
	m.state.Configure(stateInit).
		Permit(triggerIdle, stateIdle)
	m.state.OnTransitioning(func(ctx context.Context, t stateless.Transition) {
		procPkg := packager.GetProcessorPackagerFromContext(ctx)
		state := procPkg.GetState()
		switch t.Destination {
		case stateInit:
			state.GameState = pb.GameState_GameStateIdle
		case stateIdle:
			{
				state.GameState = pb.GameState_GameStateIdle
			}
		case stateMatching:
			{
				state.GameState = pb.GameState_GameStateMatching
			}
		case statePreparing:
			{
				state.GameState = pb.GameState_GameStatePreparing
			}
		case statePlay:
			{
				state.GameState = pb.GameState_GameStatePlay
			}
		case stateReward:
			{
				state.GameState = pb.GameState_GameStateReward
			}
		case stateFinish:
			{
				state.GameState = pb.GameState_GameStateFinish
			}
		}
	})

	// idle state: wait for first user, check no one and timeout
	idle := NewIdleState(fireCtx)
	m.state.Configure(stateIdle).
		OnEntry(idle.Enter).
		OnExit(idle.Exit).
		InternalTransition(triggerProcess, idle.Process).
		Permit(triggerMatching, stateMatching).
		Permit(triggerNoOne, stateFinish)

	// matching state: wait for reach min user => switch to preparing, check no one and timeout => switch to idle
	matching := NewStateMatching(fireCtx)
	m.state.Configure(stateMatching).
		OnEntry(matching.Enter).
		OnExit(matching.Exit).
		InternalTransition(triggerProcess, matching.Process).
		Permit(triggerPresenceReady, statePreparing).
		Permit(triggerIdle, stateIdle)

	// preparing state: wait for reach min user in duration => switch to play, check not enough and timeout => switch to idle
	preparing := NewStatePreparing(fireCtx)
	m.state.Configure(statePreparing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcess, preparing.Process).
		Permit(triggerPreparingDone, statePlay).
		Permit(triggerPreparingFailed, stateMatching)

	// playing state: wait for all user show card or timeout =>
	//  switch to reward
	play := NewStatePlay(fireCtx)
	m.state.Configure(statePlay).
		OnEntry(play.Enter).
		OnExit(play.Exit).
		InternalTransition(triggerProcess, play.Process).
		Permit(triggerPlayTimeout, stateReward)

	// reward state: wait for reward timeout => switch to
	reward := NewStateReward(fireCtx)
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcess, reward.Process).
		Permit(triggerRewardTimeout, stateMatching)

	m.state.ToGraph()
}

func (m *Machine) FireProcessEvent(ctx context.Context, args ...interface{}) error {
	return m.state.FireCtx(ctx, triggerProcess, args...)
}

func (m *Machine) MustState() stateless.State {
	return m.state.MustState()
}

func (m *Machine) GetPbState() pb.GameState {
	switch m.state.MustState() {
	case stateIdle:
		return pb.GameState_GameStateIdle
	case stateMatching:
		return pb.GameState_GameStateMatching
	case statePreparing:
		return pb.GameState_GameStatePreparing
	case statePlay:
		return pb.GameState_GameStatePlay
	case stateReward:
		return pb.GameState_GameStateReward
	default:
		return pb.GameState_GameStateUnknown
	}
}

func NewGameStateMachine() UseCase {
	gs := &Machine{
		state: stateless.NewStateMachine(stateInit),
	}

	gs.configure()

	return gs
}

func (m *Machine) IsMatchingState() bool {
	return m.MustState() == stateMatching
}

func (m *Machine) IsPreparingState() bool {
	return m.MustState() == statePreparing
}

func (m *Machine) IsPlayingState() bool {
	return m.MustState() == statePlay
}

func (m *Machine) IsReward() bool {
	return m.MustState() == stateReward
}

func (m *Machine) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return m.state.FireCtx(ctx, trigger, args...)
}

func (m *Machine) TriggerIdle(ctx context.Context, args ...interface{}) error {
	return m.state.FireCtx(ctx, triggerIdle, args...)
}

type FireFn func(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error

type StateBase struct {
	fireFn FireFn
}

func (s *StateBase) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return s.fireFn(ctx, trigger, args...)
}
