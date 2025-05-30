package entity

import (
	"strconv"

	"github.com/heroiclabs/nakama-common/runtime"
	pb "github.com/nakamaFramework/cgp-common/proto/whot"
)

type ArrPbPlayer []*pb.Player

func NewPlayer(presence runtime.Presence) *pb.Player {
	p := pb.Player{
		Id:       presence.GetUserId(),
		UserName: presence.GetUsername(),
	}
	m, ok := presence.(MyPrecense)
	if ok {
		p.AvatarId = m.AvatarId
		p.VipLevel = m.VipLevel
		p.Wallet = strconv.FormatInt(m.Chips, 10)
	}
	return &p
}

func NewListPlayer(presences []runtime.Presence) ArrPbPlayer {
	listPlayer := make([]*pb.Player, 0, len(presences))
	for _, presence := range presences {
		p := NewPlayer(presence)
		listPlayer = append(listPlayer, p)
	}
	return listPlayer
}

// func (arr ArrPbPlayer) ReadWallet(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
// 	listUserId := make([]string, 0, len(arr))
// 	for _, player := range arr {
// 		listUserId = append(listUserId, player.Id)
// 	}
// 	wallets, err := ReadWalletUsers(ctx, nk, logger, listUserId...)
// 	if err != nil {
// 		return err
// 	}
// 	mapWallet := make(map[string]Wallet)
// 	for _, w := range wallets {
// 		mapWallet[w.UserId] = w
// 	}
// 	for i, player := range arr {
// 		player.Wallet = strconv.FormatInt(mapWallet[player.Id].Chips, 10)
// 		arr[i] = player
// 	}
// 	return nil
// }

// func (arr ArrPbPlayer) ReadProfile(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
// 	listUserId := make([]string, 0, len(arr))
// 	for _, player := range arr {
// 		listUserId = append(listUserId, player.Id)
// 	}
// 	profiles, err := GetProfileUser(ctx, nk, listUserId...)
// 	if err != nil {
// 		return err
// 	}
// 	mapWallet := make(map[string]*pb.SimpleProfile, 0)
// 	for _, w := range profiles {
// 		mapWallet[w.UserId] = w
// 	}
// 	for i, player := range arr {
// 		profile := mapWallet[player.Id]
// 		if profile == nil {
// 			continue
// 		}
// 		player.Wallet = strconv.FormatInt(profile.GetAccountChip(), 10)
// 		player.AvatarId = profile.GetAvatarId()
// 		player.VipLevel = profile.GetVipLevel()
// 		arr[i] = player
// 	}
// 	return nil
// }
