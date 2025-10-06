package processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/cgp-common/bot"
	"github.com/nk-nigeria/cgp-common/define"
	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/whot-module/cgbdb"
	"github.com/nk-nigeria/whot-module/constant"
	"github.com/nk-nigeria/whot-module/entity"
	"github.com/nk-nigeria/whot-module/message_queue"
	"github.com/nk-nigeria/whot-module/usecase/engine"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type processor struct {
	engine      engine.UseCase
	marshaler   *proto.MarshalOptions
	unmarshaler *proto.UnmarshalOptions
}

func NewMatchProcessor(marshaler *proto.MarshalOptions, unmarshaler *proto.UnmarshalOptions, engine engine.UseCase) UseCase {
	return &processor{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
		engine:      engine,
	}
}

func (m *processor) ProcessNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {

	m.engine.NewGame(s)
	// chia bài và gửi cho từng người chơi
	if err := m.engine.Deal(s); err == nil {
		for k, v := range s.Cards {

			presence, found := s.PlayingPresences.Get(k)
			if found {
				preCardMsg := &pb.UpdateDeal{
					PresenceCard: &pb.PresenceCards{
						Presence:  k,
						WhotCards: v.WhotCards,
					},
					TopCard:  s.TopCard,
					IdDealer: s.DealerId,
				}
				err := m.broadcastMessage(logger, dispatcher,
					int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), preCardMsg,
					[]runtime.Presence{presence.(runtime.Presence)}, nil, true)
				if err != nil {
					logger.Error("failed to broadcast UpdateDeal for presence %s: %v", k, err)
				}
			}

		}
	}
	// gửi thông tin số lá bài user và bài trên bàn
	cardState := &pb.UpdateCardState{
		Event:            pb.WhotCardEvent_WHOT_EVENT_NONE,
		DeckCount:        m.engine.GetDeckCount(),
		PlayerCardCounts: m.engine.GetPlayerCardCounts(s),
	}

	err := m.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), cardState,
		nil, nil, true)
	if err != nil {
		logger.Error("failed to broadcast CountCard UpdateCardState: %v", err)
	}
	// delay 2s để client chia bài
	s.TurnReadyAt = float64(time.Now().Unix()) + 2
}

func (m *processor) UpdateTurn(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	timeTurn := s.TimeTurn
	s.IsAutoPlay = false
	if entity.BotLoader.IsBot(s.CurrentTurn) {
		botPresence, ok := s.GetPresence(s.CurrentTurn).(*bot.BotPresence)
		if ok {
			// Tạo một botTurn mới (1 lần thực thi sau random ticks)
			logger.Info("Bot %s init turn", s.CurrentTurn)
			botPresence.InitTurnWithOption(bot.TurnOpt{
				MinTick:  constant.TickRate * 1, // 1 giây = 1 * tickRate
				MaxTick:  constant.TickRate * 9, // 9 giây = 9 * tickRate
				MaxOccur: 1,                     // chỉ đánh 1 lần
			},
				func() {
					// Gọi xử lý đánh bài của bot tại đây
					m.HandleAutoPlay(logger, dispatcher, s)
				})
		} else {
			logger.Warn("Failed to cast presence to BotPresence for bot: %s", s.CurrentTurn)
		}
	} else {
		// Nếu là user đang bật trạng thái autoplay thì server tự đánh.
		if s.PresencesNoInteract[s.CurrentTurn] {
			timeTurn = 1
			s.IsAutoPlay = true
		} else {
			s.SetUserNotInteract(s.CurrentTurn, true)
		}

	}
	s.TurnExpireAt = time.Now().Unix() + int64(timeTurn)
	turnUpdate := &pb.UpdateTurn{
		UserId:    s.CurrentTurn,
		Countdown: int64(s.TimeTurn),
	}

	err := m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_TURN), turnUpdate, nil, nil, true)
	if err != nil {
		logger.Error("failed to broadcast UpdateTurn: %v", err)
	}
}

