package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nk-nigeria/cgp-common/bot"
	"github.com/nk-nigeria/whot-module/entity"
	"github.com/nk-nigeria/whot-module/pkg/packager"
)

// WhotBotIntegration implements BotIntegration for Whot game
type WhotBotIntegration struct {
	db           *sql.DB
	matchID      string
	betAmount    int64
	playerCount  int
	maxPlayers   int
	minPlayers   int
	lastResult   int
	activeTables int
	botHelper    *bot.BotIntegrationHelper
}

// NewWhotBotIntegration creates a new Whot bot integration
func NewWhotBotIntegration(db *sql.DB) *WhotBotIntegration {
	integration := &WhotBotIntegration{
		db:           db,
		maxPlayers:   4, // Whot typically has 4 players
		minPlayers:   2, // Minimum 2 players to start
		activeTables: 0, // Will be updated from game state
	}

	integration.botHelper = bot.NewBotIntegrationHelper(db, integration)
	integration.LoadBotConfig(context.Background())
	return integration
}

// GetGameCode returns the game code for Whot
func (w *WhotBotIntegration) GetGameCode() string {
	return entity.ModuleName
}

// GetMinChipBalance returns minimum chip balance for bots
func (w *WhotBotIntegration) GetMinChipBalance() int64 {
	return 75000 // 75k chips minimum
}

// GetMatchInfo returns current match information
func (w *WhotBotIntegration) GetMatchInfo(ctx context.Context) *bot.MatchInfo {
	return &bot.MatchInfo{
		MatchID:           w.matchID,
		BetAmount:         w.betAmount,
		PlayerCount:       w.playerCount,
		MaxPlayers:        w.maxPlayers,
		MinPlayers:        w.minPlayers,
		IsFull:            w.IsMatchFull(),
		LastGameResult:    w.lastResult,
		ActiveTablesCount: w.activeTables,
	}
}

// AddBotToMatch adds bots to the current match
func (w *WhotBotIntegration) AddBotToMatch(ctx context.Context, numBots int) error {

	// Add bots to match using the existing processor logic
	// This requires access to the processor and state from context
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	if procPkg == nil {
		return fmt.Errorf("processor package not found in context")
	}

	state := procPkg.GetState()

	// Add bots to match using existing processor method
	err := procPkg.GetProcessor().AddBotToMatch(
		ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state,
		numBots,
	)
	if err != nil {
		return err
	}

	// Update player count
	w.playerCount = state.GetPresenceSize()

	return nil
}

// GetMaxPlayers returns maximum players allowed
func (w *WhotBotIntegration) GetMaxPlayers() int {
	return w.maxPlayers
}

// GetMinPlayers returns minimum players required
func (w *WhotBotIntegration) GetMinPlayers() int {
	return w.minPlayers
}

// IsMatchFull returns true if match is full
func (w *WhotBotIntegration) IsMatchFull() bool {
	return w.playerCount >= w.maxPlayers
}

// GetCurrentPlayerCount returns current player count
func (w *WhotBotIntegration) GetCurrentPlayerCount() int {
	return w.playerCount
}

// GetCurrentBetAmount returns current bet amount
func (w *WhotBotIntegration) GetCurrentBetAmount() int64 {
	return w.betAmount
}

// GetLastGameResult returns last game result
func (w *WhotBotIntegration) GetLastGameResult() int {
	return w.lastResult
}

// GetActiveTablesCount returns active tables count
func (w *WhotBotIntegration) GetActiveTablesCount() int {
	return w.activeTables
}

// SetMatchState updates the match state for bot decision making
func (w *WhotBotIntegration) SetMatchState(matchID string, betAmount int64, playerCount int, lastResult int, activeTables int) {
	w.matchID = matchID
	w.betAmount = betAmount
	w.playerCount = playerCount
	w.lastResult = lastResult
	w.activeTables = activeTables
}

// ProcessBotLogic processes all bot-related logic
func (w *WhotBotIntegration) ProcessBotLogic(ctx context.Context) error {
	return w.botHelper.ProcessBotLogic(ctx)
}

// GetBotHelper returns the bot helper for direct access
func (w *WhotBotIntegration) GetBotHelper() *bot.BotIntegrationHelper {
	return w.botHelper
}

// LoadBotConfig loads bot configuration from database
func (w *WhotBotIntegration) LoadBotConfig(ctx context.Context) error {
	configLoader := bot.NewConfigLoader(w.db)
	config, err := configLoader.LoadConfigFromDatabase(ctx, w.GetGameCode())
	if err != nil {
		return fmt.Errorf("failed to load bot config: %w", err)
	}

	w.botHelper.SetBotConfig(config)
	return nil
}

// SaveBotConfig saves bot configuration to database
func (w *WhotBotIntegration) SaveBotConfig(ctx context.Context) error {
	config := w.botHelper.GetBotConfig()
	configLoader := bot.NewConfigLoader(w.db)

	err := configLoader.SaveConfigToDatabase(ctx, w.GetGameCode(), config)
	if err != nil {
		return fmt.Errorf("failed to save bot config: %w", err)
	}

	return nil
}

// DebugPendingRequests prints pending requests for debugging
func (w *WhotBotIntegration) DebugPendingRequests() {
	w.botHelper.DebugPendingRequests()
}
