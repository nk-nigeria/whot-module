package entity

import (
	"testing"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
)

func TestCard_NewCard(t *testing.T) {
	t.Logf("test card")

	card := NewCardFromPb(pb.CardRank_RANK_2, pb.CardSuit_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_3, pb.CardSuit_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_4, pb.CardSuit_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_5, pb.CardSuit_SUIT_CIRCLE)
	t.Logf("%b", card)
}