func (m *processor) PlayCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {

	var userID string
	payload := &pb.WhotCard{}

	if s.WaitingForWhotShape {
		return
	}

	if message == nil {
		userID = s.CurrentTurn
		payload = m.engine.FindPlayableCard(s, userID)
		logger.Info("Auto-playing card for user %s: %s", userID, payload.String())
		if payload == nil {
			m.DrawCard(logger, dispatcher, s, nil)
			return
		}
	} else {
		userID = message.GetUserId()
		m.unmarshaler.Unmarshal(message.GetData(), payload)
		logger.Info("User %s played card: %s", userID, payload.String())
	}

	// 1. Gọi engine xử lý đánh bài: cập nhật game state, top card, v.v...
	effect, err := m.engine.PlayCard(s, userID, payload)
	if err != nil {
		logger.Error("engine.Play error for user %s: %v", userID, err)
		return
	}

	logger.Info("User %s played card %s with effect %s", userID, payload.String(), effect)

	// Tạo thông báo cập nhật trạng thái
	cardStateMsg := &pb.UpdateCardState{
		UserId:           userID,
		Event:            pb.WhotCardEvent_WHOT_EVENT_PLAY,
		TopCard:          s.TopCard,
		Effect:           pb.WhotCardEffect(effect),
		PickPenalty:      int32(s.PickPenalty),
		TargetUserId:     s.EffectTarget,
		DeckCount:        m.engine.GetDeckCount(),
		PlayerCardCounts: m.engine.GetPlayerCardCounts(s),
		IsAutoPlay:       message == nil,
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
		cardStateMsg, nil, nil, true,
	)

	if effect == entity.EffectGeneralMarket {
		// Gọi engine xử lý General Market
		logger.Info("Handling General Market for user %s", userID)
		if err := m.engine.HandleGeneralMarket(s, userID); err != nil {
			logger.Error("Failed to handle General Market: %v", err)
			return
		}

		// Thông báo cho từng người chơi về bài mới
		for _, key := range s.PlayingPresences.Keys() {
			otherUserId := key.(string)
			if otherUserId != userID {
				presence, found := s.PlayingPresences.Get(otherUserId)
				if found {
					// Gửi bài mới cho người chơi
					dealMsg := &pb.UpdateDeal{
						PresenceCard: &pb.PresenceCards{
							Presence:  otherUserId,
							WhotCards: s.Cards[otherUserId].WhotCards,
						},
					}

					m.broadcastMessage(
						logger, dispatcher,
						int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), dealMsg,
						[]runtime.Presence{presence.(runtime.Presence)}, nil, true)

					// Thông báo công khai rằng người này đã rút bài
					drawMsg := &pb.UpdateCardState{
						UserId:           otherUserId,
						Event:            pb.WhotCardEvent_WHOT_EVENT_DRAW,
						TopCard:          s.TopCard,
						Effect:           pb.WhotCardEffect_GENERAL_MARKET,
						DeckCount:        m.engine.GetDeckCount(),
						PlayerCardCounts: m.engine.GetPlayerCardCounts(s),
					}

					m.broadcastMessage(
						logger, dispatcher,
						int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
						drawMsg, nil, nil, true,
					)
				}
			}
		}
	}

	if s.CurrentEffect != entity.EffectPickTwo && s.CurrentEffect != entity.EffectPickThree {
		s.CurrentEffect = entity.EffectNone
	}

	if s.IsEndingGame {
		s.GameState = pb.GameState_GAME_STATE_REWARD
		return
	}

	m.UpdateTurn(logger, dispatcher, s)
}

