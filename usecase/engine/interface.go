package engine

import (
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(s *entity.MatchState) error
	Organize(s *entity.MatchState, presence string, cards *pb.ListCard) error
	Combine(s *entity.MatchState, presence string) error
	Finish(s *entity.MatchState) *pb.UpdateFinish
}
