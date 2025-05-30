package entity

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
)

type MyPrecense struct {
	runtime.Presence
	AvatarId string
	Chips    int64
	VipLevel int64
}

// type ListMyPrecense []MyPrecense

func NewMyPrecense(ctx context.Context, nk runtime.NakamaModule, precense runtime.Presence) MyPrecense {
	m := MyPrecense{
		Presence: precense,
	}
	profiles, err := GetProfileUsers(ctx, nk, precense.GetUserId())
	if err != nil {
		return m
	}
	if len(profiles) == 0 {
		return m
	}
	p := profiles[0]
	m.AvatarId = p.AvatarId
	m.Chips = p.AccountChip
	m.VipLevel = p.VipLevel
	return m
}
