package game

import b "KarkasonFoollery/internal/board"

type RegionType string

const (
	RegionCity      RegionType = "city"
	RegionRoad      RegionType = "road"
	RegionMonastery RegionType = "monastery"
)

type FeatureRef struct {
	Coord    b.Coord
	Index    int
	Complete bool
}

type Region struct {
	ID        b.RegionID
	Type      RegionType
	Districts []FeatureRef
	Owner     b.OwnerID
	// Score    int
}

func MakeRegion(newCoord b.Coord, featureIndex int, featureType b.FeatureType, isMonastery bool, owner b.OwnerID) Region {
	if isMonastery {
		r := Region{
			Type:      RegionMonastery,
			Districts: []FeatureRef{{Coord: newCoord}},
			Owner:     owner,
		}
		return r
	}
	r := Region{
		Type:      RegionType(featureType),
		Districts: []FeatureRef{{Coord: newCoord, Index: featureIndex}},
		Owner:     owner,
	}
	return r
}

func GetNumberOfTilesAround(bd b.Board, newTileCoords b.Coord) int {
	var ctr int
	for _, c := range newTileCoords.GetCoordsAround() {
		_, exists := bd.GetTile(c)
		if exists {
			ctr++
		}
	}
	return ctr
}

func (r *Region) ExpandRegion(newTileCoords b.Coord, featureIndex int) {
	r.Districts = append(r.Districts, FeatureRef{Coord: newTileCoords, Index: featureIndex})
}

func (r *Region) IsComplete() bool {
	for _, feature := range r.Districts {
		if !feature.Complete {
			return false
		}
	}
	return true
}
