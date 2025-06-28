package entity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

type MyPrecense struct {
	runtime.Presence
	AvatarId string
	Chips    int64
	VipLevel int64
	DeviceID string
}

// type ListMyPrecense []MyPrecense

func NewMyPrecense(ctx context.Context, nk runtime.NakamaModule, db *sql.DB, precense runtime.Presence) MyPrecense {
	m := MyPrecense{
		Presence: precense,
	}

	profiles, err := GetProfileUsers(ctx, nk, precense.GetUserId())
	if err != nil || len(profiles) == 0 {
		return m
	}

	deviceID, err := GetDeviceIDByUserID(ctx, db, precense.GetUserId())
	if err != nil {
		return m
	}
	p := profiles[0]
	m.AvatarId = p.AvatarId
	m.Chips = p.AccountChip
	m.VipLevel = p.VipLevel
	m.DeviceID = deviceID
	return m
}

func GetDeviceIDByUserID(ctx context.Context, db *sql.DB, userID string) (string, error) {
	const query = `
		SELECT id
		FROM user_device
		WHERE user_id = $1
	`
	var deviceID string
	err := db.QueryRowContext(ctx, query, userID).Scan(&deviceID)
	fmt.Println("GetDeviceIDByUserID ", userID, deviceID, err)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Không có thiết bị nào
		}
		return "", err
	}

	return deviceID, nil
}
