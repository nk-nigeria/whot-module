package processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/cgp-common/define"
	pb1 "github.com/nk-nigeria/cgp-common/proto"
	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/cgbdb"
	"github.com/nk-nigeria/whot-module/constant"
	"github.com/nk-nigeria/whot-module/entity"
	"github.com/nk-nigeria/whot-module/message_queue"
	"github.com/nk-nigeria/whot-module/usecase/engine"
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
						Presence: k,
						Cards:    v.Cards,
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
	s.TurnReadyAt = float64(time.Now().Unix()) + 2
}

func (m *processor) UpdateTurn(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {

	s.TurnExpireAt = time.Now().Unix() + int64(s.TimeTurn)
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
	payload := &pb.Card{}

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
		UserId:       userID,
		Event:        pb.CardEvent_PLAY,
		TopCard:      s.TopCard,
		Effect:       pb.CardEffect(effect), // Convert từ entity.CardEffect sang pb.CardEffect
		PickPenalty:  int32(s.PickPenalty),
		TargetUserId: s.EffectTarget,
		IsAutoPlay:   message == nil,
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
		cardStateMsg, nil, nil, true,
	)

	// Xử lý đặc biệt với General Market
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
							Presence: otherUserId,
							Cards:    s.Cards[otherUserId].Cards,
						},
					}

					m.broadcastMessage(
						logger, dispatcher,
						int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), dealMsg,
						[]runtime.Presence{presence.(runtime.Presence)}, nil, true)

					// Thông báo công khai rằng người này đã rút bài
					drawMsg := &pb.UpdateCardState{
						UserId:  otherUserId,
						Event:   pb.CardEvent_DRAW,
						TopCard: s.TopCard,
						Effect:  pb.CardEffect_GENERAL_MARKET,
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

	// Cập nhật người chơi tiếp theo nếu không đang chờ chọn hình Whot
	// if !s.WaitingForWhotShape {

	// }

	// Kiểm tra người chiến thắng
	if s.WinnerId != "" {
		gameStateMsg := &pb.UpdateGameState{
			State: pb.GameState_GameStateReward,
		}

		m.NotifyUpdateGameState(s, logger, dispatcher, gameStateMsg)
	}
	m.UpdateTurn(logger, dispatcher, s)
}

func (m *processor) ChooseWhotShape(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {

	var userID string
	var payload pb.Card

	if message == nil {
		userID = s.CurrentTurn
		payload = *m.engine.ChooseAutomaticWhotShape()
	} else {
		userID = message.GetUserId()
		m.unmarshaler.Unmarshal(message.GetData(), &payload)
	}

	m.engine.ChooseWhotShape(s, userID, payload.Suit)

	updateMsg := &pb.UpdateCardState{
		UserId:     userID,
		Event:      pb.CardEvent_PLAY,
		Effect:     pb.CardEffect_CHOICE_SHAPE_GHOST,
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
				Presence: userID,
				Cards:    s.Cards[userID].Cards,
			},
		}
		m.broadcastMessage(logger, dispatcher,
			int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), dealMsg,
			[]runtime.Presence{playerPresence.(runtime.Presence)}, nil, true)
	}

	// Thông báo công khai về việc rút bài
	drawMsg := &pb.UpdateCardState{
		UserId:      userID,
		Event:       pb.CardEvent_DRAW,
		TopCard:     s.TopCard,
		PickPenalty: int32(cardsToDraw),
		IsAutoPlay:  message == nil,
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE),
		drawMsg, nil, nil, true,
	)

	m.UpdateTurn(logger, dispatcher, s)
}

