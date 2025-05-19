package hand

import (
	"fmt"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
)

type ChildHand struct {
	Cards    entity.ListCard
	Point    *HandPoint
	handType int
}

func (h1 ChildHand) Compare(h2 *ChildHand) int {
	return CompareHandPoint(h1.Point, h2.Point)
}

func (ch ChildHand) String() string {
	return fmt.Sprintf("Cards: %v, Point: %v", ch.Cards, ch.Point)
}

func (ch *ChildHand) calculatePoint() {
	if ch.Point != nil {
		return
	}
	ch.Point = CalculatePoint(ch)
}

func NewChildHand(cards entity.ListCard, handType int) *ChildHand {
	child := ChildHand{
		Cards:    cards[:],
		handType: handType,
	}

	return &child
}

func (ch *ChildHand) GetBonus() (pb.HandBonusType, int) {
	switch ch.handType {
	case kFronHand:
		if ch.Point.rankingType == pb.HandRanking_ThreeOfAKind {
			return pb.HandBonusType_BonusThreeOfAKindFrontHand, mapBonusPoint[pb.HandBonusType_BonusThreeOfAKindFrontHand]
		}
	case kMidHand:
		if ch.Point.rankingType == pb.HandRanking_FullHouse {
			return pb.HandBonusType_BonusFullHouseMidHand, mapBonusPoint[pb.HandBonusType_BonusFullHouseMidHand]
		}
		if ch.Point.rankingType == pb.HandRanking_FourOfAKind {
			return pb.HandBonusType_BonusFourOfAKindMidHand, mapBonusPoint[pb.HandBonusType_BonusFourOfAKindMidHand]
		}
		if ch.Point.rankingType == pb.HandRanking_StraightFlush {
			return pb.HandBonusType_BonusStraightFlushMidHand, mapBonusPoint[pb.HandBonusType_BonusStraightFlushMidHand]
		}
	case kBackHand:
		if ch.Point.rankingType == pb.HandRanking_FourOfAKind {
			return pb.HandBonusType_BonusFourOfAKindBackHand, mapBonusPoint[pb.HandBonusType_BonusFourOfAKindBackHand]
		}
		if ch.Point.rankingType == pb.HandRanking_StraightFlush {
			return pb.HandBonusType_BonusStraightFlushBackHand, mapBonusPoint[pb.HandBonusType_BonusStraightFlushBackHand]
		}
	}

	return pb.HandBonusType_None, 0
}
