package hand

import (
	"fmt"

	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/nakamaFramework/whot-module/pkg/log"
	blc "github.com/nakamaFramework/whot-module/usecase/bin_list_card"
)

//  				t1		s1		s2		s3		s4		s5
//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF

const (
	//	3				t1		m1		m2		m3
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointHighCard = uint8(0x01)
	//	3				t1		d1		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		d1		m1		m2		m3
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointPair = uint8(0x02)
	//					t1		d1		d2		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointTwoPairs = uint8(0x03)
	//	3				t1		s1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		s1		m1		m2
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointThreeOfAKind = uint8(0x04)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointStraight = uint8(0x05)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFlush = uint8(0x06)
	//					t1		s1		d1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFullHouse = uint8(0x07)
	//					t1		q1		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFourOfAKind = uint8(0x08)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointStraightFlush = uint8(0x09)

	//			n1				m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	m6		m7		m8		m9		m10		m11		m12		m13
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointNaturalThreeStraights = uint8(0x01)
	ScorePointNaturalThreeOfFlushes = uint8(0x02)
	ScorePointNaturalSixPairs       = uint8(0x03)
	ScorePointNaturalFullColors     = uint8(0x04)
	ScorePointNaturalDragon         = uint8(0x05)
	ScorePointNaturalCleanDragon    = uint8(0x06)
)

func createPoint(t, p1, p2, p3, p4, p5 uint8) uint64 {
	var point uint64 = 0
	point |= uint64(t) << (5 * 8)
	point |= uint64(p1) << (4 * 8)
	point |= uint64(p2) << (3 * 8)
	point |= uint64(p3) << (2 * 8)
	point |= uint64(p4) << (1 * 8)
	point |= uint64(p5)

	return point
}

func createPointFromList(ranking pb.HandRanking, t uint8, cards entity.ListCard) *HandPoint {
	var point uint64 = 0
	var lpoint uint64 = 0

	if len(cards) >= 13 {
		point |= uint64(t) << (6 * 8)
		point |= uint64(cards[12].GetRank()) << (4 * 8)
		point |= uint64(cards[11].GetRank()) << (3 * 8)
		point |= uint64(cards[10].GetRank()) << (2 * 8)
		point |= uint64(cards[9].GetRank()) << (1 * 8)
		point |= uint64(cards[8].GetRank())

		lpoint |= uint64(cards[7].GetRank()) << (7 * 8)
		lpoint |= uint64(cards[6].GetRank()) << (6 * 8)
		lpoint |= uint64(cards[5].GetRank()) << (5 * 8)
		lpoint |= uint64(cards[4].GetRank()) << (4 * 8)
		lpoint |= uint64(cards[3].GetRank()) << (3 * 8)
		lpoint |= uint64(cards[2].GetRank()) << (2 * 8)
		lpoint |= uint64(cards[1].GetRank()) << (1 * 8)
		lpoint |= uint64(cards[0].GetRank())

	} else if len(cards) >= 5 {
		point |= uint64(t) << (5 * 8)
		point |= uint64(cards[4].GetRank()) << (4 * 8)
		point |= uint64(cards[3].GetRank()) << (3 * 8)
		point |= uint64(cards[2].GetRank()) << (2 * 8)
		point |= uint64(cards[1].GetRank()) << (1 * 8)
		point |= uint64(cards[0].GetRank())
	} else if len(cards) >= 3 {
		point |= uint64(t) << (5 * 8)
		point |= uint64(cards[2].GetRank()) << (4 * 8)
		point |= uint64(cards[1].GetRank()) << (3 * 8)
		point |= uint64(cards[0].GetRank()) << (2 * 8)
	} else {

	}

	return &HandPoint{
		rankingType: ranking,
		point:       point,
		lpoint:      lpoint,
	}
}

func createPointNaturalCard(t uint8, cards entity.ListCard) (uint64, uint64) {
	var hpoint uint64 = 0
	var lpoint uint64 = 0

	hpoint |= uint64(t) << (6 * 8)
	hpoint |= uint64(cards[0].GetRank()) << (4 * 8)
	hpoint |= uint64(cards[1].GetRank()) << (3 * 8)
	hpoint |= uint64(cards[2].GetRank()) << (2 * 8)
	hpoint |= uint64(cards[3].GetRank()) << (1 * 8)
	hpoint |= uint64(cards[4].GetRank())

	lpoint |= uint64(cards[5].GetRank()) << (7 * 8)
	lpoint |= uint64(cards[6].GetRank()) << (6 * 8)
	lpoint |= uint64(cards[7].GetRank()) << (5 * 8)
	lpoint |= uint64(cards[8].GetRank()) << (4 * 8)
	lpoint |= uint64(cards[9].GetRank()) << (3 * 8)
	lpoint |= uint64(cards[10].GetRank()) << (2 * 8)
	lpoint |= uint64(cards[11].GetRank()) << (1 * 8)
	lpoint |= uint64(cards[12].GetRank())

	return hpoint, lpoint
}

