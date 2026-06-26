package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
	"slices"
)

type analysisResult struct {
	// Map for each feature index of newTile with its neighbours regions slice
	RegionsByFeature map[int][]b.RegionID
	Completion       completion
}

type completion struct {
	SidesToComplete     []b.Direction
	DistrictsToComplete []featureRef
}

func ResolvePlacement(state GameState, owner b.PlayerID) GameState {
	var result analysisResult
	result.RegionsByFeature = analyzeRegionsPlacement(state)
	result.Completion = analyzeCompletion(state)

	draft := state.makeStateDraft()
	draft.applyRegionsPlacement(result, owner)
	draft.applyCompletion(result, owner)

	effects := resolvePlacementEffects(draft)
	draft.applyScoring(effects.PointsToScore)
	draft.returnMeeples(effects.MeeplesToReturn)
	draft.completeRegions(effects.RegionsToComplete)

	return draft
}

func analyzeCompletion(state GameState) completion {
	tile, exists := state.Board.GetTile(state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", state.CurrCoord))
	}

	var completion completion

	for dir := b.Top; dir <= b.Left; dir++ {
		feature, index := tile.FeatureByDirection(dir)
		if feature.Type != b.FeatureCity && feature.Type != b.FeatureRoad {
			continue
		}
		neighbourCoord := state.CurrCoord.CoordByDirection(dir)
		neighbourTile, exists := state.Board.GetTile(neighbourCoord)
		if !exists {
			continue
		}
		completion.SidesToComplete = append(completion.SidesToComplete, dir)
		if feature.IsOtherSidesComplete(dir) {
			dist := featureRef{Index: index, Coord: state.CurrCoord}
			completion.DistrictsToComplete = append(completion.DistrictsToComplete, dist)
		}
		neighbourFeature, i := neighbourTile.FeatureByDirection(dir.Opposite())
		if neighbourFeature.IsOtherSidesComplete(dir.Opposite()) {
			dist := featureRef{Index: i, Coord: neighbourCoord}
			completion.DistrictsToComplete = append(completion.DistrictsToComplete, dist)
		}
	}

	return completion
}

func analyzeRegionsPlacement(state GameState) map[int][]b.RegionID {
	tile, exists := state.Board.GetTile(state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile on %+v does not exist!", state.CurrCoord))
	}

	regions := make(map[int][]b.RegionID)

	for dir := b.Top; dir <= b.Left; dir++ {
		feature, index := tile.FeatureByDirection(dir)
		if feature.Type != b.FeatureCity && feature.Type != b.FeatureRoad {
			continue
		}
		if _, ok := regions[index]; !ok {
			regions[index] = nil
		}
		neighbourCoord := state.CurrCoord.CoordByDirection(dir)
		neighbourTile, exists := state.Board.GetTile(neighbourCoord)
		if !exists {
			continue
		}
		neighbourRegion, exists := findNeighbourRegionID(*neighbourTile, dir)
		if !exists {
			continue
		}
		if !slices.Contains(regions[index], neighbourRegion) {
			regions[index] = append(regions[index], neighbourRegion)
		}
	}

	return regions
}
