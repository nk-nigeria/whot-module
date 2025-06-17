package entity

import (
	"errors"
	"math/rand"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
)

const MaxCard = 54

type Deck struct {
	Cards *pb.ListCard
	Dealt int
}

func NewDeck() *Deck {

	validCards := map[pb.CardSuit][]pb.CardRank{
		pb.CardSuit_SUIT_CIRCLE: {
			pb.CardRank_RANK_1, pb.CardRank_RANK_2, pb.CardRank_RANK_3, pb.CardRank_RANK_4, pb.CardRank_RANK_5,
			pb.CardRank_RANK_7, pb.CardRank_RANK_8,
			pb.CardRank_RANK_10, pb.CardRank_RANK_11, pb.CardRank_RANK_12, pb.CardRank_RANK_13, pb.CardRank_RANK_14,
		},
		pb.CardSuit_SUIT_TRIANGLE: {
			pb.CardRank_RANK_1, pb.CardRank_RANK_2, pb.CardRank_RANK_3, pb.CardRank_RANK_4, pb.CardRank_RANK_5,
			pb.CardRank_RANK_7, pb.CardRank_RANK_8,
			pb.CardRank_RANK_10, pb.CardRank_RANK_11, pb.CardRank_RANK_12, pb.CardRank_RANK_13, pb.CardRank_RANK_14,
		},
		pb.CardSuit_SUIT_CROSS: {
			pb.CardRank_RANK_1, pb.CardRank_RANK_2, pb.CardRank_RANK_3, pb.CardRank_RANK_5, pb.CardRank_RANK_7,
			pb.CardRank_RANK_10, pb.CardRank_RANK_11, pb.CardRank_RANK_13, pb.CardRank_RANK_14,
		},
		pb.CardSuit_SUIT_SQUARE: {
			pb.CardRank_RANK_1, pb.CardRank_RANK_2, pb.CardRank_RANK_3, pb.CardRank_RANK_5, pb.CardRank_RANK_7,
			pb.CardRank_RANK_10, pb.CardRank_RANK_11, pb.CardRank_RANK_13, pb.CardRank_RANK_14,
		},
		pb.CardSuit_SUIT_STAR: {
			pb.CardRank_RANK_1, pb.CardRank_RANK_2, pb.CardRank_RANK_3, pb.CardRank_RANK_4, pb.CardRank_RANK_5,
			pb.CardRank_RANK_7, pb.CardRank_RANK_8,
		},
	}

	cards := &pb.ListCard{}
	for suit, ranks := range validCards {
		for _, rank := range ranks {
			cards.Cards = append(cards.Cards, &pb.Card{
				Rank: rank,
				Suit: suit,
			})
		}
	}

	for i := 0; i < 5; i++ {
		cards.Cards = append(cards.Cards, &pb.Card{
			Rank: pb.CardRank_RANK_20,
			Suit: pb.CardSuit_SUIT_UNSPECIFIED,
		})
	}

	return &Deck{
		Dealt: 0,
		Cards: cards,
	}
}

// Shuffle the deck
func (d *Deck) Shuffle() {
	for i := 1; i < len(d.Cards.Cards); i++ {
		// Create a random int up to the number of Cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d.Cards.Cards[r], d.Cards.Cards[i] = d.Cards.Cards[i], d.Cards.Cards[r]
		}
	}
}

// Deal a specified amount of Cards
func (d *Deck) Deal(n int, isTopCard bool) (*pb.ListCard, error) {
	if (MaxCard - d.Dealt) < n {
		return nil, errors.New("deck.deal.error-not-enough")
	}

	if isTopCard {
		retryLimit := 10
		for retry := 0; retry < retryLimit; retry++ {
			if d.Cards.Cards[d.Dealt].Rank != pb.CardRank_RANK_20 {
				break
			}
			Shuffle(d.Cards)
		}

		if d.Cards.Cards[d.Dealt].Rank == pb.CardRank_RANK_20 {
			return nil, errors.New("deal.topcard.error-whot-retry-limit")
		}
	}
	var cards pb.ListCard
	for i := 0; i < n; i++ {
		cards.Cards = append(cards.Cards, d.Cards.Cards[d.Dealt])
		d.Dealt++
	}

	return &cards, nil
}

func (d *Deck) RemainingCards() int {
	return MaxCard - d.Dealt
}

// Shuffle any hand
func Shuffle(cards *pb.ListCard) *pb.ListCard {
	for i := 1; i < len(cards.Cards); i++ {
		// Create a random int up to the number of Cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			cards.Cards[r], cards.Cards[i] = cards.Cards[i], cards.Cards[r]
		}
	}

	return cards
}