type HandPoint struct {
	rankingType pb.HandRanking
	point       uint64
	lpoint      uint64
}

func (h HandPoint) String() string {
	return fmt.Sprintf("Rank %v, Point: %v", h.rankingType, h.point)
}

func (h HandPoint) ToHandResultPB() *pb.HandResult {
	return &pb.HandResult{
		Ranking: h.rankingType,
		Point:   h.point,
		Lpoint:  h.lpoint,
	}
}

func CompareHandPoint(p1, p2 *HandPoint) int {
	if p1.point > p2.point {
		return 1
	} else if p1.point < p2.point {
		return -1
	} else {
		if p1.lpoint > p2.lpoint {
			return 1
		} else if p1.lpoint < p2.lpoint {
			return -1
		}
	}

	return 0
}

func (h *HandPoint) IsStraight() bool {
	return h.rankingType == pb.HandRanking_Straight
}

func (h *HandPoint) IsFlush() bool {
	return h.rankingType == pb.HandRanking_Flush
}

type HandCheckFunc func(*entity.BinListCard) (*HandPoint, bool)

var HandCheckers = []HandCheckFunc{
	CheckStraightFlush,
	CheckFourOfAKind,
	CheckFullHouse,
	CheckFlush,
	CheckStraight,
	CheckThreeOfAKind,
	CheckTwoPairs,
	CheckPair,
	CheckHighCard,
}

var HandCheckerFronts = []HandCheckFunc{
	CheckThreeOfAKind,
	CheckPair,
	CheckHighCard,
}

func CalculatePoint(ch *ChildHand) *HandPoint {
	bcards := entity.NewBinListCards(ch.Cards)
	if ch.handType == kFronHand {
		for _, fn := range HandCheckerFronts {
			log.GetLogger().Info("check %v, cards %v", fn, bcards)
			if handPoint, valid := fn(bcards); valid {
				log.GetLogger().Info("check ok")
				return handPoint
			}
			log.GetLogger().Info("check not ok")
		}
	} else {
		for _, fn := range HandCheckers {
			if handPoint, valid := fn(bcards); valid {
				return handPoint
			}
		}
		log.GetLogger().Info("check not ok")
	}

	return nil
}

func CheckHighCard(bcards *entity.BinListCard) (*HandPoint, bool) {
	cards := bcards.ToList()
	return createPointFromList(pb.HandRanking_HighCard, ScorePointHighCard, cards), true
}

// CheckStraightFlush
// Thùng phá sảnh (en: Straight Flush)
// Năm lá bài cùng màu, đồng chất, cùng một chuỗi số
// Là Flush, có cùng chuỗi
func CheckStraightFlush(bcards *entity.BinListCard) (*HandPoint, bool) {
	_, valid := CheckStraight(bcards)
	if !valid {
		return nil, false
	}

	_, valid = CheckFlush(bcards)
	if !valid {
		return nil, false
	}

	hp := createPointFromList(pb.HandRanking_StraightFlush, ScorePointStraightFlush, bcards.ToList())
	return hp, true
}

// CheckFourOfAKind
// Tứ quý (en: Four of a Kind)
// Bốn lá đồng số
func CheckFourOfAKind(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineFour); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_FourOfAKind, ScorePointFourOfAKind, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckFullHouse
// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func CheckFullHouse(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineFullHouse); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_FullHouse, ScorePointFullHouse, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckFlush
// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func CheckFlush(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineFlush); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_Flush, ScorePointFlush, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckStraight
// Sảnh (en: Straight)
// Năm lá bài trong một chuỗi số (nhưng không đồng chất)
func CheckStraight(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineStraight); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_Straight, ScorePointStraight, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckThreeOfAKind
// Xám chi/Xám cô (en: Three of a Kind)
// Ba lá bài đồng số
func CheckThreeOfAKind(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombineThree); count > 0 {
		handPoint := createPointFromList(pb.HandRanking_ThreeOfAKind, ScorePointThreeOfAKind, sortedCard)
		return handPoint, true
	}
	return nil, false
}

// CheckTwoPairs
// Thú (en: Two Pairs)
// Hai đôi
func CheckTwoPairs(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombinePair); count == 2 {
		handPoint := createPointFromList(pb.HandRanking_TwoPairs, ScorePointTwoPairs, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// CheckPair
// Đôi (en: Pair)
// Hai lá bài đồng số
func CheckPair(bcards *entity.BinListCard) (*HandPoint, bool) {
	if count, sortedCard := blc.NewChinesePokerBinList().GetChain(bcards, blc.CombinePair); count == 1 {
		handPoint := createPointFromList(pb.HandRanking_Pair, ScorePointPair, sortedCard)
		return handPoint, true
	}

	return nil, false
}

// SortCard
// sort card increase by rank, equal rank will check suit
func SortCard(cards entity.ListCard) entity.ListCard {
	bl := entity.NewBinListCards(cards)
	return bl.ToList()
}
