package entity

import (
	"fmt"
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"google.golang.org/protobuf/proto"
)

var deck *Deck

func GetDeck() *Deck {
	if deck == nil {
		deck = NewDeck()
	}
	return deck
}

func TestShuffle(t *testing.T) {
	deck := GetDeck()
	originalCards := make([]*pb.Card, len(deck.Cards.Cards))
	copy(originalCards, deck.Cards.Cards)

	// Shuffle deck
	deck.Shuffle()

	// 1. Kiểm tra số lượng bài không đổi
	if len(originalCards) != len(deck.Cards.Cards) {
		t.Fatalf("❌ Số lượng lá bài thay đổi sau khi xáo")
	}

	// 2. Kiểm tra tất cả lá bài vẫn còn (không bị thiếu, thừa, hoặc sai)
	cardCount := make(map[string]int)
	for _, card := range originalCards {
		key := fmt.Sprintf("%v-%v", card.Rank, card.Suit)
		cardCount[key]++
	}
	for _, card := range deck.Cards.Cards {
		key := fmt.Sprintf("%v-%v", card.Rank, card.Suit)
		cardCount[key]--
		if cardCount[key] < 0 {
			t.Fatalf("❌ Lá bài dư hoặc không tồn tại: %v", key)
		}
	}

	// 3. Kiểm tra xem thứ tự đã thay đổi chưa
	sameOrder := true
	for i := range originalCards {
		if !proto.Equal(originalCards[i], deck.Cards.Cards[i]) {
			sameOrder = false
			break
		}
	}
	if sameOrder {
		t.Logf("⚠️ Cảnh báo: Thứ tự bài không thay đổi sau khi xáo (có thể hiếm khi xảy ra)")
	} else {
		t.Logf("✅ Thứ tự bài đã được thay đổi sau khi xáo")
	}
}

func TestDeal(t *testing.T) {
	deck := GetDeck()
	deck.Shuffle()

	players := 4
	cardsPerPlayer := 4

	for i := 0; i < players; i++ {
		cards, err := deck.Deal(cardsPerPlayer)
		if err != nil {
			t.Fatalf("❌ deal for player %d error: %v", i+1, err)
		}
		if len(cards.GetCards()) != cardsPerPlayer {
			t.Errorf("❌ player %d nhận không đủ bài, expected %d, got %d", i+1, cardsPerPlayer, len(cards.GetCards()))
		}
		t.Logf("✅ Player %d nhận: %v", i+1, cards.GetCards())
	}

	expectedDealt := players * cardsPerPlayer
	if deck.Dealt != expectedDealt {
		t.Errorf("❌ Sai số lượng bài đã chia. Expected %d, got %d", expectedDealt, deck.Dealt)
	} else {
		t.Logf("✅ Tổng số bài đã chia đúng: %d", deck.Dealt)
	}
}