func (m *processor) ChooseWhotShape(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {

	var userID string
	var payload pb.WhotCard

	if message == nil {
		userID = s.CurrentTurn
		payload = *m.engine.ChooseAutomaticWhotShape(s)
	} else {
		userID = message.GetUserId()
		m.unmarshaler.Unmarshal(message.GetData(), &payload)
	}

	m.engine.ChooseWhotShape(s, userID, payload.Suit)

	updateMsg := &pb.UpdateCardState{
		UserId:     userID,
		Event:      pb.WhotCardEvent_WHOT_EVENT_PLAY,
		Effect:     pb.WhotCardEffect_CHOICE_SHAPE_GHOST,
		TopCard:    s.TopCard,
		IsAutoPlay: message == nil,
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
		updateMsg, nil, nil, true,
	)
	m.UpdateTurn(logger, dispatcher, s)
}

func (m *processor) DrawCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	var userID string
	if message == nil {
		userID = s.CurrentTurn
	} else {
		userID = message.GetUserId()
	}

	// Gọi engine xử lý rút bài
	cardsToDraw, err := m.engine.DrawCardsFromDeck(s, userID)
	if err != nil {
		logger.Error("Failed to draw cards: %v", err)
		return
	}

	// Thông báo cho người chơi về bài mới
	playerPresence, found := s.PlayingPresences.Get(userID)
	if found {
		dealMsg := &pb.UpdateDeal{
			PresenceCard: &pb.PresenceCards{
				Presence:  userID,
				WhotCards: s.Cards[userID].WhotCards,
			},
		}
		m.broadcastMessage(logger, dispatcher,
			int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), dealMsg,
			[]runtime.Presence{playerPresence.(runtime.Presence)}, nil, true)
	}

	// Thông báo công khai về việc rút bài
	var pickPenalty int32
	if cardsToDraw != 1 {
		pickPenalty = int32(cardsToDraw)
	} else {
		pickPenalty = 0
		s.CurrentEffect = entity.EffectNone
	}
	drawMsg := &pb.UpdateCardState{
		UserId:           userID,
		Event:            pb.WhotCardEvent_WHOT_EVENT_DRAW,
		TopCard:          s.TopCard,
		PickPenalty:      pickPenalty,
		Effect:           pb.WhotCardEffect_EFFECT_NONE,
		DeckCount:        m.engine.GetDeckCount(),
		PlayerCardCounts: m.engine.GetPlayerCardCounts(s),
		IsAutoPlay:       message == nil,
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
		drawMsg, nil, nil, true,
	)
	if s.IsEndingGame {
		s.GameState = pb.GameState_GAME_STATE_REWARD
		return
	}

	m.UpdateTurn(logger, dispatcher, s)
}

func (m *processor) CheckAndHandleTurnTimeout(ctx context.Context, logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {

	userID := s.CurrentTurn
	// Nếu là bot thì gọi loop của bot
	if entity.BotLoader.IsBot(userID) {
		if bp, ok := s.GetPresence(userID).(*bot.BotPresence); ok {
			bp.Loop()
		}
		return
	}

	// Kiểm tra xem đã hết thời gian lượt của user hay chưa
	if s.TurnExpireAt <= 0 || time.Now().Unix() <= s.TurnExpireAt {
		return
	}

	logger.Info("User %s did not interact in time, auto-playing", userID)
	s.TurnExpireAt = 0
	if s.IsAutoPlay {
		m.HandleAutoPlay(logger, dispatcher, s)
	} else {
		m.DrawCard(logger, dispatcher, s, nil)
	}
}

func (m *processor) HandleAutoPlay(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) bool {
	if s.WaitingForWhotShape {
		m.ChooseWhotShape(logger, dispatcher, s, nil)
		return true
	} else {
		// Thực hiện đánh hộ - thử tìm bài phù hợp để đánh
		userID := s.CurrentTurn
		userCards := s.Cards[userID]
		if userCards != nil && len(userCards.WhotCards) > 0 {
			m.PlayCard(logger, dispatcher, s, nil)
			return true
		}
	}

	return true
}

func (m *processor) ProcessFinishGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {

	logger.Info("process finish game")
	// update finish
	updateFinish := m.engine.Finish(s)

	if updateFinish == nil {
		logger.Error("Finish game failed, no updateFinish data")
		return
	}

	m.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH), updateFinish,
		nil, nil, true)

	balanceResult := m.calcRewardForUserPlaying(ctx, nk, logger, db, dispatcher, s, updateFinish)
	if balanceResult == nil {
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		logger.
			WithField("jackpot game", entity.ModuleName).
			WithField("match id", matchId).
			// WithField("user win jackpot", updateFinish.Jackpot.GetUserId()).
			Error("calc reward failed")
		return
	}
	s.SetBalanceResult(balanceResult)

	m.updateChipByResultGameFinish(ctx, logger, nk, balanceResult) // summary balance user

	m.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_WALLET), balanceResult,
		nil, nil, true,
	)

	logger.Info("process finish game done %v", updateFinish)
}

