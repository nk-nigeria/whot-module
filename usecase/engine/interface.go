package engine

import (
	pb "github.com/nakama-nigeria/cgp-common/proto/whot"
	"github.com/nakama-nigeria/whot-module/entity"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(s *entity.MatchState) error
	PlayCard(s *entity.MatchState, presence string, card *pb.Card) (entity.CardEffect, error)
	DrawCardsFromDeck(s *entity.MatchState, userID string) (int, error)
	ChooseWhotShape(s *entity.MatchState, userID string, shape pb.CardSuit) error
	HandleGeneralMarket(s *entity.MatchState, userID string) error
	FindPlayableCard(s *entity.MatchState, userId string) *pb.Card
	Finish(s *entity.MatchState) *pb.UpdateFinish
}
