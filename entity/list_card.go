package entity

import pb "github.com/nk-nigeria/cgp-common/proto/whot"

type ListCard []Card

func NewListCard(list []*pb.Card) ListCard {
	newList := ListCard{}
	for _, card := range list {
		newList = append(newList, NewCardFromPb(card.GetRank(), card.GetSuit()))
	}

	return newList
}

func NewListCardWithSize(size uint) ListCard {
	l := make([]Card, 0, size)
	return l
}
