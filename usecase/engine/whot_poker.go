package engine

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/entity"
	mockcodegame "github.com/nk-nigeria/whot-module/mock_code_game"
	"github.com/nk-nigeria/whot-module/pkg/log"
)

type Engine struct {
	deck *entity.Deck
}

func NewWhotEngine() UseCase {
	return &Engine{}
}

func (e *Engine) NewGame(s *entity.MatchState) error {
	s.SetDealer()
	s.CurrentTurn = s.DealerId
	return nil
}

func (e *Engine) Deal(s *entity.MatchState) error {
	e.deck = entity.NewDeck()
	e.deck.Shuffle()
	// mockcodegame.InitMapMockCodeListCard(log.GetLogger())
	if list, exist := mockcodegame.MapMockCodeListCard[int(s.Label.MockCodeCard)]; exist {
		if len(list) >= s.PlayingPresences.Size() {
			log.GetLogger().Debug("[MockCard] Match has label mock code card %d " +
				"Init card for player from mock")
			idx := 0
			for _, k := range s.PlayingPresences.Keys() {
				userId := k.(string)
				s.Cards[userId] = list[idx]
				idx++
			}
			return nil
		} else {
			log.GetLogger().Debug("[MockCard] Match has label mock code card %d "+
				"but list card in mock smaller than size playert join game, fallback to normal", s.Label.MockCodeCard)
		}
	}

	card, err := e.deck.Deal(1, true)
	if err != nil {
		return err
	}

	CountCard := 0
	switch s.PlayingPresences.Size() {
	case 2:
		CountCard = entity.MaxPresenceCard
	case 3:
		CountCard = entity.MaxPresenceCard - 1
	case 4:
		CountCard = entity.MaxPresenceCard - 2
	default:
		log.GetLogger().Error("Invalid number of players: %d", s.PlayingPresences.Size())
		return nil
	}

	// loop on userid in match
	for _, k := range s.PlayingPresences.Keys() {
		userId := k.(string)
		cards, err := e.deck.Deal(CountCard, false)
		if err == nil {
			s.Cards[userId] = cards
		} else {
			return err
		}
	}

	if len(card.Cards) > 0 {
		s.TopCard = card.Cards[0]
	} else {
		return errors.New("no cards dealt for top card")
	}

	s.BuildPlayOrderFromDealer()

	return nil
}

func (e *Engine) PlayCard(s *entity.MatchState, userId string, card *pb.Card) (entity.CardEffect, error) {

	if s.CurrentTurn != userId {
		log.GetLogger().Error("not user's turn")
		return entity.EffectNone, errors.New("not user's turn")
	}

	playerCards, ok := s.Cards[userId]
	if !ok {
		log.GetLogger().Error("player cards not found")
		return entity.EffectNone, errors.New("player cards not found")
	}

	found := false
	cardIndex := -1
	for i, c := range playerCards.Cards {
		if c.GetRank() == card.GetRank() && c.GetSuit() == card.GetSuit() {
			found = true
			cardIndex = i
			break
		}
	}
	if !found {
		log.GetLogger().Error("card not in player's hand")
		return entity.EffectNone, errors.New("card not in player's hand")
	}

	playedEntityCard := entity.NewCardFromPb(card.GetRank(), card.GetSuit())
	topEntityCard := entity.NewCardFromPb(s.TopCard.GetRank(), s.TopCard.GetSuit())
	if !e.isValidPlay(playedEntityCard, topEntityCard) {
		return entity.EffectNone, errors.New("invalid card played")
	}

	effect := entity.EffectNone

	switch card.GetRank() {
	case pb.CardRank_RANK_1: // 1
		effect = entity.EffectHoldOn

	case pb.CardRank_RANK_2: // 2
		effect = entity.EffectPickTwo
		s.PickPenalty += 2
		s.EffectTarget = s.GetNextPlayerClockwise(userId)
		s.CurrentTurn = s.EffectTarget

	case pb.CardRank_RANK_5: // 5
		effect = entity.EffectPickThree
		s.PickPenalty += 3
		s.EffectTarget = s.GetNextPlayerClockwise(userId)
		s.CurrentTurn = s.EffectTarget

	case pb.CardRank_RANK_8: // 8
		effect = entity.EffectSuspension
		nextPlayer := s.GetNextPlayerClockwise(userId)
		s.EffectTarget = nextPlayer
		s.CurrentTurn = s.GetNextPlayerClockwise(nextPlayer)

	case pb.CardRank_RANK_14: // 14
		effect = entity.EffectGeneralMarket

	case pb.CardRank_RANK_20: // 20
		fmt.Printf("Whot card played by %s\n", userId)
		effect = entity.EffectWhot
		s.WaitingForWhotShape = true
	default:
		s.CurrentTurn = s.GetNextPlayerClockwise(userId)
	}

	s.TopCard = card
	s.CurrentEffect = effect

	// Xóa lá bài đã đánh khỏi bài của người chơi
	playerCards.Cards = append(playerCards.Cards[:cardIndex], playerCards.Cards[cardIndex+1:]...)
	s.Cards[userId] = playerCards

	if len(playerCards.Cards) == 0 {
		s.WinnerId = userId
		s.EffectTarget = ""
		s.PickPenalty = 0
		s.IsEndingGame = true
		log.GetLogger().Info("Player %s has won the game by playing the last card", userId)
		return entity.EffectNone, nil
	}

	return effect, nil
}

