package engine

import (
	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/whot-module/entity"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(s *entity.MatchState) error
	PlayCard(s *entity.MatchState, presence string, card *pb.WhotCard) (entity.CardEffect, error)
	CheckDoubleDeckingEligibility(s *entity.MatchState, userId string) bool
	GetPlayableCardsForDouble(s *entity.MatchState, userId string) []*pb.WhotCard
	DrawCardsFromDeck(s *entity.MatchState, userID string) (int, error)
	ChooseWhotShape(s *entity.MatchState, userID string, shape pb.WhotCardSuit) error
	ChooseAutomaticWhotShape(s *entity.MatchState) *pb.WhotCard
	HandleGeneralMarket(s *entity.MatchState, userID string) error
	FindPlayableCard(s *entity.MatchState, userId string) *pb.WhotCard
	GetPlayerCardCounts(s *entity.MatchState) map[string]int32
	GetDeckCount() int32
	Finish(s *entity.MatchState) *pb.UpdateFinish
}
