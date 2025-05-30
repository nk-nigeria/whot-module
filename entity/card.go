package entity

import (
	"fmt"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
)

// WHOT sử dụng 5 suit và các số từ 1–14 + WHOT

type Card uint8

const (
	RankStep uint8 = 0x10
	Rank1    uint8 = 0x10
	Rank2    uint8 = 0x20
	Rank3    uint8 = 0x30
	Rank4    uint8 = 0x40
	Rank5    uint8 = 0x50
	Rank6    uint8 = 0x60
	Rank7    uint8 = 0x70
	Rank8    uint8 = 0x80
	Rank9    uint8 = 0x90
	Rank10   uint8 = 0xA0
	Rank11   uint8 = 0xB0
	Rank12   uint8 = 0xC0
	Rank13   uint8 = 0xD0
	Rank14   uint8 = 0xE0
	RankWHOT uint8 = 0xF0 // Lá đặc biệt, không có chất

	SuitNone     uint8 = 0x00
	SuitCircle   uint8 = 0x01
	SuitCross    uint8 = 0x02
	SuitStar     uint8 = 0x03
	SuitTriangle uint8 = 0x04
	SuitSquare   uint8 = 0x05
)

var Ranks = []uint8{
	Rank1,    // 1
	Rank2,    // 2
	Rank3,    // 3
	Rank4,    // 4
	Rank5,    // 5
	Rank6,    // 6
	Rank7,    // 7
	Rank8,    // 8
	Rank9,    // 9
	Rank10,   // 10
	Rank11,   // 11
	Rank12,   // 12
	Rank13,   // 13
	Rank14,   // 14
	RankWHOT, //20
}

var Suits = []uint8{
	SuitCircle,
	SuitCross,
	SuitStar,
	SuitTriangle,
	SuitSquare,
	SuitNone,
}

const (
	Card1C = Card(Rank1 | SuitCircle)
	Card1X = Card(Rank1 | SuitCross)
	Card1T = Card(Rank1 | SuitTriangle)
	Card1S = Card(Rank1 | SuitStar)
	Card1R = Card(Rank1 | SuitSquare)

	Card2C = Card(Rank2 | SuitCircle)
	Card2X = Card(Rank2 | SuitCross)
	Card2T = Card(Rank2 | SuitTriangle)
	Card2S = Card(Rank2 | SuitStar)
	Card2R = Card(Rank2 | SuitSquare)

	Card3C = Card(Rank3 | SuitCircle)
	Card3X = Card(Rank3 | SuitCross)
	Card3T = Card(Rank3 | SuitTriangle)
	Card3S = Card(Rank3 | SuitStar)
	Card3R = Card(Rank3 | SuitSquare)

	Card4C = Card(Rank4 | SuitCircle)
	Card4X = Card(Rank4 | SuitCross)
	Card4T = Card(Rank4 | SuitTriangle)
	Card4S = Card(Rank4 | SuitStar)
	Card4R = Card(Rank4 | SuitSquare)

	Card5C = Card(Rank5 | SuitCircle)
	Card5X = Card(Rank5 | SuitCross)
	Card5T = Card(Rank5 | SuitTriangle)
	Card5S = Card(Rank5 | SuitStar)
	Card5R = Card(Rank5 | SuitSquare)

	Card6C = Card(Rank6 | SuitCircle)
	Card6X = Card(Rank6 | SuitCross)
	Card6T = Card(Rank6 | SuitTriangle)
	Card6S = Card(Rank6 | SuitStar)
	Card6R = Card(Rank6 | SuitSquare)

	Card7C = Card(Rank7 | SuitCircle)
	Card7X = Card(Rank7 | SuitCross)
	Card7T = Card(Rank7 | SuitTriangle)
	Card7S = Card(Rank7 | SuitStar)
	Card7R = Card(Rank7 | SuitSquare)

	Card8C = Card(Rank8 | SuitCircle)
	Card8X = Card(Rank8 | SuitCross)
	Card8T = Card(Rank8 | SuitTriangle)
	Card8S = Card(Rank8 | SuitStar)
	Card8R = Card(Rank8 | SuitSquare)

	Card9C = Card(Rank9 | SuitCircle)
	Card9X = Card(Rank9 | SuitCross)
	Card9T = Card(Rank9 | SuitTriangle)
	Card9S = Card(Rank9 | SuitStar)
	Card9R = Card(Rank9 | SuitSquare)

	Card10C = Card(Rank10 | SuitCircle)
	Card10X = Card(Rank10 | SuitCross)
	Card10T = Card(Rank10 | SuitTriangle)
	Card10S = Card(Rank10 | SuitStar)
	Card10R = Card(Rank10 | SuitSquare)

	Card11C = Card(Rank11 | SuitCircle)
	Card11X = Card(Rank11 | SuitCross)
	Card11T = Card(Rank11 | SuitTriangle)
	Card11S = Card(Rank11 | SuitStar)
	Card11R = Card(Rank11 | SuitSquare)

	Card12C = Card(Rank12 | SuitCircle)
	Card12X = Card(Rank12 | SuitCross)
	Card12T = Card(Rank12 | SuitTriangle)
	Card12S = Card(Rank12 | SuitStar)
	Card12R = Card(Rank12 | SuitSquare)

	Card13C = Card(Rank13 | SuitCircle)
	Card13X = Card(Rank13 | SuitCross)
	Card13T = Card(Rank13 | SuitTriangle)
	Card13S = Card(Rank13 | SuitStar)
	Card13R = Card(Rank13 | SuitSquare)

	Card14C = Card(Rank14 | SuitCircle)
	Card14X = Card(Rank14 | SuitCross)
	Card14T = Card(Rank14 | SuitTriangle)
	Card14S = Card(Rank14 | SuitStar)
	Card14R = Card(Rank14 | SuitSquare)

	CardWHOT = Card(RankWHOT)
)

