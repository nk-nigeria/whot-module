package processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/heroiclabs/nakama-common/runtime"
	pb1 "github.com/nakamaFramework/cgp-common/proto"
	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/cgbdb"
	"github.com/nakamaFramework/whot-module/constant"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/nakamaFramework/whot-module/message_queue"
	"github.com/nakamaFramework/whot-module/usecase/engine"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type processor struct {
	engine      engine.UseCase
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
}

func NewMatchProcessor(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions, engine engine.UseCase) UseCase {
	return &processor{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
		engine:      engine,
	}
}

// Call when client request or timeout
func (m *processor) ProcessNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	// clean up game state
	m.engine.NewGame(s)

	if err := m.engine.Deal(s); err == nil {
		for k, v := range s.Cards {
			buf, err := m.marshaler.Marshal(&pb.UpdateDeal{
				PresenceCard: &pb.PresenceCards{
					Presence: k,
					Cards:    v.Cards,
				},
			})

			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				presence, found := s.PlayingPresences.Get(k)
				if found {
					_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), buf, []runtime.Presence{presence.(runtime.Presence)}, nil, true)
				}
			}
		}
	}
}

func (m *processor) ProcessFinishGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	logger.Info("process finish game len cards %v", len(s.Cards))
	// send organize card to all
	pbGameState := pb.UpdateGameState{
		State: pb.GameState_GameStateReward,
	}
	pbGameState.PresenceCards = make([]*pb.PresenceCards, 0, len(s.Cards))
	for k, v := range s.Cards {
		organizeCards := s.OrganizeCards[k]
		if organizeCards == nil {
			logger.Warn("user %s not submit cards use deal cards", k)
			organizeCards = v
			s.OrganizeCards[k] = v
		}
		presenceCards := pb.PresenceCards{
			Presence: k,
			Cards:    organizeCards.GetCards(),
		}
		pbGameState.PresenceCards = append(pbGameState.PresenceCards, &presenceCards)
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		&pbGameState, nil, nil, true)

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
	m.broadcastMessage(
		logger,
		dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_UNSPECIFIED),
		balanceResult,
		nil, nil, true,
	)
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
		updateFinish, nil, nil, true)
	logger.Info("process finish game done %v", updateFinish)
}

func (m *processor) PlayCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	// logger.Info("User %s request play card", message.GetUserId())
	// msg := pb.UpdateGameState{
	// 	State: pb.GameState_GameStatePlay,
	// 	ArrangeCard: &pb.ArrangeCard{
	// 		Presence:  message.GetUserId(),
	// 		CardEvent: pb.CardEvent_PLAY,
	// 	},
	// }
	// m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	// m.engine.Play(s, message.GetUserId(), message.GetData())
}

func (m *processor) ChooseWhotShape(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	// if !s.WaitingForWhotShapeChoice {
	// 	return errors.New("not waiting for whot shape")
	// }
	// s.WhotChosenShape = shape
	// s.WaitingForWhotShapeChoice = false
	// return nil
}

func (m *processor) DrawCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {

}

func (m *processor) CombineCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request combineCard", message.GetUserId())
	msg := pb.UpdateGameState{
		State: pb.GameState_GameStatePlay,
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_COMBINE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.removeShowCard(logger, s, message)
}

func (m *processor) ShowCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request showCard", message.GetUserId())

	msg := pb.UpdateGameState{
		State: pb.GameState_GameStatePlay,
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_SHOW,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.saveCard(logger, s, message)
}

func (m *processor) DeclareCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request declareCard", message.GetUserId())
	// TODO: check royalties
	m.saveCard(logger, s, message)
	msg := pb.UpdateGameState{
		State: pb.GameState_GameStatePlay,
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_DECLARE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)

}

func (m *processor) saveCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
	cards := s.Cards[message.GetUserId()]
	organize := &pb.Organize{}
	err := m.unmarshaler.Unmarshal(message.GetData(), organize)
	if err != nil {
		logger.Error("Parse organize cards from client error %s", err.Error())
		return
	}
	cardsByClient := organize.Cards
	// check len card
	if len(cardsByClient.GetCards()) != len(cards.GetCards()) {
		logger.Error("Amount cards from client [%d] different amount card in server [%d]",
			len(cardsByClient.GetCards()), len(cards.GetCards()))
		return
	}
	// check card send by client is the same card in server
	if !entity.IsSameListCard(entity.NewListCard(cards.GetCards()), entity.NewListCard(cardsByClient.GetCards())) {
		logger.Error("cards from client not the same card in server, invalid action",
			len(cardsByClient.GetCards()), len(cards.GetCards()))
		return
	}

	logger.Info("update save card %v, %v", message.GetUserId(), cardsByClient)

	m.engine.Organize(s, message.GetUserId(), cardsByClient)
}

func (m *processor) removeShowCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
	m.engine.Combine(s, message.GetUserId())
}

func (m *processor) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	dataJson, err := m.marshaler.Marshal(data)
	if err != nil {
		logger.Error("Error when marshaler data for broadcastMessage")
		return err
	}
	err = dispatcher.BroadcastMessage(opCode, dataJson, presences, sender, true)

	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(dataJson))
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataJson))
		return err
	}
	return nil
}

func (m *processor) NotifyUpdateGameState(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
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
		Bet:            int64(s.Label.Bet),
		JoinPlayers:    pjoins,
		LeavePlayers:   pleaves,
		Players:        players,
		PlayingPlayers: playing_players,
	}
	{
		// mapPlayging := make(map[string]bool, 0)

		for _, p := range msg.Players {
			// check playing
			mapUserPlaying := s.PlayingPresences
			_, p.IsPlaying = mapUserPlaying.Get(p.GetId())
			// status hold card
			if _, exist := s.OrganizeCards[p.GetId()]; exist {
				p.CardStatus = pb.CardStatus(pb.CardEvent_DECLARE)
				// p.Cards = s.OrganizeCards[p.GetId()]
			} else {
				p.CardStatus = pb.CardStatus(pb.CardEvent_COMBINE)
			}
		}
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
	m.notifyUpdateTable(ctx, logger, nk, dispatcher, s, presences, nil)
	//send cards for player rejoin
	for _, presence := range presences {
		if _, found := s.PlayingPresences.Get(presence.GetUserId()); found {
			card := s.Cards[presence.GetUserId()]
			if card == nil {
				continue
			}
			buf, _ := m.marshaler.Marshal(&pb.UpdateDeal{
				PresenceCard: &pb.PresenceCards{
					Presence: presence.GetUserId(),
					Cards:    card.Cards,
				},
			})

			_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL),
				buf, []runtime.Presence{presence.(runtime.Presence)}, nil, true)

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
	defer func() {
		players := entity.NewListPlayer(s.GetPresences())
		// players.ReadWallet(ctx, nk, logger)

		playing_players := entity.NewListPlayer(s.GetPlayingPresences())
		// playing_players.ReadWallet(ctx, nk, logger)

		msg := &pb.UpdateTable{
			Bet:            int64(s.Label.Bet),
			Players:        players,
			PlayingPlayers: playing_players,
			JpTreasure:     s.GetJackpotTreasure(),
		}

		m.NotifyUpdateTable(s, logger, dispatcher, msg)
	}()
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
	maxJP := int64(bet) * 100 * vipLv
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
