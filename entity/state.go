package entity

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/cgp-common/bot"
	pb1 "github.com/nk-nigeria/cgp-common/proto"
	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/pkg/log"
)

const (
	MinPresences = 2
	MaxPresences = 4
)

var BotLoader = bot.NewBotLoader(nil, "", 0)

// type MatchLabel struct {
// 	Open         int32  `json:"open"`
// 	Bet          int32  `json:"bet"`
// 	Code         string `json:"code"`
// 	Name         string `json:"name"`
// 	Password     string `json:"password"`
// 	MaxSize      int32  `json:"max_size"`
// 	MockCodeCard int32  `json:"mock_code_card"`
// }

type MatchState struct {
	Random       *rand.Rand
	Label        *pb1.Match
	MinPresences int

	// Currently connected users, or reserved spaces.
	Presences        *linkedhashmap.Map
	PlayingPresences *linkedhashmap.Map
	LeavePresences   *linkedhashmap.Map
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress     int
	PresencesNoInteract map[string]bool // Map of userId to a boolean indicating if the user is not interacting

	CurrentTurn      string   // userId của người đang đánh
	DealerId         string   // id người chia bài, người này sẽ đánh đầu tiên
	PlayOrder        []string // Danh sách người chơi theo thứ tự đánh bài, bắt đầu từ người chia bài
	PreviousWinnerId string   // id người thắng ván trước, nếu có
	WinnerId         string   // id người thắng ván này, nếu có

	PickPenalty         int        // Số lá phải rút khi bị phạt
	CurrentEffect       CardEffect // Hiệu ứng đang áp dụng
	EffectTarget        string     // ID người chơi bị ảnh hưởng
	WaitingForWhotShape bool       // Đang chờ người chơi chọn hình sau Whot

	//bot
	Bots []*bot.BotPresence
	// BotLogic *BotLogic

	// The top card on the table.
	TopCard *pb.Card
	// Mark assignments to player user IDs.
	Cards map[string]*pb.ListCard
	// Delay for the first turn.
	TurnReadyAt float64

	//time turn play
	TimeTurn     int
	TurnExpireAt int64

	IsEndingGame       bool
	CountDownReachTime time.Time
	LastCountDown      int
	GameState          pb.GameState
	// save balance result in state reward
	// using for send noti to presence join in state reward
	balanceResult   *pb.BalanceResult
	jackpotTreasure *pb.Jackpot
}

func NewMatchState(label *pb1.Match) MatchState {
	m := MatchState{
		Random:              rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:               label,
		MinPresences:        MinPresences,
		Presences:           linkedhashmap.New(),
		PlayingPresences:    linkedhashmap.New(),
		LeavePresences:      linkedhashmap.New(),
		PresencesNoInteract: make(map[string]bool),
		Cards:               make(map[string]*pb.ListCard),
		TimeTurn:            10,
		// balanceResult:       nil,
	}

	if bots, err := BotLoader.GetFreeBot(int(label.NumBot)); err != nil {
		fmt.Printf("\r\n load bot failed %s  \r\n", err.Error())
	} else {
		m.Bots = bots
	}

	for _, bot := range m.Bots {
		m.Presences.Put(bot.GetUserId(), bot)
		m.Label.Size += 1
		fmt.Printf("\r\n add bot %s to match \r\n", bot.GetUserId())
		// m.addOrderPresence(bot)
	}
	return m
}

func (s *MatchState) SetDealer() {
	// Nếu đã có người thắng ván trước, set làm dealer
	if s.PreviousWinnerId != "" {
		if val, found := s.PlayingPresences.Get(s.PreviousWinnerId); found && val != nil {
			s.DealerId = s.PreviousWinnerId
			return
		}
	}
	// Nếu là bàn mới (chưa có người thắng) hoặc người thắng ván trước không còn trong bàn,
	// chọn ngẫu nhiên một người chơi làm dealer
	if s.PlayingPresences.Size() == 0 {
		log.GetLogger().Warn("No players in the match to set as dealer")
		s.DealerId = ""
		return
	}
	keys := s.PlayingPresences.Keys()
	if len(keys) == 0 {
		s.DealerId = ""
		return
	}
	randomIndex := s.Random.Intn(len(keys))
	s.DealerId = keys[randomIndex].(string)
	log.GetLogger().Info("Set dealer to %s", s.DealerId)
}

func (s *MatchState) ResetMatch() {
	s.PlayingPresences = linkedhashmap.New()
	s.Cards = make(map[string]*pb.ListCard)
	s.TopCard = nil
	s.CurrentTurn = ""
	s.DealerId = ""
	s.TurnExpireAt = 0
	s.TurnReadyAt = 0
	s.PreviousWinnerId = s.WinnerId
	s.ResetBalanceResult()
	s.WinnerId = ""
	s.EffectTarget = ""
	s.PickPenalty = 0
	s.CurrentEffect = EffectNone
	s.WaitingForWhotShape = false
	s.IsEndingGame = false
}

func (s *MatchState) GetPresenceNotBotSize() int {
	count := 0
	s.Presences.Each(func(index any, value interface{}) {
		presence, ok := value.(runtime.Presence)
		if !ok {
			return
		}
		if BotLoader.IsBot(presence.GetUserId()) {
			return
		}
		count++
	})
	return count
}

func (s *MatchState) BuildPlayOrderFromDealer() {
	// Xây dựng danh sách người chơi theo thứ tự đánh bài, bắt đầu từ dealer
	keys := s.PlayingPresences.Keys()
	values := make([]string, 0, len(keys))
	startIndex := 0

	for i, key := range keys {
		val, _ := s.PlayingPresences.Get(key)
		presence := val.(runtime.Presence)
		userID := presence.GetUserId()
		values = append(values, userID)

		if userID == s.DealerId {
			startIndex = i
		}
	}

	rotated := append(values[startIndex:], values[:startIndex]...)
	s.PlayOrder = rotated
}