var mapStringRanks = map[uint8]string{
	0x10: "1", 0x20: "2", 0x30: "3", 0x40: "4",
	0x50: "5", 0x60: "6", 0x70: "7", 0x80: "8",
	0x90: "9", 0xA0: "10", 0xB0: "11", 0xC0: "12",
	0xD0: "13", 0xE0: "14", RankWHOT: "WHOT",
}

var mapStringSuits = map[uint8]string{
	SuitCircle:   "CIRCLE",
	SuitCross:    "CROSS",
	SuitStar:     "STAR",
	SuitTriangle: "TRIANGLE",
	SuitSquare:   "SQUARE",
	SuitNone:     "NONE",
}

var mapRanks = map[pb.CardRank]uint8{
	pb.CardRank_RANK_1:  Ranks[0],
	pb.CardRank_RANK_2:  Ranks[1],
	pb.CardRank_RANK_3:  Ranks[2],
	pb.CardRank_RANK_4:  Ranks[3],
	pb.CardRank_RANK_5:  Ranks[4],
	pb.CardRank_RANK_7:  Ranks[6],
	pb.CardRank_RANK_8:  Ranks[7],
	pb.CardRank_RANK_10: Ranks[9],
	pb.CardRank_RANK_11: Ranks[10],
	pb.CardRank_RANK_12: Ranks[11],
	pb.CardRank_RANK_13: Ranks[12],
	pb.CardRank_RANK_14: Ranks[13],
	pb.CardRank_RANK_20: RankWHOT,
}

var mapSuits = map[pb.CardSuit]uint8{
	pb.CardSuit_SUIT_CIRCLE:   SuitCircle,
	pb.CardSuit_SUIT_CROSS:    SuitCross,
	pb.CardSuit_SUIT_SQUARE:   SuitSquare,
	pb.CardSuit_SUIT_STAR:     SuitStar,
	pb.CardSuit_SUIT_TRIANGLE: SuitTriangle,
}

type CardEffect int

const (
	EffectNone          CardEffect = iota
	EffectHoldOn                   // Số 1 - Hold On
	EffectPickTwo                  // Số 2 - Pick Two
	EffectPickThree                // Số 5 - Pick Three
	EffectSuspension               // Số 8 - Suspension
	EffectGeneralMarket            // Số 14 - General Market
	EffectWhot                     // Số 20 - Whot
)

// Giá trị của các lá đặc biệt
const (
	CardValueHoldOn        = 1
	CardValuePickTwo       = 2
	CardValuePickThree     = 5
	CardValueSuspension    = 8
	CardValueGeneralMarket = 14
	CardValueWhot          = 20
)

func NewCardFromPb(rank pb.CardRank, suit pb.CardSuit) Card {
	card := uint8(0)
	card |= mapRanks[rank]
	card |= mapSuits[suit]
	return Card(card)
}

func NewCardFromUint(c uint) Card {
	return Card(c)
}

func NewCard(rank uint8, suit uint8) Card {
	if rank == RankWHOT {
		return Card(rank)
	}
	return Card(rank | suit)
}

func (c Card) GetRank() uint8 {
	return uint8(c & 0xF0)
}

func (c Card) GetSuit() uint8 {
	if c.GetRank() == RankWHOT {
		return SuitNone
	}
	return uint8(c & 0x0F)
}

func CalculateCardValue(card *pb.Card) int {
	value := int(card.GetRank())

	// Lá Star có giá trị gấp đôi
	if card.GetSuit() == pb.CardSuit(SuitStar) {
		value *= 2
	}

	return value
}

func (c Card) String() string {
	return fmt.Sprintf("Rank: %s, Suit: %s", mapStringRanks[c.GetRank()], mapStringSuits[c.GetSuit()])
}
