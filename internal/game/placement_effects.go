package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
	"os"
)

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
		fmt.Fprintf(os.Stdout, "Reg %d COMPLETE\n", region.ID)
		effects.RegionsToComplete = append(effects.RegionsToComplete, region.ID)
		if region.Type == RegionMonastery {
			effects.PointsToScore[region.Owner] += 9
			continue
		}

		if region.Owner == b.NoOwner && region.Contested == false {
			continue
		}

		tileCtr := 0
		tilesToCount := make(map[b.Coord]struct{})
		for _, district := range region.Districts {
			tilesToCount[district.Coord] = struct{}{}

			tile, _ := state.Board.GetTile(district.Coord)
			for _, feature := range tile.Features {
				if feature.RegionID != region.ID {
					continue
				}

				if feature.Shield {
					tileCtr++
				}

				owner := feature.Meeple.Owner
				if owner == b.NoOwner {
					continue
				}
				effects.MeeplesToReturn[owner] = append(effects.MeeplesToReturn[owner], feature.Meeple.Type)
			}
		}

		tileCtr += len(tilesToCount)
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