func (e *Engine) DrawCardsFromDeck(s *entity.MatchState, userID string) (int, error) {
	// Kiểm tra lượt chơi
	if s.CurrentTurn != userID && s.CurrentEffect != entity.EffectGeneralMarket {
		return 0, errors.New("not user's turn")
	}

	// Xác định số lá cần rút
	cardsToDraw := 1

	if s.PickPenalty > 0 && s.EffectTarget == userID {
		cardsToDraw = s.PickPenalty
		s.PickPenalty = 0
		s.EffectTarget = ""
		s.CurrentEffect = entity.EffectNone
	}

	// Rút bài từ bộ bài
	card, err := e.deck.Deal(cardsToDraw, false)
	if err != nil {
		return 0, err
	}
	cardsToDraw = len(card.Cards)

	s.Cards[userID].Cards = append(s.Cards[userID].Cards, card.Cards...)

	if e.deck.RemainingCards() == 0 {
		log.GetLogger().Info("Deck is empty, handle game reward")
		s.IsEndingGame = true
		return cardsToDraw, nil
	}

	// Xác định người chơi tiếp theo
	if s.CurrentEffect != entity.EffectGeneralMarket {
		s.CurrentTurn = s.GetNextPlayerClockwise(userID)
	}
	return cardsToDraw, nil
}

func (e *Engine) HandleGeneralMarket(s *entity.MatchState, userID string) error {
	for _, key := range s.PlayingPresences.Keys() {
		otherUserId := key.(string)
		if otherUserId != userID {
			_, err := e.DrawCardsFromDeck(s, otherUserId)
			if err != nil {
				return err
			}
			if s.IsEndingGame {
				return nil
			}

		}

	}
	return nil
}

func (e *Engine) ChooseWhotShape(s *entity.MatchState, userID string, shape pb.CardSuit) error {
	// Kiểm tra xem có đang chờ chọn hình không
	if !s.WaitingForWhotShape {
		return errors.New("not waiting for Whot shape choice")
	}

	// Kiểm tra lượt chơi
	if s.CurrentTurn != userID {
		return errors.New("not user's turn")
	}

	// Cập nhật hình được chọn
	s.TopCard.Suit = shape
	s.WaitingForWhotShape = false
	s.CurrentTurn = s.GetNextPlayerClockwise(userID)
	return nil
}

func (e *Engine) FindPlayableCard(s *entity.MatchState, userId string) *pb.Card {
	userCards := s.Cards[userId]
	if userCards == nil || len(userCards.Cards) == 0 {
		return nil
	}

	topCard := s.TopCard

	// Thực hiện logic tìm bài phù hợp:
	// 1. Ưu tiên lá bài có cùng số (rank)
	// 2. Ưu tiên lá bài có cùng suit
	// 3. Nếu có lá Whot (joker), chọn Whot

	for _, card := range userCards.Cards {
		if card.Rank == topCard.Rank {
			return card
		}
	}

	for _, card := range userCards.Cards {
		if card.Suit == topCard.Suit && s.CurrentEffect != entity.EffectPickTwo && s.CurrentEffect != entity.EffectPickThree {
			return card
		}
	}

	for _, card := range userCards.Cards {
		if card.Rank == pb.CardRank_RANK_20 && s.CurrentEffect != entity.EffectPickTwo && s.CurrentEffect != entity.EffectPickThree {
			return card
		}
	}
	return nil
}

