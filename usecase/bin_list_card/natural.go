package bin_list_card

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/nk-nigeria/whot-module/entity"
)

func (s service) lookupFullColor(b *entity.BinListCard) (uint, entity.ListCard) {
	var black *bitset.BitSet
	var red *bitset.BitSet

	red = b.GetBitSet().Intersection(BitSetColor[kRed])
	black = b.GetBitSet().Intersection(BitSetColor[kBlack])

	if red.Count() >= 12 {
		remain := b.GetBitSet().Difference(red)
		return 1, createResult(remain.Count()+red.Count(), remain, red)
	}

	if black.Count() >= 12 {
		remain := b.GetBitSet().Difference(black)
		return 1, createResult(remain.Count()+black.Count(), remain, black)
	}

	return 0, nil
}
