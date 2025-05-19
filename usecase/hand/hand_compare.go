package hand

import (
	"context"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
)

var (
	mapNaturalPoint = map[pb.HandBonusType]int{
		pb.HandBonusType_BonusNaturalCleanDragon:    60,
		pb.HandBonusType_BonusNaturalDragon:         50,
		pb.HandBonusType_BonusNaturalFullColors:     25,
		pb.HandBonusType_BonusNaturalSixPairs:       15,
		pb.HandBonusType_BonusNaturalThreeStraights: 15,
		pb.HandBonusType_BonusNaturalThreeOfFlushes: 10,
	}

	mapBonusPoint = map[pb.HandBonusType]int{
		pb.HandBonusType_BonusThreeOfAKindFrontHand: 3,
		pb.HandBonusType_BonusFullHouseMidHand:      2,
		pb.HandBonusType_BonusFourOfAKindMidHand:    8,
		pb.HandBonusType_BonusStraightFlushMidHand:  10,
		pb.HandBonusType_BonusFourOfAKindBackHand:   4,
		pb.HandBonusType_BonusStraightFlushBackHand: 5,

		pb.HandBonusType_MisSet: 6,
		pb.HandBonusType_Scoop:  3,
		// pb.HandBonusType_ScoopAll: 6,
	}

	// mapRatioScoop = map[pb.HandBonusType]int64{
	// 	// pb.HandBonusType_Scoop:    2,
	// 	pb.HandBonusType_ScoopAll: 4,
	// }
	ratioScoopApp = []int64{1, 1, 2, 3}

	baseHandPoint = 1
)

func rankingTypeToBonusType(ranking pb.HandRanking) pb.HandBonusType {
	switch ranking {
	case pb.HandRanking_NaturalCleanDragon:
		return pb.HandBonusType_BonusNaturalCleanDragon
	case pb.HandRanking_NaturalDragon:
		return pb.HandBonusType_BonusNaturalDragon
	case pb.HandRanking_NaturalFullColors:
		return pb.HandBonusType_BonusNaturalFullColors
	case pb.HandRanking_NaturalSixPairs:
		return pb.HandBonusType_BonusNaturalSixPairs
	case pb.HandRanking_NaturalThreeOfFlushes:
		return pb.HandBonusType_BonusNaturalThreeOfFlushes
	case pb.HandRanking_NaturalThreeStraights:
		return pb.HandBonusType_BonusNaturalThreeStraights
	}

	return pb.HandBonusType_None
}

func (h *Hand) CompareHand(h2 *Hand) *ComparisonResult {
	result := ComparisonResult{
		//WinType: entity.WinType_WIN_TYPE_UNSPECIFIED,
	}

	return &result
}

var kCmp = "cmp"

type CompareContext struct {
	PresenceCount  int
	ScoopAllUser   string
	ScoopAllResult *pb.ComparisonResult
	Bonuses        []*pb.HandBonus
}

func NewCompareContext(pc int) context.Context {
	return context.WithValue(context.TODO(), kCmp, &CompareContext{
		PresenceCount: pc,
	})
}

func GetCompareContext(ctx context.Context) *CompareContext {
	return ctx.Value(kCmp).(*CompareContext)
}

