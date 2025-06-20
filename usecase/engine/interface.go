package engine

import (
	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/entity"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(s *entity.MatchState) error
	PlayCard(s *entity.MatchState, presence string, card *pb.Card) (entity.CardEffect, error)
	DrawCardsFromDeck(s *entity.MatchState, userID string) (int, error)
	ChooseWhotShape(s *entity.MatchState, userID string, shape pb.CardSuit) error
	ChooseAutomaticWhotShape() *pb.Card
	HandleGeneralMarket(s *entity.MatchState, userID string) error
	FindPlayableCard(s *entity.MatchState, userId string) *pb.Card
	GetPlayerCardCounts(s *entity.MatchState) map[string]int32
	GetDeckCount() int32
	Finish(s *entity.MatchState) *pb.UpdateFinish
}
