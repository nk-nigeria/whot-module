package engine

import (
	"errors"
	"math"

	pb "github.com/nakama-nigeria/cgp-common/proto/whot"
	"github.com/nakama-nigeria/whot-module/entity"
	mockcodegame "github.com/nakama-nigeria/whot-module/mock_code_game"
	"github.com/nakama-nigeria/whot-module/pkg/log"
)

type Engine struct {
	deck *entity.Deck
}

func NewWhotPokerEngine() UseCase {
	return &Engine{}
}

func (e *Engine) NewGame(s *entity.MatchState) error {
	s.Cards = make(map[string]*pb.ListCard)
	if s.WinnerId != "" {
		log.GetLogger().Info("Resetting match state for new game")
		s.PreviousWinnerId = s.WinnerId
		s.WinnerId = ""
	}
	s.SetDealer()
	s.CurrentTurn = s.DealerId
	return nil
}

func (e *Engine) Deal(s *entity.MatchState) error {
	e.deck = entity.NewDeck()
	e.deck.Shuffle()
	mockcodegame.InitMapMockCodeListCard(log.GetLogger())
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
		cards, err := e.deck.Deal(CountCard)
		if err == nil {
			s.Cards[userId] = cards
		} else {
			return err
		}
	}
	card, err := e.deck.Deal(1)
	if err != nil {
		return err
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

	if !e.IsValidPlay(playedEntityCard, topEntityCard) {
		log.GetLogger().Error("invalid card played: %v on top %v", card, s.TopCard)
		return entity.EffectNone, errors.New("invalid card played")
	}

	effect := entity.EffectNone

	switch card.GetRank() {
	case entity.CardValueHoldOn: // 1
		effect = entity.EffectHoldOn
		s.IsHoldOn = true

	case entity.CardValuePickTwo: // 2
		effect = entity.EffectPickTwo
		s.PickPenalty += 2
		s.EffectTarget = s.GetNextPlayerClockwise(userId)
		s.CurrentTurn = s.EffectTarget

	case entity.CardValuePickThree: // 5
		effect = entity.EffectPickThree
		s.PickPenalty += 3
		s.EffectTarget = s.GetNextPlayerClockwise(userId)
		s.CurrentTurn = s.EffectTarget

	case entity.CardValueSuspension: // 8
		effect = entity.EffectSuspension
		s.IsSuspension = true
		nextPlayer := s.GetNextPlayerClockwise(userId)
		s.CurrentTurn = s.GetNextPlayerClockwise(nextPlayer)

	case entity.CardValueGeneralMarket: // 14
		effect = entity.EffectGeneralMarket

	case entity.CardValueWhot: // 20
		effect = entity.EffectWhot
		s.WaitingForWhotShape = true
	}
	// Cập nhật bài trên bàn
	s.TopCard = card
	s.CurrentEffect = effect

	// Xóa lá bài đã đánh khỏi bài của người chơi
	playerCards.Cards = append(playerCards.Cards[:cardIndex], playerCards.Cards[cardIndex+1:]...)
	s.Cards[userId] = playerCards

	// Kiểm tra người chơi đã hết bài chưa
	if len(playerCards.Cards) == 0 {
		s.WinnerId = userId
	}

	return effect, nil
}

func (e *Engine) DrawCardsFromDeck(s *entity.MatchState, userID string) (int, error) {
	// Kiểm tra lượt chơi
	if s.CurrentTurn != userID {
		return 0, errors.New("not user's turn")
	}

	// Xác định số lá cần rút
	cardsToDraw := 1
	drawingPenalty := false

	if s.PickPenalty > 0 && s.EffectTarget == userID {
		cardsToDraw = s.PickPenalty
		s.PickPenalty = 0
		s.EffectTarget = ""
		drawingPenalty = true
	}

	// Rút bài từ bộ bài
	card, err := e.deck.Deal(cardsToDraw)
	if err != nil {
		return 0, err
	}

	s.Cards[userID].Cards = append(s.Cards[userID].Cards, card.Cards...)

	// Xác định người chơi tiếp theo
	if !drawingPenalty {
		s.CurrentTurn = s.GetNextPlayerClockwise(userID)
	}

	return cardsToDraw, nil
}

func (e *Engine) HandleGeneralMarket(s *entity.MatchState, userID string) error {
	for _, key := range s.PlayingPresences.Keys() {
		otherUserId := key.(string)
		if otherUserId != userID {
			cardsToDraw, err := e.DrawCardsFromDeck(s, otherUserId)
			if cardsToDraw == 1 {
				continue
			}
			if err != nil {
				return err
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

	return nil
}

func (e *Engine) Finish(s *entity.MatchState) *pb.UpdateFinish {
	updateFinish := pb.UpdateFinish{}

	// Nếu có người thắng trực tiếp (đánh hết bài)
	if s.WinnerId != "" {
		// Người thắng nhận toàn bộ tiền
		winResult := &pb.WhotPlayerResult{
			UserId: s.WinnerId,
			Score: &pb.WhotScoreResult{
				WinFactor: 1, // Người thắng nhận toàn bộ tiền
			},
		}
		updateFinish.Results = append(updateFinish.Results, winResult)

		// Những người khác thua
		for _, val := range s.PlayingPresences.Keys() {
			uid := val.(string)
			if uid != s.WinnerId {
				loseResult := &pb.WhotPlayerResult{
					UserId: uid,
					Score: &pb.WhotScoreResult{
						WinFactor: 0, // Người thua không nhận gì
					},
				}
				updateFinish.Results = append(updateFinish.Results, loseResult)
			}
		}
		return &updateFinish
	}

	// Nếu bộ bài hết mà chưa có ai hết bài, tính điểm
	// Tính điểm cho mỗi người chơi
	scores := make(map[string]int)
	lowestScore := math.MaxInt32

	for _, val := range s.PlayingPresences.Keys() {
		uid := val.(string)
		playerCards := s.Cards[uid]

		// Tính tổng điểm
		totalScore := 0
		for _, card := range playerCards.Cards {
			cardValue := entity.CalculateCardValue(card)
			totalScore += cardValue
		}

		scores[uid] = totalScore
		if totalScore < lowestScore {
			lowestScore = totalScore
		}

		// Lưu điểm vào kết quả
		result := &pb.WhotPlayerResult{
			UserId: uid,
			Score: &pb.WhotScoreResult{
				TotalPoints: int64(totalScore),
			},
		}
		updateFinish.Results = append(updateFinish.Results, result)
	}

	// Tìm những người có điểm thấp nhất
	winners := []string{}
	for uid, score := range scores {
		if score == lowestScore {
			winners = append(winners, uid)
		}
	}

	// Phân bổ tiền thắng
	winFactor := 1.0 / float64(len(winners))

	// Đặt kết quả cuối cùng
	for i, result := range updateFinish.Results {
		isWinner := false
		for _, winnerID := range winners {
			if result.UserId == winnerID {
				isWinner = true
				break
			}
		}

		if isWinner {
			updateFinish.Results[i].Score.WinFactor = float64(winFactor)
			updateFinish.Results[i].Score.IsWinner = true
		} else {
			updateFinish.Results[i].Score.WinFactor = 0
			updateFinish.Results[i].Score.IsWinner = false
		}
	}

	return &updateFinish
}

func (e *Engine) IsValidPlay(playedCard, topCard entity.Card) bool {

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
