package game

import b "KarkasonFoollery/internal/board"

type PlacementEffects struct {
	PointsToScore     map[b.PlayerID]int
	MeeplesToReturn   map[b.PlayerID][]b.MeepleType
	RegionsToComplete []b.RegionID
}

func NewPlacementEffects() PlacementEffects {
	return PlacementEffects{PointsToScore: make(map[b.PlayerID]int),
		MeeplesToReturn: make(map[b.PlayerID][]b.MeepleType)}
}

func ResolvePlacementEffects(state GameState) PlacementEffects {
	effects := NewPlacementEffects()

	return effects
}
