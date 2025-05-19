package entity

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
)

type BinListCard struct {
	b *bitset.BitSet
}

func (b BinListCard) GetBitSet() *bitset.BitSet {
	return b.b
}

func NewBinListCards(cards ListCard) *BinListCard {
	b := bitset.New(MaxCard)
	for _, card := range cards {
		b.Set(uint(card))
	}
	return &BinListCard{
		b: b,
	}
}

func (b BinListCard) String() string {
	var str = "[\n"

	for i, e := b.b.NextSet(0); e; i, e = b.b.NextSet(i + 1) {
		str += fmt.Sprintf("%d\n", i)
	}
	str += "]"

	return str
}

func (b BinListCard) ToList() ListCard {
	return BitSetToListCard(b.b)
}

func BitSetToListCard(b *bitset.BitSet) ListCard {
	cards := ListCard{}
	for i, e := b.NextSet(0); e; i, e = b.NextSet(i + 1) {
		cards = append(cards, NewCardFromUint(i))
	}
	return cards
}

func IsSameListCard(cards1, cards2 ListCard) bool {
	bl1 := NewBinListCards(cards1)
	bl2 := NewBinListCards(cards2)
	return bl1.b.DifferenceCardinality(bl2.b) == 0
}