func CompareHand(ctx context.Context, h1, h2 *Hand) *ComparisonResult {
	// A (NA) vs
	//			B(NA) => case 1
	//			B(MS) => case 2
	//			B(NM) => case 3
	// A (MS) vs
	//			B(NA) => case 2
	//			B(MS) => case 4
	//			B(NM) => case 5
	// A (NM) vs
	//			B(NA) => case 3
	//			B(MS) => case 5
	//			B(NM) => case 6

	// validate scoop all user
	result := &ComparisonResult{}
	h1.calculatePoint()
	h2.calculatePoint()

	// case 1
	if h1.IsNatural() && h2.IsNatural() {
		compareNaturalWithNatural(h1, h2, result)
		return result
	}

	// case 2
	if h1.IsNatural() && h2.IsMisSet() {
		compareNaturalWithMisset(h1, h2, result)
		return result
	}

	if h2.IsNatural() && h1.IsMisSet() {
		compareNaturalWithMisset(h2, h1, result)
		result.swap()
		return result
	}

	// case 3
	if h1.IsNatural() && h2.IsNormal() {
		compareNaturalWithNormal(h1, h2, result)
		return result
	}

	if h2.IsNatural() && h1.IsNormal() {
		compareNaturalWithNormal(h2, h1, result)
		result.swap()
		return result
	}

	// case 4
	if h1.IsMisSet() && h2.IsMisSet() {
		compareMissetWithMisset(h1, h2, result)
		return result
	}

	// case 5
	if h1.IsNormal() && h2.IsMisSet() {
		compareNormalWithMisset(h1, h2, result)
		return result
	}

	if h2.IsNormal() && h1.IsMisSet() {
		compareNormalWithMisset(h2, h1, result)
		result.swap()
		return result
	}

	// case 6
	if h1.IsNormal() && h2.IsNormal() {
		compareNormalWithNormal(h1, h2, result)
		return result
	}

	return result
}

// compareNaturalWithNatural
// case 1
func compareNaturalWithNatural(h1, h2 *Hand, result *ComparisonResult) {
	var score = 0
	if cmp := CompareHandPoint(h1.naturalPoint, h2.naturalPoint); cmp > 0 {
		score = mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
		result.addHandBonus(h1.owner, h2.owner, rankingTypeToBonusType(h1.naturalPoint.rankingType), int64(score))
	} else if cmp < 0 {
		score = mapNaturalPoint[rankingTypeToBonusType(h2.naturalPoint.rankingType)]
		result.addHandBonus(h2.owner, h1.owner, rankingTypeToBonusType(h2.naturalPoint.rankingType), int64(score))
	}

	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

// compareNaturalWithMisset
// case2
func compareNaturalWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	score := mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
	result.addHandBonus(h1.owner, h2.owner, rankingTypeToBonusType(h1.naturalPoint.rankingType), int64(score))
	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

// compareNaturalWithNormal
// case3
func compareNaturalWithNormal(h1, h2 *Hand, result *ComparisonResult) {
	score := mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
	result.addHandBonus(h1.owner, h2.owner, rankingTypeToBonusType(h1.naturalPoint.rankingType), int64(score))
	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

// compareMissetWithMisset
// case4
func compareMissetWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	// Don't need to do anything
}

// compareNormalWithMisset
func compareNormalWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	// check special case bonus only
	if bonus, bonusScore := h1.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.FrontFactor += bonusScore
		result.r2.FrontFactor += -bonusScore

		result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
	}

	if bonus, bonusScore := h1.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.MiddleFactor += bonusScore
		result.r2.MiddleFactor += -bonusScore

		result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
	}

	if bonus, bonusScore := h1.backHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.BackFactor += bonusScore
		result.r2.BackFactor += -bonusScore

		result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
	}

	bonusMisset := mapBonusPoint[pb.HandBonusType_MisSet]
	result.r1.BonusFactor = bonusMisset
	result.r2.BonusFactor = -bonusMisset
	result.addHandBonus(h1.owner, h2.owner, pb.HandBonusType_MisSet, int64(bonusMisset))

	result.r1.Scoop = kWinMisset
	result.r2.Scoop = kLoseMisset

	// scoopScore := mapBonusPoint[pb.HandBonusType_Scoop]
	// result.addHandBonus(h1.owner, h2.owner, pb.HandBonusType_Scoop, int64(scoopScore))
}

