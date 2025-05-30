package bin_list_card

import "github.com/nakama-nigeria/whot-module/entity"

var (
	CombinePair      = 1
	CombineThree     = 2
	CombineFour      = 3
	CombineStraight  = 4
	CombineFullHouse = 5
	CombineFlush     = 6
	CombineFullColor = 7
)

type ChinesePokerBinList interface {
	GetChain(b *entity.BinListCard, comb int) (uint, entity.ListCard)
}