func (m *processor) AddBotToMatch(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, count int) error {
	if count <= 0 {
		return nil
	}
	bJoin := s.AddBotToMatch(count)

	if len(bJoin) > 0 {
		listUserId := make([]string, 0, len(bJoin))
		for _, p := range bJoin {
			listUserId = append(listUserId, p.GetUserId())
		}
		m.emitNkEvent(ctx, define.NakEventMatchJoin, nk, listUserId, s)
		m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, bJoin, nil)
		matchJson, err := protojson.Marshal(s.Label)
		if err != nil {
			logger.Error("update json label failed ", err)
			return nil
		}
		dispatcher.MatchLabelUpdate(string(matchJson))
		return nil
	}
	return fmt.Errorf("no bot join")
}

func (m *processor) RemoveBotFromMatch(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, botUserID string) error {
	logger.Info("RemoveBotFromMatch %s", botUserID)
	if botUserID == "" {
		return nil
	}

	err, botPresence := s.RemoveBotFromMatch(botUserID)
	if err != nil {
		return err
	}

	// Emit leave event
	listUserId := []string{botUserID}
	m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, listUserId, s)

	// Notify table update
	leavePresences := []runtime.Presence{botPresence}
	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, nil, leavePresences)

	// Update match label
	matchJson, err := protojson.Marshal(s.Label)
	if err != nil {
		logger.Error("update json label failed ", err)
		return nil
	}
	dispatcher.MatchLabelUpdate(string(matchJson))

	logger.Info("Bot %s removed from match", botUserID)
	return nil
}

func (m *processor) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	// tạo json cho logging
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal data to JSON for logging: %v", err)
		jsonData = []byte("{}")
	}
	// tạo byte gửi tới client
	dataByte, err := m.marshaler.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal data: %v", err)
		return err
	}

	err = dispatcher.BroadcastMessage(opCode, dataByte, presences, sender, true)
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s, err : %v", string(jsonData), err)
		return err
	}

	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(jsonData))
	return nil
}

func (m *processor) NotifyUpdateGameState(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
	logger.Info("notify update game state %v", updateState)
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		updateState, nil, nil, true)
}

func (m *processor) NotifyUpdateTable(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
	logger.Info("notify update table data %v", updateState)
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		updateState, nil, nil, true)

}

