package game

import b "KarkasonFoollery/internal/board"

type RegionType string

const (
	RegionCity      RegionType = "city"
	RegionRoad      RegionType = "road"
	RegionMonastery RegionType = "monastery"
)

type featureRef struct {
	Coord    b.Coord
	Index    int
	Complete bool
}

type Region struct {
	ID        b.RegionID
	Type      RegionType
	Districts []featureRef
	Owner     b.PlayerID
	Contested bool
	// Score    int
}

func MakeRegion(newCoord b.Coord, featureIndex int, featureType b.FeatureType, owner b.PlayerID) Region {
	r := Region{
		Type:      RegionType(featureType),
		Districts: []featureRef{{Coord: newCoord, Index: featureIndex}},
		Owner:     owner,
	}
	return r
}

func (r *Region) Clone() Region {
	clone := Region{
		ID:        r.ID,
		Type:      r.Type,
		Owner:     r.Owner,
		Contested: r.Contested,
	}

	if r.Districts != nil {
		clone.Districts = make([]featureRef, len(r.Districts))
		copy(clone.Districts, r.Districts)
	}

	return clone
}

func getNumberOfTilesAround(bd b.Board, newTileCoords b.Coord) int {
	var ctr int
	for _, c := range newTileCoords.GetCoordsAround() {
		_, exists := bd.GetTile(c)
		if exists {
			ctr++
		}
	}
	return ctr
}

func (r *Region) expandRegion(newTileCoords b.Coord, featureIndex int) {
	r.Districts = append(r.Districts, featureRef{Coord: newTileCoords, Index: featureIndex})
}

func (r *Region) isComplete() bool {
	for _, feature := range r.Districts {
		if !feature.Complete {
			return false
		}
	}
	return true
}
