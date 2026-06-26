package game

import b "KarkasonFoollery/internal/board"

type placementEffects struct {
	PointsToScore     map[b.PlayerID]int
	MeeplesToReturn   map[b.PlayerID][]b.MeepleType
	RegionsToComplete []b.RegionID
}

func newPlacementEffects() placementEffects {
	return placementEffects{PointsToScore: make(map[b.PlayerID]int),
		MeeplesToReturn: make(map[b.PlayerID][]b.MeepleType)}
}

func resolvePlacementEffects(state GameState) placementEffects {
	effects := newPlacementEffects()

	return effects
}
