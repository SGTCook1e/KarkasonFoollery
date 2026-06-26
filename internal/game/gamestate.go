package game

import (
	b "KarkasonFoollery/internal/board"
	"fmt"
)

type GameState struct {
	Board   b.Board
	Regions Regions

	Deck      Deck
	TopTile   b.Tile
	CurrCoord b.Coord

	PlacedMeeple b.MeepleType
	MeepleFeatId int
}

func NewState(tiles []*b.Tile) *GameState {
	s := &GameState{
		Deck:    *NewDeck(tiles),
		Regions: *NewRegions(),
		Board:   *b.NewBoard(),
	}
	s.TopTile = s.Deck.Draw()
	return s
}

func (s *GameState) Clone() GameState {
	return GameState{
		Board:     s.Board.Clone(),
		Regions:   s.Regions.Clone(),
		CurrCoord: s.CurrCoord,
	}
}

func (s *GameState) CompleteDistrict(dist FeatureRef) {
	t, _ := s.Board.GetTile(dist.Coord)
	id := t.Features[dist.Index].RegionID
	r := s.Regions.byID[id]
	for i := range r.Districts {
		if r.Districts[i].Index == dist.Index {
			r.Districts[i].Complete = true
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

func (s *GameState) ApplyPlacement(draft GameState, owner b.PlayerID) {
	s.Board = draft.Board.Clone()
	s.Regions = draft.Regions.Clone()
}

func (s *GameState) ApplyCompletion(res AnalysisResult, owner b.PlayerID) {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile on %+v does not exist!", s.CurrCoord))
	}

	for _, dir := range res.Completion.SidesToComplete {
		tile.CompleteSide(dir)
		neighbourTile, _ := s.Board.GetTile(s.CurrCoord.CoordByDirection(dir))
		neighbourTile.CompleteSide(dir.Opposite())
	}
	for _, dist := range res.Completion.DistrictsToComplete {
		s.CompleteDistrict(dist)
	}
}

func (s *GameState) ApplyRegionsPlacement(res AnalysisResult, owner b.PlayerID) {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile on %+v does not exist!", s.CurrCoord))
	}

	if tile.Monastery {
		newReg := MakeRegion(s.CurrCoord, 0, b.FeatureMonastery, true, owner)
		s.Regions.AppendRegion(newReg)
	}
	for featId, regIds := range res.RegionsByFeature {
		switch len(regIds) {
		case 0: //If 0 regions found for feature, make a new region
			newReg := MakeRegion(s.CurrCoord, featId, tile.Features[featId].Type, false, owner)
			id := s.Regions.AppendRegion(newReg)
			tile.UpdateRegionId(featId, id)
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			s.Regions.byID[regIds[0]].ExpandRegion(s.CurrCoord, featId)
			tile.UpdateRegionId(featId, regIds[0])
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			targetRegId := regIds[0]
			s.UpdateTilesRegionIds(regIds[1:], targetRegId)
			s.Regions.MergeRegions(s.CurrCoord, featId, regIds)
			tile.UpdateRegionId(featId, targetRegId)
		}
	}
}

func (s *GameState) ApplyScoring(pts map[b.PlayerID]int) {

}

func (s *GameState) ReturnMeeples(mpr map[b.PlayerID][]b.MeepleType) {

}

func (s *GameState) CompleteRegions(rtc []b.RegionID) {

}
