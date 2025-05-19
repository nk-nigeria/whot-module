package mock

import (
	"context"
	"os"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MockModule struct {
}

func (m *MockModule) GetSatori() runtime.Satori {
	panic("implement me")
}

func (m MockModule) FriendsBlock(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("implement me")
}

func (m MockModule) ChannelMessageRemove(ctx context.Context, channelId, messageId, userId, username string, authoritative bool) (*rtapi.ChannelMessageAck, error) {
	panic("implement me")
}

func (m MockModule) ChannelMessagesList(ctx context.Context, channelId string, limit int, forward bool, cursor string) ([]*api.ChannelMessage, string, string, error) {
	panic("implement me")
}

func (m MockModule) UsersGetRandom(ctx context.Context, count int) ([]*api.User, error) {
	panic("implement me")
}

func (m MockModule) MatchSignal(ctx context.Context, id string, data string) (string, error) {
	panic("implement me")
}

func (m MockModule) NotificationSendAll(ctx context.Context, subject string, content map[string]interface{}, code int, persistent bool) error {
	panic("implement me")
}

func (m MockModule) LeaderboardRecordsHaystack(ctx context.Context, id, ownerID string, limit int, cursor string, expiry int64) (*api.LeaderboardRecordList, error) {
	panic("implement me")
}

func (m MockModule) LeaderboardRecordsListCursorFromRank(id string, rank, overrideExpiry int64) (string, error) {
	panic("implement me")
}

func (m MockModule) PurchaseValidateApple(ctx context.Context, userID, receipt string, persist bool, passwordOverride ...string) (*api.ValidatePurchaseResponse, error) {
	panic("implement me")
}

func (m MockModule) PurchaseValidateGoogle(ctx context.Context, userID, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidatePurchaseResponse, error) {
	panic("implement me")
}

func (m MockModule) PurchaseValidateHuawei(ctx context.Context, userID, signature, inAppPurchaseData string, persist bool) (*api.ValidatePurchaseResponse, error) {
	panic("implement me")
}

func (m MockModule) TournamentCreate(ctx context.Context, id string, authoritative bool, sortOrder, operator, resetSchedule string, metadata map[string]interface{}, title, description string, category, startTime, endTime, duration, maxSize, maxNumScore int, joinRequired bool) error {
	panic("implement me")
}

func (m MockModule) GroupUsersBan(ctx context.Context, callerID, groupID string, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) FriendsAdd(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("implement me")
}

func (m MockModule) FriendsDelete(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("implement me")
}

func (m MockModule) ChannelIdBuild(ctx context.Context, sender string, target string, chanType runtime.ChannelType) (string, error) {
	panic("implement me")
}

func (m MockModule) ChannelMessageSend(ctx context.Context, channelID string, content map[string]interface{}, senderId, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("implement me")
}

func (m MockModule) ChannelMessageUpdate(ctx context.Context, channelID, messageID string, content map[string]interface{}, senderId, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateApple(ctx context.Context, token, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateCustom(ctx context.Context, id, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateDevice(ctx context.Context, id, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateEmail(ctx context.Context, email, password, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateFacebook(ctx context.Context, token string, importFriends bool, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateFacebookInstantGame(ctx context.Context, signedPlayerInfo string, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateGameCenter(ctx context.Context, playerID, bundleID string, timestamp int64, salt, signature, publicKeyUrl, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateGoogle(ctx context.Context, token, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateSteam(ctx context.Context, token, username string, create bool) (string, string, bool, error) {
	panic("implement me")
}

func (m MockModule) AuthenticateTokenGenerate(userID, username string, exp int64, vars map[string]string) (string, int64, error) {
	panic("implement me")
}

func (m MockModule) AccountGetId(ctx context.Context, userID string) (*api.Account, error) {
	panic("implement me")
}

func (m MockModule) AccountsGetId(ctx context.Context, userIDs []string) ([]*api.Account, error) {
	return []*api.Account{}, nil
}

func (m MockModule) AccountUpdateId(ctx context.Context, userID, username string, metadata map[string]interface{}, displayName, timezone, location, langTag, avatarUrl string) error {
	panic("implement me")
}

func (m MockModule) AccountDeleteId(ctx context.Context, userID string, recorded bool) error {
	panic("implement me")
}

func (m MockModule) AccountExportId(ctx context.Context, userID string) (string, error) {
	panic("implement me")
}

func (m MockModule) UsersGetId(ctx context.Context, userIDs []string, facebookIDs []string) ([]*api.User, error) {
	panic("implement me")
}

func (m MockModule) UsersGetUsername(ctx context.Context, usernames []string) ([]*api.User, error) {
	panic("implement me")
}

func (m MockModule) UsersBanId(ctx context.Context, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) UsersUnbanId(ctx context.Context, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) LinkApple(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) LinkCustom(ctx context.Context, userID, customID string) error {
	panic("implement me")
}

func (m MockModule) LinkDevice(ctx context.Context, userID, deviceID string) error {
	panic("implement me")
}

func (m MockModule) LinkEmail(ctx context.Context, userID, email, password string) error {
	panic("implement me")
}

func (m MockModule) LinkFacebook(ctx context.Context, userID, username, token string, importFriends bool) error {
	panic("implement me")
}

func (m MockModule) LinkFacebookInstantGame(ctx context.Context, userID, signedPlayerInfo string) error {
	panic("implement me")
}

func (m MockModule) LinkGameCenter(ctx context.Context, userID, playerID, bundleID string, timestamp int64, salt, signature, publicKeyUrl string) error {
	panic("implement me")
}

func (m MockModule) LinkGoogle(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) LinkSteam(ctx context.Context, userID, username, token string, importFriends bool) error {
	panic("implement me")
}

func (m MockModule) ReadFile(path string) (*os.File, error) {
	panic("implement me")
}

func (m MockModule) UnlinkApple(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) UnlinkCustom(ctx context.Context, userID, customID string) error {
	panic("implement me")
}

func (m MockModule) UnlinkDevice(ctx context.Context, userID, deviceID string) error {
	panic("implement me")
}

func (m MockModule) UnlinkEmail(ctx context.Context, userID, email string) error {
	panic("implement me")
}

func (m MockModule) UnlinkFacebook(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) UnlinkFacebookInstantGame(ctx context.Context, userID, signedPlayerInfo string) error {
	panic("implement me")
}

func (m MockModule) UnlinkGameCenter(ctx context.Context, userID, playerID, bundleID string, timestamp int64, salt, signature, publicKeyUrl string) error {
	panic("implement me")
}

func (m MockModule) UnlinkGoogle(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) UnlinkSteam(ctx context.Context, userID, token string) error {
	panic("implement me")
}

func (m MockModule) StreamUserList(mode uint8, subject, subcontext, label string, includeHidden, includeNotHidden bool) ([]runtime.Presence, error) {
	panic("implement me")
}

func (m MockModule) StreamUserGet(mode uint8, subject, subcontext, label, userID, sessionID string) (runtime.PresenceMeta, error) {
	panic("implement me")
}

func (m MockModule) StreamUserJoin(mode uint8, subject, subcontext, label, userID, sessionID string, hidden, persistence bool, status string) (bool, error) {
	panic("implement me")
}

func (m MockModule) StreamUserUpdate(mode uint8, subject, subcontext, label, userID, sessionID string, hidden, persistence bool, status string) error {
	panic("implement me")
}

func (m MockModule) StreamUserLeave(mode uint8, subject, subcontext, label, userID, sessionID string) error {
	panic("implement me")
}

func (m MockModule) StreamUserKick(mode uint8, subject, subcontext, label string, presence runtime.Presence) error {
	panic("implement me")
}

func (m MockModule) StreamCount(mode uint8, subject, subcontext, label string) (int, error) {
	panic("implement me")
}

func (m MockModule) StreamClose(mode uint8, subject, subcontext, label string) error {
	panic("implement me")
}

func (m MockModule) StreamSend(mode uint8, subject, subcontext, label, data string, presences []runtime.Presence, reliable bool) error {
	panic("implement me")
}

func (m MockModule) StreamSendRaw(mode uint8, subject, subcontext, label string, msg *rtapi.Envelope, presences []runtime.Presence, reliable bool) error {
	panic("implement me")
}

func (m MockModule) SessionDisconnect(ctx context.Context, sessionID string, reason ...runtime.PresenceReason) error {
	panic("implement me")
}

func (m MockModule) SessionLogout(userID, token, refreshToken string) error {
	panic("implement me")
}

func (m MockModule) MatchCreate(ctx context.Context, module string, params map[string]interface{}) (string, error) {
	panic("implement me")
}

func (m MockModule) MatchGet(ctx context.Context, id string) (*api.Match, error) {
	panic("implement me")
}

func (m MockModule) MatchList(ctx context.Context, limit int, authoritative bool, label string, minSize, maxSize *int, query string) ([]*api.Match, error) {
	panic("implement me")
}

func (m MockModule) NotificationSend(ctx context.Context, userID, subject string, content map[string]interface{}, code int, sender string, persistent bool) error {
	panic("implement me")
}

func (m MockModule) NotificationsSend(ctx context.Context, notifications []*runtime.NotificationSend) error {
	panic("implement me")
}

func (m MockModule) NotificationsDelete(ctx context.Context, notifications []*runtime.NotificationDelete) error {
	panic("implement me")
}

func (m MockModule) WalletUpdate(ctx context.Context, userID string, changeset map[string]int64, metadata map[string]interface{}, updateLedger bool) (map[string]int64, map[string]int64, error) {
	panic("implement me")
}

func (m MockModule) WalletsUpdate(ctx context.Context, updates []*runtime.WalletUpdate, updateLedger bool) ([]*runtime.WalletUpdateResult, error) {
	return []*runtime.WalletUpdateResult{
		{
			UserID: "123",
		},
	}, nil
}

func (m MockModule) WalletLedgerUpdate(ctx context.Context, itemID string, metadata map[string]interface{}) (runtime.WalletLedgerItem, error) {
	panic("implement me")
}

func (m MockModule) WalletLedgerList(ctx context.Context, userID string, limit int, cursor string) ([]runtime.WalletLedgerItem, string, error) {
	panic("implement me")
}

func (m MockModule) SubscriptionGetByProductId(ctx context.Context, userID, productID string) (*api.ValidatedSubscription, error) {
	panic("implement me")
}

func (m MockModule) SubscriptionValidateApple(ctx context.Context, userID, receipt string, persist bool, passwordOverride ...string) (*api.ValidateSubscriptionResponse, error) {
	panic("implement me")
}

func (m MockModule) SubscriptionValidateGoogle(ctx context.Context, userID, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidateSubscriptionResponse, error) {
	panic("implement me")
}

func (m MockModule) SubscriptionsList(ctx context.Context, userID string, limit int, cursor string) (*api.SubscriptionList, error) {
	panic("implement me")
}

func (m MockModule) StorageList(ctx context.Context, callerID, userID, collection string, limit int, cursor string) ([]*api.StorageObject, string, error) {
	panic("implement me")
}

func (m MockModule) StorageRead(ctx context.Context, reads []*runtime.StorageRead) ([]*api.StorageObject, error) {
	panic("implement me")
}

func (m MockModule) StorageWrite(ctx context.Context, writes []*runtime.StorageWrite) ([]*api.StorageObjectAck, error) {
	panic("implement me")
}

func (m MockModule) StorageDelete(ctx context.Context, deletes []*runtime.StorageDelete) error {
	panic("implement me")
}

func (m MockModule) MultiUpdate(ctx context.Context, accountUpdates []*runtime.AccountUpdate, storageWrites []*runtime.StorageWrite, walletUpdates []*runtime.WalletUpdate, updateLedger bool) ([]*api.StorageObjectAck, []*runtime.WalletUpdateResult, error) {
	panic("implement me")
}

func (m MockModule) LeaderboardCreate(ctx context.Context, id string, authoritative bool, sortOrder, operator, resetSchedule string, metadata map[string]interface{}) error {
	panic("implement me")
}

func (m MockModule) LeaderboardDelete(ctx context.Context, id string) error {
	panic("implement me")
}

func (m MockModule) LeaderboardList(limit int, cursor string) (*api.LeaderboardList, error) {
	panic("implement me")
}

func (m MockModule) LeaderboardRecordsList(ctx context.Context, id string, ownerIDs []string, limit int, cursor string, expiry int64) ([]*api.LeaderboardRecord, []*api.LeaderboardRecord, string, string, error) {
	panic("implement me")
}

func (m MockModule) LeaderboardRecordWrite(ctx context.Context, id, ownerID, username string, score, subscore int64, metadata map[string]interface{}, overrideOperator *int) (*api.LeaderboardRecord, error) {
	panic("implement me")
}

func (m MockModule) LeaderboardRecordDelete(ctx context.Context, id, ownerID string) error {
	panic("implement me")
}

func (m MockModule) LeaderboardsGetId(ctx context.Context, ids []string) ([]*api.Leaderboard, error) {
	panic("implement me")
}

func (m MockModule) PurchasesList(ctx context.Context, userID string, limit int, cursor string) (*api.PurchaseList, error) {
	panic("implement me")
}

func (m MockModule) PurchaseGetByTransactionId(ctx context.Context, transactionID string) (*api.ValidatedPurchase, error) {
	panic("implement me")
}

func (m MockModule) PurchaseValidateFacebookInstant(ctx context.Context, userID, signedRequest string, persist bool) (*api.ValidatePurchaseResponse, error) {
	panic("implement me")
}

func (m *MockModule) StorageIndexList(ctx context.Context, collection, indexName, query string, limit int) (*api.StorageObjects, error) {
	panic("implement me")
}

func (m MockModule) TournamentDelete(ctx context.Context, id string) error {
	panic("implement me")
}

func (m MockModule) TournamentRecordDelete(ctx context.Context, id, ownerID string) error {
	panic("implement me")
}

func (m MockModule) TournamentAddAttempt(ctx context.Context, id, ownerID string, count int) error {
	panic("implement me")
}

func (m MockModule) TournamentJoin(ctx context.Context, id, ownerID, username string) error {
	panic("implement me")
}

func (m MockModule) TournamentsGetId(ctx context.Context, tournamentIDs []string) ([]*api.Tournament, error) {
	panic("implement me")
}

func (m MockModule) TournamentList(ctx context.Context, categoryStart, categoryEnd, startTime, endTime, limit int, cursor string) (*api.TournamentList, error) {
	panic("implement me")
}

func (m MockModule) TournamentRecordsList(ctx context.Context, tournamentId string, ownerIDs []string, limit int, cursor string, overrideExpiry int64) ([]*api.LeaderboardRecord, []*api.LeaderboardRecord, string, string, error) {
	panic("implement me")
}

func (m MockModule) TournamentRecordWrite(ctx context.Context, id, ownerID, username string, score, subscore int64, metadata map[string]interface{}, operatorOverride *int) (*api.LeaderboardRecord, error) {
	panic("implement me")
}

func (m MockModule) TournamentRecordsHaystack(ctx context.Context, id, ownerID string, limit int, cursor string, expiry int64) (*api.TournamentRecordList, error) {
	panic("implement me")
}

func (m MockModule) GroupsGetId(ctx context.Context, groupIDs []string) ([]*api.Group, error) {
	panic("implement me")
}

func (m MockModule) GroupCreate(ctx context.Context, userID, name, creatorID, langTag, description, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) (*api.Group, error) {
	panic("implement me")
}

func (m MockModule) GroupUpdate(ctx context.Context, id, userID, name, creatorID, langTag, description, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) error {
	panic("implement me")
}

func (m MockModule) GroupsGetRandom(ctx context.Context, count int) ([]*api.Group, error) {
	panic("implement me")
}

func (m MockModule) GroupDelete(ctx context.Context, id string) error {
	panic("implement me")
}

func (m MockModule) GroupUserJoin(ctx context.Context, groupID, userID, username string) error {
	panic("implement me")
}

func (m MockModule) GroupUserLeave(ctx context.Context, groupID, userID, username string) error {
	panic("implement me")
}

func (m MockModule) GroupUsersAdd(ctx context.Context, callerID, groupID string, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) GroupUsersKick(ctx context.Context, callerID, groupID string, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) GroupUsersPromote(ctx context.Context, callerID, groupID string, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) GroupUsersDemote(ctx context.Context, callerID, groupID string, userIDs []string) error {
	panic("implement me")
}

func (m MockModule) GroupUsersList(ctx context.Context, id string, limit int, state *int, cursor string) ([]*api.GroupUserList_GroupUser, string, error) {
	panic("implement me")
}

func (m MockModule) GroupsList(ctx context.Context, name, langTag string, members *int, open *bool, limit int, cursor string) ([]*api.Group, string, error) {
	panic("implement me")
}

func (m MockModule) UserGroupsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.UserGroupList_UserGroup, string, error) {
	panic("implement me")
}

func (m MockModule) FriendsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.Friend, string, error) {
	panic("implement me")
}

func (m MockModule) Event(ctx context.Context, evt *api.Event) error {
	panic("implement me")
}

func (m MockModule) MetricsCounterAdd(name string, tags map[string]string, delta int64) {
	panic("implement me")
}

func (m MockModule) MetricsGaugeSet(name string, tags map[string]string, value float64) {
	panic("implement me")
}

func (m MockModule) MetricsTimerRecord(name string, tags map[string]string, value time.Duration) {
	panic("implement me")
}