// caculator amount chips user win or lose on this match
// with amount chip before and after apply reward
// and add jackpot if user win
func (m *processor) calcRewardForUserPlaying(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, updateFinish *pb.UpdateFinish) *pb.BalanceResult {
	listUserId := make([]string, 0, len(updateFinish.ResultWhots))
	for _, uf := range updateFinish.ResultWhots {
		listUserId = append(listUserId, uf.UserId)
	}

	logger.Info("update Chips For User Playing users %v, label %v", listUserId, s.Label)

	wallets, err := m.readWalletUsers(ctx, nk, logger, listUserId...)
	if err != nil {
		updateFinishData, _ := m.marshaler.Marshal(updateFinish)
		logger.
			WithField("users", strings.Join(listUserId, ",")).
			WithField("data", string(updateFinishData)).
			WithField("err", err).
			Error("read wallet error")
		return nil
	}
	mapUserWallet := make(map[string]entity.Wallet)
	for _, w := range wallets {
		mapUserWallet[w.UserId] = w
	}

	balanceResult := pb.BalanceResult{}
	listFeeGame := make([]entity.FeeGame, 0)
	for _, uf := range updateFinish.ResultWhots {
		balance := &pb.BalanceUpdate{
			UserId:           uf.UserId,
			AmountChipBefore: mapUserWallet[uf.UserId].Chips,
		}

		myPrecense, ok := s.GetPresence(uf.GetUserId()).(entity.MyPrecense)
		percentFreeGame := entity.GetFeeGameByLevel(0)
		if ok {
			percentFreeGame = entity.GetFeeGameByLevel(int(myPrecense.VipLevel))
		}
		percentFee := percentFreeGame

		if (uf.WinFactor) > 0 {
			// win
			balance.AmoutChipFee = int64(s.Label.MarkUnit) / 100 * int64(percentFee)
			balance.AmoutChipAddPrefee = int64(uf.WinFactor * float64(s.Label.MarkUnit))
			balance.AmountChipAdd = balance.AmoutChipAddPrefee - balance.AmoutChipFee
			listFeeGame = append(listFeeGame, entity.FeeGame{
				UserID: balance.UserId,
				Fee:    balance.AmoutChipFee,
			})
		} else {
			// lose
			balance.AmountChipAdd = int64(uf.WinFactor * float64(s.Label.MarkUnit))
		}
		balance.AmountChipCurrent = balance.AmountChipBefore + balance.AmountChipAdd
		balanceResult.Updates = append(balanceResult.Updates, balance)
		logger.Info("update user %v, change %s", uf.UserId, balance)
	}
	cgbdb.AddNewMultiFeeGame(ctx, logger, db, listFeeGame)
	return &balanceResult

}

func (m *processor) readWalletUsers(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userIds ...string) ([]entity.Wallet, error) {
	return entity.ReadWalletUsers(ctx, nk, logger, userIds...)
}

func (m *processor) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, balanceResult *pb.BalanceResult) {
	logger.Info("updateChipByResultGameFinish %v", balanceResult)
	walletUpdates := make([]*runtime.WalletUpdate, 0, len(balanceResult.Updates))
	for _, result := range balanceResult.Updates {
		amountChip := result.AmountChipCurrent - result.AmountChipBefore
		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_reward": entity.ModuleName,
		}
		if amountChip > 0 {
			leaderBoardRecord := &pb.CommonApiLeaderBoardRecord{
				GameCode: constant.GameCode,
				UserId:   result.UserId,
				Score:    amountChip,
			}
			message_queue.GetNatsService().Publish(constant.TopicLeaderBoardAddScore, leaderBoardRecord)
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}

	// add chip for user win jackpot
	if jp := balanceResult.Jackpot; jp != nil && jp.UserId != "" {
		changeset := map[string]int64{
			"chips": jp.Chips, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_reward": entity.ModuleName,
			"action":      entity.WalletActionWinGameJackpot,
		}
		wu := &runtime.WalletUpdate{
			UserID:    jp.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		}
		walletUpdates = append(walletUpdates, wu)
	}
	logger.Info("wallet update ctx %v, walletUpdates %v", ctx, walletUpdates)
	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		payload, _ := json.Marshal(walletUpdates)
		logger.
			WithField("payload", string(payload)).
			WithField("err", err).
			Error("Wallets update error.")
		return
	}
}

