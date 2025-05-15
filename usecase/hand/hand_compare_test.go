package hand

import (
	"reflect"
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/stretchr/testify/assert"
)

func mockHandNatural1() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func mockHandNatural2() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})
}

func mockHandNatural3Flush() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})
}

func mockHandNormal() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_A,
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
		},
	})
}

func TestCompareHand(t *testing.T) {

	h1, _ := mockHandNatural1()
	h2, _ := mockHandNatural2()
	ctx := NewCompareContext(2)
	result := CompareHand(ctx, h1, h2)
	r1 := pb.ComparisonResult{
		ScoreResult: &pb.ScoreResult{},
		PointResult: &pb.PointResult{},
	}
	r2 := pb.ComparisonResult{
		ScoreResult: &pb.ScoreResult{},
		PointResult: &pb.PointResult{},
	}

	ProcessCompareResult(ctx, &r1, result.GetR1())
	ProcessCompareResult(ctx, &r2, result.GetR2())
	t.Logf("result %v", result)
	t.Logf("result %v", result.bonuses)
	t.Logf("r1 %v", r1.ScoreResult)
	t.Logf("r1 %v", r1.PointResult)
	t.Logf("r2 %v", r2.ScoreResult)
	t.Logf("r2 %v", r2.PointResult)

}

func TestHand_CompareHand(t *testing.T) {
	type fields struct {
		cards        entity.ListCard
		ranking      pb.HandRanking
		frontHand    *ChildHand
		middleHand   *ChildHand
		backHand     *ChildHand
		naturalPoint *HandPoint
		pointType    pb.PointType
		calculated   bool
		owner        string
	}
	type args struct {
		h2 *Hand
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ComparisonResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Hand{
				cards:        tt.fields.cards,
				ranking:      tt.fields.ranking,
				frontHand:    tt.fields.frontHand,
				middleHand:   tt.fields.middleHand,
				backHand:     tt.fields.backHand,
				naturalPoint: tt.fields.naturalPoint,
				pointType:    tt.fields.pointType,
				calculated:   tt.fields.calculated,
				owner:        tt.fields.owner,
			}
			if got := h.CompareHand(tt.args.h2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hand.CompareHand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareHand3FlushWithNormal(t *testing.T) {

	h1, _ := mockHandNatural3Flush()
	h2, _ := mockHandNormal()

	result := CompareHand(nil, h1, h2)

	t.Logf("result %v", result)
	assert.Equal(t, pb.PointType_Point_Natural, h1.pointType, "should be natural, 3flush")
	assert.Equal(t, pb.PointType_Point_Normal, h2.pointType, "should be normal")
	assert.Equal(t, mapNaturalPoint[pb.HandBonusType_BonusNaturalThreeOfFlushes], result.GetR1().NaturalFactor, "should be natural, 3flush")
	assert.Equal(t, -mapNaturalPoint[pb.HandBonusType_BonusNaturalThreeOfFlushes], result.GetR2().NaturalFactor, "should be normal")

	assert.Equal(t, 0, result.GetR1().BackFactor, "should be 0")
	assert.Equal(t, 0, result.GetR1().MiddleFactor, "should be 0")
	assert.Equal(t, 0, result.GetR1().FrontFactor, "should be 0")
	assert.Equal(t, 0, result.GetR1().BackBonusFactor, "should be 0")
	assert.Equal(t, 0, result.GetR1().MiddleBonusFactor, "should be 0")
	assert.Equal(t, 0, result.GetR1().FrontBonusFactor, "should be 0")

	assert.Equal(t, 0, result.GetR2().BackFactor, "should be 0")
	assert.Equal(t, 0, result.GetR2().MiddleFactor, "should be 0")
	assert.Equal(t, 0, result.GetR2().FrontFactor, "should be 0")
	assert.Equal(t, 0, result.GetR2().BackBonusFactor, "should be 0")
	assert.Equal(t, 0, result.GetR2().MiddleBonusFactor, "should be 0")
	assert.Equal(t, 0, result.GetR2().FrontBonusFactor, "should be 0")
}

func TestCompareNormalWithNormal(t *testing.T) {
	h1, _ := NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})

	h2, _ := NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})
	result := CompareHand(nil, h1, h2)
	assert.Equal(t, 1, result.r1.BackFactor, "win")
	assert.Equal(t, 1, result.r1.MiddleFactor, "win")
	assert.Equal(t, 1, result.r1.FrontFactor, "win")
	assert.Equal(t, 4, result.r1.BackBonusFactor, "Four of a Kind bonus")
	assert.Equal(t, 0, result.r1.MiddleBonusFactor, "win no bobus")
	assert.Equal(t, 0, result.r1.FrontBonusFactor, "win no bobus")

	assert.Equal(t, -1, result.r2.BackFactor, "lose")
	assert.Equal(t, -1, result.r2.MiddleFactor, "lose")
	assert.Equal(t, -1, result.r2.FrontFactor, "lose")
	assert.Equal(t, -4, result.r2.BackBonusFactor, "Four of a Kind lose bonus")
	assert.Equal(t, 0, result.r2.MiddleBonusFactor, "lose no bobus")
	assert.Equal(t, 0, result.r2.FrontBonusFactor, "lose no bobus")
	t.Logf("result %v", result)
}
