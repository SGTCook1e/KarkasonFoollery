package game

import b "KarkasonFoollery/internal/board"

type placementEffects struct {
	PointsToScore     map[b.PlayerID]int
	MeeplesToReturn   map[b.PlayerID][]b.MeepleType
	RegionsToComplete []b.RegionID
}

func newPlacementEffects() placementEffects {
	return placementEffects{PointsToScore: make(map[b.PlayerID]int),
		MeeplesToReturn: make(map[b.PlayerID][]b.MeepleType)}
}

var multiplier = map[RegionType]int{
	RegionCity: 2,
	RegionRoad: 1,
}

const scoreMonastery = 9

func resolvePlacementEffects(state GameState) placementEffects {
	effects := newPlacementEffects()

	for _, region := range state.Regions.ByID {
		if !region.isComplete(state.Board) {
			continue
		}
		if region.Type == RegionMonastery {
			effects.PointsToScore[region.Owner] += 9
			continue
		}
		effects.RegionsToComplete = append(effects.RegionsToComplete, region.ID)

		tileCtr := 0
		for _, district := range region.Districts {
			tile, _ := state.Board.GetTile(district.Coord)
			for _, feature := range tile.Features {

				if feature.RegionID != region.ID {
					continue
				}

				if shouldCountTile(region.Districts, district) {
					tileCtr++
					if feature.Shield {
						tileCtr++
					}
				}

				owner := feature.Meeple.Owner
				if owner == b.NoOwner {
					continue
				}
				effects.MeeplesToReturn[owner] = append(effects.MeeplesToReturn[owner], feature.Meeple.Type)
			}
		}

		regionScore := tileCtr * multiplier[region.Type]
		if region.Contested {
			owners := state.getContestingOwners(*region)
			for _, owner := range owners {
				effects.PointsToScore[owner] += regionScore
			}
		} else {
			effects.PointsToScore[region.Owner] += regionScore
		}
	}

	return effects
}

func shouldCountTile(districts []featureRef, district featureRef) bool {
	for _, d := range districts {
		if d.Coord == district.Coord && d.Index != district.Index {
			return false
		}
	}
	return true
}