func (m *processor) notifyUpdateTable(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, s *entity.MatchState, joins, leaves []runtime.Presence) {
	players := entity.NewListPlayer(s.GetPresences())
	players.ReadProfile(ctx, nk, logger)

	playing_players := entity.NewListPlayer(s.GetPlayingPresences())
	// playing_players.ReadWallet(ctx, nk, logger)

	var pjoins, pleaves []*pb.Player
	if joins != nil {
		pjoins = entity.NewListPlayer(joins)
	}

	if leaves != nil {
		pleaves = entity.NewListPlayer(leaves)
	}

	msg := &pb.UpdateTable{
		Bet:            int64(s.Label.Bet.GetMarkUnit()),
		JoinPlayers:    pjoins,
		LeavePlayers:   pleaves,
		Players:        players,
		PlayingPlayers: playing_players,
	}

	// msg.JpTreasure = s.GetJackpotTreasure()
	msg.RemainTime = int64(s.GetRemainCountDown())
	msg.GameState = s.GameState

	m.NotifyUpdateTable(s, logger, dispatcher, msg)
}

func (m *processor) ProcessPresencesJoin(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule, db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
	presences []runtime.Presence,
) {
	logger.Info("process presences join %v", presences)

	newJoins := make([]runtime.Presence, 0, len(presences))
	listUserId := make([]string, 0, len(presences))
	for _, p := range presences {
		uid := p.GetUserId()
		if _, found := s.LeavePresences.Get(uid); found {
			s.RemoveLeavePresence(uid)
		} else if _, found := s.Presences.Get(uid); !found {
			newJoins = append(newJoins, p)
			listUserId = append(listUserId, uid)
		}
	}

	// Cập nhật danh sách người chơi mới
	s.AddPresence(ctx, nk, db, presences)
	s.JoinsInProgress -= len(presences)
	// Cập nhật playing_match vào DB
	// if len(listUserId) > 0 {
	// 	matchID, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
	// 	playingMatch := &pb.PlayingMatch{
	// 		Code:    entity.ModuleName,
	// 		MatchId: matchID,
	// 	}
	// 	if data, err := json.Marshal(playingMatch); err == nil {
	// 		cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, string(data))
	// 	}
	// }

	if len(listUserId) > 0 {
		m.emitNkEvent(ctx, define.NakEventMatchJoin, nk, listUserId, s)
	}

	// Cập nhật bàn chơi
	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, presences, nil)

	switch s.GameState {
	case pb.GameState_GAME_STATE_REWARD:
		// Nếu đã hết game thì gửi balance
		if result := s.GetBalanceResult(); result != nil {
			m.broadcastMessage(
				logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_WALLET),
				result, presences, nil, true,
			)
		}

	case pb.GameState_GAME_STATE_PLAY:
		// Gửi trạng thái số lá bài trên bàn
		cardState := &pb.UpdateCardState{
			Event:            pb.WhotCardEvent_WHOT_EVENT_NONE,
			DeckCount:        m.engine.GetDeckCount(),
			PlayerCardCounts: m.engine.GetPlayerCardCounts(s),
			TopCard:          s.TopCard,
		}
		if err := m.broadcastMessage(
			logger, dispatcher,
			int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
			cardState, presences, nil, true,
		); err != nil {
			logger.Error("failed to broadcast card state: %v", err)
		}

		// Gửi lượt hiện tại
		if s.CurrentTurn != "" {
			turnUpdate := &pb.UpdateTurn{
				UserId:    s.CurrentTurn,
				Countdown: int64(s.TurnExpireAt - time.Now().Unix()),
			}
			if err := m.broadcastMessage(
				logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_TURN),
				turnUpdate, nil, nil, true,
			); err != nil {
				logger.Error("failed to broadcast turn update: %v", err)
			}
		}

		// Gửi bài riêng cho từng người
		for _, p := range presences {
			uid := p.GetUserId()

			if _, found := s.PlayingPresences.Get(uid); !found {
				continue
			}
			s.AddPlayingPresences(p)

			cards := s.Cards[uid]
			if cards == nil || len(cards.WhotCards) == 0 {
				continue
			}

			msg := &pb.UpdateDeal{
				PresenceCard: &pb.PresenceCards{
					Presence:  uid,
					WhotCards: cards.WhotCards,
				},
				TopCard: s.TopCard,
			}
			if err := m.broadcastMessage(
				logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL),
				msg, []runtime.Presence{p}, nil, true,
			); err != nil {
				logger.Error("failed to send cards to %s: %v", uid, err)
			}
		}
	}
}

