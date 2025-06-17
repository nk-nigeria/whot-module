package bin_list_card

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/nk-nigeria/whot-module/entity"
)

var (
	kRed   = uint8(0)
	kBlack = uint8(1)
)

var (
	BitSetRankMap map[uint8]*bitset.BitSet
	BitSetSuitMap map[uint8]*bitset.BitSet
	BitSetColor   map[uint8]*bitset.BitSet
)

func init() {
	BitSetRankMap = make(map[uint8]*bitset.BitSet)
	for _, rank := range entity.Ranks {
		BitSet := bitset.New(4)
		for _, suit := range entity.Suits {
			BitSet.Set(uint(entity.NewCard(rank, suit)))
		}

		BitSetRankMap[rank] = BitSet
	}

	BitSetSuitMap = make(map[uint8]*bitset.BitSet)
	for _, suit := range entity.Suits {
		BitSet := bitset.New(16)
		for _, rank := range entity.Ranks {
			BitSet.Set(uint(entity.NewCard(rank, suit)))
		}
		BitSetSuitMap[suit] = BitSet
	}

	BitSetColor = make(map[uint8]*bitset.BitSet)
	BitSetColor[kRed] = BitSetSuitMap[entity.SuitCircle].Union(BitSetSuitMap[entity.SuitCross])
	BitSetColor[kBlack] = BitSetSuitMap[entity.SuitSquare].Union(BitSetSuitMap[entity.SuitTriangle])
}

type service struct {
}

func (s service) GetChain(b *entity.BinListCard, comb int) (uint, entity.ListCard) {
	switch comb {
	case CombineFour:
		return s.lookupFour(b)
	case CombineThree:
		return s.lookupThree(b)
	case CombinePair:
		return s.lookupTwo(b)
	case CombineStraight:
		return s.lookupStraight(b)
	case CombineFullHouse:
		return s.lookupFullHouse(b)
	case CombineFlush:
		return s.lookupFlush(b)
	case CombineFullColor:
		return s.lookupFullColor(b)
	}

	return 0, nil
}

var s = &service{}

func NewChinesePokerBinList() ChinesePokerBinList {
	return s
}
