package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
)

type GameState struct {
	Deck    Deck
	Board   b.Board
	Regions Regions
}

func NewState(tiles []*b.Tile) GameState {
	return GameState{Deck: *NewDeck(tiles),
		Regions: *NewRegions(),
		Board:   *b.NewBoard()}
}

func (s *GameState) CompleteDistrict(dist FeatureRef) {
	t, _ := s.Board.GetTile(dist.Coord)
	id := t.Features[dist.Index].RegionID
	r := s.Regions.byID[id]
	for _, d := range r.Districts {
		if d.Index == dist.Index {
			d.Complete = true
			return
		}
	}
}

func (s *GameState) UpdateTilesRegionIds(ids []b.RegionID, newId b.RegionID) {
	for _, id := range ids {
		region := s.Regions.byID[id]
		for _, district := range region.Districts {
			t, _ := s.Board.GetTile(district.Coord)
			t.Features[district.Index].RegionID = newId
		}
	}
}

func (s *GameState) ApplyPlacement(res PlacementResult) {
	tile, exists := s.Board.GetTile(res.Coord)
	if !exists {
		panic(fmt.Errorf("Tile not present!"))
	}

	if tile.Monastery {
		s.Regions.AppendRegion(MakeRegion(res.Coord, 0, b.FeatureMonastery, true, res.Owner))
	}
	for featId, regIds := range res.RegionsByFeature {
		switch len(regIds) {
		case 0: //If 0 regions found for feature, make a new region
			id := s.Regions.AppendRegion(MakeRegion(res.Coord, featId, tile.Features[featId].Type, false, res.Owner))
			tile.Features[featId].RegionID = id
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			s.Regions.byID[regIds[0]].ExpandRegion(res.Coord, featId)
			tile.Features[featId].RegionID = regIds[0]
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			targetRegId := regIds[0]
			s.UpdateTilesRegionIds(regIds[1:], targetRegId)
			s.Regions.UniteRegions(res.Coord, featId, regIds)
			tile.Features[featId].RegionID = targetRegId
		}
	}

	for _, dir := range res.SidesToComplete {
		tile.CompleteSide(dir)
		neighbourTile, _ := s.Board.GetTile(res.Coord.CoordByDirection(dir))
		neighbourTile.CompleteSide(dir.Opposite())
	}
	for _, dist := range res.DistrictsToComplete {
		s.CompleteDistrict(dist)
	}
}