func (m *processor) ProcessPresencesLeave(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence) {
	logger.Info("process presences leave %v", presences)
	s.RemovePresence(presences...)
	listUserId := make([]string, 0, len(presences))
	for _, p := range presences {
		listUserId = append(listUserId, p.GetUserId())
	}
	m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, listUserId, s)
	// cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, "")
	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, nil, presences)
}

func (m *processor) ProcessPresencesLeavePending(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence) {
	logger.Info("process presences leave pending %v", presences)
	listUserId := make([]string, 0, len(presences))
	for _, presence := range presences {
		_, found := s.PlayingPresences.Get(presence.GetUserId())
		if found {
			s.AddLeavePresence(presence)
		} else {
			s.RemovePresence(presence)
			m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, nil, []runtime.Presence{presence})
			listUserId = append(listUserId, presence.GetUserId())
		}
	}
	if len(listUserId) > 0 {
		m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, listUserId, s)
	}

}

func (m *processor) ProcessApplyPresencesLeave(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
) {
	pendingLeaves := s.GetLeavePresences()

	if len(pendingLeaves) <= 0 {
		return
	}
	logger.Info("process apply presences")

	s.ApplyLeavePresence() // cật nhật lại danh sách người chơi

	listUserId := make([]string, 0, len(pendingLeaves))
	for _, p := range pendingLeaves {
		listUserId = append(listUserId, p.GetUserId())
	}

	m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, listUserId, s)
	// cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, "")

	logger.Info("notify to player kick off %s", strings.Join(listUserId, ","))
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_KICK_OFF_THE_TABLE),
		nil, pendingLeaves, nil, true)

	players := entity.NewListPlayer(s.GetPresences())
	players.ReadProfile(ctx, nk, logger)

	playing_players := entity.NewListPlayer(s.GetPlayingPresences())
	// playing_players.ReadWallet(ctx, nk, logger)

	leaves_player := entity.NewListPlayer(pendingLeaves)

	msg := &pb.UpdateTable{
		Bet:            int64(s.Label.Bet.GetMarkUnit()),
		Players:        players,
		PlayingPlayers: playing_players,
		LeavePlayers:   leaves_player,
	}

	m.NotifyUpdateTable(s, logger, dispatcher, msg)
}

func (m *processor) ProcessKickUserNotInterac(logger runtime.Logger,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
) {
	leaves := s.GetPresenceNotInteract()
	if len(leaves) > 0 {
		err := dispatcher.MatchKick(leaves)
		if err != nil {
			logger.Error("Failed to kick users not interact: %v", err)
			return
		}
		m.broadcastMessage(logger, dispatcher,
			int64(pb.OpCodeUpdate_OPCODE_KICK_OFF_THE_TABLE),
			nil, leaves, nil, true)
	}
}

func (m *processor) ProcessMatchTerminate(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
) {
	listUserId := make([]string, 0, len(s.GetPresences()))
	for _, presence := range s.GetPresences() {
		listUserId = append(listUserId, presence.GetUserId())
	}
	m.emitNkEvent(ctx, define.NakEventMatchEnd, nk, listUserId, s)
}