// compareNormalWithNormal
func compareNormalWithNormal(h1, h2 *Hand, result *ComparisonResult) {
	winHand := 0
	cmpHand := 0
	// front hand
	cmpHand++
	if cmp := CompareHandPoint(h1.frontHand.Point, h2.frontHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.FrontBonusFactor = bonusScore
			result.r2.FrontBonusFactor = -bonusScore

			result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
		}

		result.r1.FrontFactor = baseHandPoint
		result.r2.FrontFactor = -baseHandPoint

		winHand++
	} else if cmp < 0 {
		if bonus, bonusScore := h2.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.FrontBonusFactor = bonusScore
			result.r1.FrontBonusFactor = -bonusScore

			result.addHandBonus(h2.owner, h1.owner, bonus, int64(bonusScore))
		}

		result.r2.FrontFactor = baseHandPoint
		result.r1.FrontFactor = -baseHandPoint

		winHand--
	}

	// middle hand
	cmpHand++
	if cmp := CompareHandPoint(h1.middleHand.Point, h2.middleHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.MiddleBonusFactor = bonusScore
			result.r2.MiddleBonusFactor = -bonusScore

			result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
		}

		result.r1.MiddleFactor = baseHandPoint
		result.r2.MiddleFactor = -baseHandPoint

		winHand++
	} else if cmp < 0 {
		if bonus, bonusScore := h2.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.MiddleBonusFactor = bonusScore
			result.r1.MiddleBonusFactor = -bonusScore

			result.addHandBonus(h2.owner, h1.owner, bonus, int64(bonusScore))
		}

		result.r2.MiddleFactor = baseHandPoint
		result.r1.MiddleFactor = -baseHandPoint

		winHand--
	}

	// backhand
	cmpHand++
	if cmp := CompareHandPoint(h1.backHand.Point, h2.backHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.backHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.BackBonusFactor = bonusScore
			result.r2.BackBonusFactor = -bonusScore

			result.addHandBonus(h1.owner, h2.owner, bonus, int64(bonusScore))
		}

		result.r1.BackFactor = baseHandPoint
		result.r2.BackFactor = -baseHandPoint

		winHand++
	} else if cmp < 0 {
		if bonus, bonusScore := h2.backHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.BackBonusFactor = bonusScore
			result.r1.MiddleBonusFactor = -bonusScore

			result.addHandBonus(h2.owner, h1.owner, bonus, int64(bonusScore))
		}

		result.r2.BackFactor = baseHandPoint
		result.r1.BackFactor = -baseHandPoint

		winHand--
	}

	scoopScore := mapBonusPoint[pb.HandBonusType_Scoop]
	if winHand == cmpHand {
		result.r1.Scoop = kWinScoop
		result.r2.Scoop = kLoseScoop

		result.addHandBonus(h1.owner, h2.owner, pb.HandBonusType_Scoop, int64(scoopScore))
	} else if -winHand == cmpHand {
		result.r2.Scoop = kWinScoop
		result.r1.Scoop = kLoseScoop

		result.addHandBonus(h2.owner, h1.owner, pb.HandBonusType_Scoop, int64(scoopScore))
	}
}

func ProcessCompareResult(ctx context.Context, cmpResult *pb.ComparisonResult, cresult Result) {
	result := cmpResult.ScoreResult
	result.FrontFactor += int64(cresult.FrontFactor)
	result.MiddleFactor += int64(cresult.MiddleFactor)
	result.BackFactor += int64(cresult.BackFactor)

	result.FrontBonusFactor += int64(cresult.FrontBonusFactor)
	result.MiddleBonusFactor += int64(cresult.MiddleBonusFactor)
	result.BackBonusFactor += int64(cresult.BackBonusFactor)

	scoopScore := mapBonusPoint[pb.HandBonusType_Scoop]
	missetScore := mapBonusPoint[pb.HandBonusType_MisSet]
	if cresult.Scoop != 0 {
		if cresult.Scoop == kWinScoop {
			result.BonusFactor += int64(scoopScore)
		} else if cresult.Scoop == kLoseScoop {
			result.BonusFactor -= int64(scoopScore)
		}
	}
	if cresult.Scoop > 0 {
		result.Scoop++
	} else if cresult.Scoop < 0 {
		result.Scoop--
	}

	if cresult.Scoop == kWinMisset {
		result.BonusFactor += int64(missetScore)
	} else if cresult.Scoop == kLoseMisset {
		result.BonusFactor -= int64(missetScore)
	}

	result.NaturalFactor += int64(cresult.NaturalFactor)

	cmpCtx := GetCompareContext(ctx)
	if cmpCtx.PresenceCount > 2 && result.Scoop >= int64(cmpCtx.PresenceCount)-1 {
		cmpCtx.ScoopAllUser = cmpResult.UserId
		cmpCtx.ScoopAllResult = cmpResult
	}

	if cresult.NaturalFactor > 0 {
		result.NumHandWin += 3
	}
	if cresult.FrontFactor > 0 {
		result.NumHandWin++
	}
	if cresult.MiddleFactor > 0 {
		result.NumHandWin++
	}
	if cresult.BackFactor > 0 {
		result.NumHandWin++
	}

}

