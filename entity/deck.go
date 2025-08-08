package entity

import (
	"errors"
	"math/rand"

	pb "github.com/nk-nigeria/cgp-common/proto"
)

const MaxCard = 54

type Deck struct {
	Cards *pb.ListCard
	Dealt int
}

func NewDeck() *Deck {

	validCards := map[pb.WhotCardSuit][]pb.WhotCardRank{
		pb.WhotCardSuit_WHOT_SUIT_CIRCLE: {
			pb.WhotCardRank_WHOT_RANK_1, pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardRank_WHOT_RANK_4, pb.WhotCardRank_WHOT_RANK_5,
			pb.WhotCardRank_WHOT_RANK_7, pb.WhotCardRank_WHOT_RANK_8,
			pb.WhotCardRank_WHOT_RANK_10, pb.WhotCardRank_WHOT_RANK_11, pb.WhotCardRank_WHOT_RANK_12, pb.WhotCardRank_WHOT_RANK_13, pb.WhotCardRank_WHOT_RANK_14,
		},
		pb.WhotCardSuit_WHOT_SUIT_TRIANGLE: {
			pb.WhotCardRank_WHOT_RANK_1, pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardRank_WHOT_RANK_4, pb.WhotCardRank_WHOT_RANK_5,
			pb.WhotCardRank_WHOT_RANK_7, pb.WhotCardRank_WHOT_RANK_8,
			pb.WhotCardRank_WHOT_RANK_10, pb.WhotCardRank_WHOT_RANK_11, pb.WhotCardRank_WHOT_RANK_12, pb.WhotCardRank_WHOT_RANK_13, pb.WhotCardRank_WHOT_RANK_14,
		},
		pb.WhotCardSuit_WHOT_SUIT_CROSS: {
			pb.WhotCardRank_WHOT_RANK_1, pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardRank_WHOT_RANK_5, pb.WhotCardRank_WHOT_RANK_7,
			pb.WhotCardRank_WHOT_RANK_10, pb.WhotCardRank_WHOT_RANK_11, pb.WhotCardRank_WHOT_RANK_13, pb.WhotCardRank_WHOT_RANK_14,
		},

		pb.WhotCardSuit_WHOT_SUIT_SQUARE: {
			pb.WhotCardRank_WHOT_RANK_1, pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardRank_WHOT_RANK_5, pb.WhotCardRank_WHOT_RANK_7,
			pb.WhotCardRank_WHOT_RANK_10, pb.WhotCardRank_WHOT_RANK_11, pb.WhotCardRank_WHOT_RANK_13, pb.WhotCardRank_WHOT_RANK_14,
		},
		pb.WhotCardSuit_WHOT_SUIT_STAR: {
			pb.WhotCardRank_WHOT_RANK_1, pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardRank_WHOT_RANK_4, pb.WhotCardRank_WHOT_RANK_5,
			pb.WhotCardRank_WHOT_RANK_7, pb.WhotCardRank_WHOT_RANK_8,
		},
	}

	cards := &pb.ListCard{}
	for suit, ranks := range validCards {
		for _, rank := range ranks {
			cards.WhotCards = append(cards.WhotCards, &pb.WhotCard{
				Rank: rank,
				Suit: suit,
			})
		}
	}

	for i := 0; i < 5; i++ {
		cards.WhotCards = append(cards.WhotCards, &pb.WhotCard{
			Rank: pb.WhotCardRank_WHOT_RANK_20,
			Suit: pb.WhotCardSuit_WHOT_SUIT_UNSPECIFIED,
		})
	}

	return &Deck{
		Dealt: 0,
		Cards: cards,
	}
}

// Shuffle the deck
func (d *Deck) Shuffle() {
	for i := 1; i < len(d.Cards.WhotCards); i++ {
		// Create a random int up to the number of Cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d.Cards.WhotCards[r], d.Cards.WhotCards[i] = d.Cards.WhotCards[i], d.Cards.WhotCards[r]
		}
	}
}

// Deal a specified amount of Cards
func (d *Deck) Deal(n int, isTopCard bool) (*pb.ListCard, error) {
	remainCountCards := d.RemainingCards()
	if (remainCountCards) <= 0 {
		return nil, errors.New("deal.no-cards-left")
	}
	if remainCountCards < n {
		n = MaxCard - d.Dealt
	}

	if isTopCard {
		retryLimit := 10
		for retry := 0; retry < retryLimit; retry++ {
			if d.Cards.WhotCards[d.Dealt].Rank != pb.WhotCardRank_WHOT_RANK_20 {
				break
			}
			Shuffle(d.Cards)
		}

		if d.Cards.WhotCards[d.Dealt].Rank == pb.WhotCardRank_WHOT_RANK_20 {
			return nil, errors.New("deal.topcard.error-whot-retry-limit")
		}
	}
	var cards pb.ListCard
	for i := 0; i < n; i++ {
		cards.WhotCards = append(cards.WhotCards, d.Cards.WhotCards[d.Dealt])
		d.Dealt++
	}

	return &cards, nil
}

func (d *Deck) RemainingCards() int {
	return MaxCard - d.Dealt
}

// Shuffle any hand
func Shuffle(cards *pb.ListCard) *pb.ListCard {
	for i := 1; i < len(cards.WhotCards); i++ {
		// Create a random int up to the number of Cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			cards.WhotCards[r], cards.WhotCards[i] = cards.WhotCards[i], cards.WhotCards[r]
		}
	}

	return cards
}
