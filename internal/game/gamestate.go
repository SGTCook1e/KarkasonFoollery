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

	Players      []*Player
	CurrPlayer   b.PlayerID
	NextPlayerId b.PlayerID
}

func NewState(tiles []*b.Tile) *GameState {
	s := &GameState{
		Deck:         *NewDeck(tiles),
		Regions:      *NewRegions(),
		Board:        *b.NewBoard(),
		CurrPlayer:   1,
		NextPlayerId: 1,
	}
	s.TopTile = s.Deck.Draw()
	return s
}

func (s *GameState) makeStateDraft() GameState {
	return GameState{
		Board:      s.Board.Clone(),
		Regions:    s.Regions.Clone(),
		CurrCoord:  s.CurrCoord,
		CurrPlayer: s.CurrPlayer,
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
	r := s.Regions.ByID[id]
	for i := range r.Districts {
		if r.Districts[i].Index == dist.Index {
			r.Districts[i].Complete = true
			return
		}
	}
}

func (s *GameState) updateTilesRegionIds(ids []b.RegionID, newId b.RegionID) {
	for _, id := range ids {
		region := s.Regions.ByID[id]
		for _, district := range region.Districts {
			t, _ := s.Board.GetTile(district.Coord)
			t.Features[district.Index].RegionID = newId
		}
	}
}

func (s *GameState) applyCompletion(res analysisResult) {
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

func (s *GameState) applyRegionsPlacement(res analysisResult) {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}

	for featId, regIds := range res.RegionsByFeature {
		owner := tile.Features[featId].Meeple.Owner

		switch len(regIds) {
		case 0: //If 0 regions found for feature, make a new region
			newReg := MakeRegion(s.CurrCoord, featId, tile.Features[featId].Type, owner)
			id := s.Regions.addRegion(newReg)
			tile.UpdateRegionId(featId, id)
		case 1: //If 1 regions found for feature, append this feature to existing neighbour region
			s.Regions.ByID[regIds[0]].expandRegion(s.CurrCoord, featId)
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

func (s *GameState) CanPlaceMeepleOnTile() bool {
	tile, exists := s.Board.GetTile(s.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}

	for i := range tile.Features {
		if s.IsValidMeeplePlacement(i) {
			return true
		}
	}
	return false
}

func (s *GameState) IsValidMeeplePlacement(featId int) bool {
	tile, _ := s.Board.GetTile(s.CurrCoord)
	feature := tile.Features[featId]

	if feature.Type != b.FeatureCity &&
		feature.Type != b.FeatureRoad &&
		feature.Type != b.FeatureMonastery {
		return false
	}

	for _, side := range feature.Sides {
		rotDir := side.Direction.Rotate(tile.Orientation)
		neighbourCoord := s.CurrCoord.CoordByDirection(rotDir)
		neighbourTile, exists := s.Board.GetTile(neighbourCoord)
		if !exists {
			continue
		}
		neighbourRegion, exists := findNeighbourRegionID(*neighbourTile, rotDir)
		if !exists {
			continue
		}

		reg := s.Regions.ByID[neighbourRegion]
		if reg.Contested {
			return false
		}
		if reg.Owner != s.CurrPlayer && reg.Owner != b.NoOwner {
			return false
		}
	}

	return true
}

func (s *GameState) AdvanceTurn() {
	var index int
	for i, p := range s.Players {
		if p.Id == s.CurrPlayer {
			if i+1 == len(s.Players) {
				index = 0
			} else {
				index = i + 1
			}
			break
		}
	}
	s.CurrPlayer = s.Players[index].Id
}

func (s *GameState) AddPlayer() b.PlayerID {
	p := NewPlayer(s.NextPlayerId)
	s.Players = append(s.Players, &p)
	return s.NextPlayerId + 1
}