func (e *Engine) Finish(s *entity.MatchState) *pb.UpdateFinish {

	updateFinish := pb.UpdateFinish{}

	// 1. Tính điểm cho tất cả người chơi
	playerScores := e.calculatePlayerScores(s) // map[uid]int
	lowestScore := math.MaxInt32
	for _, score := range playerScores {
		if score < lowestScore {
			lowestScore = score
		}
	}

	// 2. Tìm người chơi có điểm thấp nhất
	winners := []string{}
	losers := []string{}
	for uid, score := range playerScores {
		if score == lowestScore {
			winners = append(winners, uid)
		} else {
			losers = append(losers, uid)
		}
	}

	// 3. Xác định WinnerId nếu chưa có (trường hợp hết bài rút)
	if s.WinnerId == "" {
		if len(winners) == 1 {
			s.WinnerId = winners[0]
		} else if len(winners) > 1 {
			s.WinnerId = winners[rand.Intn(len(winners))]
			log.GetLogger().Info("Multiple winners found, randomly selected: %s", s.WinnerId)
		}
	}

	numLosers := float64(len(losers))
	numWinners := float64(len(winners))
	winFactorPerWinner := 0.0
	if numWinners > 0 {
		winFactorPerWinner = numLosers / numWinners
	}

	// 4. Gán kết quả cho từng người chơi
	for uid, total := range playerScores {
		isWinner := false
		for _, w := range winners {
			if uid == w {
				isWinner = true
				break
			}
		}

		var winFactor float64
		if isWinner {
			winFactor = winFactorPerWinner
		} else {
			winFactor = -1.0
		}

		updateFinish.Results = append(updateFinish.Results, &pb.WhotPlayerResult{
			UserId:         uid,
			TotalPoints:    int64(total),
			IsWinner:       isWinner,
			WinFactor:      winFactor,
			RemainingCards: s.Cards[uid].Cards,
		})
	}

	return &updateFinish
}

func (e *Engine) calculatePlayerScores(s *entity.MatchState) map[string]int {
	scores := make(map[string]int)
	for _, val := range s.PlayingPresences.Keys() {
		uid := val.(string)
		playerCards := s.Cards[uid]
		total := 0
		for _, card := range playerCards.Cards {
			total += entity.CalculateCardValue(card)
		}
		scores[uid] = total
		log.GetLogger().Info("User %s has total score: %d", uid, total)
	}
	return scores
}

func (e *Engine) isValidPlay(playedCard, topCard entity.Card) bool {

	if playedCard.GetSuit() == entity.SuitNone && playedCard.GetRank() == entity.RankWHOT {
		return true
	}

	if playedCard.GetSuit() == topCard.GetSuit() {
		return true
	}

	if playedCard.GetRank() == topCard.GetRank() {
		return true
	}
	return false
}

func (e *Engine) ChooseAutomaticWhotShape() *pb.Card {
	// Phương án 1: Chọn shape ngẫu nhiên
	shapes := []pb.CardSuit{
		pb.CardSuit_SUIT_CIRCLE,
		pb.CardSuit_SUIT_SQUARE,
		pb.CardSuit_SUIT_TRIANGLE,
		pb.CardSuit_SUIT_STAR,
		pb.CardSuit_SUIT_CROSS,
	}
	randomIndex := rand.Intn(len(shapes))

	return &pb.Card{
		Rank: pb.CardRank_RANK_20,
		Suit: shapes[randomIndex],
	}
}

func (e *Engine) GetPlayerCardCounts(s *entity.MatchState) map[string]int32 {
	counts := make(map[string]int32)
	for _, key := range s.PlayingPresences.Keys() {
		userId := key.(string)
		if cards, ok := s.Cards[userId]; ok {
			counts[userId] = int32(len(cards.Cards))
		}
	}
	return counts
}

func (e *Engine) GetDeckCount() int32 {
	return int32(e.deck.RemainingCards())
}
