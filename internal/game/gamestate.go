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

func (s *GameState) makeStateDraft() *GameState {
	draft := GameState{
		Board:      s.Board.Clone(),
		Regions:    s.Regions.Clone(),
		Players:    make([]*Player, len(s.Players)),
		CurrCoord:  s.CurrCoord,
		CurrPlayer: s.CurrPlayer,
	}
	for i := range s.Players {
		draft.Players[i] = s.Players[i].Clone()
	}
	return &draft
}

func (s *GameState) ApplyPlacement(draft GameState, owner b.PlayerID) {
	s.Board = draft.Board.Clone()
	s.Regions = draft.Regions.Clone()
	s.Players = draft.Players
}

func (s *GameState) completeDistrict(dist featureRef) {
	t, exists := s.Board.GetTile(dist.Coord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
	}
	id := t.Features[dist.Index].RegionID
	r := s.Regions.ByID[id]
	for i := range r.Districts {
		if r.Districts[i].Coord == dist.Coord {
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

	for _, dir := range res.SidesToComplete {
		tile.CompleteSide(dir)

		neigbourCoord := s.CurrCoord.CoordByDirection(dir)
		neighbourTile, exists := s.Board.GetTile(neigbourCoord)
		if !exists {
			panic(fmt.Sprintf("Tile at %+v does not exist!", s.CurrCoord))
		}
		neighbourTile.CompleteSide(dir.Opposite())

		feature, index := tile.FeatureByDirection(dir)
		if feature.IsComplete() {
			currDist := featureRef{Coord: s.CurrCoord, Index: index}
			s.completeDistrict(currDist)
		}
		feature, index = neighbourTile.FeatureByDirection(dir.Opposite())
		if feature.IsComplete() {
			neighbourDist := featureRef{Coord: neigbourCoord, Index: index}
			s.completeDistrict(neighbourDist)
		}
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
			s.Regions.ByID[regIds[0]].expandRegion(s.CurrCoord, featId, owner)
			tile.UpdateRegionId(featId, regIds[0])
		default: //If more than 1 regions found for feature, unite this regions and add feature to it
			targetRegId := regIds[0]
			s.updateTilesRegionIds(regIds[1:], targetRegId)
			s.mergeRegions(s.CurrCoord, featId, regIds)
			tile.UpdateRegionId(featId, targetRegId)
		}
	}
}

func (s *GameState) applyScoring(pts map[b.PlayerID]int) {
	for i := range s.Players {
		s.Players[i].Score += pts[s.Players[i].Id]
	}
}

func (s *GameState) returnMeeples(mtr map[b.PlayerID][]b.MeepleType) {
	for i := range s.Players {
		for _, m := range mtr[s.Players[i].Id] {
			s.Players[i].Meeples[m]++
		}
	}
}

func (s *GameState) completeRegions(rtc []b.RegionID) {
	for _, id := range rtc {
		reg := s.Regions.ByID[id]
		for _, district := range reg.Districts {
			tile, _ := s.Board.GetTile(district.Coord)
			feature := &tile.Features[district.Index]

			feature.RegionID = b.NoRegion
			feature.Meeple.Owner = b.NoOwner
			feature.Meeple.Type = b.NoUnit
		}

		s.Regions.deleteRegion(id)
	}
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

func (s *GameState) mergeRegions(coords b.Coord, newFeature int, regionIds []b.RegionID) {
	targetReg := s.Regions.ByID[regionIds[0]]

	for i := 1; i < len(regionIds); i++ {
		reg := s.Regions.ByID[regionIds[i]]
		targetReg.Districts = append(targetReg.Districts, reg.Districts...)
		s.Regions.deleteRegion(regionIds[i])
	}
	targetReg.expandRegion(coords, newFeature, b.NoOwner)

	targetReg.Owner, targetReg.Contested = s.calcOwner(*targetReg)
}

func (s *GameState) calcOwner(reg Region) (b.PlayerID, bool) {
	uniqueOwners := make(map[b.PlayerID]struct{})

	for _, district := range reg.Districts {
		tile, _ := s.Board.GetTile(district.Coord)
		for i, feature := range tile.Features {
			if i != district.Index {
				continue
			}
			owner := feature.Meeple.Owner
			if owner != b.NoOwner {
				uniqueOwners[owner] = struct{}{}
			}
		}
	}

	switch len(uniqueOwners) {
	case 0:
		return b.NoOwner, false
	case 1:
		for k := range uniqueOwners {
			return k, false
		}
	}
	return b.NoOwner, true
}

func (s *GameState) getContestingOwners(reg Region) []b.PlayerID {
	var owners []b.PlayerID

	max := 0
	playersByMeeples := make(map[b.PlayerID]int)
	countedDistricts := make(map[b.Coord]struct{})
	for _, district := range reg.Districts {
		tile, _ := s.Board.GetTile(district.Coord)
		for _, feature := range tile.Features {
			if feature.RegionID != reg.ID {
				continue
			}
			if _, ok := countedDistricts[district.Coord]; ok {
				continue
			}

			owner := feature.Meeple.Owner
			if owner == b.NoOwner {
				continue
			}

			playersByMeeples[owner]++

			if max < playersByMeeples[owner] {
				max = playersByMeeples[owner]
			}
			countedDistricts[district.Coord] = struct{}{}
			continue
		}
	}

	for player, num := range playersByMeeples {
		if max == num {
			owners = append(owners, player)
		}
	}

	return owners
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
