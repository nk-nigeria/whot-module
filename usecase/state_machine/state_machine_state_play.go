package state_machine

import (
	"context"
	"time"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	log "github.com/nk-nigeria/whot-module/pkg/log"
	"github.com/nk-nigeria/whot-module/pkg/packager"
)

type StatePlay struct {
	StateBase
}

func NewStatePlay(fn FireFn) *StatePlay {
	return &StatePlay{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StatePlay) Enter(ctx context.Context, agrs ...interface{}) error {
	log.GetLogger().Info("[play] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	// Setup count down
	// state.SetUpCountDown(playTimeout)

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State: pb.GameState_GameStatePlay,
			// CountDown: int64(state.GetRemainCountDown()),
		},
	)
	// thêm user đang chờ matching sang playing
	state.SetupMatchPresence()

	// New game here
	procPkg.GetProcessor().ProcessNewGame(procPkg.GetLogger(), procPkg.GetDispatcher(), state)

	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[play] exit")
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[play] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)

	state := procPkg.GetState()
	// if remain := state.GetRemainCountDown(); remain > 0 {
	// log.GetLogger().Info("[play] not timeout %v, message %v", remain, procPkg.GetMessages())
	messages := procPkg.GetMessages()
	processor := procPkg.GetProcessor()
	logger := procPkg.GetLogger()
	dispatcher := procPkg.GetDispatcher()
	if state.TurnReadyAt > 0 && float64(time.Now().Unix()) >= state.TurnReadyAt {
		processor.UpdateTurn(logger, dispatcher, state)
		state.TurnReadyAt = 0
	}
	processor.CheckAndHandleTurnTimeout(ctx, logger, dispatcher, state)

	for _, message := range messages {
		switch pb.OpCodeRequest(message.GetOpCode()) {
		case pb.OpCodeRequest_OPCODE_REQUEST_PLAY_CARD:
			logger.Info("User %s play card", message.GetUserId())
			processor.PlayCard(logger, dispatcher, state, message)
		case pb.OpCodeRequest_OPCODE_REQUEST_DRAW_CARD:
			logger.Info("User %s draw card", message.GetUserId())
			processor.DrawCard(logger, dispatcher, state, message)
		case pb.OpCodeRequest_OPCODE_REQUEST_CALL_WHOT:
			logger.Info("User %s call whot", message.GetUserId())
			processor.ChooseWhotShape(logger, dispatcher, state, message)
		case pb.OpCodeRequest_OPCODE_USER_INTERACT_CARDS:
			logger.Info("User %s interact with card", message.GetUserId())
			state.ResetUserNotInteract(message.GetUserId())
		}
	}

	// log.GetLogger().Info("[play] not timeout show %v, play %v", state.GetShowCardCount(), state.GetPlayingCount())
	if state.GameState == pb.GameState_GameStateReward {
		s.Trigger(ctx, triggerPlayTimeout)
	}
	// } else {
	// 	log.GetLogger().Info("[play] timeout reach %v", remain)
	// 	s.Trigger(ctx, triggerPlayTimeout)
	// }
	return nil
}
