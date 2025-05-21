package engine

import (
	"testing"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
)

func TestGame(t *testing.T) {
	t.Logf("Test Game")
	processor := NewWhotPokerEngine()

	// mock presense
	presense := linkedhashmap.New()
	presense.Put("user1", nil)
	presense.Put("user2", nil)
	// presense.Put("user3", nil)

	// mock state
	state := &entity.MatchState{
		Presences:        presense,
		PlayingPresences: linkedhashmap.New(),
		Cards:            make(map[string]*pb.ListCard),
		OrganizeCards:    make(map[string]*pb.ListCard),
		Label: &entity.MatchLabel{
			Open:         1,
			Bet:          100,
			Code:         "test",
			Name:         "name",
			Password:     "password",
			MaxSize:      4,
			MockCodeCard: 1,
		},
	}

	// var err = processor.NewGame(state)
	// if err != nil {
	// 	t.Errorf("new game error %v", err)
	// }

	t.Logf("new game success")
	processor.Deal(state)

	// check dealt cards
	for u, cards := range state.Cards {
		t.Logf("card %v, %v", u, cards)
	}

	card1 := &pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_1,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},

			{
				Rank: pb.CardRank_RANK_13,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_14,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_10,
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
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
		},
	}
	// card2 := state.Cards["user2"]
	card2 := &pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_1,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CROSS,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_12,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_14,
				Suit: pb.CardSuit_SUIT_STAR,
			},

			{
				Rank: pb.CardRank_RANK_12,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_20,
				Suit: pb.CardSuit_SUIT_UNSPECIFIED,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
		},
	}
	// card3 := state.Cards["user3"]
	state.Cards["user1"] = card1
	state.Cards["user2"] = card2
	// cardOrganize1 := entity.Shuffle(card1)
	// cardOrganize2 := entity.Shuffle(card2)
	// cardOrganize3 := entity.Shuffle(card3)
	cardOrganize1 := (card1)
	cardOrganize2 := (card2)
	processor.Organize(state, "user1", cardOrganize1)
	processor.Organize(state, "user2", cardOrganize2)
	// processor.Organize(state, "user3", cardOrganize3)
	state.PlayingPresences.Put("user1", "user1")
	state.PlayingPresences.Put("user2", "user2")
	result := processor.Finish(state)
	t.Logf("%v", result)
	// check dealt cards
	// for u, cards := range state.OrganizeCards {
	// 	t.Logf("card organize %v, %v", u, cards)
	// }

	// for u, cards := range state.Cards {
	// 	t.Logf("card2 %v, %v", u, cards)
	// }
}

func TestGameNormalWithNorma(t *testing.T) {
	listCard1 := (&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_1,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_11,
				Suit: pb.CardSuit_SUIT_CROSS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_20,
				Suit: pb.CardSuit_SUIT_UNSPECIFIED,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_STAR,
			},
		},
	})

	listCard2 := (&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_1,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_STAR,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_12,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_12,
				Suit: pb.CardSuit_SUIT_SQUARE,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_STAR,
			},

			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_TRIANGLE,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CIRCLE,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_CROSS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_STAR,
			},
		},
	})

	processor := NewWhotPokerEngine()

	// mock presense
	presense := linkedhashmap.New()
	presense.Put("user1", nil)
	presense.Put("user2", nil)

	// mock state
	state := &entity.MatchState{
		Presences:        presense,
		PlayingPresences: presense,
		OrganizeCards:    make(map[string]*pb.ListCard),
		Cards:            make(map[string]*pb.ListCard),
	}

	//var err = processor.NewGame(state)
	//if err != nil {
	//	t.Errorf("new game error %v", err)
	//}

	processor.Organize(state, "user1", listCard1)
	processor.Organize(state, "user2", listCard2)

	result := processor.Finish(state)
	t.Logf("%v", result)
}
