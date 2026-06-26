package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
	"slices"
)

type AnalysisResult struct {
	// Map for each feature index of newTile with its neighbours regions slice
	RegionsByFeature map[int][]b.RegionID
	Completion       Completion
}

type Completion struct {
	SidesToComplete     []b.Direction
	DistrictsToComplete []FeatureRef
}

func ResolvePlacement(state GameState, owner b.PlayerID) GameState {
	var result AnalysisResult
	result.RegionsByFeature = analyzeRegionsPlacement(state)
	result.Completion = analyzeCompletion(state)

	draft := state.Clone()
	draft.ApplyRegionsPlacement(result, owner)
	draft.ApplyCompletion(result, owner)

	effects := ResolvePlacementEffects(draft)
	draft.ApplyScoring(effects.PointsToScore)
	draft.ReturnMeeples(effects.MeeplesToReturn)
	draft.CompleteRegions(effects.RegionsToComplete)

	return draft
}

func analyzeCompletion(state GameState) Completion {
	tile, exists := state.Board.GetTile(state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile on %+v does not exist!", state.CurrCoord))
	}

	var completion Completion

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
			dist := FeatureRef{Index: index, Coord: state.CurrCoord}
			completion.DistrictsToComplete = append(completion.DistrictsToComplete, dist)
		}
		neighbourFeature, i := neighbourTile.FeatureByDirection(dir.Opposite())
		if neighbourFeature.IsOtherSidesComplete(dir.Opposite()) {
			dist := FeatureRef{Index: i, Coord: neighbourCoord}
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
		neighbourRegion, exists := FindNeighbourRegionID(*neighbourTile, dir)
		if !exists {
			continue
		}
		if !slices.Contains(regions[index], neighbourRegion) {
			regions[index] = append(regions[index], neighbourRegion)
		}
	}

	return regions
}
