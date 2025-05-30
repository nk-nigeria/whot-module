package api

import (
	"context"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	pb "github.com/nakamaFramework/cgp-common/proto"
)

func (m *MatchHandler) addChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": int64(amountChip), // Add amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
	}

}

func (m *MatchHandler) subtractChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": -int64(amountChip), // Substract amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.
			WithField("err", err).
			WithField("userID", userID).
			Error("Wallet update error.")
	}
}

func (m *MatchHandler) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, resultGame *pb.UpdateFinish) {
	logger.Info("update chip result %v, size %v", resultGame, len(resultGame.Results))
	walletUpdates := make([]*runtime.WalletUpdate, 0, len(resultGame.Results))
	for _, result := range resultGame.Results {
		amountChip := int64(0)

		amountChip = 200*(result.ScoreResult.FrontFactor+result.ScoreResult.MiddleFactor+result.ScoreResult.BackFactor) +
			(result.ScoreResult.FrontBonusFactor + result.ScoreResult.MiddleBonusFactor + result.ScoreResult.BackBonusFactor)

		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_topup": "topup",
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}

	logger.Info("update wallet ctx %v, walletUpdates %v", ctx, walletUpdates)
	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		payload, _ := json.Marshal(walletUpdates)
		logger.
			WithField("err", err).WithField("payload", string(payload)).Error("Wallets update error.")
	}
}
