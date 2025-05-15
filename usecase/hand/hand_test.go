package hand

import (
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/stretchr/testify/assert"
)

func mockHand1() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func mockHand2() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			// Front
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			// Middle
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			// Back
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func TestHand(t *testing.T) {
	t.Logf("Test Hand")

	h1, err := mockHand1()
	if err != nil {
		t.Errorf("invalid hand %v", err)
	}

	for _, card := range h1.GetCards() {
		t.Logf("hand %v", card)
	}

	// test calculate
	h1.calculatePoint()
	t.Logf("caculate front %s", h1.frontHand.Point)
	t.Logf("caculate middle %s", h1.middleHand.Point)
	t.Logf("caculate back %s", h1.backHand.Point)

	// test compare
	h2, err := mockHand2()
	if err != nil {
		t.Errorf("invalid hand %v", err)
	}

	for _, card := range h2.GetCards() {
		t.Logf("hand2 %v", card)
	}

	// test calculate
	h2.calculatePoint()
	t.Logf("caculate front %v", h2.frontHand.Point)
	t.Logf("caculate middle %v", h2.middleHand.Point)
	t.Logf("caculate back %v", h2.backHand.Point)

	//t.Logf("compare result: %v", comp)
}

func TestCheck(t *testing.T) {
	t.Logf("check begin")

	unsortCard := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	sortedCard := SortCard(entity.NewListCard(unsortCard))
	t.Logf("sorted %v", sortedCard)

	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckFlush(entity.NewBinListCards(entity.NewListCard(cards))); !ok {
		t.Errorf("wrong check flush")
	} else {
		t.Logf("check flush ok")
	}

	if _, ok := CheckStraight(entity.NewBinListCards(entity.NewListCard(cards))); !ok {
		t.Errorf("wrong check straight")
	} else {
		t.Logf("check straight ok")
	}

	if _, ok := CheckStraightFlush(entity.NewBinListCards(entity.NewListCard(cards))); !ok {
		t.Errorf("wrong check straight flush")
	} else {
		t.Logf("check straight flush ok")
	}

	fourOfAKindCards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckFourOfAKind(entity.NewBinListCards(entity.NewListCard(fourOfAKindCards))); !ok {
		t.Errorf("wrong check four of a kind")
	} else {
		t.Logf("check four of a kind ok")
	}

	fullHouseCards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckFullHouse(entity.NewBinListCards(entity.NewListCard(fullHouseCards))); !ok {
		t.Errorf("wrong check full house card")
	} else {
		t.Logf("check full house ok")
	}
}

func TestTwoPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckTwoPairs(entity.NewBinListCards(entity.NewListCard(cards))); !ok {
		t.Errorf("wrong check two pairs")
	} else {
		t.Logf("check two pairs ok")
	}
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// Same level card
func TestCompareBasicStraightFlushVsStraightFlushDraw(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// list card 1 higher
func TestCompareBasicStraightFlushHigherStraightFlush(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// list card 1 lower
func TestCompareBasicStraightFlushLowerStraightFlushLower(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(-1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Tứ quý (en: Four of a Kind)
func TestCompareBasicStraightFlushVsFourOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Cù lũ (en: Full House
func TestCompareBasicStraightFlushVsFullhouse(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng (en: Flush)
func TestCompareBasicStraightFlushVsFlush(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Sảnh (en: Straight)
func TestCompareBasicStraightFlushVsStraight(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Xám chi/Xám cô (en: Three of a Kind)
func TestCompareBasicStraightFlushVsThreeOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thú (en: Two Pairs)
func TestCompareBasicStraightFlushVsTwoPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Đôi (en: Pair)
func TestCompareBasicStraightFlushVsPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Mậu Thầu (en: High Card)
func TestCompareBasicStraightFlushVsHighCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Tứ quý (en: Four of a Kind)
// Same level card
func TestCompareFourOfAKindVsFourOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Tứ quý (en: Four of a Kind) vs Cù lũ (en: Full House)
func TestCompareFourOfAKindVsFullHouse(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng (en: Flush) vs Thùng (en: Flush)
func TestCompareFlushVsFlushHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_6,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// No contain A card
func TestCompareStraightVsStraightNoACardEqual(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

func TestCompareStraightVsStraightNoACardHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card
func TestCompareStraightVsStraightContainACardEqual(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card, No card K
func TestCompareStraightVsStraightContainACardNotCardKLower(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_6,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(-1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card, contain K card
func TestCompareStraightVsStraightContainACardKCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Xám chi/Xám cô (en: Three of a Kind) vs Xám chi/Xám cô (en: Three of a Kind)
func TestCompareThreeOfAKindVsThreeOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairDraw(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher1(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher2(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher3(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Đôi (en: Pair) vs Đôi (en: Pair)
func TestComparePairVsPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Đôi (en: Pair) vs Đôi (en: Pair)
func TestComparePairVsPairHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Mậu Thầu (en: High Card) vs Mậu Thầu (en: High Card)
func TestCompareHighCardVsHighCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

func TestCompareHighCardVsHighCardHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards), kBackHand)
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards), kBackHand)
	point := strainghtFlush1.Compare(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

//
