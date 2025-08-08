package entity

import (
	"testing"

	pb "github.com/nk-nigeria/cgp-common/proto"
)

func TestCard_NewCard(t *testing.T) {
	t.Logf("test card")

	card := NewCardFromPb(pb.WhotCardRank_WHOT_RANK_2, pb.WhotCardSuit_WHOT_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.WhotCardRank_WHOT_RANK_3, pb.WhotCardSuit_WHOT_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.WhotCardRank_WHOT_RANK_4, pb.WhotCardSuit_WHOT_SUIT_CIRCLE)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.WhotCardRank_WHOT_RANK_5, pb.WhotCardSuit_WHOT_SUIT_CIRCLE)
	t.Logf("%b", card)
}