// func ProcessCompareBonusResult(ctx context.Context, cmpResult []*pb.ComparisonResult, bonuses *[]*pb.HandBonus) {
// 	lbonuses := *bonuses
// 	cmpCtx := GetCompareContext(ctx)
// 	if cmpCtx.ScoopAllUser != "" {
// 		resultScoopAll := cmpCtx.ScoopAllResult
// 		bonus := int64(mapBonusPoint[pb.HandBonusType_ScoopAll])
// 		for _, result := range cmpResult {
// 			if result.UserId != cmpCtx.ScoopAllUser {
// 				resultScoopAll.ScoreResult.BonusFactor += bonus
// 				result.ScoreResult.BonusFactor -= bonus

// 				lbonuses = append(lbonuses, &pb.HandBonus{
// 					Win:  cmpCtx.ScoopAllUser,
// 					Lose: result.UserId,
// 					Type: pb.HandBonusType_ScoopAll,
// 				})
// 			}
// 		}
// 	}
// }

func ProcessCompareBonusResult(ctx context.Context, cmpResult []*pb.ComparisonResult, mapRCPair map[string]*ComparisonResult, bonuses *[]*pb.HandBonus) {
	cmpCtx := GetCompareContext(ctx)
	if cmpCtx.ScoopAllUser == "" {
		return
	}
	// lbonuses := *bonuses
	resultScoopAll := cmpCtx.ScoopAllResult
	// resultScoopAllUserId := cmpCtx.ScoopAllUser
	bonusRatioScoopAll := ratioScoopApp[cmpCtx.PresenceCount-1]
	// bonusRatioScoop := mapRatioScoop[pb.HandBonusType_Scoop]
	for _, result := range cmpResult {
		if result.UserId == cmpCtx.ScoopAllUser {
			continue
		}
		rc, exist := mapRCPair[result.UserId+cmpCtx.ScoopAllUser]
		if !exist {
			rc = mapRCPair[cmpCtx.ScoopAllUser+result.UserId]
		}
		r1 := rc.GetR1()
		totalFactor := int64(r1.BackBonusFactor +
			r1.BackFactor + r1.BonusFactor +
			r1.FrontBonusFactor + r1.FrontFactor +
			r1.MiddleBonusFactor + r1.MiddleFactor +
			r1.NaturalFactor)
		if totalFactor < 0 {
			totalFactor = -totalFactor
		}

		resultScoopAll.ScoreResult.BonusFactor += totalFactor * bonusRatioScoopAll
		result.ScoreResult.BonusFactor -= totalFactor * (bonusRatioScoopAll)

		*bonuses = append(*bonuses, &pb.HandBonus{
			Win:    cmpCtx.ScoopAllUser,
			Lose:   result.UserId,
			Type:   pb.HandBonusType_ScoopAll,
			Factor: totalFactor * bonusRatioScoopAll,
		})
	}
	// lbonuses = append(lbonuses, lbonuses...)
}

func CalcTotalFactor(cmpResult []*pb.ComparisonResult) {
	for _, result := range cmpResult {
		result.ScoreResult.TotalFactor = getTotalFactor(result)
	}
}

func getTotalFactor(result *pb.ComparisonResult) int64 {
	return result.ScoreResult.BackBonusFactor +
		result.ScoreResult.BackFactor + result.ScoreResult.BonusFactor +
		result.ScoreResult.FrontBonusFactor + result.ScoreResult.FrontFactor +
		result.ScoreResult.MiddleBonusFactor + result.ScoreResult.MiddleFactor +
		result.ScoreResult.NaturalFactor
}
