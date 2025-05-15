package hand

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

var randinst *rand.Rand

func init() {
	sourceRand := rand.NewSource(time.Now().UnixNano())
	randinst = rand.New(sourceRand)
}

func mockDragon() []*pb.Card {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_6,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	return cards
}

func mockCard(rank, suit int) *pb.Card {
	card := pb.Card{
		Rank: pb.CardRank(rank),
		Suit: pb.CardSuit(suit),
	}
	return &card
}

func mockPair(rank, suit int) []*pb.Card {
	card := pb.Card{
		Rank: pb.CardRank(rank),
		Suit: pb.CardSuit(suit),
	}
	cards := make([]*pb.Card, 0, 2)
	for i := 0; i < 2; i++ {
		cards = append(cards, &card)
	}
	return cards
}

func mockFivePair() []*pb.Card {

	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	return cards
}

func mockThreeOfAKind() []*pb.Card {
	startRank := (randinst.Intn(100) + 1) % 11
	cards := make([]*pb.Card, 0, 3)
	for i := 0; i < 3; i++ {
		card := mockCard(startRank, i)
		cards = append(cards, card)
	}
	return cards

}

func mockJackpot() []*pb.Card {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
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
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	return cards
}
func TestIsDragonSuccess(t *testing.T) {
	listCard := entity.NewListCard(mockDragon())
	handCard, isDragon := CheckDragon(entity.NewBinListCards(listCard))
	assert.Equal(t, true, isDragon)
	assert.NotNil(t, handCard)
}

func TestIsDragonFail(t *testing.T) {
	cards := mockDragon()
	cards[4].Rank = pb.CardRank_RANK_K
	listCard := entity.NewListCard(cards)
	handCard, isDragon := CheckDragon(entity.NewBinListCards(listCard))
	assert.Equal(t, false, isDragon)
	assert.Nil(t, handCard)
}

func TestIsCleanDragonSuccess(t *testing.T) {
	cards := mockDragon()
	for idx, _ := range cards {
		cards[idx].Suit = pb.CardSuit_SUIT_HEARTS
	}
	listCard := entity.NewListCard(cards)
	handCard, isCleanDragon := CheckCleanDragon(entity.NewBinListCards(listCard))
	assert.Equal(t, true, isCleanDragon)
	assert.NotNil(t, handCard)
}

func TestIsCleanDragonFailed(t *testing.T) {
	cards := mockDragon()
	listCard := entity.NewListCard(cards)
	handCard, isCleanDragon := CheckCleanDragon(entity.NewBinListCards(listCard))
	assert.Equal(t, false, isCleanDragon)
	assert.Nil(t, handCard)
}

func TestIsFullColoredSuccess(t *testing.T) {
	cards := mockDragon()
	for idx, _ := range cards {
		cards[idx].Suit = pb.CardSuit_SUIT_HEARTS
	}
	cards[10].Rank = pb.CardRank_RANK_2
	listCard := entity.NewListCard(cards)
	handCard, isFullColor := CheckFullColor(entity.NewBinListCards(listCard))
	assert.Equal(t, true, isFullColor)
	assert.NotNil(t, handCard)
}

func TestIsFullColoredFailed(t *testing.T) {
	cards := mockDragon()

	cards[10].Rank = pb.CardRank_RANK_2
	listCard := entity.NewListCard(cards)
	handCard, isFullColor := CheckFullColor(entity.NewBinListCards(listCard))
	assert.Equal(t, false, isFullColor)
	assert.Nil(t, handCard)
}

func TestIsSixPairSuccess(t *testing.T) {
	cards := mockFivePair()
	cards = append(cards, mockPair(0, 0)...)
	cards = append(cards, &pb.Card{
		Rank: pb.CardRank_RANK_10,
		Suit: pb.CardSuit_SUIT_DIAMONDS,
	})
	listCard := entity.NewListCard(cards)
	handCard, isvalid := CheckSixPairs(entity.NewBinListCards(listCard))
	assert.Equal(t, true, isvalid)
	assert.NotNil(t, handCard)
}

func TestIsSixPairFailed(t *testing.T) {
	cards := mockFivePair()
	cards = append(cards, mockThreeOfAKind()...)
	cards[12].Rank = (cards[12].Rank + 1) % 14
	listCard := entity.NewListCard(cards)
	handCard, isvalid := CheckSixPairs(entity.NewBinListCards(listCard))
	assert.Equal(t, false, isvalid)
	assert.Nil(t, handCard)
}

func TestCheckJackpot(t *testing.T) {
	type args struct {
		childHand *ChildHand
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "TestCheckJackpot",
			args: args{
				childHand: NewChildHand(entity.NewListCard(mockJackpot()), kMidHand),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckJackpot(tt.args.childHand); got != tt.want {
				t.Errorf("CheckJackpot() = %v, want %v", got, tt.want)
			}
		})
	}
}
