package hand

import (
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/stretchr/testify/assert"
)

func mockHandDontMissets() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},

			{
				Rank: pb.CardRank_RANK_11,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_12,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_STAR,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
		},
	})
}

func mockHandMissets1() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CROSS,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},

			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
		},
	})
}

func mockHandMissets2() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CROSS,
			},

			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_11,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CROSS,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
		},
	})
}

func mockHandDontMissets2() (*Hand, error) {
	return NewHand(entity.ListCard{
		// entity.Card8H,
		// entity.Card5C,
		// entity.Card10H,

		// entity.Card8C,
		// entity.Card5S,
		// entity.Card7D,
		// entity.Card6C,
		// entity.Card4H,

		// entity.Card2H,
		// entity.Card2C,
		// entity.Card2D,
		// entity.Card3D,
		// entity.Card3S,
	})
}

func TestIsMisSets(t *testing.T) {
	t.Logf("test is mis sets")
	var h1 *Hand
	var mis bool
	// h1, _ = mockHandMissets1()
	// h1.calculatePoint()

	// var mis bool
	// mis = IsMisSets(h1)
	// assert.Equal(t, true, mis)

	// h1, _ = mockHandMissets2()
	// h1.calculatePoint()
	// assert.Equal(t, true, mis)

	// h1, _ = mockHandDontMissets()
	// h1.calculatePoint()
	// mis = IsMisSets(h1)
	// assert.Equal(t, false, mis)

	h1, _ = mockHandDontMissets2()
	h1.calculatePoint()

	t.Logf("front %s", h1.frontHand.Point)
	t.Logf("middle %s", h1.middleHand.Point)
	t.Logf("back %s", h1.backHand.Point)

	mis = IsMisSets(h1)
	assert.Equal(t, false, mis)
}
