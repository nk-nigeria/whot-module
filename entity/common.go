package entity

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

const (
	ModuleName        = "whot-poker"
	MaxPresenceCard   = 13
	JackpotPercentTax = 1 // 1%
)

// free game by lv
// [level]=%
// https://docs.google.com/spreadsheets/d/1OKPtCzTGe5Da-HRUKe37rS3bIGYw4F_B/edit#gid=1754766987
var feeGameByLvPercent = map[int]int{0: 7, 1: 7, 2: 6, 3: 5, 4: 5, 6: 4, 7: 4, 8: 4, 9: 4, 10: 4}

func GetFeeGameByLevel(lv int) int {
	val, exist := feeGameByLvPercent[lv]
	if exist {
		return val
	}
	return 5
}

var SnowlakeNode, _ = snowflake.NewNode(1)

type WalletAction string

const (
	WalletActionWinGameJackpot WalletAction = "win_game_jackpot"
)

func InterfaceToString(inf interface{}) string {
	if inf == nil {
		return ""
	}
	str, ok := inf.(string)
	if !ok {
		return ""
	}
	return str
}

func ToInt64(inf interface{}, def int64) int64 {
	if inf == nil {
		return def
	}
	switch v := inf.(type) {
	case int:
		return int64(inf.(int))
	case int64:
		return inf.(int64)
	case string:
		str := inf.(string)
		i, _ := strconv.ParseInt(str, 10, 64)
		return i
	case float64:
		return int64(inf.(float64))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	return def
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
