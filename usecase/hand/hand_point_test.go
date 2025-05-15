package hand

import (
	"testing"

	"github.com/nakamaFramework/whot-module/entity"
	"github.com/stretchr/testify/assert"
)

func mockHighCard() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card5C,
		entity.Card9D,
		entity.Card7H,
		entity.Card10C,
		entity.Card2H,
	})
}

func mockStraightFlushSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card5C,
		entity.Card6C,
		entity.Card7C,
		entity.Card8C,
		entity.Card9C,
	})
}

func TestCheckStraightFlush(t *testing.T) {
	t.Logf("test TestCheckStraightFlush")

	var ok bool
	_, ok = CheckStraightFlush(mockStraightFlushSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckStraightFlush(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockFourOfAKindSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card5C,
		entity.Card5D,
		entity.Card5S,
		entity.Card5H,
		entity.Card9C,
	})
}

func TestCheckFourOfAKind(t *testing.T) {
	t.Logf("test TestCheckFourOfAKind")

	var ok bool
	_, ok = CheckFourOfAKind(mockFourOfAKindSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckFourOfAKind(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockFullHouseSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card3D,
		entity.Card2S,
		entity.Card2H,
		entity.Card2C,
	})
}

func TestCheckFullHouse(t *testing.T) {
	t.Logf("test TestCheckFullHouse")

	var ok bool
	_, ok = CheckFullHouse(mockFullHouseSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckFullHouse(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockFlushSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card5C,
		entity.Card7C,
		entity.Card9C,
		entity.CardJC,
	})
}

func TestFlush(t *testing.T) {
	t.Logf("test TestCheckFlush")

	var ok bool
	_, ok = CheckFlush(mockFlushSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckFlush(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockStraightSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card4H,
		entity.Card5S,
		entity.Card6C,
		entity.Card7S,
	})
}

func mockStraightSuccess2() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.CardAC,
		entity.Card2H,
		entity.Card3S,
		entity.Card4C,
		entity.Card5S,
	})
}

func TestStraight(t *testing.T) {
	t.Logf("test TestStraight")

	var ok bool
	_, ok = CheckStraight(mockStraightSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckStraight(mockStraightSuccess2())
	assert.Equal(t, true, ok)

	_, ok = CheckStraight(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockThreeOfAKindSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card3H,
		entity.Card3S,
		entity.Card6C,
		entity.Card7S,
	})
}

func TestThreeOfAKind(t *testing.T) {
	t.Logf("test TestThreeOfAKind")

	var ok bool
	_, ok = CheckThreeOfAKind(mockThreeOfAKindSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckThreeOfAKind(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockTwoPairsSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card3H,
		entity.Card4S,
		entity.Card4C,
		entity.Card7S,
	})
}

func TestTwoPairs(t *testing.T) {
	t.Logf("test TestTwoPairs")

	var ok bool
	_, ok = CheckTwoPairs(mockTwoPairsSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckTwoPairs(mockHighCard())
	assert.Equal(t, false, ok)
}

func mockPairSuccess() *entity.BinListCard {
	return entity.NewBinListCards(entity.ListCard{
		entity.Card3C,
		entity.Card3H,
		entity.Card4S,
		entity.Card5C,
		entity.Card7S,
	})
}

func TestPair(t *testing.T) {
	t.Logf("test TestPairs")

	var ok bool
	_, ok = CheckPair(mockPairSuccess())
	assert.Equal(t, true, ok)

	_, ok = CheckTwoPairs(mockHighCard())
	assert.Equal(t, false, ok)
}