func (m *processor) CheckAndHandleTurnTimeout(ctx context.Context, logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) bool {
	// Kiểm tra xem đã hết thời gian lượt chưa
	if s.TurnExpireAt <= 0 || time.Now().Unix() <= s.TurnExpireAt {
		return false
	}

	userID := s.CurrentTurn
	logger.Info("Turn timeout for user %s, auto-playing", userID)

	// Tăng số lần đánh hộ liên tiếp
	if s.AutoPlayCounts == nil {
		s.AutoPlayCounts = make(map[string]int)
	}
	s.AutoPlayCounts[userID]++

	if s.WaitingForWhotShape {
		m.ChooseWhotShape(logger, dispatcher, s, nil)
		return true
	} else {
		// Thực hiện đánh hộ - thử tìm bài phù hợp để đánh
		userCards := s.Cards[userID]
		if userCards != nil && len(userCards.Cards) > 0 {
			m.PlayCard(logger, dispatcher, s, nil)
			return true
		}
	}

	return true
}

func (m *processor) ProcessFinishGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	logger.Info("process finish game len cards %v", len(s.Cards))
	pbGameState := pb.UpdateGameState{
		State: pb.GameState_GameStateReward,
	}
	pbGameState.PresenceCards = make([]*pb.PresenceCards, 0, len(s.Cards))
	for k := range s.Cards {

		presenceCards := pb.PresenceCards{
			Presence: k,
		}
		pbGameState.PresenceCards = append(pbGameState.PresenceCards, &presenceCards)
	}

	m.NotifyUpdateGameState(s, logger, dispatcher, &pbGameState)
	// update finish
	updateFinish := m.engine.Finish(s)
	m.readJackpotTreasure(ctx, nk, logger, db, dispatcher, s, updateFinish)
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
	m.handlerJackpotProcess(ctx, logger, nk, db, s, updateFinish, balanceResult)
	// balanceResult.Jackpot = updateFinish.Jackpot
	// read new treasure after update chips win to jp treasure
	m.readJackpotTreasure(ctx, nk, logger, db, dispatcher, s, updateFinish)
	// s.SetJackpotTreasure(updateFinish.JpTreasure)
	m.updateChipByResultGameFinish(ctx, logger, nk, balanceResult) // summary balance ủe
	// summary balance user if win jackpot
	// if updateFinish.Jackpot != nil {
	// 	for _, b := range balanceResult.GetUpdates() {
	// 		if b.GetUserId() == updateFinish.Jackpot.UserId {
	// 			b.AmountChipAdd += updateFinish.Jackpot.Chips
	// 			b.AmountChipCurrent += updateFinish.Jackpot.Chips
	// 			break
	// 		}
	// 	}
	// }
	// s.SetBalanceResult(balanceResult)

	m.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_UNSPECIFIED), balanceResult,
		nil, nil, true,
	)
	m.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH), updateFinish,
		nil, nil, true)
	logger.Info("process finish game done %v", updateFinish)
}

