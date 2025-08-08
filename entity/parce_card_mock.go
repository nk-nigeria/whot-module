package entity

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	pb "github.com/nk-nigeria/cgp-common/proto"
)

func ParseCard(str string) *pb.WhotCard {
	l := len(str)
	if l < 2 {
		return nil
	}
	rankStr := str[:l-1]
	suitStr := strings.ToLower(str[l-1:])

	rankInt, err := strconv.Atoi(rankStr)
	if err != nil || (rankInt < 1 || (rankInt > 14 && rankInt != 20)) {
		return nil
	}

	card := &pb.WhotCard{
		Rank: pb.WhotCardRank(rankInt),
	}

	switch suitStr {
	case "c":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_CIRCLE
	case "x":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_CROSS
	case "s":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_STAR
	case "t":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_TRIANGLE
	case "q":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_SQUARE
	case "w":
		card.Suit = pb.WhotCardSuit_WHOT_SUIT_UNSPECIFIED
	default:
		return nil
	}

	return card
}

func ParseListCard(str string) []*pb.WhotCard {
	ml := make([]*pb.WhotCard, 0)
	arrCardMock := strings.Split(str, ";")
	for _, s := range arrCardMock {
		s = strings.TrimSpace(s)
		ml = append(ml, ParseCard(s))
	}
	return ml
}

func ParseMockCard(fileMock string) [][]*pb.WhotCard {
	f, err := os.Open(fileMock)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	list := make([][]*pb.WhotCard, 0)
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