func (s *MatchState) GetNextPlayerClockwise(current string) string {
	for i, userID := range s.PlayOrder {
		if userID == current {
			return s.PlayOrder[(i+1)%len(s.PlayOrder)]
		}
	}
	return current
}

func (s *MatchState) GetBalanceResult() *pb.BalanceResult {
	return s.balanceResult
}

func (s *MatchState) SetBalanceResult(u *pb.BalanceResult) {
	s.balanceResult = u
}

func (s *MatchState) ResetBalanceResult() {
	s.SetBalanceResult(nil)
}

func (s *MatchState) SetJackpotTreasure(jp *pb.Jackpot) {
	s.jackpotTreasure = &pb.Jackpot{
		GameCode: jp.GetGameCode(),
		Chips:    jp.GetChips(),
	}
}

func (s *MatchState) GetJackpotTreasure() *pb.Jackpot {
	return s.jackpotTreasure
}

func (s *MatchState) AddPresence(ctx context.Context, nk runtime.NakamaModule, db *sql.DB, presences []runtime.Presence) {
	for _, presence := range presences {
		m := NewMyPrecense(ctx, nk, db, presence)
		log.GetLogger().Info("Add presence %s to match %s", m.DeviceID)
		s.Presences.Put(presence.GetUserId(), m)
		s.ResetUserNotInteract(presence.GetUserId())
	}
}

func (s *MatchState) RemovePresence(presences ...runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Remove(presence.GetUserId())
		if _, ok := s.PlayingPresences.Get(presence.GetUserId()); ok {
			s.PlayingPresences.Remove(presence.GetUserId())
		}
		delete(s.PresencesNoInteract, presence.GetUserId())
	}
}

func (s *MatchState) AddLeavePresence(presences ...runtime.Presence) {
	for _, presence := range presences {
		s.LeavePresences.Put(presence.GetUserId(), presence)
	}
}

func (s *MatchState) RemoveLeavePresence(userId string) {
	s.LeavePresences.Remove(userId)
}

func (s *MatchState) ApplyLeavePresence() {
	s.LeavePresences.Each(func(key interface{}, value interface{}) {
		s.Presences.Remove(key)
		s.PlayingPresences.Remove(key)
		delete(s.PresencesNoInteract, key.(string))
	})

	s.LeavePresences = linkedhashmap.New()
}

func (s *MatchState) SetupMatchPresence() {
	s.PlayingPresences = linkedhashmap.New()
	presences := make([]runtime.Presence, 0, s.Presences.Size())
	s.Presences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})
	s.AddPlayingPresences(presences...)
}

func (s *MatchState) AddPlayingPresences(presences ...runtime.Presence) {
	for _, p := range presences {
		s.PlayingPresences.Put(p.GetUserId(), p)
		keyStr := p.GetUserId()
		s.PresencesNoInteract[keyStr] = false
	}
}

func (s *MatchState) GetPresenceNotInteract() []runtime.Presence {
	listPresence := make([]runtime.Presence, 0)
	s.Presences.Each(func(key interface{}, value interface{}) {
		if isAutoPlay, exist := s.PresencesNoInteract[key.(string)]; exist && isAutoPlay {
			listPresence = append(listPresence, value.(runtime.Presence))
		}
	})
	return listPresence
}

func (s *MatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
	s.LastCountDown = -1
}

func (s *MatchState) GetRemainCountDown() int {
	currentTime := time.Now()
	difference := s.CountDownReachTime.Sub(currentTime)
	return int(difference.Seconds())
}

func (s *MatchState) SetLastCountDown(countDown int) {
	s.LastCountDown = countDown
}

func (s *MatchState) GetLastCountDown() int {
	return s.LastCountDown
}

func (s *MatchState) IsNeedNotifyCountDown() bool {
	return s.GetRemainCountDown() != s.LastCountDown || s.LastCountDown == -1
}

func (s *MatchState) GetPresenceSize() int {
	return s.Presences.Size()
}

func (s *MatchState) GetPlayingPresenceSize() int {
	return s.PlayingPresences.Size()
}

func (s *MatchState) IsReadyToPlay() bool {
	return s.Presences.Size() >= s.MinPresences
}

func (s *MatchState) GetPlayingCount() int {
	return s.PlayingPresences.Size()
}

func (s *MatchState) GetPresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.Presences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) GetPresence(userID string) runtime.Presence {
	_, value := s.Presences.Find(func(key, value interface{}) bool {
		if key == userID {
			return true
		}
		return false
	})
	if value == nil {
		return nil
	}
	return value.(runtime.Presence)
}

func (s *MatchState) GetPlayingPresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.PlayingPresences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) GetLeavePresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.LeavePresences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) ResetUserNotInteract(userId string) {
	s.PresencesNoInteract[userId] = false
}

func (s *MatchState) SetUserNotInteract(userId string) {
	log.GetLogger().Info("Set user %s not interact in match", userId)
	s.PresencesNoInteract[userId] = true
}

func (s *MatchState) AddUserNotInteractToLeaves() {
	listPrecense := s.GetPresenceNotInteract()
	if len(listPrecense) > 0 {
		listUserId := make([]string, len(listPrecense))
		for _, p := range listPrecense {
			listUserId = append(listUserId, p.GetUserId())
		}
		log.GetLogger().Info("Kick %d user from match %s",
			len(listPrecense), strings.Join(listUserId, ","))
		s.AddLeavePresence(listPrecense...)
	}
}
