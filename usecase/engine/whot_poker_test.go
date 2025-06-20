package engine

import (
	"testing"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/entity"
	"github.com/stretchr/testify/assert"
)

// Test cho các hàm khởi tạo và xử lý bài ban đầu
func TestNewGameAndDeal(t *testing.T) {
	// Khởi tạo game engine
	engine := NewWhotEngine()

	// Mock presences
	presences := linkedhashmap.New()
	presences.Put("user1", nil)
	presences.Put("user2", nil)
	// Mock state
	state := &entity.MatchState{
		Presences:        presences,
		PlayingPresences: presences,
		// Label: &entity.MatchLabel{
		// 	Open:         1,
		// 	Bet:          100,
		// 	Code:         "test",
		// 	Name:         "test game",
		// 	Password:     "",
		// 	MaxSize:      4,
		// 	MockCodeCard: 0, // Không dùng mock card
		// },
		DealerId: "user1", // Giả sử user1 là người chia bài
	}
	// Khởi tạo game mới
	err := engine.NewGame(state)
	assert.NoError(t, err)
	assert.NotNil(t, state.Cards)

	// Chia bài
	err = engine.Deal(state)
	assert.NoError(t, err)

	// Kiểm tra chia bài cho cả 2 user
	assert.Equal(t, 2, len(state.Cards))
	assert.NotNil(t, state.Cards["user1"])
	assert.NotNil(t, state.Cards["user2"])

	// Kiểm tra số lá bài của người chơi (theo số người chơi)
	expectedCardCount := entity.MaxPresenceCard // Với 2 người chơi
	assert.Equal(t, expectedCardCount, len(state.Cards["user1"].Cards))
	assert.Equal(t, expectedCardCount, len(state.Cards["user2"].Cards))

	// Kiểm tra lá bài trên bàn
	assert.NotNil(t, state.TopCard)
}

// Test cho việc đánh bài thông thường
func TestPlayCard_Normal(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CIRCLE},
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập bài trên bàn và lượt chơi
	state.TopCard = &pb.Card{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_TRIANGLE}
	state.CurrentTurn = "user1"

	// Đánh lá cùng số (5)
	cardToPlay := &pb.Card{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CIRCLE}
	effect, err := engine.PlayCard(state, "user1", cardToPlay)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, entity.EffectNone, effect)
	assert.Equal(t, cardToPlay, state.TopCard)
	assert.Equal(t, 1, len(state.Cards["user1"].Cards))
}

// Test cho việc đánh lá bài đặc biệt Hold On (1)
func TestPlayCard_HoldOn(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_1, Suit: pb.CardSuit_SUIT_CIRCLE}, // Hold On
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập bài trên bàn và lượt chơi
	state.TopCard = &pb.Card{Rank: pb.CardRank_RANK_1, Suit: pb.CardSuit_SUIT_TRIANGLE}
	state.CurrentTurn = "user1"

	// Đánh lá Hold On
	cardToPlay := &pb.Card{Rank: pb.CardRank_RANK_1, Suit: pb.CardSuit_SUIT_CIRCLE}
	effect, err := engine.PlayCard(state, "user1", cardToPlay)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, entity.EffectHoldOn, effect)
	assert.Equal(t, cardToPlay, state.TopCard)
	assert.Equal(t, 1, len(state.Cards["user1"].Cards))
}

// Test cho việc đánh lá bài đặc biệt Pick Two (2)
func TestPlayCard_PickTwo(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Thêm user2 vào presences
	state.PlayingPresences.Put("user2", nil)

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_CIRCLE}, // Pick Two
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập bài trên bàn và lượt chơi
	state.TopCard = &pb.Card{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_TRIANGLE}
	state.CurrentTurn = "user1"

	// Đánh lá Pick Two
	cardToPlay := &pb.Card{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_CIRCLE}
	effect, err := engine.PlayCard(state, "user1", cardToPlay)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, entity.EffectPickTwo, effect)
	assert.Equal(t, cardToPlay, state.TopCard)
	assert.Equal(t, 2, state.PickPenalty)
	assert.Equal(t, "user2", state.EffectTarget) // Người chơi tiếp theo sẽ bị phạt
}

// Test cho việc đánh lá Whot (20)
func TestPlayCard_Whot(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_20, Suit: pb.CardSuit_SUIT_UNSPECIFIED}, // Whot
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập bài trên bàn và lượt chơi
	state.TopCard = &pb.Card{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_TRIANGLE}
	state.CurrentTurn = "user1"

	// Đánh lá Whot
	cardToPlay := &pb.Card{Rank: pb.CardRank_RANK_20, Suit: pb.CardSuit_SUIT_UNSPECIFIED}
	effect, err := engine.PlayCard(state, "user1", cardToPlay)

	// Kiểm tra kết quả
	assert.NoError(t, err)
	assert.Equal(t, entity.EffectWhot, effect)
	assert.Equal(t, cardToPlay, state.TopCard)
	assert.True(t, state.WaitingForWhotShape) // Đang chờ người chơi chọn hình
}

