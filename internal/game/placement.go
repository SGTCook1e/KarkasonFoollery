package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
	"slices"
)

type PlacementResult struct {
	// Map for each feature index of newTile with its neighbours regions slice
	RegionsByFeature    map[int][]b.RegionID
	Coord               b.Coord
	SidesToComplete     []b.Direction
	DistrictsToComplete []FeatureRef
	Owner               b.PlayerID
}

func NewPlacementResult(coord b.Coord, tile b.Tile) PlacementResult {
	return PlacementResult{Coord: coord,
		RegionsByFeature: make(map[int][]b.RegionID)}
}

func AnalyzePlacement(state GameState, newCoord b.Coord) PlacementResult {
	newTile, exists := state.Board.GetTile(newCoord)
	if !exists {
		panic(fmt.Errorf("NewTile not present!"))
	}

	result := NewPlacementResult(newCoord, *newTile)

	for dir := b.Top; dir <= b.Left; dir++ {
		feature, index := newTile.FeatureByDirection(dir)
		if feature.Type != b.FeatureCity && feature.Type != b.FeatureRoad {
			continue
		}
		if _, ok := result.RegionsByFeature[index]; !ok {
			result.RegionsByFeature[index] = nil
		}
		neighbourCoord := newCoord.CoordByDirection(dir)
		neighbourTile, exists := state.Board.GetTile(neighbourCoord)
		if !exists {
			continue
		}

		result.SidesToComplete = append(result.SidesToComplete, dir)
		if feature.IsOtherSidesComplete(dir) {
			dist := FeatureRef{Index: index, Coord: newCoord}
			result.DistrictsToComplete = append(result.DistrictsToComplete, dist)
		}
		neighbourFeature, i := neighbourTile.FeatureByDirection(dir.Opposite())
		if neighbourFeature.IsOtherSidesComplete(dir.Opposite()) {
			dist := FeatureRef{Index: i, Coord: neighbourCoord}
			result.DistrictsToComplete = append(result.DistrictsToComplete, dist)
		}

		neighbourRegion, exists := FindNeighbourRegionID(*neighbourTile, dir)
		if !exists {
			continue
		}
		if !slices.Contains(result.RegionsByFeature[index], neighbourRegion) {
			result.RegionsByFeature[index] = append(result.RegionsByFeature[index], neighbourRegion)
		}
	}

	return result
}
