package state_machine

import (
	"context"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	log "github.com/nakamaFramework/whot-module/pkg/log"
	"github.com/nakamaFramework/whot-module/pkg/packager"
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
	state.SetUpCountDown(playTimeout)

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePlay,
			CountDown: int64(state.GetRemainCountDown()),
		},
	)
	// Setup match presences
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
	// log.GetLogger().Info("[play] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	if remain := state.GetRemainCountDown(); remain > 0 {
		// log.GetLogger().Info("[play] not timeout %v, message %v", remain, procPkg.GetMessages())
		messages := procPkg.GetMessages()
		processor := procPkg.GetProcessor()
		logger := procPkg.GetLogger()
		dispatcher := procPkg.GetDispatcher()
		for _, message := range messages {
			switch pb.OpCodeRequest(message.GetOpCode()) {
			case pb.OpCodeRequest_OPCODE_REQUEST_COMBINE_CARDS:
				processor.PlayCard(logger, dispatcher, state, message)
			case pb.OpCodeRequest_OPCODE_REQUEST_SHOW_CARDS:
				processor.ShowCard(logger, dispatcher, state, message)
			case pb.OpCodeRequest_OPCODE_REQUEST_DECLARE_CARDS:
				processor.DeclareCard(logger, dispatcher, state, message)
				state.ResetUserNotInteract(message.GetUserId())
			case pb.OpCodeRequest_OPCODE_USER_INTERACT_CARDS:
				logger.Info("User %s interact with card", message.GetUserId())
				state.ResetUserNotInteract(message.GetUserId())
			}
		}

		// log.GetLogger().Info("[play] not timeout show %v, play %v", state.GetShowCardCount(), state.GetPlayingCount())
		// Check all user show card
		if state.GetShowCardCount() >= state.GetPlayingCount() {
			s.Trigger(ctx, triggerPlayCombineAll)
		}
	} else {
		log.GetLogger().Info("[play] timeout reach %v", remain)
		s.Trigger(ctx, triggerPlayTimeout)
	}
	return nil
}