func (m *processor) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	dataByte, err := m.marshaler.Marshal(data)
	if err != nil {
		return err
	}
	err = dispatcher.BroadcastMessage(opCode, dataByte, presences, sender, true)

	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(dataByte))
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataByte))
		return err
	}
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
	// listUserId := make([]string, 0, len(updateFinish.Results))
	// for _, uf := range updateFinish.Results {
	// 	listUserId = append(listUserId, uf.UserId)
	// }

	// logger.Info("update Chips For User Playing users %v, label %v", listUserId, s.Label)

	// wallets, err := m.readWalletUsers(ctx, nk, logger, listUserId...)
	// if err != nil {
	// 	updateFinishData, _ := m.marshaler.Marshal(updateFinish)
	// 	logger.
	// 		WithField("users", strings.Join(listUserId, ",")).
	// 		WithField("data", string(updateFinishData)).
	// 		WithField("err", err).
	// 		Error("read wallet error")
	// 	return nil
	// }
	// mapUserWallet := make(map[string]entity.Wallet)
	// for _, w := range wallets {
	// 	mapUserWallet[w.UserId] = w
	// }

	balanceResult := pb.BalanceResult{}
	// listFeeGame := make([]entity.FeeGame, 0)
	// for _, uf := range updateFinish.Results {
	// 	balance := &pb.BalanceUpdate{
	// 		UserId:           uf.UserId,
	// 		AmountChipBefore: mapUserWallet[uf.UserId].Chips,
	// 	}

	// 	myPrecense, ok := s.GetPresence(uf.GetUserId()).(entity.MyPrecense)
	// 	percentFreeGame := entity.GetFeeGameByLevel(0)
	// 	if ok {
	// 		percentFreeGame = entity.GetFeeGameByLevel(int(myPrecense.VipLevel))
	// 	}
	// 	percentFee := percentFreeGame

	// 	fee := int64(uf.ScoreResult.NumHandWin) * int64(s.Label.Bet) / 100 * int64(percentFee)
	// 	balance.AmountChipAdd = uf.ScoreResult.TotalFactor * int64(s.Label.Bet)
	// 	if (balance.AmountChipAdd) > 0 {
	// 		// win
	// 		balance.AmountChipCurrent = balance.AmountChipBefore + balance.AmountChipAdd - fee
	// 		listFeeGame = append(listFeeGame, entity.FeeGame{
	// 			UserID: balance.UserId,
	// 			Fee:    fee,
	// 		})
	// 	} else {
	// 		// lose
	// 		balance.AmountChipCurrent = balance.AmountChipBefore + balance.AmountChipAdd
	// 	}
	// 	balanceResult.Updates = append(balanceResult.Updates, balance)
	// 	// logger.Info("update user %v, fee %d change %s", uf.UserId, fee, balance)
	// }
	// cgbdb.AddNewMultiFeeGame(ctx, logger, db, listFeeGame)
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
			leaderBoardRecord := &pb1.CommonApiLeaderBoardRecord{
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
	// players.ReadProfile(ctx, nk, logger)

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
	{
		// mapPlayging := make(map[string]bool, 0)

		// for _, p := range msg.Players {
		// 	// check playing
		// 	mapUserPlaying := s.PlayingPresences
		// 	_, p.IsPlaying = mapUserPlaying.Get(p.GetId())
		// 	// status hold card
		// 	if _, exist := s.OrganizeCards[p.GetId()]; exist {
		// 		p.CardStatus = pb.CardStatus(pb.CardEvent_DRAW)
		// 		// p.Cards = s.OrganizeCards[p.GetId()]
		// 	} else {
		// 		p.CardStatus = pb.CardStatus(pb.CardEvent_DRAW)
		// 	}
		// }
	}
	msg.JpTreasure = s.GetJackpotTreasure()
	msg.RemainTime = int64(s.GetRemainCountDown())
	msg.GameState = s.GameState

	m.NotifyUpdateTable(s, logger, dispatcher, msg)
}

func (m *processor) ProcessPresencesJoin(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule, db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
	presences []runtime.Presence,
) {
	logger.Info("process presences join %v", presences)
	// update new presence
	newJoins := make([]runtime.Presence, 0)

	for _, presence := range presences {
		// check in list leave pending
		{
			_, found := s.LeavePresences.Get(presence.GetUserId())
			if found {
				s.LeavePresences.Remove(presence.GetUserId())
			} else {
				newJoins = append(newJoins, presence)
			}
		}
	}

	s.AddPresence(ctx, nk, newJoins)
	s.JoinsInProgress -= len(newJoins)

	// update match profile user
	{
		var listUserId []string
		for _, p := range newJoins {
			listUserId = append(listUserId, p.GetUserId())
		}
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		playingMatch := &pb1.PlayingMatch{
			Code:    entity.ModuleName,
			MatchId: matchId,
		}
		playingMatchJson, _ := json.Marshal(playingMatch)
		cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, string(playingMatchJson))
	}

	for _, presence := range newJoins {
		m.emitNkEvent(ctx, define.NakEventMatchJoin, nk, presence.GetUserId(), s)
	}

	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, presences, nil)
	//send cards for player rejoin
	for _, presence := range presences {
		if _, found := s.PlayingPresences.Get(presence.GetUserId()); found {
			card := s.Cards[presence.GetUserId()]
			if card == nil {
				continue
			}
			dealMsg := &pb.UpdateDeal{
				PresenceCard: &pb.PresenceCards{
					Presence: presence.GetUserId(),
					Cards:    card.Cards,
				},
				TopCard: s.TopCard,
			}
			m.broadcastMessage(
				logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), dealMsg,
				[]runtime.Presence{presence}, nil, true)
		}
	}
	// send update wallet for new user join
	switch s.GameState {
	case pb.GameState_GameStateReward:
		{
			balanceResult := s.GetBalanceResult()
			if balanceResult != nil {
				m.broadcastMessage(
					logger,
					dispatcher,
					int64(pb.OpCodeUpdate_OPCODE_UPDATE_WALLET),
					balanceResult,
					presences,
					nil,
					true,
				)
			}
		}
	default:
		{
		}
	}
}

