package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
	"slices"
)

type analysisResult struct {
	// Map for each feature index of newTile with its neighbours regions slice
	RegionsByFeature map[int][]b.RegionID
	SidesToComplete  []b.Direction
}

func ResolvePlacement(state GameState) *GameState {
	var result analysisResult
	result.RegionsByFeature = analyzeRegionsPlacement(state)
	result.SidesToComplete = analyzeCompletion(state)

	draft := state.makeStateDraft()
	draft.applyRegionsPlacement(result)
	draft.applyCompletion(result)

	effects := resolvePlacementEffects(*draft)
	draft.applyScoring(effects.PointsToScore)
	draft.returnMeeples(effects.MeeplesToReturn)
	draft.completeRegions(effects.RegionsToComplete)

	return draft
}

func analyzeCompletion(state GameState) []b.Direction {
	tile, exists := state.Board.GetTile(state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", state.CurrCoord))
	}

	var sides []b.Direction

	for dir := b.Top; dir <= b.Left; dir++ {
		feature, _ := tile.FeatureByDirection(dir)
		if feature.Type != b.FeatureCity && feature.Type != b.FeatureRoad {
			continue
		}
		neighbourCoord := state.CurrCoord.CoordByDirection(dir)
		_, exists := state.Board.GetTile(neighbourCoord)
		if !exists {
			continue
		}
		sides = append(sides, dir)
	}

	return sides
}

func analyzeRegionsPlacement(state GameState) map[int][]b.RegionID {
	tile, exists := state.Board.GetTile(state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile on %+v does not exist!", state.CurrCoord))
	}

	regions := make(map[int][]b.RegionID)

	for index, feature := range tile.Features {
		if feature.Type != b.FeatureCity &&
			feature.Type != b.FeatureRoad &&
			feature.Type != b.FeatureMonastery {
			continue
		}

		if feature.Type == b.FeatureMonastery {
			if feature.Meeple.Owner != b.NoOwner {
				regions[index] = nil
			}
			continue
		}

		for _, side := range feature.Sides {
			rotDir := side.Direction.Rotate(tile.Orientation)
			neighbourCoord := state.CurrCoord.CoordByDirection(rotDir)
			neighbourTile, exists := state.Board.GetTile(neighbourCoord)
			if !exists {
				continue
			}
			neighbourRegion, exists := findNeighbourRegionID(*neighbourTile, rotDir)
			if !exists {
				continue
			}
			if !slices.Contains(regions[index], neighbourRegion) {
				regions[index] = append(regions[index], neighbourRegion)
			}
		}
	}

	return regions
}