// check win jackpot, and always get jackpot treasure before exit
// if user win. update jackpot, jackpot history
func (m *processor) handlerJackpotProcess(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule, db *sql.DB,
	s *entity.MatchState,
	updateFinish *pb.UpdateFinish,
	balanceResult *pb.BalanceResult,
) {
	// add chips to jackpot treasure
	defer func() {
		totalChipsWin := int64(0)
		for _, v := range balanceResult.Updates {
			if v.AmountChipAdd > 0 {
				totalChipsWin += v.AmountChipAdd
			}
		}
		totalJpChipTax := totalChipsWin / 100 * entity.JackpotPercentTax
		cgbdb.AddOrUpdateChipJackpot(ctx, logger, db, entity.ModuleName, int64(totalJpChipTax))
	}()
	// update chip if have user win jackpot

	if updateFinish.Jackpot == nil || updateFinish.Jackpot.UserId == "" {
		// no user win
		return
	}
	jackpotTreasure, err := cgbdb.GetJackpot(ctx, logger, db, entity.ModuleName)
	if err != nil {
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		logger.
			WithField("jackpot game", entity.ModuleName).
			WithField("match id", matchId).
			WithField("err", err.Error()).Error("get jackpot treasure error")
		return
	}
	if jackpotTreasure.Chips <= 0 {
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		logger.
			WithField("jackpot game", entity.ModuleName).
			WithField("match id", matchId).
			WithField("user win jackpot", updateFinish.Jackpot.GetUserId()).
			Debug("No chips in jackpot treasure, ignore this win jackpot")
		return
	}
	myPrecense := s.GetPresence(updateFinish.Jackpot.UserId).(entity.MyPrecense)
	// JACKPOT PUSOY
	// Công thức tính tiền max khi JP: JP = MCB x 100 x hệ số Vip
	bet := s.Label.Bet
	vipLv := entity.MaxInt64(myPrecense.VipLevel, 1)
	maxJP := int64(bet.GetMarkUnit()) * 100 * vipLv
	maxJP = entity.MinInt64(maxJP, jackpotTreasure.Chips)
	err = cgbdb.AddOrUpdateChipJackpot(ctx, logger, db, entity.ModuleName, -maxJP)
	if err != nil {
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		logger.
			WithField("jackpot game", entity.ModuleName).
			WithField("match id", matchId).
			WithField("err", err.Error()).Error("update jackpot treasure error")
		return
	}
	updateFinish.Jackpot.Chips = maxJP
	cgbdb.AddJackpotHistoryUserWin(ctx, logger, db, updateFinish.Jackpot.GameCode,
		updateFinish.Jackpot.UserId, -updateFinish.Jackpot.Chips)

}

// read jackpot treasure and set to updateFinish
func (m *processor) readJackpotTreasure(
	ctx context.Context,
	nk runtime.NakamaModule,
	logger runtime.Logger,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
	updateFinish *pb.UpdateFinish,
) {
	updateFinish.JpTreasure = &pb.Jackpot{}
	jpTreasure, _ := cgbdb.GetJackpot(ctx, logger, db, entity.ModuleName)
	if jpTreasure != nil {
		updateFinish.JpTreasure = &pb.Jackpot{
			GameCode: jpTreasure.GetGameCode(),
			Chips:    jpTreasure.Chips,
		}
	}
}

func (m *processor) emitNkEvent(ctx context.Context, eventNk define.NakEvent, nk runtime.NakamaModule, userIds []string, s *entity.MatchState) {

	matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
	gameCode := entity.ModuleName
	endMatchUnix := strconv.FormatInt(time.Now().Unix(), 10)
	mcbValue := strconv.FormatInt(int64(s.Label.Bet.GetMarkUnit()), 10)

	nk.Event(ctx, &api.Event{
		Name:      string(eventNk),
		Timestamp: timestamppb.Now(),
		Properties: map[string]string{
			"user_id":        strings.Join(userIds, ","),
			"game_code":      gameCode,
			"end_match_unix": endMatchUnix,
			"match_id":       matchId,
			"mcb":            mcbValue,
		},
	})
}