func (m *processor) ProcessPresencesLeave(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence) {
	logger.Info("process presences leave %v", presences)
	s.RemovePresence(presences...)
	var listUserId []string
	for _, p := range presences {
		listUserId = append(listUserId, p.GetUserId())
		m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, p.GetUserId(), s)
	}
	cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, "")
	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, nil, presences)
}

func (m *processor) ProcessPresencesLeavePending(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, s *entity.MatchState, presences []runtime.Presence) {
	logger.Info("process presences leave pending %v", presences)
	for _, presence := range presences {
		_, found := s.PlayingPresences.Get(presence.GetUserId())
		if found {
			s.AddLeavePresence(presence)
		} else {
			s.RemovePresence(presence)
			m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, nil, []runtime.Presence{presence})
		}
		m.emitNkEvent(ctx, define.NakEventMatchLeave, nk, presence.GetUserId(), s)
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

	if len(pendingLeaves) == 0 {
		return
	}
	logger.Info("process apply presences")

	s.RemovePresence(pendingLeaves...)

	if len(pendingLeaves) > 0 {
		listUserId := make([]string, 0)
		for _, p := range pendingLeaves {
			listUserId = append(listUserId, p.GetUserId())
		}
		cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, "")
		logger.Info("notify to player kick off %s", strings.Join(listUserId, ","))
		m.broadcastMessage(
			logger, dispatcher,
			int64(pb.OpCodeUpdate_OPCODE_KICK_OFF_THE_TABLE),
			nil, pendingLeaves, nil, true)
	}
	s.ApplyLeavePresence()

	players := entity.NewListPlayer(s.GetPresences())
	// players.ReadWallet(ctx, nk, logger)

	playing_players := entity.NewListPlayer(s.GetPlayingPresences())
	// playing_players.ReadWallet(ctx, nk, logger)

	msg := &pb.UpdateTable{
		Bet:            int64(s.Label.Bet.GetMarkUnit()),
		Players:        players,
		PlayingPlayers: playing_players,
		JpTreasure:     s.GetJackpotTreasure(),
	}

	m.NotifyUpdateTable(s, logger, dispatcher, msg)
}

func (m *processor) ProcessMatchTerminate(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	s *entity.MatchState,
) {
	for _, presence := range s.GetPresences() {
		m.emitNkEvent(ctx, define.NakEventMatchEnd, nk, presence.GetUserId(), s)
	}
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

	// if updateFinish.Jackpot == nil || updateFinish.Jackpot.UserId == "" {
	// 	// no user win
	// 	return
	// }
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

func (m *processor) emitNkEvent(ctx context.Context, eventNk define.NakEvent, nk runtime.NakamaModule, userId string, s *entity.MatchState) {
	matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
	nk.Event(ctx, &api.Event{
		Name:      string(eventNk),
		Timestamp: timestamppb.Now(),
		Properties: map[string]string{
			"user_id":        userId,
			"game_code":      entity.ModuleName,
			"end_match_unix": strconv.FormatInt(time.Now().Unix(), 10),
			"match_id":       matchId,
			"mcb":            strconv.FormatInt(int64(s.Label.Bet.GetMarkUnit()), 10),
		},
	})
}
