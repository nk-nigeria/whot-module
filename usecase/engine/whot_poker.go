package engine

import (
	combinations "github.com/mxschmitt/golang-combinations"
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	mockcodegame "github.com/nakamaFramework/whot-module/mock_code_game"
	"github.com/nakamaFramework/whot-module/pkg/log"
	"github.com/nakamaFramework/whot-module/usecase/hand"
)

type Engine struct {
	deck *entity.Deck
}

func NewWhotPokerEngine() UseCase {
	return &Engine{}
}

func (c *Engine) NewGame(s *entity.MatchState) error {
	s.Cards = make(map[string]*pb.ListCard)
	s.OrganizeCards = make(map[string]*pb.ListCard)

	return nil
}

func (c *Engine) Deal(s *entity.MatchState) error {
	c.deck = entity.NewDeck()
	c.deck.Shuffle()
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

	// loop on userid in match
	for _, k := range s.PlayingPresences.Keys() {
		userId := k.(string)
		cards, err := c.deck.Deal(entity.MaxPresenceCard)
		if err == nil {
			s.Cards[userId] = cards
		} else {
			return err
		}
	}

	return nil
}

func (c *Engine) Organize(s *entity.MatchState, presence string, cards *pb.ListCard) error {
	s.UpdateShowCard(presence, cards)
	return nil
}

func (c *Engine) Combine(s *entity.MatchState, presence string) error {
	s.RemoveShowCard(presence)
	return nil
}

func (c *Engine) Finish(s *entity.MatchState) *pb.UpdateFinish {
	// Check every user
	updateFinish := pb.UpdateFinish{}
	presenceCount := s.PlayingPresences.Size()
	ctx := hand.NewCompareContext(presenceCount)

	log.GetLogger().Info("Finish presence %v, size %v", s.PlayingPresences, presenceCount)

	// prepare for compare data
	userIds := make([]string, 0, presenceCount)
	hands := make(map[string]*hand.Hand)
	results := make(map[string]*pb.ComparisonResult)
	userJackpot := ""
	for _, val := range s.PlayingPresences.Keys() {
		uid := val.(string)
		userIds = append(userIds, uid)

		cards := s.OrganizeCards[uid]
		var h *hand.Hand
		var err error
		h, err = hand.NewHandFromPb(cards)
		h.SetOwner(uid)
		if err != nil {
			continue
		}

		hands[uid] = h

		result := &pb.ComparisonResult{
			UserId:      uid,
			PointResult: h.GetPointResult(),
			ScoreResult: &pb.ScoreResult{},
		}

		results[uid] = result
		if h.IsJackpot() {
			userJackpot = uid
		}

		updateFinish.Results = append(updateFinish.Results, result)

		log.GetLogger().Info("prepare for %s, hand %v, result %v", uid, h, result)
	}

	pairs := combinations.Combinations(userIds, 2)
	mapRcPair := make(map[string]*hand.ComparisonResult, 0)
	log.GetLogger().Info("combination %v of %v", pairs, len(userIds))
	for _, pair := range pairs {
		uid1 := pair[0]
		uid2 := pair[1]
		log.GetLogger().Info("compare %v with %v", pair[0], pair[1])

		// calculate natural point, normal point, hand bonus case
		rc := hand.CompareHand(ctx, hands[uid1], hands[uid2])
		mapRcPair[uid1+uid2] = rc
		hand.ProcessCompareResult(ctx, results[uid1], rc.GetR1())
		hand.ProcessCompareResult(ctx, results[uid2], rc.GetR2())
		for _, bonus := range rc.GetBonuses() {
			if bonus.Type != pb.HandBonusType_Scoop {
				updateFinish.Bonuses = append(updateFinish.Bonuses, bonus)
			}

		}

	}

	// hand.ProcessCompareBonusResult(ctx, updateFinish.Results, &updateFinish.Bonuses)
	hand.ProcessCompareBonusResult(ctx,
		updateFinish.Results, mapRcPair, &updateFinish.Bonuses)

	for _, rc := range mapRcPair {
		for _, bonus := range rc.GetBonuses() {
			if bonus.Type == pb.HandBonusType_Scoop {
				updateFinish.Bonuses = append(updateFinish.Bonuses, bonus)
			}
		}

	}
	hand.CalcTotalFactor(updateFinish.Results)
	if userJackpot != "" {
		updateFinish.Jackpot = &pb.Jackpot{
			UserId:   userJackpot,
			GameCode: entity.ModuleName,
		}
	} else {
		updateFinish.Jackpot = &pb.Jackpot{}
	}
	return &updateFinish
}