// Test cho việc đánh bài không hợp lệ
func TestPlayCard_Invalid(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_CIRCLE},
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập bài trên bàn và lượt chơi
	state.TopCard = &pb.Card{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_TRIANGLE}
	state.CurrentTurn = "user1"

	// Đánh lá không hợp lệ (không cùng số, không cùng hình)
	cardToPlay := &pb.Card{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_CIRCLE}
	_, err := engine.PlayCard(state, "user1", cardToPlay)

	// Kiểm tra kết quả
	assert.Error(t, err) // Phải trả về lỗi
}

// Test cho việc rút bài
func TestDrawCardsFromDeck(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Khởi tạo và chia bài
	engine.NewGame(state)
	engine.Deal(state)

	// Lưu số lá bài ban đầu
	// initialCardCount := len(state.Cards["user1"].Cards)

	// Rút thêm 2 lá
	// err := engine.DrawCardsFromDeck(state, "user1", 2)

	// // Kiểm tra kết quả
	// assert.NoError(t, err)
	// assert.Equal(t, initialCardCount+2, len(state.Cards["user1"].Cards))
}

// Test cho việc kết thúc game khi có người thắng trực tiếp
func TestFinish_DirectWinner(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Thêm user2 vào presences
	state.PlayingPresences.Put("user2", nil)

	// Cài đặt bài trên tay người chơi
	state.Cards["user1"] = &pb.ListCard{Cards: []*pb.Card{}} // user1 hết bài
	state.Cards["user2"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CIRCLE},
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},
		},
	}

	// Thiết lập người thắng
	state.WinnerId = "user1"

	// Tính điểm kết thúc
	result := engine.Finish(state)

	// Kiểm tra kết quả
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Results)) // Có kết quả cho cả 2 người chơi

	// Kiểm tra người thắng
	var winnerFound bool
	for _, playerResult := range result.Results {
		if playerResult.UserId == "user1" {
			assert.Equal(t, float64(1), playerResult.WinFactor)
			winnerFound = true
		} else {
			assert.Equal(t, float64(0), playerResult.WinFactor)
		}
	}
	assert.True(t, winnerFound)
}

// Test cho việc kết thúc game khi cần tính điểm (không có người thắng trực tiếp)
func TestFinish_ScoreCalculation(t *testing.T) {
	engine := NewWhotEngine()
	state := createTestMatchState()

	// Thêm user2 vào presences
	state.PlayingPresences.Put("user2", nil)

	// Cài đặt bài trên tay người chơi với điểm khác nhau
	// user1: 5 + (3*2) = 11 điểm (lá Star tính gấp đôi)
	state.Cards["user1"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CIRCLE}, // 5 điểm
			{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_STAR},   // 3*2 = 6 điểm
		},
	}

	// user2: 2 + 10 = 12 điểm
	state.Cards["user2"] = &pb.ListCard{
		Cards: []*pb.Card{
			{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_CIRCLE},    // 2 điểm
			{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_TRIANGLE}, // 10 điểm
		},
	}

	// Không thiết lập người thắng trực tiếp
	state.WinnerId = ""

	// Tính điểm kết thúc
	result := engine.Finish(state)

	// Kiểm tra kết quả
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Results)) // Có kết quả cho cả 2 người chơi

	// Người có điểm thấp nhất (user1) thắng
	var winnerFound bool
	for _, playerResult := range result.Results {
		if playerResult.UserId == "user1" {
			assert.Equal(t, float64(1), playerResult.WinFactor)
			assert.Equal(t, int64(11), playerResult.TotalPoints)
			assert.True(t, playerResult.IsWinner)
			winnerFound = true
		} else {
			assert.Equal(t, float64(0), playerResult.WinFactor)
			assert.Equal(t, int64(12), playerResult.TotalPoints)
			assert.False(t, playerResult.IsWinner)
		}
	}
	assert.True(t, winnerFound)
}

// Hàm tiện ích để tạo MatchState cho các test
func createTestMatchState() *entity.MatchState {
	presences := linkedhashmap.New()
	presences.Put("user1", nil)

	playingPresences := linkedhashmap.New()
	playingPresences.Put("user1", nil)

	return &entity.MatchState{
		Presences:        presences,
		PlayingPresences: playingPresences,
		Cards:            make(map[string]*pb.ListCard),
	}
}
