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

func (s *GameState) makeStateDraft() GameState {
	return GameState{
		Board:     s.Board.Clone(),
		Regions:   s.Regions.Clone(),
		CurrCoord: s.CurrCoord,
	}
}

func (s *GameState) ApplyPlacement(draft GameState, owner b.PlayerID) {
	s.Board = draft.Board.Clone()
	s.Regions = draft.Regions.Clone()
}

func (s *GameState) completeDistrict(dist featureRef) {
	t, exists := s.Board.GetTile(dist.Coord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}
	id := t.Features[dist.Index].RegionID
	r := s.Regions.byID[id]
	for i := range r.Districts {
		if r.Districts[i].Index == dist.Index {
			r.Districts[i].Complete = true
			return
		}
	}
}

func (s *GameState) updateTilesRegionIds(ids []b.RegionID, newId b.RegionID) {
	for _, id := range ids {
		region := s.Regions.byID[id]
		for _, district := range region.Districts {
			t, _ := s.Board.GetTile(district.Coord)
			t.Features[district.Index].RegionID = newId
		}
	}
}

func (s *GameState) applyCompletion(res analysisResult, owner b.PlayerID) {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}

	for _, dir := range res.Completion.SidesToComplete {
		tile.CompleteSide(dir)
		neighbourTile, exists := s.Board.GetTile(s.CurrCoord.CoordByDirection(dir))
		if !exists {
			panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
		}
		neighbourTile.CompleteSide(dir.Opposite())
	}
	for _, dist := range res.Completion.DistrictsToComplete {
		s.completeDistrict(dist)
	}
}

func (s *GameState) applyRegionsPlacement(res analysisResult, owner b.PlayerID) {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}

	for featId, regIds := range res.RegionsByFeature {
		switch len(regIds) {
		case 0: //If 0 regions found for feature, make a new region
			newReg := MakeRegion(s.CurrCoord, featId, tile.Features[featId].Type, owner)
			id := s.Regions.addRegion(newReg)
			tile.UpdateRegionId(featId, id)
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			s.Regions.byID[regIds[0]].expandRegion(s.CurrCoord, featId)
			tile.UpdateRegionId(featId, regIds[0])
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			targetRegId := regIds[0]
			s.updateTilesRegionIds(regIds[1:], targetRegId)
			s.Regions.mergeRegions(s.CurrCoord, featId, regIds)
			tile.UpdateRegionId(featId, targetRegId)
		}
	}
}

func (s *GameState) applyScoring(pts map[b.PlayerID]int) {

}

func (s *GameState) returnMeeples(mpr map[b.PlayerID][]b.MeepleType) {

}

func (s *GameState) completeRegions(rtc []b.RegionID) {

}
