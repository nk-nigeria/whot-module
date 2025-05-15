package hand

import (
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	blc "github.com/nakamaFramework/whot-module/usecase/bin_list_card"
)

var naturalCardChecker = []HandCheckFunc{
	CheckCleanDragon,
	CheckDragon,
	CheckSixPairs,
	CheckFullColor,
}
var naturalHandChecker = []func(*Hand) (*HandPoint, bool){
	CheckThreeFlushes,
	CheckThreeStraights,
}

func CheckNaturalCards(h *Hand) (*HandPoint, bool) {
	// check natural win
	bcards := entity.NewBinListCards(h.GetCards())
	for _, checkerFn := range naturalCardChecker {
		handPoint, valid := checkerFn(bcards)
		if valid {
			return handPoint, valid
		}
	}

	return nil, false
}

func CheckNaturalHands(hand *Hand) (*HandPoint, bool) {
	// check natural win
	for _, checkerFn := range naturalHandChecker {
		handPoint, valid := checkerFn(hand)
		if valid {
			return handPoint, valid
		}
	}

	return nil, false
}

// CheckCleanDragon
// Sảnh rồng đồng màu
func CheckCleanDragon(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineFlush); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_NaturalCleanDragon, ScorePointNaturalCleanDragon, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckFullColor
// Đồng màu 12 lá
func CheckFullColor(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineFullColor); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_NaturalFullColors, ScorePointNaturalFullColors, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckDragon
// Sảnh rồng
func CheckDragon(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineStraight); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_NaturalDragon, ScorePointNaturalDragon, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckSixPairs
// 6 đôi
func CheckSixPairs(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombinePair); count >= 6 {
		handPoint := createPointFromList(pb.HandRanking_NaturalSixPairs, ScorePointNaturalSixPairs, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckThreeStraight
// 3 sảnh
func CheckThreeStraights(hand *Hand) (*HandPoint, bool) {
	bcards := entity.NewBinListCards(hand.frontHand.Cards)
	_, isFontHandStraight := CheckStraight(bcards)
	threeStraight := isFontHandStraight && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	if threeStraight {
		var listCard entity.ListCard
		listCard = append(listCard, SortCard(hand.backHand.Cards)...)
		listCard = append(listCard, SortCard(hand.middleHand.Cards)...)
		listCard = append(listCard, SortCard(hand.frontHand.Cards)...)

		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalThreeStraights, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalThreeStraights,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

// CheckThreeFlushes
// 3 cái thùng
func CheckThreeFlushes(hand *Hand) (*HandPoint, bool) {
	// isFontHandFlush :=
	bcards := entity.NewBinListCards(hand.frontHand.Cards)
	_, isFontHandFlush := CheckFlush(bcards)
	threeFlush := isFontHandFlush && hand.middleHand.Point.IsFlush() && hand.backHand.Point.IsFlush()
	if threeFlush {
		var listCard entity.ListCard
		listCard = append(listCard, SortCard(hand.backHand.Cards)...)
		listCard = append(listCard, SortCard(hand.middleHand.Cards)...)
		listCard = append(listCard, SortCard(hand.frontHand.Cards)...)

		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalThreeOfFlushes, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalThreeOfFlushes,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

func CheckJackpot(childHand *ChildHand) bool {
	if len(childHand.Cards) != 5 {
		return false
	}
	card := SortCard(childHand.Cards)
	if card[0].GetRank() != entity.Rank10 {
		return false
	}
	if card[0].GetSuit() != entity.SuitSpades {
		return false
	}
	bcards := entity.NewBinListCards(childHand.Cards)
	_, isStraightFlush := CheckStraightFlush(bcards)
	if !isStraightFlush {
		return false
	}
	return true
}
