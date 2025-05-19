package hand

import (
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
)

func TestCheckCleanDragon(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_11,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_12,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_13,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_14,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
	}

	if _, ok := CheckCleanDragon(entity.NewBinListCards(entity.NewListCard(cards))); ok {
		t.Logf("check clean dragon ok")
	} else {
		t.Logf("check clean dragon failed")
	}
}

func TestCheckFullColor(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
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
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_11,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_12,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_13,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_1,
			Suit: pb.CardSuit_SUIT_STAR,
		},
	}

	if _, ok := CheckFullColor(entity.NewBinListCards(entity.NewListCard(cards))); ok {
		t.Logf("check full color ok")
	} else {
		t.Logf("check full color failed")
	}
}

func TestCheckDragon(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_11,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_12,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_13,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_1,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
	}

	if _, ok := CheckDragon(entity.NewBinListCards(entity.NewListCard(cards))); ok {
		t.Logf("check dragon ok")
	} else {
		t.Logf("check dragon failed")
	}
}

func TestCheckSixPairs(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_STAR,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_STAR,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_STAR,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_STAR,
		},
		{
			Rank: pb.CardRank_RANK_11,
			Suit: pb.CardSuit_SUIT_TRIANGLE,
		},
		{
			Rank: pb.CardRank_RANK_11,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_13,
			Suit: pb.CardSuit_SUIT_CROSS,
		},
		{
			Rank: pb.CardRank_RANK_13,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
		{
			Rank: pb.CardRank_RANK_1,
			Suit: pb.CardSuit_SUIT_CIRCLE,
		},
	}

	if _, ok := CheckSixPairs(entity.NewBinListCards(entity.NewListCard(cards))); ok {
		t.Logf("check six pairs ok")
	} else {
		t.Logf("check six pairs failed")
	}
}
