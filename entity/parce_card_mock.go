package entity

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
)

func ParseCard(str string) *pb.Card {
	l := len(str)
	if l < 2 {
		return nil
	}
	rank := strings.ToLower(str[:l-1])
	suit := str[l-1 : l]
	card := &pb.Card{}
	switch rank {
	case "1", "a":
		card.Rank = pb.CardRank_RANK_1
	case "11", "j":
		card.Rank = pb.CardRank_RANK_11
	case "12", "q":
		card.Rank = pb.CardRank_RANK_12
	case "13", "k":
		card.Rank = pb.CardRank_RANK_13
	default:
		rankInt, _ := strconv.Atoi(rank)
		if rankInt < 2 || rankInt > 13 {
			return nil
		}
		card.Rank = pb.CardRank(rankInt)
	}

	switch suit {
	case "c":
		card.Suit = pb.CardSuit_SUIT_CIRCLE
	case "s":
		card.Suit = pb.CardSuit_SUIT_CROSS
	case "d":
		card.Suit = pb.CardSuit_SUIT_SQUARE
	case "h":
		card.Suit = pb.CardSuit_SUIT_STAR
	default:
		suitInt, _ := strconv.Atoi(rank)
		if suitInt <= 0 || suitInt > 4 {
			return nil
		}
		card.Suit = pb.CardSuit(suitInt)
	}
	return card
}
func ParseListCard(str string) []*pb.Card {
	ml := make([]*pb.Card, 0)
	arrCardMock := strings.Split(str, ";")
	for _, s := range arrCardMock {
		s = strings.TrimSpace(s)
		ml = append(ml, ParseCard(s))
	}
	return ml
}

func ParseMockCard(fileMock string) [][]*pb.Card {
	f, err := os.Open(fileMock)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	list := make([][]*pb.Card, 0)
	for scanner.Scan() {
		lineText := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(lineText, "#") || len(lineText) == 0 {
			continue
		}
		list = append(list, ParseListCard(lineText))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return list
}
